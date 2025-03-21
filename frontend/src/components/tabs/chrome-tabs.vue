<template>
  <div class="chrome-tabs">
    <div class="tabs-content" :ref="setContentRef">
      <div
        class="tabs-item"
        v-for="(tab, i) in tabs"
        :class="{ active: tab.key === modelValue }"
        :key="tab.key"
        :style="{ width: tabWidth + 'px' }"
        :ref="(e: Element | ComponentPublicInstance | null) => setTabRef(e, tab)"
        @contextmenu="(e: Event) => handleContextMenu(e, tab, i)"
        @click="(e: Event) => handleNativeClick(e, tab, i)"
      >
        <div class="tabs-background">
          <div class="tabs-background-divider"></div>
          <div class="tabs-background-content"></div>
          <svg class="tabs-background-before" width="8" height="8">
            <path d="M 0 8 A 8 8 0 0 0 8 0 L 8 8 Z"></path>
          </svg>
          <svg class="tabs-background-after" width="8" height="8">
            <path d="M 0 0 A 8 8 0 0 0 8 8 L 0 8 Z"></path>
          </svg>
        </div>
        <div class="tabs-close" @click.stop="handleDelete(tab, i)" v-show="showTabCloseIcon(tab)">
          <svg class="tabs-close-icon" width="16" height="16" stroke="#595959">
            <path d="M 4 4 L 12 12 M 12 4 L 4 12"></path>
          </svg>
        </div>
        <div class="tabs-main" :title="tab.label">
          <span class="tabs-favicon" v-if="tab.favicon">
            <render-temp v-if="typeof tab.favicon === 'function'" :render="tab.favicon" :params="[tab, i]" />
            <img v-else-if="tab.favicon" :src="tab.favicon" alt="" />
          </span>
          <span class="tabs-label" :class="{ 'no-close': !showTabCloseIcon(tab), 'no-icon': !tab.favicon }">
            <render-temp v-if="typeof renderLabel === 'function'" :render="renderLabel" :params="[tab, i]" />
            <template v-else>{{ tab.title }}</template>
          </span>
        </div>
      </div>
      <span
          class="tabs-after"
          :ref="setAfterRef"
          :style="{ left: (tabWidth - gap * 2) * tabs.length + gap * 2 + 'px' }"
      >
        <slot name="after"/>
      </span>
      <span
        class="tabs-end"
        :ref="setEndRef"
      >
        <slot name="end"/>
      </span>
    </div>
  </div>
</template>

<!--:style="{ left: (tabWidth - gap * 2) * tabs.length + gap * 2 + 'px' }"-->
<script lang="ts">
import RenderTemp from './render-temp.vue';
import Draggabilly from 'draggabilly';
import {
  defineComponent,
  ref,
  reactive,
  onMounted,
  PropType,
  nextTick,
  h,
  onUnmounted,
  VNode
} from 'vue'

import type { ComponentPublicInstance } from 'vue';

export type FaviconType = ((...args: unknown[]) => VNode) | ((...args: unknown[]) => string) | NodeRequire | string
export interface Tab {
  /** 显示名称 */
  label: string
  /** 显示标题 */
  title: string
  /** 唯一 key */
  key: string
  favicon?: FaviconType
  /**
   * 是否可关闭
   */
  closable?: boolean
  /**
   * 是否可被交换
   */
  swappable?: boolean
  /**
   * 是否可拖拽
   */
  dragable?: boolean
  $el?: HTMLElement
  // eslint-disable-next-line
  _instance?: any
  _x?: number
}

export interface Refs {
  [key: string]: Element | null
}

export default defineComponent({
  name: 'VueTabsChrome',
  components: { RenderTemp },
  emits: ['click', 'update:modelValue', 'remove', 'dragstart', 'dragging', 'dragend', 'swap', 'contextmenu'],
  props: {
    modelValue: {
      type: [String, Number],
      default: ''
    },
    tabs: {
      type: Array as PropType<Tab[]>,
      default: () => []
    },
    /**
     * 当宽度小于设置的值时，会自动隐藏关闭按钮
     */
    autoHiddenCloseIconWidth: {
      type: Number,
      default: 120
    },
    /**
     * tab 的最小宽度
     */
    minWidth: {
      type: Number,
      default: 40
    },
    /**
     * tab 的最大宽度
     */
    maxWidth: {
      type: Number,
      default: 245
    },
    /**
     * 两个相邻的 tab 的空隙大小
     */
    gap: {
      type: Number,
      default: 7
    },
    /**
     * 关闭事件
     */
    onClose: {
      type: Function
    },
    /**
     * 新 tab 追加时，是否追加到当前 tab 之后
     */
    insertToAfter: {
      type: Boolean,
      default: false
    },
    /**
     * 鼠标按下时，是否自动将命中的 tab 设置为激活状态
     */
    isMousedownActive: {
      type: Boolean,
      default: true
    },
    /**
     * 自定义渲染 label
     */
    renderLabel: {
      type: Function
    }
  },
  setup(props, context) {
    const $refs = reactive<Refs>({})
    const tabWidth = ref<number>(0)

    /**
     * 计算单个 tab 的宽度
     */
    const calcTabWidth = () => {
      const { tabs, minWidth, maxWidth, gap } = props
      const { $content } = $refs
      const endWidth = $refs.$end?.getBoundingClientRect().width || 0
      const afterWidth = $refs.$after?.getBoundingClientRect().width || 0
      if (!$content) return Math.max(maxWidth, minWidth)
      const contentWidth: number = $content.clientWidth - gap * 4 - endWidth - afterWidth - 1
      let width: number = contentWidth / tabs.length
      width += gap * 2
      if (width > maxWidth) width = maxWidth
      if (width < minWidth) width = minWidth
      tabWidth.value = width
    }

    /**
     * 拖拽开始事件
     * @param e 拖拽事件
     * @param tab 当前正在拖拽的 tab
     * @param i 当前拖拽的下标
     */
    const handlePointerDown = (e: Event, tab: Tab, i: number) => {
      const { emit } = context
      const { isMousedownActive } = props
      // 如果允许按下就 active，才命中
      if (isMousedownActive) {
        emit('update:modelValue', tab.key)
      }
      emit('dragstart', e, tab, i)
    }

    /**
     * 拖拽事件监听
     * @param e 拖拽事件
     * @param tab 当前正在拖拽的 tab
     * @param i 当前拖拽的下标
     */
    const handleDragMove = (e: Event, tab: Tab, i: number) => {
      const { tabs, gap } = props
      const { emit } = context

      if (tab.swappable === false) {
        return
      }

      // 获取一半 tab 宽度
      const halfWidth = (tabWidth.value - gap) / 2
      // 获取 tab 当前的 x 值
      const { x } = tab._instance.position
      let swapTab: Tab | null = null
      for (let i = 0; i < tabs.length; i++) {
        const currentTab: Tab = tabs[i]
        const targetX: number = (currentTab._x || 1) - 1

        // 如果命中自己本身，则无需交换
        if (tab.key === currentTab.key) {
          // eslint-disable-next-line no-continue
          continue
        }
        // 判断是否有重叠的 tab，只需要判定是否在前半部分即可
        if (targetX <= x && x < targetX + halfWidth) {
          swapTab = currentTab
          swapTabs(tab, swapTab)
          break
        }
      }
      emit('dragging', e, tab, i)
    }

    /**
     * 交换俩 tab
     * @param tab 当前 tab
     * @param swapTab 需要交换的 tab
     */
    const swapTabs = (tab: Tab, swapTab: Tab) => {
      if (swapTab.swappable === false) {
        return
      }
      const { tabs } = props
      const { emit } = context

      let index = -1
      let swapIndex = -1

      for (let i = 0; i < tabs.length; i++) {
        const obj: Tab = tabs[i]
        if (obj.key === tab.key) {
          index = i
        }
        if (obj.key === swapTab.key) {
          swapIndex = i
        }
      }

      if (index < 0 || swapIndex < 0 || index === swapIndex) {
        return
      }

      // eslint-disable-next-line
      ;[tabs[index], tabs[swapIndex]] = [tabs[swapIndex], tabs[index]]

      // swap x
      const { _x } = tab
      tab._x = swapTab._x
      swapTab._x = _x

      // swap position
      const { _instance } = swapTab
      setTimeout(() => {
        _instance.element.classList.add('move')
        _instance.setPosition(_x, _instance.position.y)
      }, 50)
      setTimeout(() => {
        _instance.element.classList.remove('move')
        emit('swap', tab, swapTab)
      }, 200)
    }

    /**
     * 拖拽完成监听
     * @param e 拖拽事件
     * @param tab 命中的 tab
     * @param i 当前拖拽的下标
     */
    const handleDragEnd = (e: Event, tab: Tab, i: number) => {
      const { _instance } = tab
      const { emit } = context

      if (_instance.position.x === 0) return
      setTimeout(() => {
        _instance.element.classList.add('move')
        _instance.setPosition(tab._x, 0)
      }, 50)
      setTimeout(() => {
        _instance.element.classList.remove('move')
        emit('dragend', e, tab, i)
      }, 200)
      return false
    }

    /**
     * 单击事件监听
     * @param e 单击事件
     * @param tab 命中的 tab
     * @param i 当前单击的下标
     */
    const handleClick = (e: Event, tab: Tab, i: number) => {
      const { emit } = context
      emit('click', e, tab, i)
    }

    /**
     * 原生点击事件
     * @param e 单击事件
     * @param tab 命中的 tab
     * @param i 当前单击的下标
     */
    const handleNativeClick = (e: Event, tab: Tab, i: number) => {
      if (tab.dragable === false) {
        handleClick(e, tab, i)
      }
    }

    /**
     * 右键事件监听
     * @param e 右键事件
     * @param tab 命中的 tab
     * @param i 当前右键的下标
     */
    const handleContextMenu = (e: Event, tab: Tab, i: number) => {
      const { emit } = context
      emit('contextmenu', e, tab, i)
    }

    /**
     * 删除事件
     * @param tab 当前命中 tab
     * @param i 当前命中 tab 的下标
     */
    const handleDelete = (tab: Tab, i: number) => {
      const { tabs, modelValue, onClose } = props
      const { emit } = context
      const index = tabs.findIndex((item) => item.key === modelValue)

      // 可以通过 onClose 返回 false 来主动阻止事件
      if (typeof onClose === 'function' && onClose(tab, tab.key, i) === false) {
        return false
      }

      let after, before
      if (i === index) {
        after = tabs[i + 1]
        before = tabs[i - 1]
      }

      if (after) {
        emit('update:modelValue', after.key)
      } else if (before) {
        emit('update:modelValue', before.key)
      } else if (tabs.length <= 1) {
        emit('update:modelValue', null)
      }
      tabs.splice(i, 1)
      // emit('update:tabs',tabs)
      emit('remove', tab, i)

      nextTick(() => {
        doLayout()
      })
    }

    /**
     * 主动添加 tab
     * @param newTabs 用户需要添加的 tab
     */
    const addTab = (...newTabs: Array<Tab>) => {
      const { insertToAfter, modelValue, tabs } = props
      if (insertToAfter) {
        const i = tabs.findIndex((tab) => tab.key === modelValue)
        tabs.splice(i + 1, 0, ...newTabs)
      } else {
        tabs.push(...newTabs)
      }
      nextTick(() => {
        init()
        doLayout()
      })
    }

    /**
     * 主动移除 tab
     * @param tabKey 如果为数字则判定为用下标删除
     */
    const removeTab = (tabKey: string | number) => {
      const { tabs } = props

      if (typeof tabKey === 'number') {
        const index: number = tabKey
        const tab = tabs[index]
        handleDelete(tab, index)
      } else {
        const index: number = tabs.findIndex((item) => item.key === tabKey)
        const tab: Tab | undefined = tabs.find((item) => item.key === tabKey)
        if (tab) {
          handleDelete(tab, index)
        }
      }
    }

    // 计时器
    let timer: number
    /**
     * 窗口改变，重新布局
     */
    const handleResize = () => {
      if (timer) window.clearTimeout(timer)
      timer = window.setTimeout(() => {
        doLayout()
      }, 100)
    }

    /**
     * 判断关闭按钮是否展示
     */
    const showTabCloseIcon = (tab: Tab) => {
      const { modelValue, autoHiddenCloseIconWidth } = props
      if (tab.closable === false) {
        return false
      }

      if (tab.key === modelValue) {
        return true
      }

      return autoHiddenCloseIconWidth <= tabWidth.value;
    }

    /**
     * 渲染文本
     */
    const renderLabelText = (tab: Tab) => {
      const { renderLabel } = props
      if (renderLabel) {
        return renderLabel(tab)
      }
      return h('span', tab.label)
    }

    /**
     * 重新调整 tab 位置
     */
    const doLayout = () => {
      calcTabWidth()
      const { tabs, gap } = props
      tabs.forEach((tab, i) => {
        const instance = tab._instance
        const _x = (tabWidth.value - gap * 2) * i
        tab._x = _x
        instance.setPosition(_x, 0)
      })
    }

    /**
     * 添加 tab 实例
     * @param tab 当前命中 tab
     * @param i 当前命中 tab 的下标
     */
    const addInstance = (tab: Tab, i: number) => {
      const { gap } = props

      // 如果已经存在实例，则重新设置位置
      if (tab._instance) {
        tab._instance.setPosition(tab._x, 0)
        return
      }
      // 如果不存在 dom 元素，则无需设置
      if (!tab.$el || !$refs.$content) {
        return
      }
      // 添加实例
      tab._instance = new Draggabilly(tab.$el, {
        axis: 'x',
        containment: $refs.$content,
        handle: '.tabs-main'
      })
      if (tab.dragable === false) {
        tab._instance.disable()
      }
      // 计算实际 x 值
      const x = (tabWidth.value - gap * 2) * i
      // 记录 x 位置到 tab 上
      tab._x = x
      // 设置位置
      tab._instance.setPosition(x, 0)
      // 绑定拖拽事件
      tab._instance.on('pointerDown', (e: Event) => handlePointerDown(e, tab, i))
      tab._instance.on('dragMove', (e: Event) => handleDragMove(e, tab, i))
      tab._instance.on('dragEnd', (e: Event) => handleDragEnd(e, tab, i))
      tab._instance.on('staticClick', (e: Event) => handleClick(e, tab, i))
    }

    /**
     * 初始化，为 tab 添加实例
     */
    const init = () => {
      props.tabs.forEach((tab: Tab, i: number) => {
        addInstance(tab, i)
      })
    }

    /**
     * 为 Tab 添加 dom 节点
     * @param el 当前 tab 对应的 dom 元素
     * @param tab 当前命中 tab
     */
    const setTabRef = (el: Element | ComponentPublicInstance | null, tab: Tab) => {
      if (el) {
        tab.$el = el as HTMLElement
      }
    }

    /**
     * 添加容器 dom 节点
     * @param el tab 对应的 dom 父元素
     */
    const setContentRef = (el: Element | ComponentPublicInstance | null) => {
      if (el) {
        $refs.$content = el as Element
      }
    }

    /**
     * 添加后缀元素 dom 节点
     * @param el 在 tab 后面的元素
     */
    const setAfterRef = (el: Element | ComponentPublicInstance | null) => {
      if (el) {
        $refs.$after = el as Element
      }
    }

    /**
     * 添加末尾元素 dom 节点
     * @param el 在 tab 后面的元素
     */
    const setEndRef = (el: Element | ComponentPublicInstance | null) => {
      if (el) {
        $refs.$end = el as Element
      }
    }

    onMounted(() => {
      calcTabWidth()
      init()
      window.addEventListener('resize', handleResize)
    })

    onUnmounted(() => {
      window.removeEventListener('resize', handleResize)
      if (timer) window.clearTimeout(timer)
    })

    return {
      setTabRef,
      setContentRef,
      setAfterRef,
      setEndRef,
      tabWidth,
      handleDelete,
      handleContextMenu,
      showTabCloseIcon,
      renderLabelText,
      handleNativeClick,
      doLayout,
      addTab,
      removeTab
    }
  }
})
</script>

<style scoped lang="less">
.chrome-tabs {
  @bg: #303335;
  @gap: 7px;
  @divider: #1f1f1f;
  @speed: 150ms;
  --wails-draggable:drag;
  cursor: default;

  padding-top: 7px;
  background-color: @bg;
  position: relative;
  color: rgba(255, 255, 255, 0.85);
  font-weight: bold;

  .tabs-content {
    height: 34px;
    position: relative;
    overflow: hidden;
  }

  /* divider */
  .tabs-divider {
    left: 0;
    top: 50%;
    width: 1px;
    height: 20px;
    background-color: @divider;
    position: absolute;
    transform: translateY(-50%);
  }

  .tabs-item {
    height: 100%;
    display: flex;
    align-items: center;
    user-select: none;
    box-sizing: border-box;
    transition: width @speed;
    position: absolute;

    &:hover {
      z-index: 1;

      .tabs-background-divider {
        display: none;
      }

      .tabs-background-content {
        //background-color: #1A1B1E;
        background-color: #202225;
      }

      .tabs-background-before,
      .tabs-background-after {
        //fill: #1A1B1E;
        fill: #202225;
      }
    }

    &.move {
      transition: @speed;
    }

    &.is-dragging {
      z-index: 3;

      .tabs-background-content {
        background-color: #1A1B1E;
      }

      .tabs-background-divider {
        display: none;
      }

      .tabs-background-before,
      .tabs-background-after {
        fill: #1A1B1E;
      }
    }

    &.active {
      z-index: 2;

      .tabs-close {
        background-color: #1A1B1E;
        cursor: pointer;
      }

      .tabs-background-divider {
        display: none;
      }

      .tabs-background-content {
        background-color: #1A1B1E;
      }

      .tabs-background-before,
      .tabs-background-after {
        fill: #1A1B1E;
      }
    }

    &:first-of-type {
      .tabs-background-divider::before {
        display: none;
      }
    }
  }

  .tabs-main {
    height: 100%;
    left: 0;
    right: 0;
    margin: 0 @gap * 2;
    border-top-left-radius: .5rem;
    border-top-right-radius: .5rem;
    transition: @speed;
    display: flex;
    align-items: center;
    position: absolute;
    box-sizing: border-box;
    overflow: hidden;
  }

  .tabs-close {
    top: 50%;
    right: @gap * 2;
    width: 16px;
    height: 16px;
    z-index: 1;
    position: absolute;
    transform: translateY(-50%);
  }

  .tabs-close-icon {
    width: 100%;
    height: 100%;
    border-radius: 50%;

    &:hover {
      stroke: white;
      background-color: rgba(214, 200, 200, 0.3);
      border-radius: 20%;
    }
  }

  .tabs-favicon {
    height: 16px;
    margin-left: 3%;
    display: flex;
    align-items: center;
    overflow: hidden;

    img {
      height: 100%;
    }
  }

  .tabs-label {
    flex: 1;
    margin-left: 5%;
    margin-right: 16px;
    box-sizing: border-box;
    overflow: hidden;
    white-space: nowrap;
    position: relative;

    &.no-close {
      margin-right: 0;
    }

    &.no-icon {
      margin-left: 0;
    }
  }

  .tabs-background {
    width: 100%;
    height: 100%;
    padding: 0 @gap - 1px;
    position: absolute;
    box-sizing: border-box;
  }

  .tabs-background-divider {
    left: 0;
    width: calc(100% - 14px);
    height: 100%;
    margin: 0 7px;
    position: absolute;

    &::before {
      content: '';
      top: 20%;
      right: 100%;
      width: 1px;
      height: 60%;
      background-color: #81878c;
      position: absolute;
    }

    &::after {
      content: '';
      top: 20%;
      left: calc(100% - 1px);
      width: 1px;
      height: 60%;
      background-color: #81878c;
      position: absolute;
    }
  }

  .tabs-background-content {
    height: 100%;
    border-top-left-radius: .5rem;
    border-top-right-radius: .5rem;
    transition: @speed;
  }

  .tabs-background-before,
  .tabs-background-after {
    bottom: 0;
    position: absolute;
    fill: transparent;
    transition: @speed;
  }

  .tabs-background-before {
    left: -2px;
  }

  .tabs-background-after {
    right: -2px;
  }

  .tabs-footer {
    height: 4px;
    background-color: #fff;
  }

  .tabs-after {
    top: 50%;
    display: flex;
    position: absolute;
    overflow: hidden;
    transform: translateY(-50%);
  }

  .tabs-end {
    right: 14px;
    top: 50%;
    display: flex;
    position: absolute;
    overflow: hidden;
    transform: translateY(-50%);
  }

  @keyframes tab-show {
    from {
      transform: scaleX(0);
    }
    to {
      transform: scaleX(1);
    }
  }
}
</style>

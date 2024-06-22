<template>
  <div ref="resizeDivRefDom" class="resize-div">
    <slot></slot>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, onBeforeUnmount } from 'vue';

interface DomInfo {
  width: number;
  height: number;
}

const props = defineProps({
  // 防抖时间
  t: {
    type: Number,
    default: 50
  },
  // 初始第一次是否触发 resize 事件
  firstTrigger: {
    type: Boolean,
    default: true
  }
});

const emit = defineEmits<{
  (event: 'resize', payload: {
    resizeBox: HTMLElement;
    oldDomInfo: DomInfo;
    newDomInfo: DomInfo;
  }): void;
}>();

const resizeDivRefDom = ref<HTMLElement | null>(null);
const resizeObserver = ref<ResizeObserver | null>(null);
const tickTimer = ref<number | null>(null);
const firstTime = ref(true);
const domWh = reactive<DomInfo>({
  width: 0,
  height: 0
});

function getDomInfo(dom: HTMLElement): DomInfo {
  return {
    width: dom.clientWidth,
    height: dom.clientHeight
  };
}

function emitEvent() {
  const resizeBox = resizeDivRefDom.value;
  if (!resizeBox) return;

  const oldDomInfo = { ...domWh };
  const newDomInfo = getDomInfo(resizeBox);
  domWh.width = newDomInfo.width;
  domWh.height = newDomInfo.height;

  emit('resize', {
    resizeBox,
    oldDomInfo,
    newDomInfo
  });
}

function eventResize() {
  if (tickTimer.value) clearTimeout(tickTimer.value);
  tickTimer.value = setTimeout(() => {
    const resizeBox = resizeDivRefDom.value;
    if (!resizeBox) return;

    if (firstTime.value && !props.firstTrigger) {
      const initialDomInfo = getDomInfo(resizeBox);
      domWh.width = initialDomInfo.width;
      domWh.height = initialDomInfo.height;
    } else {
      emitEvent();
    }
    firstTime.value = false;
  }, props.t) as unknown as number; // TypeScript 类型断言
}

function addEvents() {
  removeEvents();
  const resizeBox = resizeDivRefDom.value;
  if (!resizeBox) return;

  resizeObserver.value = new ResizeObserver(eventResize);
  resizeObserver.value.observe(resizeBox);
}

function removeEvents() {
  const resizeBox = resizeDivRefDom.value;
  if (resizeBox && resizeObserver.value) {
    resizeObserver.value.unobserve(resizeBox);
  }
  if (tickTimer.value) clearTimeout(tickTimer.value);
}

onMounted(() => {
  addEvents();
});

onBeforeUnmount(() => {
  removeEvents();
});
</script>

<style scoped>
.resize-div {
  width: 100%;
  height: 100%;
}
</style>

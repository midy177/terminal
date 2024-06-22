<script lang="ts" setup>
import TerminalTabs, {Tab} from "./components/tabs/chrome-tabs.vue";
import {ComponentPublicInstance, nextTick, onMounted, onUnmounted, reactive, ref, watch} from 'vue';
import Terminal from "./components/terminal/terminal.vue";
import {nanoid} from "nanoid";
const tabRef = ref();
const fileBrowserRef = ref();
const  terminalLayoutRef = ref();
import Dropdown  from "./components/dropdown/dropdown.vue";
import {ClosePty, CreateLocalPty, CreateSshPty} from "../wailsjs/go/logic/Logic";
import Hosts from "./components/hosts/hosts.vue";
import {termx} from "../wailsjs/go/models";
import More from "./components/more/more.vue";
import FileBrowser from "./components/terminal/file_browser.vue";
import zhCN from "ant-design-vue/es/locale/zh_CN";
import {ConfigProvider, theme, Space, Spin, notification, message} from "ant-design-vue";
import {EventsOff} from "../wailsjs/runtime";

const state = reactive({
  tabs: <Array<Tab>>[],
  tab: '',
  termRefMap: new Map<string, Element | ComponentPublicInstance | null>(),
  resizeObserver: <ResizeObserver | null>null,
  tickTimer: <number | null>null,
})

function addLocalTab(data: termx.SystemShell) {
  const hide = message.loading('打开本地终端中...', 0);
  let key = nanoid()
  data.id = key
  CreateLocalPty(data).then(()=>{
    let newTab = {
      label: data.name,
      key: key,
    }
    tabRef.value.addTab(newTab)
    state.tab = key
  }).catch(e=>{
    notification.error({
      message: '创建本地终端失败',
      description: e,
      duration: null
    })
  }).finally(hide)
}

function handleOpenSshTerminal(id:number,label:string){
  const hide = message.loading('连接到ssh服务器中...', 0);
  let tid = nanoid()
  CreateSshPty(tid,id,70,40).then(()=>{
    let newTab = {
      label: label,
      key: tid,
    }
    tabRef.value.addTab(newTab)
    state.tab = tid
  }).catch(e=>{
    notification.error({
      message: '创建ssh连接失败',
      description: e,
      duration: null
    })
  }).finally(hide)
}
function openFileBrowser(){
  if (fileBrowserRef.value) {
    fileBrowserRef.value.openModel()
  }
}

function closePty(tab: Tab,key: string,i: number){
  ClosePty(key).then().catch();
  EventsOff(key);
  state.termRefMap.delete(key);
}
function setTerminalRef(tabKey: string,el: Element | ComponentPublicInstance | null) {
  if (el) {
    state.termRefMap.set(tabKey,el);
  }
}

function resizeTerminal(newKey: string) {
  let newRef = state.termRefMap.get(newKey)
  if (newRef) {
    (newRef as InstanceType<typeof Terminal>).fitWithHeightWidth()
  }
}

function resizeHandle() {
  let currTermRef = state.termRefMap.get(state.tab)
  if (currTermRef) {
    (currTermRef as InstanceType<typeof Terminal>).fitWithHeightWidth()
  }
}

function eventResize() {
  if (state.tickTimer) clearTimeout(state.tickTimer);
  state.tickTimer = setTimeout(() => {
    const resizeBox = terminalLayoutRef.value;
    if (!resizeBox) return;
    // console.log('layout',resizeBox.clientWidth,resizeBox.clientHeight)
    // 计算行和列数
    resizeHandle()
  }, 100) as unknown as number; // TypeScript 类型断言
}

function addEvents() {
  removeEvents();
  const resizeBox = terminalLayoutRef.value;
  if (!resizeBox) return;
  // addEventListener('resize',eventResize)
  state.resizeObserver = new ResizeObserver(eventResize);
  state.resizeObserver.observe(resizeBox);
}

function removeEvents() {
  // removeEventListener('resize',eventResize)
  const resizeBox = terminalLayoutRef.value;
  if (resizeBox && state.resizeObserver) {
    state.resizeObserver.unobserve(resizeBox);
  }
  if (state.tickTimer) clearTimeout(state.tickTimer);
}

onMounted(()=>{
  nextTick(()=>{
    message.config({top: '40%'});
    addEvents();
  })
})
onUnmounted(()=>{
  removeEvents();
})

watch(
    () => state.tab,
    (newVal, oldVal) => {
      eventResize();
    },
    // {immediate: true}
);

const shouldShowTerminal = (key: string) => {
  return key === state.tab
};
</script>
<template>
  <ConfigProvider
      :locale="zhCN"
      :theme="{
              algorithm: theme.darkAlgorithm,
              }"
  >
  <terminal-tabs
      ref="tabRef"
      :tabs="state.tabs"
      v-model="state.tab"
      :on-close="closePty"
  >
    <template v-slot:after>
      <Space :size="1">
        <dropdown :at-click="addLocalTab"/>
      </Space>
    </template>
    <template v-slot:end>
        <Space :size="1">
<!--        <dropdown :at-click="addLocalTab"/>-->
          <hosts :open-ssh-terminal="handleOpenSshTerminal"/>
          <more :file-browser="openFileBrowser"/>
        </Space>
    </template>
  </terminal-tabs>
      <div ref="terminalLayoutRef" class="terminal-layout">
        <TransitionGroup name="fade">
          <terminal
              v-for="item in state.tabs"
              :key="item.key"
              :id="item.key"
              v-show="shouldShowTerminal(item.key)"
              v-model:title="item.label"
              :ref="(el: Element | ComponentPublicInstance | null)=> setTerminalRef(item.key,el)"
          />
        </TransitionGroup>
      </div>
  <FileBrowser ref="fileBrowserRef" :tid="state.tab"></FileBrowser>
  </ConfigProvider>
</template>

<style lang="less">
input[type=search]::-webkit-search-cancel-button{
  -webkit-appearance: none;
}

.header-btn-bar {
  width: 34px;
  border-radius: 5px;
  padding: 0;
  cursor: pointer;
  display: flex;
  justify-content: center;
  align-items: center;
  transition: background 300ms;
  height: 34px;
  line-height: 34px;
  &:hover {
    background-color: rgba(0, 0, 0, .1);
  }
}
.terminal-layout {
  background-color: transparent;
  height: 100%;
  width: 100%;
  max-height: 100%;
  max-width: 100%;
  display: flex;
  flex: 1;
  overflow: hidden; /* 防止子元素内容撑开 */
  justify-content: center; /* 水平居中对齐内容 */
  align-items: center; /* 垂直居中对齐内容 */
}

.fade-box {
  padding: 10px;
  margin: 5px 0;
  background-color: lightblue;
  transition: opacity 0.5s ease; /* 过渡效果设置 */
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.5s ease; /* 过渡效果设置 */
}

.fade-enter, .fade-leave-to {
  opacity: 0; /* 初始和结束状态的透明度 */
}
</style>

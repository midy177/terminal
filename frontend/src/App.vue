<script lang="ts" setup>
import TerminalTabs, {Tab} from "./components/tabs/chrome-tabs.vue";
import {ComponentPublicInstance, nextTick, onMounted, onUnmounted, reactive, ref, watch} from 'vue';
import Terminal from "./components/terminal/terminal.vue";
import {nanoid} from "nanoid";
const tabRef = ref();
const fileBrowserRef = ref();
const  terminalLayoutRef = ref();
import Dropdown  from "./components/dropdown/dropdown.vue";
import {ClosePty, CreateLocalPty, CreateSshPty, CreateSshPtyWithJumper, ResizePty} from "../wailsjs/go/logic/Logic";
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
  width: 0,
  height: 0,
})

function addLocalTab(data: termx.SystemShell) {
  const hide = message.loading('打开本地终端中...', 0);
  const tid = nanoid()
  data.id = tid
  CreateLocalPty(data).then(()=>{
    let newTab = {
      label: data.name,
      title: data.name,
      key: tid,
    }
    tabRef.value.addTab(newTab);
    state.tab = tid;
    resizeHandle(tid);
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
  const tid = nanoid()
  CreateSshPty(tid, id,22,60).then((e)=>{
    let newTab = {
      label: label,
      title: label,
      key: tid,
    }
    tabRef.value.addTab(newTab);
    state.tab = tid;
    resizeHandle(tid);
  }).catch(e=>{
    notification.error({
      message: '创建ssh连接失败',
      description: e,
      duration: null
    })
  }).finally(hide)
}

function handleOpenSshTerminalWithJumper(tid:number,jid:number,label:string){
  const hide = message.loading('连接到ssh服务器中...', 0);
  const id = nanoid()
  CreateSshPtyWithJumper(id, tid, jid,22,60).then((e)=>{
    let newTab = {
      label: label,
      title: label,
      key: id,
    }
    tabRef.value.addTab(newTab);
    state.tab = id;
    resizeHandle(id);
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
  ClosePty(key).finally(()=>{
    EventsOff(key);
    state.termRefMap.delete(key);
  });
}
function setTerminalRef(tabKey: string,el: Element | ComponentPublicInstance | null) {
  if (el) {
    state.termRefMap.set(tabKey,el);
  }
}

function resizeHandle(tabKey?: string) {
  setTimeout(() => {
    if (!tabKey) tabKey = state.tab;
    const currTermRef = state.termRefMap.get(tabKey)
    if (currTermRef) {
      (currTermRef as InstanceType<typeof Terminal>).fitWithHeightWidth(state.width, state.height)
    }
  }, 100)
}

function focusHandle(tabKey?: string) {
  setTimeout(() => {
    if (!tabKey) tabKey = state.tab;
    const currTermRef = state.termRefMap.get(tabKey)
    if (currTermRef) {
      (currTermRef as InstanceType<typeof Terminal>).focusTerminal();
    }
  }, 100)
}

function eventResize() {
  if (state.tickTimer) clearTimeout(state.tickTimer);
  state.tickTimer = setTimeout(() => {
    const resizeBox = terminalLayoutRef.value;
    if (!resizeBox) return;
    const currentWidth = resizeBox.clientWidth-24;
    const currentHeight = resizeBox.clientHeight-16;
    if (currentWidth <= 0 && currentHeight <=0) return;
    if (state.width == currentWidth && state.height == currentHeight) return;
    state.width = currentWidth;
    state.height = currentHeight;
    // 计算行和列数
    resizeHandle();
  }, 100) as unknown as number; // TypeScript 类型断言
}

function addEvents() {
  removeEvents();
  const resizeBox = terminalLayoutRef.value;
  if (!resizeBox) return;
  state.resizeObserver = new ResizeObserver(eventResize);
  state.resizeObserver.observe(resizeBox);
}

function removeEvents() {
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
      resizeHandle(newVal);
      focusHandle(newVal);
    }
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
          <hosts
              :open-ssh-terminal="handleOpenSshTerminal"
              :open-ssh-terminal-with-jumper="handleOpenSshTerminalWithJumper"
          />
          <more :file-browser="openFileBrowser"/>
        </Space>
    </template>
  </terminal-tabs>
      <div ref="terminalLayoutRef" :style="{backgroundColor: state.tabs.length>0 ? 'rgb(26, 27, 30)': 'transparent'}" class="terminal-layout">
          <terminal
              v-for="item in state.tabs"
              :key="item.key"
              :id="item.key"
              v-show="shouldShowTerminal(item.key)"
              v-model:title="item.title"
              :ref="(el: Element | ComponentPublicInstance | null)=> setTerminalRef(item.key,el)"
              :width="state.width"
              :height="state.height"
          />
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
  //background-color: transparent;
  height: 100%;
  width: 100%;
  max-height: 100%;
  max-width: 100%;
  display: flex;
  flex: 1;
  overflow: hidden; /* 防止子元素内容撑开 */
  justify-content: center; /* 水平居中对齐内容 */
  align-items: center; /* 垂直居中对齐内容 */
  padding: 8px;
  //background-color: rgb(26, 27, 30);
}
</style>

<script lang="ts" setup>
import TerminalTabs, {Tab} from "./components/tabs/chrome-tabs.vue";
import {ComponentPublicInstance, nextTick, onMounted, onUnmounted, reactive, ref, watch} from 'vue';
import Terminal from "./components/terminal/terminal.vue";
import {nanoid} from "nanoid";
const tabRef = ref();
const fileBrowserRef = ref();
const asciiPlayerRef = ref();
const  terminalLayoutRef = ref();
import Dropdown  from "./components/dropdown/dropdown.vue";
import {
  ClosePty,
  CreateLocalPty,
  CreateSshPty,
  CreateSshPtyWithJumper,
  ResizePty
} from "../wailsjs/go/logic/Logic";
import Hosts from "./components/hosts/hosts.vue";
import {termx} from "../wailsjs/go/models";
import More from "./components/more/more.vue";
import FileBrowser from "./components/terminal/file_browser.vue";
import zhCN from "ant-design-vue/es/locale/zh_CN";
import {ConfigProvider, theme, Space, Spin, notification, message} from "ant-design-vue";
import {EventsOff} from "../wailsjs/runtime";
import AsciiPlayer from "./components/terminal/AsciiPlayer.vue";

const state = reactive({
  tabs: <Array<Tab>>[],
  tab: '',
  termRefMap: new Map<string, Element | ComponentPublicInstance | null>(),
  resizeObserver: <ResizeObserver | null>null,
  tickTimer: <number | null>null,
  width: 0,
  height: 0,
  enableMarquee: false,
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

function recordingCurrentTab() {
  if (state.tab) {
    const currTermRef = state.termRefMap.get(state.tab)
    if (currTermRef) {
      (currTermRef as InstanceType<typeof Terminal>).startRecording();
    }
  } else {
    notification.error({
      message: '当前没有打开的终端',
      duration: null
    })
  }
}

function recordingReview() {
  if (asciiPlayerRef.value) {
    asciiPlayerRef.value.openModel()
  }
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
          <more
              :file-browser="openFileBrowser"
              :start-recording="recordingCurrentTab"
              :record-review="recordingReview"/>
        </Space>
    </template>
  </terminal-tabs>
    <div
        ref="terminalLayoutRef"
        :style="{backgroundColor: state.tabs.length>0 ? 'rgb(26, 27, 30)': 'transparent'}"
        class="terminal-layout"
    >
      <div class="terminal-stack">
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
    </div>
    <FileBrowser ref="fileBrowserRef" :tid="state.tab"/>
    <AsciiPlayer ref="asciiPlayerRef" />
  </ConfigProvider>
</template>

<style lang="less">
input[type=search]::-webkit-search-cancel-button{
  -webkit-appearance: none;
}

.terminal-layout {
  height: 100%;
  width: 100%;
  max-height: 100%;
  max-width: 100%;
  display: flex;
  flex: 1;
  overflow: hidden;
  justify-content: center;
  align-items: center;
  padding: 8px;
}

.terminal-stack {
  position: relative;
  width: 100%;
  height: 100%;
}
</style>

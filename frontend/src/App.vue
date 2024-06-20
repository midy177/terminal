<script lang="ts" setup>
import TerminalTabs, {Tab} from "./components/tabs/chrome-tabs.vue";
import {ComponentPublicInstance, onMounted, onUnmounted, reactive, ref, watch} from 'vue';
import Terminal from "./components/terminal/terminal.vue";
import {nanoid} from "nanoid";
const tabRef = ref();
const fileBrowserRef = ref();
import Dropdown  from "./components/dropdown/dropdown.vue";
import {ClosePty, CreateLocalPty, CreateSshPty} from "../wailsjs/go/logic/Logic";
import Hosts from "./components/hosts/hosts.vue";
import {logic, termx} from "../wailsjs/go/models";
import {NotificationService, LoadingService, Message} from "vue-devui";
import More from "./components/more/more.vue";
import FileBrowser from "./components/terminal/file_browser.vue";
import zhCN from "ant-design-vue/es/locale/zh_CN";
import { ConfigProvider, theme, Space } from "ant-design-vue";
import {EventsOff} from "../wailsjs/runtime";

const state = reactive({
  tabs: <Array<Tab>>[],
  tab: '',
  loading: <any>null,
  termRefMap: new Map<string, Element | ComponentPublicInstance | null>()
})

function addLocalTab(data: termx.SystemShell) {
  state.loading = LoadingService.open({
    message: '打开本地终端中...',
  })
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
    NotificationService.open({
      type: 'error',
      title: '创建本地终端失败',
      content: e,
      duration: 5000,
    })
  }).finally(closeLoading)
}

function handleOpenSshTerminal(id:number,label:string){
  state.loading = LoadingService.open({
    message: '连接到ssh服务器中...',
  })
  let tid = nanoid()
  CreateSshPty(tid,id,70,40).then(()=>{
    let newTab = {
      label: label,
      key: tid,
    }
    tabRef.value.addTab(newTab)
    state.tab = tid
  }).catch(e=>{
    NotificationService.open({
      type: 'error',
      title: '创建ssh连接失败',
      content: e,
      duration: 5000,
    })
  }).finally(closeLoading)
}
function openFileBrowser(){
  if (fileBrowserRef.value) {
    fileBrowserRef.value.openModel()
  }
}
function closeLoading() {
  setTimeout(()=>{
    state.loading?.loadingInstance?.close()
    state.loading = null
  },10)
}
function closePty(tab: Tab,key: string,i: number){
  ClosePty(key).then().catch(e=>{
    console.log(e)
  }).finally(()=>{
    EventsOff(key);
    state.termRefMap.delete(key);
  });
}
function setTerminalRef(tabKey: string,el: Element | ComponentPublicInstance | null) {
  if (el) {
    state.termRefMap.set(tabKey,el);
  }
}

function resizeTerminal(newKey: string) {
  let newRef = state.termRefMap.get(newKey)
  if (newRef) {
    (newRef as InstanceType<typeof Terminal>).fitTerminal()
  }
}

function resizeHandle() {
  let currTermRef = state.termRefMap.get(state.tab)
  if (currTermRef) {
    (currTermRef as InstanceType<typeof Terminal>).fitTerminal()
  }
}

onMounted(()=>{
  addEventListener("resize", resizeHandle);
})
onUnmounted(()=>{
  removeEventListener("resize", resizeHandle);
})

watch(
    () => state.tab,
    (newVal, oldVal) => {
      setTimeout(() => resizeTerminal(newVal),200)
      // resizeTerminal(newVal)
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
          <hosts :open-ssh-terminal="handleOpenSshTerminal"/>
          <more :file-browser="openFileBrowser"/>
        </Space>
    </template>
  </terminal-tabs>
  <div class="terminal-layout" v-if="state.tabs.length>0">
    <template v-for="item in state.tabs" :key="item.key">
      <terminal
          :id="item.key"
          v-show="shouldShowTerminal(item.key)"
          v-model:title="item.label"
          :ref="(el: Element | ComponentPublicInstance | null)=> setTerminalRef(item.key,el)"
      />
    </template>
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
  padding-right: .3rem;
  padding-left: .3rem;
  background-color: #1A1B1E;
  height: 100%;
  width: 100%;
  flex: 1;
}
</style>

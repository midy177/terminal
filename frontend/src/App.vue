<template>
  <terminal-tabs ref="tabRef" style="--wails-draggable:drag" :tabs="tabs" v-model="tab">
    <template v-slot:after>
      <span class="header-btn-bar">
        <dropdown :at-click="addLocalTab"/>
      </span>
      <span class="header-btn-bar">
        <hosts :open-ssh-terminal="handleOpenSshTerminal"/>
      </span>
      <span class="header-btn-bar">
        <more/>
      </span>
    </template>
  </terminal-tabs>
  <div class="terminal-layout" v-if="tabs.length>0">
    <template v-for="item in tabs" :key="item.key">
      <terminal :id="item.key" :item="item" v-show="item.key === tab" v-model:title="item.label"/>
    </template>
  </div>
</template>

<script lang="ts" setup>
import TerminalTabs, {Tab} from "./components/tabs/chrome-tabs.vue";
import {reactive, ref} from 'vue';
import Terminal from "./components/terminal/terminal.vue";
import {nanoid} from "nanoid";
const tab = ref('')
const tabRef = ref()
import Dropdown  from "./components/dropdown/dropdown.vue";
import {CreateLocalPty, CreateSshPty} from "../wailsjs/go/logic/Logic";
import Hosts from "./components/hosts/hosts.vue";
import {logic, termx} from "../wailsjs/go/models";
import {NotificationService} from "vue-devui";
import More from "./components/more/more.vue";
const tabs = <Array<Tab>>reactive([])

function addLocalTab(data: termx.SystemShell) {
  let key = nanoid()
  data.id = key
  CreateLocalPty(data).then(res=>{
    let newTab = {
      label: data.name,
      key: key,
    }
    // tabs.push(newTab)
    tabRef.value.addTab(newTab)
    tab.value = key
  })
}

function handleOpenSshTerminal(id:number,label:string){
  let tid = nanoid()
  CreateSshPty(tid,id,70,40).then(()=>{
    let newTab = {
      label: label,
      key: tid,
    }
    tabRef.value.addTab(newTab)
    tab.value = tid
  }).catch(e=>{
    NotificationService.open({
      type: 'error',
      title: '创建ssh连接失败',
      content: e,
      duration: 3000,
    })
  })

}
</script>

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

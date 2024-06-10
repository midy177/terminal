<template>
  <terminal-tabs ref="tabRef" style="--wails-draggable:drag" :tabs="tabs" v-model="tab">
    <template v-slot:after>
      <span class="header-btn-bar">
        <dropdown :at-click="addLocalTab"/>
<!--          <i class="iconfont icon-add"></i>-->
      </span>
      <span class="header-btn-bar">
        <hosts/>
      </span>
      <span class="header-btn-bar" @click="handleMore">
            <i class="icon-more-operate"></i>
          </span>
    </template>
  </terminal-tabs>
  <div class="terminal-layout" v-if="tabs.length>0">
    <template v-for="item in tabs">
      <terminal :id="item.key" :item="item" v-show="item.key === tab" v-model:title="item.label"></terminal>
    </template>
  </div>
</template>

<script lang="ts" setup>
import TerminalTabs, {Tab} from "./components/tabs/chrome-tabs.vue";
import {onMounted, reactive, ref} from 'vue';
import Terminal from "./components/terminal/terminal.vue";
import {nanoid} from "nanoid";
const tab = ref('google')
const tabRef = ref()
import {termx} from "./types/models";
import Dropdown  from "./components/dropdown/dropdown.vue";
import {CreateLocalPty} from "../wailsjs/go/main/App";
import Hosts from "./components/hosts/hosts.vue";
const tabs = <Array<Tab>>reactive([])

function addLocalTab(data: termx.SystemShell) {
  let key = nanoid()
  data.ID = key
  CreateLocalPty(data).then(res=>{
    console.log(res)
    let newTab = {
      label: data.Name,
      key: key,
    }
    tabRef.value.addTab(newTab)
    tab.value = key
  })
}

function handleSearch(){

}
function handleMore(){

}
function fitTerm(){

}
function addTerminalResize() {
  window.addEventListener("resize", fitTerm);
}
function removeResizeListener() {
  window.removeEventListener("resize", fitTerm);
}
onMounted(()=>{

})

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

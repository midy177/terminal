<template>
  <terminal-tabs ref="tabRef" style="--wails-draggable:drag" :tabs="tabs" v-model="tab">
    <template v-slot:after>
        <span class="btn-add" @click="addTab">
          <i class="iconfont icon-add"></i>
      </span>
      <span class="btn-add">
            <i class="iconfont icon-tree"></i>
          </span>
      <span class="btn-add" @click="handleMore">
            <i class="iconfont icon-more"></i>
          </span>
    </template>
  </terminal-tabs>
  <div class="terminal-layout" v-if="tabs.length>0">
    <template v-for="item in tabs">
      <terminal :id="item.key" v-show="item.key === tab"></terminal>
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
import './theme-dark.css';
import {GetLocalPtyList} from "../wailsjs/go/main/App";
const tabs = <Array<Tab>>reactive([])

function addTab() {
  // let key = nanoid()
  // let newTab = {
  //   label: key,
  //   key: key,
  // }
  // tabRef.value.addTab(newTab)
  // tab.value = key
  GetLocalPtyList().then(res=>{
    console.log(res[1].Name)
    let key = nanoid()
    let newTab = {
      label: res[1].Name,
      key: key,
    }
    tabRef.value.addTab(newTab)
    tab.value = key
  }).catch(e=>{
    console.log(e);
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

.btn-add {
  width: 34px;
  border-radius: 5px;
  padding: 0;
  //color: #333;
  cursor: pointer;
  display: flex;
  justify-content: center;
  align-items: center;
  transition: background 300ms;
  height: 34px;
  line-height: 34px;
  margin-left: .3rem;
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
}
</style>

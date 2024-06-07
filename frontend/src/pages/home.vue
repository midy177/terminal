<template>
  <vue3-tabs-chrome ref="tabRef" data-tauri-drag-region :tabs="tabs" v-model="tab">
    <template v-slot:after>
        <span slot="after" class="btn-add" @click="addTab">
          <i class="iconfont icon-add"></i>
      </span>
    </template>
  </vue3-tabs-chrome>
  <div class="nav">
    <div class="nav-btns">
          <span class="nav-btn">
            <i class="iconfont icon-tree"></i>
          </span>
    </div>
    <div class="nav-location">
      <input type="search" ref="location" class="nav-location-input" @keydown.enter="handleSearch">
    </div>
    <div class="nav-btns">
          <span class="nav-btn" @click="handleMore">
            <i class="iconfont icon-more"></i>
          </span>
    </div>
  </div>
</template>

<script lang="ts" setup>
import Vue3TabsChrome from 'vue3-tabs-chrome'
import 'vue3-tabs-chrome/dist/vue3-tabs-chrome.css'
import {onMounted, reactive, ref} from 'vue'
import {GetLocalPtyList} from "../../wailsjs/go/main/App";
const tab = ref('google')
const tabRef = ref()
const tabs = reactive([
  {
    label: 'google',
    key: 'google',
  }
])
function addTab () {
  let item = Date.now().toString()
  let newTab = {
    label: 'New Tab' + item,
    key: item
  }
  tabRef.value.addTab(newTab)
  tab.value = item
}

function handleSearch(){

}
function handleMore(){

}

onMounted(()=>{
  GetLocalPtyList().then(res=>{
    console.log(res);
  }).catch(e=>{
    console.log(e);
  })
})

</script>

<style lang="less">
input[type=search]::-webkit-search-cancel-button{
  -webkit-appearance: none;
}

.btn-add {
  width: 20px;
  border-radius: 50%;
  padding: 0;
  //color: #333;
  cursor: pointer;
  display: flex;
  justify-content: center;
  align-items: center;
  transition: background 300ms;
  height: 20px;
  line-height: 20px;
  margin-left: .5rem;
  &:hover {
    background-color: rgba(0, 0, 0, .1);
  }
}

.nav {
  padding: 8px;
  background-color: #fff;
  border-bottom: 1px solid #d5d7db;
  display: flex;
  align-items: center;
  position: relative;
}

.nav-btns {
  display: flex;
}

.nav-btn {
  width: 28px;
  height: 28px;
  margin-left: 2px;
  border-radius: 14px;
  color: #666;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 300ms;
  cursor: pointer;

  &:first-of-type {
    margin-left: 0;
  }

  &:hover {
    background-color: rgba(0, 0, 0, .1);
  }
}

.nav-location {
  flex: 1;
  height: 28px;
  margin: 0 4px;
  position: relative;
}

.nav-location-input {
  width: 100%;
  height: 100%;
  border: none;
  background-color: #eff1f2;
  border-radius: 14px;
  outline: none;
  padding-left: 16px;
  transition: background 300ms;

  &:hover {
    background-color: #e6e8e9;
  }

  &:focus {
    box-shadow: 0 0 0 3px Highlight;
    background-color: #fff;
  }
}

.nav-collection {
  top: 50%;
  right: 2px;
  width: 32px;
  height: 24px;
  border-radius: 12px;
  position: absolute;
  transform: translateY(-50%);
}
</style>

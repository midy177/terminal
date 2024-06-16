<script setup lang="ts">
import {Dropdown, List, ListItem, Button, NotificationService} from 'vue-devui';
import {
  Quit, WindowCenter,
  WindowFullscreen, WindowIsFullscreen, WindowIsMaximised, WindowIsMinimised,
  WindowMinimise, WindowToggleMaximise, WindowUnfullscreen, WindowUnminimise
} from "../../../wailsjs/runtime";
import {reactive} from "vue";
const state = reactive({
  isMax: false,
  isFull: false,
  isMin: false,
})
function toggleMin(){
  WindowIsMinimised().then((res:boolean)=>{
    if (res) {
      WindowUnminimise()
    } else {
      WindowMinimise()
    }
  }).catch(e=>{
    NotificationService.open({
      type: 'error',
      title: '获取窗口是否最小化失败',
      content: e,
      duration: 3000,
    })
  })
}
function toggleFull(){
  WindowIsFullscreen().then((res:boolean)=>{
    state.isFull = res
    if (res) {
      WindowUnfullscreen()
    } else {
      WindowFullscreen()
    }
  }).catch(e=>{
    NotificationService.open({
      type: 'error',
      title: '获取窗口是否全屏失败',
      content: e,
      duration: 3000,
    })
  })
}
</script>

<template>
  <Dropdown style="min-width: 60px;" :position="['bottom-end']">
    <i class="icon-more-operate"></i>
    <template #menu>
      <List style="min-width: 60px;padding-bottom: 0;">
        <ListItem :key="0">
          <Button icon="icon-exit" variant="text" @click="WindowToggleMaximise">最大化切换</Button>
        </ListItem>
        <ListItem :key="1">
          <Button icon="icon-exit" variant="text" @click="WindowCenter">窗口居中</Button>
        </ListItem>
        <ListItem :key="2">
          <Button icon="icon-exit" variant="text" @click="toggleMin">
            {{state.isMin ? '取消最小化': '最小化'}}
          </Button>
        </ListItem>
        <ListItem :key="3">
          <Button icon="icon-exit" variant="text" @click="toggleFull">
            {{ state.isFull ? '取消全屏': '全屏' }}
          </Button>
        </ListItem>
        <ListItem :key="4">
          <Button icon="icon-exit" variant="text" @click="Quit()">退出</Button>
        </ListItem>
      </List>
    </template>
  </Dropdown>
</template>

<style scoped lang="less">
</style>

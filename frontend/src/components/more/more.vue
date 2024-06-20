<script setup lang="ts">
import { NotificationService, Icon} from 'vue-devui';
import {
  Dropdown, Button, Menu, MenuItem, Tooltip
} from "ant-design-vue";
import {
  Quit, WindowCenter,
  WindowFullscreen, WindowIsFullscreen, WindowIsMinimised,
  WindowMinimise, WindowToggleMaximise, WindowUnfullscreen, WindowUnminimise
} from "../../../wailsjs/runtime";
import {reactive} from "vue";
const props = defineProps({
  fileBrowser: {
    type: Function,
  }
})
const state = reactive({
  isMax: false,
  isFull: true,
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

function openFileBrowser() {
  if (props.fileBrowser) props.fileBrowser()
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
  <Dropdown :position="['bottom-end']">
    <Button
        type="text"
        ghost
        size="small"
    >
      <template #icon>
        <Tooltip placement="left" title="更多">
          <Icon name="icon-drag-small" color="#f2f3f5"/>
        </Tooltip>
      </template>
    </Button>
    <template #overlay>
      <Menu>
        <MenuItem :key="1" @click="openFileBrowser">
          <template #icon>
            <Icon name="icon-folder" color="#f2f3f5" size="1rem">
            </Icon>
          </template>
          管理文件
        </MenuItem>
        <MenuItem :key="2" @click="WindowToggleMaximise">
          <template #icon>
            <Icon name="icon-maxmize" color="#f2f3f5" size="1rem">
            </Icon>
          </template>
          最大化切换
        </MenuItem>
        <MenuItem :key="3" @click="WindowCenter">
          <template #icon>
            <Icon name="icon-text-align-center" color="#f2f3f5" size="1rem">
            </Icon>
          </template>
          窗口居中
        </MenuItem>
        <MenuItem :key="4" @click="toggleMin">
          <template #icon>
            <Icon name="icon-minimize" color="#f2f3f5" size="1rem">
            </Icon>
          </template>
          {{state.isMin ? '取消最小化': '最小化'}}
        </MenuItem>
        <MenuItem :key="5" @click="toggleFull">
          <template #icon>
            <Icon name="icon-ue-expand" color="#f2f3f5" size="1rem">
            </Icon>
          </template>
          {{ state.isFull ? '全屏':'取消全屏' }}
        </MenuItem>
        <MenuItem :key="6" @click="Quit()">
          <template #icon>
            <Icon name="icon-exit" color="#f2f3f5" size="1rem">
            </Icon>
          </template>
          退出
        </MenuItem>
      </Menu>
    </template>
  </Dropdown>
</template>

<style scoped lang="less">
/deep/.devui-icon__container {
  display: block;
}
</style>

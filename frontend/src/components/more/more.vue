<script setup lang="ts">
import { Icon } from 'vue-devui';
import {
  Dropdown, Button, Menu, MenuItem, Tooltip, notification
} from "ant-design-vue";
import {
  EventsOff,
  EventsOn,
  Quit, WindowCenter,
  WindowFullscreen, WindowIsFullscreen, WindowIsMinimised,
  WindowMinimise, WindowSetAlwaysOnTop, WindowToggleMaximise,
  WindowUnfullscreen, WindowUnminimise
} from "../../../wailsjs/runtime";
import {onMounted, reactive} from "vue";
import {ExportData, ImportData, IsRunAsAdmin, OsGoos, RunAsAdmin} from "../../../wailsjs/go/logic/Logic";
const props = defineProps({
  fileBrowser: {
    type: Function,
  }
})
const state = reactive({
  goos: '',
  isMax: false,
  isFull: true,
  isRunAsAdmin: false,
  alwaysOnTop: false,
})
function toggleMin(){
  WindowIsMinimised().then((res:boolean)=>{
    if (res) {
      WindowUnminimise()
    } else {
      WindowMinimise()
    }
  }).catch(e=>{
    notification.error({
      message: '获取窗口是否最小化失败',
      description: '错误信息：'+ e,
      duration: null
    });
  })
}

function toggleMax(){
  state.isMax = !state.isMax
  WindowToggleMaximise();
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
    notification.error({
      message: '获取窗口是否全屏失败',
      description: '错误信息：'+ e,
      duration: null
    });
  })
}

function SetAlwaysOnTop() {
  state.alwaysOnTop = !state.alwaysOnTop
  WindowSetAlwaysOnTop(state.alwaysOnTop)
}

function ExportConfig(){
  ExportData().then((resp: string)=>{
    notification.info({
      message: '导出成功',
      description: '导出文件：'+ resp,
      duration: null
    });
  }).catch(err=>{
    notification.error({
      message: '导出失败',
      description: '错误信息：'+ err,
      duration: null
    });
  })
}

function ImportConfig(){
  const key = 'import_data_event';
  EventsOn(key,resp => {
        notification.info({
          key: key,
          message: '导入中',
          description: resp,
          duration: null
        });
  });
  ImportData().then(()=>{
    notification.info({
      key: key,
      message: '✅导入完成',
      description: '导入完成，如果配置存在则不会导入,如果使用了key,需要手动重新关联!',
      duration: null
    });
  }).catch(err=>{
    notification.error({
      key: key,
      message: '导入错误',
      description: err,
      duration: null
    });
  }).finally(()=>{
    EventsOff(key)
  })
}

onMounted(()=>{
  IsRunAsAdmin().then(resp=>{
    state.isRunAsAdmin = resp
  });
  WindowSetAlwaysOnTop(state.alwaysOnTop);
  OsGoos().then(resp=>{
    state.goos = resp
  })
})
</script>

<template>
  <Dropdown :position="['bottom-end']">
    <Button
        type="text"
        size="small"
    >
      <template #icon>
          <Icon name="icon-drag-small" color="#f2f3f5"/>
      </template>
    </Button>
    <template #overlay>
      <Menu>
        <MenuItem :key="-1" @click="SetAlwaysOnTop">
          <template #icon>
            <Icon name="icon-group-submit" color="#f2f3f5" size="1rem">
            </Icon>
          </template>
         {{ state.alwaysOnTop ? '取消置顶':'窗口置顶' }}
        </MenuItem>
        <MenuItem v-if="!state.isRunAsAdmin" :key="0" @click="RunAsAdmin">
          <template #icon>
            <Icon name="icon-op-mine" color="#f2f3f5" size="1rem">
            </Icon>
          </template>
          管理员运行
        </MenuItem>
        <MenuItem :key="1" @click="openFileBrowser">
          <template #icon>
            <Icon name="icon-folder-2" color="#f2f3f5" size="1rem">
            </Icon>
          </template>
          管理文件
        </MenuItem>
        <MenuItem :key="2" @click="toggleMax">
          <template #icon>
            <Icon :name="state.isMax ? 'icon-minimize':'icon-maxmize'" color="#f2f3f5" size="1rem">
            </Icon>
          </template>
          {{ state.isMax ? '取消最大化':'窗口最大化' }}
        </MenuItem>
        <MenuItem :key="3" @click="WindowCenter">
          <template #icon>
            <Icon name="icon-location" color="#f2f3f5" size="1rem">
            </Icon>
          </template>
          窗口居中
        </MenuItem>
        <MenuItem :key="4" @click="toggleMin">
          <template #icon>
            <Icon name="icon-minimize" color="#f2f3f5" size="1rem">
            </Icon>
          </template>
          最小化
        </MenuItem>
        <MenuItem :key="5" @click="toggleFull" v-if="state.goos != 'darwin'">
          <template #icon>
            <Icon name="icon-ue-expand" color="#f2f3f5" size="1rem">
            </Icon>
          </template>
          {{ state.isFull ? '全屏':'取消全屏' }}
        </MenuItem>
        <MenuItem :key="6" @click="ExportConfig">
          <template #icon>
            <Icon name="icon-export" color="#f2f3f5" size="1rem">
            </Icon>
          </template>
          导出ssh配置
        </MenuItem>
        <MenuItem :key="7" @click="ImportConfig">
          <template #icon>
            <Icon name="icon-import" color="#f2f3f5" size="1rem">
            </Icon>
          </template>
          导入ssh配置
        </MenuItem>
        <MenuItem :key="8" @click="Quit()">
          <template #icon>
            <Icon name="icon-op-exit-2" color="#f2f3f5" size="1rem">
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

<script setup lang="ts">

import {message, notification, Tooltip} from "ant-design-vue";
import { StopRec } from "../../../wailsjs/go/logic/Logic";

const props = defineProps({
  id: {
    type: String,
    required: true
  },
  filename: {
    type: String,
    required: true
  },
  stopRecording: {
    type: Function,
    required: true
  }
})

function stopRecording(){
  StopRec(props.id).then(()=>{
    notification.success({
      message: '录屏已停止',
      description: '录屏已保存到:'+ props.filename,
      duration: null
    })
    props.stopRecording();
  }).catch(err=>{
    message.error(err);
  });
}

</script>

<template>
  <Tooltip title="停止录屏">
    <div class="rec-indicator" @click="stopRecording">
      <span class="rec-dot"></span>
      <span class="rec-text">REC</span>
    </div>
  </Tooltip>
</template>

<style scoped lang="less">
.rec-indicator {
  z-index: 1000;
  position: absolute;
  top: 5px;
  right: 5px;
  display: flex;
  align-items: center;
  font-family: Arial, sans-serif;
  font-weight: bold;
  background-color: rgba(0, 0, 0, 0.5);
  padding: 5px 10px;
  border-radius: 15px;
  cursor: pointer;
}

.rec-dot {
  width: 12px;
  height: 12px;
  background-color: red;
  border-radius: 50%;
  margin-right: 5px;
  animation: blink 1s infinite;
}

.rec-text {
  color: red;
  font-size: 14px;
}

@keyframes blink {
  0% { opacity: 1; }
  50% { opacity: 0.5; }
  100% { opacity: 1; }
}
</style>

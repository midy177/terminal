<script setup lang="ts">
import { Terminal } from '@xterm/xterm';
import { WebglAddon } from '@xterm/addon-webgl';
import "./xterm.css";
import {ComponentPublicInstance, nextTick, onMounted, onUnmounted, reactive, ref, VNodeRef} from 'vue';
import {
  ClosePty,
  ResizePty,
  SetClipTextToClipboard,
  StartRec,
  WriteClipboardToPty,
  WriteToPty
} from "../../../wailsjs/go/logic/Logic";
import {EventsOff, EventsOn} from "../../../wailsjs/runtime";
import {IRenderDimensions} from "@xterm/xterm/src/browser/renderer/shared/Types";
import {message, notification} from "ant-design-vue";
import Recording from "./recording.vue";
const webgl_addon = new WebglAddon();
webgl_addon.onContextLoss(e=>{
  webgl_addon.dispose();
  message.error('Webgl 上下文丢失: '+ e);
});

const props = defineProps({
  id: {
    type: String,
    required: true
  },
  width: {
    type: Number,
    required: true
  },
  height: {
    type: Number,
    required: true
  }
});

const state = reactive({
  term: null as unknown as Terminal,
  cols: 0,
  rows: 0,
  recording: false,
  recordingFile: ''
});
const currentRef = ref<VNodeRef | null>(null);
// 赋值动态ref到变量
function setItemRef(vn: Element | ComponentPublicInstance | VNodeRef | undefined | null) {
  if (vn) {
    currentRef.value = vn
  }
}

function NewTerminal(){
  state.term = new Terminal({
    theme: {
      background: '#1A1B1E',
      cursor: '#90f64c',
      cursorAccent: '#10B98100',
    },
    fontFamily: 'JetBrainsMono, monaco, Consolas, Lucida Console, monospace',
    fontWeight: 'bold',
    fontSize: 18,
    cursorBlink: true,
    convertEol: true, // 确保换行符被正确处理
    disableStdin: false,
    cursorStyle: 'bar',
    allowTransparency: true,
    allowProposedApi: true,
    overviewRulerWidth: 8,
    scrollback: 1000,
    logLevel: 'off',
  });
  state.term.open(currentRef.value);
  state.term.loadAddon(webgl_addon);
  state.term.onTitleChange((title)=>{
    emit('update:title', title);
  })
  state.term.onResize(size=>{
    ResizePty(props.id,size.rows,size.cols).catch(err=>{
      message.error(err);
    });
  })
  state.term.attachCustomKeyEventHandler(event => {
    // 处理粘贴操作
    if (event.type === 'keydown' && event.ctrlKey && event.key === 'v') {
      event.preventDefault(); // 阻止默认的粘贴行为
      WriteClipboardToPty(props.id);
      return false; // 阻止 xterm.js 进一步处理这个事件
    }

    // 处理复制操作
    if (event.type === 'keydown' && event.ctrlKey && event.key === 'c') {
      const copiedText = state.term.getSelection();
      if (copiedText.length > 0) {
        SetClipTextToClipboard(copiedText).then(()=>{
          state.term.clearSelection();
        }).catch( e => {
          message.error(e)
        })
        return false; // 阻止 xterm.js 进一步处理这个事件
      }
    }

    // 对于其他按键事件，让 xterm.js 正常处理
    return true;
  })
}

const emit = defineEmits(['update:title']);

interface ColsRows {
  cols: number,
  rows: number
}

function getColsRows() {
  return <ColsRows>{
    cols: state.term.cols,
    rows :state.term.rows
  }
}

function fitWithHeightWidth(width:number = props.width,height:number = props.height) {
  if (width == 0 || height == 0) return;
  const core =  (state.term as any)._core;
  const dims: IRenderDimensions = core._renderService.dimensions;
  if (dims.css.cell.width === 0 || dims.css.cell.height === 0) return;
  const cols = Math.floor(width / dims.css.cell.width);
  const rows = Math.floor(height / dims.css.cell.height);
  if (state.cols == cols && state.rows == rows) return;
  state.cols = cols;
  state.rows = rows;
  if (Number.isFinite(rows) && Number.isFinite(cols)){
    state.term.resize(cols, rows);
  }
}

// Write data from pty into the terminal
function writeToTerminal(data: string | Uint8Array | ArrayBuffer | Blob) {
  toUint8Array(data).then(res=>{
    state.term.write(res);
  })
}

// Write data from the terminal to the pty
function writeToPty(data: string | Uint8Array | ArrayBuffer | Blob) {
  toUint8Array(data).then(resp=>{
    WriteToPty(props.id,Array.from(resp)).catch(e=>{
      message.error(e);
    });
  })
}

function ptyStdoutListener(){
  EventsOn(props.id,(resp: string)=>{
    writeToTerminal(resp);
  })
}

function initShell() {
  NewTerminal();
  ptyStdoutListener();
}

async function toUint8Array(input: string | Uint8Array | ArrayBuffer | Blob): Promise<Uint8Array> {
  if (typeof input === 'string') {
    // Convert string to Uint8Array
    return new TextEncoder().encode(input);
  } else if (input instanceof Uint8Array) {
    // Return Uint8Array directly
    return input;
  } else if (input instanceof ArrayBuffer) {
    // Convert ArrayBuffer to Uint8Array
    return new Uint8Array(input);
  } else if (input instanceof Blob) {
    // Convert Blob to Uint8Array
    const arrayBuffer = await input.arrayBuffer();
    return new Uint8Array(arrayBuffer);
  } else {
    throw new Error('Unsupported input type');
  }
}

function rightMouseDown(event: any) {
  if (event.button === 2) {
    if (state.term.hasSelection()) {
      const copiedText = state.term.getSelection();
      if (copiedText.length > 0) {
        SetClipTextToClipboard(copiedText).then(()=>{
          state.term.clearSelection();
        }).catch( e => {
          message.error(e)
        })
      }
    } else {
      WriteClipboardToPty(props.id).catch(e => {
        message.error(e)
      });
    }
  }
}

function focusTerminal(){
  state.term.focus()
}

defineExpose({
  fitWithHeightWidth,
  getColsRows,
  focusTerminal,
  startRecording
})

function onDragover(event: DragEvent){
  event.preventDefault();
}

function onDrop(event: DragEvent){
  event.preventDefault();
  if (event?.dataTransfer?.items && event?.dataTransfer?.items?.length > 0) {
    writeToPty('trz\r');
  }
}

function startRecording(){
  StartRec(props.id, state.term.rows, state.term.cols).then(res => {
    state.recording = true;
    state.recordingFile = res;
    notification.warning({
      message: '录屏开始',
      description: '录制中...',
      duration: 3
    })
  }).catch(err=>{
    message.error(err);
  });
}

function stopRecording(){
  state.recording = false;
  state.recordingFile = '';
}

onMounted(()=>{
  nextTick(() => {
    initShell();
    state.term.onData(writeToPty);
    state.term.onBinary(writeToPty);
    fitWithHeightWidth();
  });
})

onUnmounted( () => {
  ClosePty(props.id).catch(err=>{
    message.error(err);
  });
  EventsOff(props.id);
})
</script>

<template>
  <div class="terminal-item">
    <div
        :ref="setItemRef"
        class="xterm-layout"
        @contextmenu.prevent="rightMouseDown"
        @dragover="onDragover"
        @drop="onDrop"
    />
    <recording v-if="state.recording" :id="props.id" :filename="state.recordingFile" :stop-recording="stopRecording"/>
  </div>
</template>

<style scoped lang="less">
.terminal-item {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}
.xterm-layout {
  background-color: #1d1e21;
  height: 100%;
  width: 100%;
  max-height: 100%;
  max-width: 100%;
}
/deep/ .terminal {
  height: 100%;
  width: 100%;
  max-height: 100%;
  max-width: 100%;
  justify-content: center; /* 水平居中对齐内容 */
  align-items: center; /* 垂直居中对齐内容 */
  background-color: rgb(26, 27, 30);
}
</style>

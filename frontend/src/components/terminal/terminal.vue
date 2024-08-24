<script setup lang="ts">
import { Terminal } from '@xterm/xterm';
import "./xterm.css";
import {ComponentPublicInstance, nextTick, onMounted, onUnmounted, reactive, ref, VNodeRef} from 'vue';
import {ClosePty, GetStats, ResizePty, WriteToPty} from "../../../wailsjs/go/logic/Logic";
import {EventsOff, EventsOn} from "../../../wailsjs/runtime";
import {IRenderDimensions} from "@xterm/xterm/src/browser/renderer/shared/Types";
import {message} from "ant-design-vue";

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

// const fitAddon = new FitAddon();
const state = reactive({
  term: null as unknown as Terminal,
  cols: 0,
  rows: 0
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
    scrollback: 10000,
  });
  state.term.open(currentRef.value);
  state.term.onTitleChange((title)=>{
    emit('update:title', title);
  })
  state.term.onResize(size=>{
    ResizePty(props.id,size.rows,size.cols).catch(event=>{
      console.error(event);
    });
  })
  state.term.attachCustomKeyEventHandler(event => {
    if (event.ctrlKey && event.key === 'v') {
      return false; // 阻止默认行为
    } else if (event.ctrlKey && event.key === 'c') {
      const copiedText = state.term.getSelection();
      if (copiedText.length > 0) {
        navigator.clipboard.writeText(copiedText).then(() => {
          state.term.clearSelection();
        });
        return false; // 阻止默认行为
      }
    }
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
    // Todo 从后端读取数据，通过调用写入xterm
    state.term.write(res);
  })
}

// Write data from the terminal to the pty
function writeToPty(data: string | Uint8Array | ArrayBuffer | Blob) {
  toUint8Array(data).then(resp=>{
    // Todo 通过调用func写入后端
    WriteToPty(props.id,Array.from(resp)).catch(e=>{
      console.error(e);
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
    handleSelectToClipboardOrClipboardToTerm();
  }
}

function handleSelectToClipboardOrClipboardToTerm() {
  try {
    if (state.term.hasSelection()) {
      navigator.clipboard.writeText(state.term.getSelection()).then(()=>{
        state.term.clearSelection();
      });
    } else {
      navigator.clipboard.readText().then(clipText => {
        if (clipText.length>0){
          writeToPty(clipText.replace(/\r\n/g, "\r"));
        }
      })
    }
  } catch (e) {

  }
}

function focusTerminal(){
  state.term.focus()
}

defineExpose({
  fitWithHeightWidth,
  getColsRows,
  focusTerminal,
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
    <div
        :ref="setItemRef"
        class="xterm-layout"
        @contextmenu.prevent="rightMouseDown"
        @dragover="onDragover"
        @drop="onDrop"
    />
</template>

<style scoped lang="less">
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

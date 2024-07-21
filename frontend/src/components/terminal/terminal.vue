<script setup lang="ts">
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import "./xterm.css";
import {ComponentPublicInstance, nextTick, onMounted, onUnmounted, reactive, ref, VNodeRef} from 'vue';
import {ClosePty, GetStats, ResizePty, WriteToPty} from "../../../wailsjs/go/logic/Logic";
import {EventsOff, EventsOn} from "../../../wailsjs/runtime";
import {IRenderDimensions} from "xterm/src/browser/renderer/shared/Types";
import {DropEvent} from "vue-devui/dragdrop-new";
import {TrzszAddon, TrzszFilter} from "trzsz";

const props = defineProps({
  id: {
    type: String,
    required: true
  },
});
const trzszFilter = new TrzszFilter({
  // 将服务器的输出转发给终端进行显示，当用户在服务器上执行 trz / tsz 命令时，输出则会被接管。
  writeToTerminal: (data) => writeToTerminal(data),
  // 将服务器的输出转发给终端进行显示，当用户在服务器上执行 trz / tsz 命令时，输出则会被接管。
  sendToServer: (data) => writeToPty(data),
});

const fitAddon = new FitAddon();
const state = reactive({
  term: null as unknown as Terminal,
  width: 0,
  height: 0,
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
    disableStdin: false,
    cursorStyle: 'bar',
    allowTransparency: true,
    allowProposedApi: true,
    overviewRulerWidth: 8,
    scrollback: 5000
  });
  state.term.loadAddon(fitAddon);
  state.term.open(currentRef.value);
  state.term.onKey(async (e) => {
    if (e.key === '\x03') {
      const copiedText = state.term.getSelection();
      if (copiedText.length > 0) {
        navigator.clipboard.writeText(copiedText).then(() => {
          state.term.clearSelection();
        });
        e.domEvent.preventDefault();
      }
    } else if (e.key === '\x16') {
      const clipText = await navigator.clipboard.readText();
      if (clipText.length > 0) {
        writeToPty(clipText.replace(/\r\n/g, "\n"));
        e.domEvent.preventDefault();
      }
    }
  })
  state.term.onTitleChange((title)=>{
    emit('update:title', title);
  })
  state.term.onResize(size=>{
    ResizePty(props.id,size.rows,size.cols).catch(e=>{
      console.error(e);
    });
  })
}

const emit = defineEmits(['update:title']);

interface ColsRows {
  cols: number,
  rows: number
}

function getColsRows() {
  // fitWithHeightWidth()
  return <ColsRows>{
    cols: state.term.cols,
    rows :state.term.rows
  }
}

function fitWithHeightWidth(width:number = state.width,height:number = state.height) {
  if (width == 0 || height == 0) return;
  const core =  (state.term as any)._core;
  const dims: IRenderDimensions = core._renderService.dimensions;
  if (dims.css.cell.width === 0 || dims.css.cell.height === 0) return;
  const cols = Math.floor(width / dims.css.cell.width);
  const rows = Math.floor(height / dims.css.cell.height);
  if (width == state.width && height == state.height && state.cols == cols && state.rows == rows) return;
  state.cols = cols;
  state.rows = rows;
  state.width = width;
  state.height = height;
  if (Number.isFinite(rows) && Number.isFinite(cols)){
    state.term.resize(cols, rows);
    trzszFilter.setTerminalColumns(cols);
    // ResizePty(props.id,rows,cols).catch(e=>{
    //   console.error(e);
    // });
  }
}

// Make the terminal fit all the window size
function fitTerminal() {
  fitAddon.fit();
  // getHelperRect().then(()=>{
  //   fitWithHeightWidth();
    // Todo 从后端读取数据，通过调用func写入后端
    ResizePty(props.id,state.term.rows,state.term.cols).catch(e=>{
      console.error(e);
    });
  // });
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
    trzszFilter.processServerOutput(resp);
    // writeToTerminal(resp);
  })
}

function initShell() {
  NewTerminal();
  // Todo 请求pty或者ssh
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
          writeToPty(clipText.replace(/\r\n/g, "\n"));
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
    trzszFilter.uploadFiles(event.dataTransfer.items)
        .then(() => console.log("upload success"))
        .catch((err) => console.log(err));
  }
}

onMounted(()=>{
  nextTick(() => {
    initShell();
    // state.term.onData(writeToPty);
    // state.term.onBinary(writeToPty);
    state.term.onData((data) => trzszFilter.processTerminalInput(data));
    state.term.onBinary((data) => trzszFilter.processBinaryInput(data));
    // 初次渲染时调整大小
    fitAddon.fit();
    trzszFilter.setTerminalColumns(state.term.cols);
    // fitTerminal();
    // 获取监控信息
    // GetStats(props.id).then(resp=>{
    //   console.log(resp)
    // }).catch(e=>{
    //   console.log(e)
    // })
  });
})

onUnmounted( () => {
  ClosePty(props.id).catch(e=>{
    console.error(e);
  });
  EventsOff(props.id);
  // state.term.dispose();
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
///deep/ .xterm .xterm-viewport::-webkit-scrollbar-thumb {
//  border-radius: 10px;
//  border-color: transparent;
//  //-webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.13);
//  background-color: rgba(184, 184, 184, 0.1);
//  background-clip: padding-box;
//}
///deep/ .xterm .xterm-viewport::-webkit-scrollbar-thumb:hover {
//  background-color: rgb(224, 225, 227);
//}
</style>

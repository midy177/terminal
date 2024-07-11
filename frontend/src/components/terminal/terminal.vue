<script setup lang="ts">
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import "./xterm.css";
import {ComponentPublicInstance, nextTick, onMounted, onUnmounted, reactive, ref, VNodeRef} from 'vue';
import {ClosePty, ResizePty, WriteToPty} from "../../../wailsjs/go/logic/Logic";
import {EventsOff, EventsOn} from "../../../wailsjs/runtime";

const props = defineProps({
  id: {
    type: String,
    required: true
  },
});
const fitAddon = new FitAddon();
const state = reactive({
  term: null as unknown as Terminal,
  width: 0,
  height: 0,
  helperRect:{
    width: 0,
    height: 0
  }
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
  // state.term.onResize((size) => {
  //   ResizePty(props.id,size.rows,size.cols).then();
  // })
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
  if (width == state.width && height == state.height) return;
  state.width = width;
  state.height = height;
  if (state.helperRect.width == 0 || state.helperRect.height == 0) {
    // console.log('getHelperRect');
    if (!currentRef.value) return;
    const xtermElement = currentRef.value;
    if (xtermElement.style.display == 'none') return;
    // const xtermRect = xtermElement.getBoundingClientRect();
    const xtermHelperElement = xtermElement.querySelector('.xterm-helper-textarea');
    if (!xtermHelperElement) return;
    const helperRect = xtermHelperElement.getBoundingClientRect();
    state.helperRect.height =helperRect.height;
    state.helperRect.width = helperRect.width;
  }
  // console.log('window size change');
  const cols = Math.floor(width / state.helperRect.width);
  const rows = Math.round(height / state.helperRect.height);
  if (Number.isFinite(rows) && Number.isFinite(cols)){
    ResizePty(props.id,rows,cols).then(()=>{
      state.term.resize(cols, rows);
    });
  }
}

// Make the terminal fit all the window size
function fitTerminal() {
  fitAddon.fit();
  // getHelperRect().then(()=>{
  //   fitWithHeightWidth();
    // Todo 从后端读取数据，通过调用func写入后端
    ResizePty(props.id,state.term.rows,state.term.cols).then();
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
  toUint8Array(data).then(res=>{
    // Todo 通过调用func写入后端
    WriteToPty(props.id,Array.from(res)).then()
  })
}

function ptyStoutListener(){
  EventsOn(props.id,(res: string)=>{
    writeToTerminal(res);
  })
}

function initShell() {
  NewTerminal();
  // Todo 请求pty或者ssh
  ptyStoutListener()
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


onMounted(()=>{
  nextTick(() => {
    initShell();
    state.term.onData(writeToPty);
    state.term.onBinary(writeToPty);
    // 初次渲染时调整大小
    fitTerminal();
  });
})
onUnmounted( () => {
  ClosePty(props.id).then().catch();
  EventsOff(props.id);
})
</script>

<template>
    <div :ref="setItemRef"
       class="xterm-layout"
       @contextmenu.prevent="rightMouseDown"
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

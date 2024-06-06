<script setup lang="ts">
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import "./xterm.css";
import { TrzszFilter } from 'trzsz';
import {ComponentPublicInstance, onMounted, reactive, ref, VNodeRef} from 'vue';

const props = defineProps({
  id: {
    type: String,
    required: true
  }
});
const fitAddon = new FitAddon();
const state = reactive({
  trzszFilter: null as unknown as TrzszFilter,
  term: null as unknown as Terminal,
});
const currentRef = ref<VNodeRef | null>(null);
// 赋值动态ref到变量
function setItemRef(vn: Element | ComponentPublicInstance | null) {
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
  });
  state.term.loadAddon(fitAddon);
  state.term.open(currentRef.value);
}

// Make the terminal fit all the window size
async function fitTerminal() {
  fitAddon.fit();
  // Todo 从后端读取数据，通过调用func写入后端
}

// Write data from pty into the terminal
function writeToTerminal(data: string | Uint8Array | ArrayBuffer | Blob) {
  toUint8Array(data).then(res=>{
    // Todo 从后端读取数据，通过调用写入xterm
    state.term.write(res);
  }).catch(e=>console.log(e))
}

// Write data from the terminal to the pty
function writeToPty(data: string | Uint8Array | ArrayBuffer | Blob) {
  toUint8Array(data).then(res=>{
    // Todo 从后端读取数据，通过调用func写入后端
  })
}
function initShell() {
  NewTerminal();
  // Todo 请求pty或者ssh
}


async function readFromPty() {
  // Todo 从后端读取数据，通过事件的方式
  window.requestAnimationFrame(readFromPty);
}
function dragover(event: any) {
  event.preventDefault();
}

function drop(event: any){
  event.preventDefault();
  state.trzszFilter.uploadFiles(event.dataTransfer.items)
      .then(() => console.log("upload success"))
      .catch((err:any) => console.error(err));
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

defineExpose({
  fitTerminal,
  writeToTerminal,
})

onMounted(()=>{
  state.trzszFilter = new TrzszFilter({
    writeToTerminal: (data: string | Uint8Array | ArrayBuffer | Blob) => {
      toUint8Array(data).then(res=>{
        writeToTerminal(res);
      }).catch(e=>console.log(e))
    },
    sendToServer: (data: string | Uint8Array | ArrayBuffer | Blob) => {
      toUint8Array(data).then(res=>{
        writeToPty(res);
      }).catch(e=>console.log(e))
    },
  });
  initShell();
  state.term.onData(writeToPty);
  addEventListener("resize", fitTerminal);
  fitTerminal();
  writeToTerminal(new TextEncoder().encode(props.id));
  // window.requestAnimationFrame(readFromPty);
})

</script>

<template>
  <div :ref="setItemRef"
       @dragover="dragover"
       @drop="drop"
       class="xterm-layout"/>
</template>

<style scoped lang="less">
.xterm-layout {
  display: flex;
  background-color: #1d1e21;
  height: 100%;
  width: 100%;
}
</style>

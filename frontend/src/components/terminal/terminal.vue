<script setup lang="ts">
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import "./xterm.css";
import { TrzszFilter } from 'trzsz';
import {ComponentPublicInstance, onMounted, onUnmounted, reactive, ref, VNodeRef} from 'vue';
import {ClosePty, ResizePty, WriteToPty} from "../../../wailsjs/go/main/App";
import { Tab } from "../tabs/chrome-tabs.vue";
import {EventsOff, EventsOn} from "../../../wailsjs/runtime";

const props = defineProps({
  id: {
    type: String,
    required: true
  },
  item: {
    type: Object,
    required: true
  },
});
const fitAddon = new FitAddon();
const state = reactive({
  trzszFilter: null as unknown as TrzszFilter,
  term: null as unknown as Terminal,
});
const currentRef = ref<VNodeRef | null>(null);
// 赋值动态ref到变量
function setItemRef(vn: Element | ComponentPublicInstance | VNodeRef | undefined) {
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
  state.term.onTitleChange((title)=>{
    emit('update:title', title);
  })
}
const emit = defineEmits(['update:title']);
// Make the terminal fit all the window size
async function fitTerminal() {
  fitAddon.fit();
  // Todo 从后端读取数据，通过调用func写入后端
  ResizePty(props.id,state.term.rows,state.term.cols).then(res => {
    if (res !== null && res !== undefined) console.log(res)
  }).catch(e=>{
    console.log(e);
  })
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
    // Todo 通过调用func写入后端
    WriteToPty(props.id,Array.from(res)).then().catch(e=>{
      console.log(e);
    })
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
})
onUnmounted(()=>{
  ClosePty(props.id)
  EventsOff(props.id)
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
  //display: flex;
  background-color: #1d1e21;
  height: 100%;
  width: 100%;
}
</style>

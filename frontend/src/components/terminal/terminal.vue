<script setup lang="ts">
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import "./xterm.css";
import { TrzszFilter } from 'trzsz';
import {ComponentPublicInstance, onMounted, onUnmounted, reactive, ref, VNodeRef} from 'vue';
import {ClosePty, ResizePty, WriteToPty} from "../../../wailsjs/go/logic/Logic";
import { Tab } from "../tabs/chrome-tabs.vue";
import {EventsOff, EventsOn} from "../../../wailsjs/runtime";
import {logic} from "../../../wailsjs/go/models";
import {NotificationService} from "vue-devui";

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
}
const emit = defineEmits(['update:title']);
// Make the terminal fit all the window size
async function fitTerminal() {
  fitAddon.fit();
  // Todo 从后端读取数据，通过调用func写入后端
  ResizePty(props.id,state.term.rows,state.term.cols).then().catch(e=>{
    NotificationService.open({
      type: 'error',
      title: '重置窗口大小失败',
      content: e,
      duration: 3000,
    })
  })
}

// Write data from pty into the terminal
function writeToTerminal(data: string | Uint8Array | ArrayBuffer | Blob) {
  toUint8Array(data).then(res=>{
    // Todo 从后端读取数据，通过调用写入xterm
    state.term.write(res);
  }).catch(e=> {
    NotificationService.open({
      type: 'error',
      title: '写入前端终端失败',
      content: e,
      duration: 3000,
    })
  })
}

// Write data from the terminal to the pty
function writeToPty(data: string | Uint8Array | ArrayBuffer | Blob) {
  toUint8Array(data).then(res=>{
    // Todo 通过调用func写入后端
    WriteToPty(props.id,Array.from(res)).then().catch(e=>{
      NotificationService.open({
        type: 'error',
        title: '写入后端失败',
        content: e,
        duration: 3000,
      })
    })
  })
}
// 监听pty返回的消息，并写入前端
function ptyStoutListener(){
  EventsOn(props.id,(res: string)=>{
    state.trzszFilter.processServerOutput(res);
    // writeToTerminal(res);
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
      .then(() => {
        NotificationService.open({
          type: 'success',
          title: 'Trzsz拖拽上传文件成功',
          duration: 3000,
        })
      })
      .catch((e) => {
        NotificationService.open({
          type: 'error',
          title: 'Trzsz拖拽上传文件失败',
          content: e,
          duration: 3000,
        })
      });
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
}

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
  state.term.onData(state.trzszFilter.processTerminalInput);
  state.term.onBinary(state.trzszFilter.processTerminalInput);
  initShell();
  addEventListener("resize", fitTerminal);
  fitTerminal();
})
onUnmounted( () => {
  ClosePty(props.id).then().catch(e=>{
    console.log(e)
  });
  EventsOff(props.id);
  removeEventListener("resize", fitTerminal);
})
</script>

<template>
  <div :ref="setItemRef"
       @dragover="dragover"
       @drop="drop"
       class="xterm-layout"
       @contextmenu.prevent="rightMouseDown"
  />
</template>

<style scoped lang="less">
.xterm-layout {
  //display: flex;
  background-color: #1d1e21;
  height: 100%;
  width: 100%;
}
</style>

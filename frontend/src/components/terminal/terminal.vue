<script setup lang="ts">
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
// import "xterm/css/xterm.css";
import { nanoid } from "nanoid";
import {ComponentPublicInstance, onMounted, ref, VNodeRef} from "vue";

const props = defineProps({
  id: {
    type: Number,
    required: true
  }
});
const currentId = nanoid();
const fitAddon = new FitAddon();
const currentRef = ref<VNodeRef | null>(null);
let term = ref();
// 赋值动态ref到变量
function setItemRef(vn: Element | ComponentPublicInstance | null) {
  if (vn) {
    currentRef.value = vn
  }
}

function NewTerminal(){
  term.value = new Terminal({
    // TODO: Add possibility to customize theme
    theme: {
      background: '#1A1B1E',
      cursor: '#10B981',
      cursorAccent: '#10B98100',
    },
    fontFamily: 'Cascadia Mono, MesloLGS NF, Monospace',
    fontWeight: 'normal',
    fontSize: 14,
    cursorBlink: true,
    allowTransparency: true,
    allowProposedApi: true,
    overviewRulerWidth: 8,
  });
  term.value.loadAddon(fitAddon);
  term.value.open(currentRef.value);
}


// Make the terminal fit all the window size
async function fitTerminal() {
  fitAddon.fit();
  void invoke<string>("async_resize_pty", {
    tid: currentId,
    rows: term.value.rows,
    cols: term.value.cols,
  });
}

// Write data from pty into the terminal
function writeToTerminal(data: string) {
  return new Promise<void>((r) => {
    term.value.write(data, () => r());
  });
}

// Write data from the terminal to the pty
function writeToPty(data: string) {
  // void invoke("async_write_to_pty", {
  //   tid: currentId,
  //   data: data,
  // });
}
function initShell() {
  // invoke("async_create_shell", {tid: currentId, id: props.id}).catch((error) => {
  //   // on linux it seem to to "Operation not permitted (os error 1)" but it still works because echo $SHELL give /bin/bash
  //   console.error("Error creating shell:", error);
  // });
  NewTerminal();
}


async function readFromPty() {
  // const data = await invoke<string>("async_read_from_pty",{tid: currentId});
  // if (data) {
  //   await writeToTerminal(data);
  // }
  window.requestAnimationFrame(readFromPty);
}

onMounted(()=>{
  initShell();
  term.value.onData(writeToPty);
  addEventListener("resize", fitTerminal);
  fitTerminal();
  window.requestAnimationFrame(readFromPty);
})

</script>

<template>
  <div :ref="setItemRef"></div>
</template>

<style scoped lang="less">

</style>

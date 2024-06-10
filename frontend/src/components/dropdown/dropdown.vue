<script setup lang="ts">
// import Dropdown from 'v-dropdown';
import {onMounted, PropType, reactive} from "vue";
import {GetLocalPtyList} from "../../../wailsjs/go/main/App";
import {termx} from "../../types/models";
import {Button, Dropdown, List, ListItem} from "vue-devui";
const props = defineProps({
  atClick: {
    type: Function as PropType<(item: termx.SystemShell) => void>,
    required: false
  }
});
const state = reactive({
  items: null as unknown as Array<termx.SystemShell>,
})

function dropClick(item: termx.SystemShell) {
  if (props.atClick) {
    props.atClick(item);
  } else {
    console.warn("atClick function is not provided");
  }
}
 onMounted(()=>{
   GetLocalPtyList().then((res: Array<termx.SystemShell>)=>{
     state.items = res;
   }).catch(e=>console.log(e))
 })
</script>

<template>
  <Dropdown trigger="hover" :position="['bottom-end']">
    <i class="icon-op-add"></i>
    <template #menu>
      <List style="width: 100px;padding-bottom: 0;">
        <ListItem
            v-for="(item,index) of state.items"
            :key="index"
            @click="dropClick(item)">
          {{item.Name}}
        </ListItem>
      </List>
    </template>
  </Dropdown>
</template>

<style lang="less">
</style>

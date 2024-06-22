<script setup lang="ts">
import {onMounted, PropType, reactive} from "vue";
import {GetLocalPtyList} from "../../../wailsjs/go/logic/Logic";
import {Icon} from "vue-devui";

import {Dropdown, Button, Menu, MenuItem, Tooltip
} from "ant-design-vue";
import {logic, termx} from "../../../wailsjs/go/models";
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
  <Dropdown :position="['bottom-end']">
    <Button
        type="text"
        size="small"
    >
      <template #icon>
        <Icon name="icon-copy-to-new" color="#f2f3f5"/>
      </template>
    </Button>
    <template #overlay>
      <Menu>
        <MenuItem
            v-for="(item,index) of state.items"
            :key="index"
            @click="dropClick(item)"
        >
          {{item.name}}
        </MenuItem>
      </Menu>
    </template>
  </Dropdown>
</template>

<style scoped lang="less">
/deep/.devui-icon__container {
  display: block;
}
</style>

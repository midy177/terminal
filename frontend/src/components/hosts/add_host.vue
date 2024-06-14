<script setup lang="ts">
import {Button, Modal, Form, FormItem, FormOperation, Input, FixedOverlay} from "vue-devui";
import {reactive} from "vue";
import {main} from "../../../wailsjs/go/models";
const state = reactive({
  visible: false,
  formModel: <main.HostEntry>{}
})
function openModel() {
  state.visible = true
}
function closeModel() {
  state.visible = false
}
</script>

<template>
  <Button
      icon="add"
      variant="solid"
      title="Add"
      @click="openModel"
  />
  <FixedOverlay v-model="state.visible" class="hosts-fixed-overlay" :close-on-click-overlay="false">
  <Modal
      v-model="state.visible"
      style="min-width: 60%;"
      :show-close="false"
      :draggable="false"
      :show-overlay="false"
      :close-on-click-overlay="false"
  >
    <Form layout="vertical" :data="state.formModel">
      <FormItem field="name" label="Name">
        <Input v-model="state.formModel.label" />
      </FormItem>
      <FormOperation>
        <Button variant="solid">提交</Button>
        <Button @click="closeModel">取消</Button>
      </FormOperation>
    </Form>
  </Modal>
  </FixedOverlay>
</template>

<style scoped lang="less">
.hosts-fixed-overlay {
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.43);
}
</style>

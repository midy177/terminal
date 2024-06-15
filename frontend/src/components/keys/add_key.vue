<script setup lang="ts">
import {Button, Modal, Form, FormItem,
  FormOperation, Input, FixedOverlay,
  Textarea, Row,Col,NotificationService} from "vue-devui";
import { onUnmounted, reactive} from "vue";
import {logic} from "../../../wailsjs/go/models";
import {AddKey} from "../../../wailsjs/go/logic/Logic";
const props = defineProps({
  onSuccess: {
    type: Function
  }
})
const initState = () => ({
  visible: false,
  formModel: <logic.KeyEntry>{
    id: 0,
    label: '',
    content: '',
  }
})

const state = reactive(initState())

const rules = {
  label: [{ required: true, min: 3, max: 16, message: '标签需要不小于3个字符，不大于16个字符', trigger: 'blur' }],
  content: [{ required: true, message: '私钥信息不能为空', trigger: 'blur' }],
};
function openModel() {
  state.visible = true
}
function closeModel() {
  // reset reactive
  Object.assign(state, initState());
  state.visible = false;
}
function submitData() {
  AddKey(state.formModel).then(()=>{
    NotificationService.open({
      type: 'success',
      title: '添加私钥成功',
      duration: 1000,
    })
    closeModel()
  }).catch(e=>{
    NotificationService.open({
      type: 'error',
      title: '添加私钥失败',
      content: e,
      duration: 1000,
    })
  })
}
onUnmounted(()=>{
  // reset reactive
  Object.assign(state, initState());
})
</script>

<template>
  <Button
      variant="solid"
      @click="openModel"
  >添加私钥</Button>
  <FixedOverlay v-model="state.visible" class="hosts-fixed-overlay" :close-on-click-overlay="false">
    <Modal
        v-model="state.visible"
        style="min-width: 60%;"
        title="添加私钥"
        :show-close="false"
        :draggable="false"
        :show-overlay="false"
        :close-on-click-overlay="false"
    >
      <Form
          layout="horizontal"
          :data="state.formModel"
          label-size="sm"
          label-align="center"
          :rules="rules"
          :pop-position="['top-start','bottom-start']"
      >
        <FormItem field="label" label="标签">
          <Input v-model="state.formModel.label" placeholder="请设置标签名"/>
        </FormItem>
        <FormItem field="content" label="私钥">
          <Textarea v-model="state.formModel.content" placeholder="请输入私钥"/>
        </FormItem>
        <FormOperation>
          <Row justify="end" style="width: 100%;">
            <Col :span="4">
              <Button @click="closeModel">取消</Button>
            </Col>
            <Col :span="4">
              <Button variant="solid" @click="submitData">提交</Button>
            </Col>
          </Row>
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

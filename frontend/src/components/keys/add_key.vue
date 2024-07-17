<script setup lang="ts">
import {onUnmounted, reactive, ref} from "vue";
import {logic} from "../../../wailsjs/go/models";
import {AddKey} from "../../../wailsjs/go/logic/Logic";
import {
  Modal, Input, Button, FormItem, Textarea, Form, message, notification,
} from "ant-design-vue";
import {Rule} from "ant-design-vue/es/form";

const formRef = ref();
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

const rules: Record<string, Rule[]> = {
  label: [{ required: true, min: 1, max: 64, message: '标签需要不小于3个字符，不大于16个字符', trigger: 'blur' }],
  content: [{ required: true, message: '私钥信息不能为空', trigger: 'blur' }],
};
function openModel() {
  state.visible = true
}
function closeModel() {
  if (props.onSuccess) props.onSuccess()
  // reset reactive
  Object.assign(state, initState());
  state.visible = false;
}
function submitData() {
  AddKey(state.formModel).then(()=>{
    message.success('添加私钥成功',1)
    closeModel()
  }).catch(e=>{
    notification.error({
      message: '添加私钥失败',
      description: e,
      duration: null
    });
  })
}

function onSubmit(){
  formRef.value.validate()
      .then(() => {
        submitData();
      });
}
</script>

<template>
  <Button
      variant="solid"
      @click="openModel"
  >添加私钥</Button>
    <Modal
        title="添加私钥"
        v-model:open="state.visible"
        width="80%"
        centered
        :closable="false"
        :destroyOnClose="true"
        :maskClosable="false"
        @ok="onSubmit"
        :mask="true"
        :maskStyle="{borderRadius: '.5rem'}"
    >
      <Form
          layout="horizontal"
          ref="formRef"
          :model="state.formModel"
          :rules="rules"
          name="add_host"
          :label-col="{ span: 4 }"
          :wrapper-col="{ span: 20 }"
          autocomplete="off"
      >
        <FormItem name="label" label="标签">
          <Input v-model:value="state.formModel.label" placeholder="请设置标签名"/>
        </FormItem>
        <FormItem name="content" label="私钥">
          <Textarea v-model:value="state.formModel.content" :rows="4" placeholder="请输入私钥"/>
        </FormItem>
      </Form>
    </Modal>
</template>

<style scoped lang="less">
.hosts-fixed-overlay {
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.43);
}
</style>

<script setup lang="ts">
import {NotificationService} from "vue-devui";
import {onUnmounted, reactive, ref} from "vue";
import {logic} from "../../../wailsjs/go/models";
import {AddKey} from "../../../wailsjs/go/logic/Logic";
import {
  Modal, Input, Button, FormItem, Textarea, Form,
} from "ant-design-vue";

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

const rules = {
  label: [{ required: true, min: 3, max: 16, message: '标签需要不小于3个字符，不大于16个字符', trigger: 'blur' }],
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

function onSubmit(){
  formRef.value.validate()
      .then(() => {
        submitData();
      }).catch((e:any) => {
    console.log('error', e);
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
        :mask="false"
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

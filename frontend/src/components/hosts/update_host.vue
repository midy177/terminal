<script setup lang="ts">
import { reactive, ref } from "vue";
import {logic} from "../../../wailsjs/go/models";
import { DelKey, GetFolds, GetKeyList, UpdFoldOrHost } from "../../../wailsjs/go/logic/Logic";
import Add_key from "../keys/add_key.vue";
import {
  Button, Modal, Form, FormItem, Input, InputPassword,
  Switch, InputNumber, Select, SelectOption,
  Popover, Row, Col, message, notification,
} from "ant-design-vue";
const formRef = ref();
const props= defineProps({
  data: {
    type: logic.HostEntry,
    require: true,
  },
})

const initState = () => ({
  visible: false,
  formModel: <logic.HostEntry>{},
  useKey: false,
  keyList: <Array<logic.KeyEntry>>[],
  foldList: <Array<logic.HostEntry>>[],
  title: '修改目录'
})

const state = reactive(initState())

const rules = {
  label: [{ required: true, min: 1, max: 64, message: '标签需不小于1个字符，不大于6个字符', trigger: 'blur' }],
  username: [{ required: true, message: '用户信息不能为空', trigger: 'blur' }],
  address: [{ required: true, message: '地址不能为空', trigger: 'blur' }],
  port: [{ required: true, message: '端口必须填写', trigger: 'blur' }],
  password: [{ required: false, message: '请填写密码', trigger: 'blur' }],
  key_id: [{ required: true, message: '请选择私钥', trigger: 'blur' }],
};

function openModel() {
  if (props.data) {
    state.formModel = props.data
  } else {
    closeModel()
  }
  if (!props.data?.is_folder) {
    state.title = '修改主机'
    getKeys()
  }
  GetFolds().then(res=>{
    state.foldList = res
  }).catch(e=>{
    notification.error({
      message: '获取目录列表失败',
      description: e,
      duration: null
    });
  })
  state.visible = true
}

function delKey(id:number) {
  DelKey(id).then(()=>{
    message.success('成功删除私钥',1)
    getKeys();
  }).catch(e=>{
    notification.error({
      message: '删除私钥失败',
      description: e,
      duration: null
    });
  })
}

function closeModel() {
  // reset reactive
  Object.assign(state, initState());
  state.visible = false;
}
function getKeys() {
  GetKeyList(false).then((res: Array<logic.KeyEntry>)=>{
    if (res) {
      state.keyList = res;
    }
  }).catch(e=>{
    notification.error({
      message: '获取主机列表失败',
      description: e,
      duration: null
    });
  })
}
function setData(data: logic.HostEntry) {
  state.formModel = data;
  if (state.formModel.key_id>0) state.useKey =true
  if (state.formModel.is_folder) state.title = '修改目录'
}

function onSubmit(){
  formRef.value.validate()
      .then(() => {
        addHost();
      });
}

function addHost(){
  UpdFoldOrHost(state.formModel).then(()=>{
    message.success('修改主机或目录成功',1)
    closeModel()
  }).catch(e=>{
    notification.error({
      message: '添加主机或目录失败',
      description: e,
      duration: null
    });
  })
}
defineExpose({
  setData,
  openModel,
})
</script>

<template>
  <Modal
      :title="state.title"
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
      <FormItem name="folder_id" label="上级目录">
        <Select
            v-model:value="state.formModel.folder_id"
            placeholder="请选择上级目录"
            :options="state.foldList"
            :field-names="{ label: 'label', value: 'id'}"
        />
      </FormItem>
      <FormItem name="label" label="标签">
        <Input v-model:value="state.formModel.label" placeholder="请设置标签名"/>
      </FormItem>
      <template v-if="!state.formModel.is_folder">
        <FormItem name="username" label="用户">
          <Input v-model:value="state.formModel.username" placeholder="请输入用户名"/>
        </FormItem>
        <FormItem name="address" label="地址">
          <Input v-model:value="state.formModel.address" placeholder="请输入用户名"/>
        </FormItem>
        <FormItem name="port" label="端口">
          <InputNumber id="sshPort" v-model:value="state.formModel.port" :min="0" :max="65535"/>
        </FormItem>
        <FormItem v-if="state.useKey" name="key_id">
          <template #label>
            <Popover content="是否使用私钥?" trigger="hover" style="background-color: #7693f5; color: #fff">
              <Switch v-model:checked="state.useKey" checked-children="是" un-checked-children="否"/>
            </Popover>
          </template>
          <Row :gutter="8" style="width: 98%;">
            <Col :span="18" style="flex: 1;">
              <Select
                  v-model:value="state.formModel.key_id"
                  placeholder="请选择私钥"
                  allowClear
              >

                <SelectOption
                    v-for="(item, index) in state.keyList"
                    :key="index"
                    :value="item.id"
                >
                  <Row justify="space-between">
                    <Col>{{item.label}}</Col>
                    <Col>
                      <Button v-if="state.formModel.key_id !== item.id"
                              type="text"
                              danger
                              size="small"
                              @click="delKey(<number>item.id)"
                      >
                        删除
                      </Button>
                    </Col>
                  </Row>
                </SelectOption>
              </Select>
            </Col>
            <Col :span="4">
              <add_key :on-success="getKeys"/>
            </Col>
          </Row>
        </FormItem>
        <FormItem v-else name="password">
          <template #label>
            <Popover content="是否使用私钥?" trigger="hover" style="background-color: #7693f5; color: #fff">
              <Switch v-model:checked="state.useKey" checked-children="是" un-checked-children="否"/>
            </Popover>
          </template>
          <InputPassword v-model:value="state.formModel.password" show-password placeholder="请输入ssh密码"/>
        </FormItem>
      </template>
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

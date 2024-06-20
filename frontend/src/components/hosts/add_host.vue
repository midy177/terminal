<script setup lang="ts">
import {
  NotificationService, Message, Icon
} from "vue-devui";
import {onMounted, reactive, ref} from "vue";
import {logic} from "../../../wailsjs/go/models";
import {AddFoldOrHost, DelFoldOrHost, DelKey, GetFolds, GetKeyList} from "../../../wailsjs/go/logic/Logic";
import Add_key from "../keys/add_key.vue";
import {
  Modal, Form, FormItem, InputPassword, Input,
  Switch, InputNumber, Select, SelectOption,
  Row, Col, Button, Popover, Tooltip,
} from "ant-design-vue";
const formRef = ref();
const props= defineProps({
  folder_id: {
    type: Number,
    require: true,
  },
  callback: {
    type: Function,
    require: true,
  }
})

const initState = () => ({
  visible: false,
  formModel: <logic.HostEntry>{
    label: '',
    username: '',
    address: '',
    port: 22,
    password: '',
    folder_id: 0,
    key_id: 0,
  },
  useKey: false,
  keyList: <Array<logic.KeyEntry>>[],
  foldList: <Array<logic.HostEntry>>[]
})

const state = reactive(initState())

const rules = {
  label: [{ required: true, min: 3, max: 26, message: '用户名需不小于3个字符，不大于6个字符', trigger: 'blur' }],
  username: [{ required: true, message: '用户信息不能为空', trigger: 'blur' }],
  address: [{ required: true, message: '用户信息不能为空', trigger: 'blur' }],
  port: [{ required: true, message: '端口必须填写', trigger: 'blur' }],
  password: [{ required: false, message: '必须填写密码', trigger: 'blur' }],
  key_id: [{ required: true, message: '请选择私钥', trigger: 'blur' }],
};

function openModel() {
  if (props.folder_id) {
    state.formModel.folder_id = props.folder_id
  } else {
    closeModel()
  }
  getKeys()
  GetFolds().then(res=>{
    state.foldList = res
  }).catch(e=>{
    NotificationService.open({
      type: 'error',
      title: '获取目录列表失败',
      content: e,
      duration: 1000,
    })
  })
  state.visible = true
}
function getKeys() {
  GetKeyList(false).then((res: Array<logic.KeyEntry>)=>{
    if (res.length>0) {
      state.keyList = res;
    }
  }).catch(e=>{
    NotificationService.open({
      type: 'error',
      title: '获取主机列表失败',
      content: e,
      duration: 1000,
    })
  })
}

function delKey(id:number) {
  DelKey(id).then(()=>{
    Message({
      message: '成功删除私钥',
      type: 'success'
    })
    getKeys()
  }).catch(e=>{
    NotificationService.open({
      type: 'error',
      title: '删除私钥失败',
      content: e,
      duration: 3000,
    })
  })
}
function closeModel() {
  // reset reactive
  Object.assign(state, initState());
  state.visible = false;
}

function onSubmit(){
  formRef.value.validate()
      .then(() => {
        addHost();
      }).catch((e:any) => {
        console.log('error', e);
      });
}

function addHost(){
  AddFoldOrHost(state.formModel).then(()=>{
    NotificationService.open({
      type: 'success',
      title: '添加主机或目录成功',
      duration: 1000,
    })
    closeModel()
    if (props.callback) {
      props.callback()
    }
  }).catch(e=>{
    NotificationService.open({
      type: 'error',
      title: '添加主机或目录失败',
      content: e,
      duration: 3000,
    })
  })
}

</script>

<template>
  <Button
      type="primary"
      ghost
      size="small"
      @click="openModel"
  >
    <template #icon>
      <Tooltip placement="right" title="添加主机终端">
        <Icon name="add" color="#f2f3f5"/>
      </Tooltip>
    </template>
  </Button>
  <Modal
      title="添加主机"
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
      <FormItem name="is_folder" label="目录?">
        <Switch v-model:checked="state.formModel.is_folder" checked-children="是" un-checked-children="否"/>
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
                        ghost
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
        <InputPassword  v-model:value="state.formModel.password" show-password placeholder="请输入ssh密码"/>
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

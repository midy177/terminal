<script setup lang="ts">
import {
  Button, Modal, Form, FormItem,
  FormOperation, Input, FixedOverlay,
  Switch, InputNumber, Select, Option,
  Popover, Row, Col, NotificationService, Message
} from "vue-devui";
import {onMounted, reactive} from "vue";
import {logic} from "../../../wailsjs/go/models";
import {DelKey, GetFolds, GetKeyList, UpdFoldOrHost} from "../../../wailsjs/go/logic/Logic";
import Add_key from "../keys/add_key.vue";
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
  label: [{ required: true, min: 3, max: 6, message: '用户名需不小于3个字符，不大于6个字符', trigger: 'blur' }],
  username: [{ required: true, message: '用户信息不能为空', trigger: 'blur' }],
  address: [{ required: true, message: '用户信息不能为空', trigger: 'blur' }],
  port: [{ required: true, message: '端口必须填写', trigger: 'blur' }],
  password: [{ required: false, message: '必须填写密码', trigger: 'blur' }],
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
    NotificationService.open({
      type: 'error',
      title: '获取目录列表失败',
      content: e,
      duration: 1000,
    })
  })
  state.visible = true
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
function setData(data: logic.HostEntry) {
  state.formModel = data;
  if (state.formModel.key_id>0) state.useKey =true
}

function addHost(){
  UpdFoldOrHost(state.formModel).then(()=>{
    NotificationService.open({
      type: 'success',
      title: '添加主机或目录成功',
      duration: 1000,
    })
    closeModel()
  }).catch(e=>{
    NotificationService.open({
      type: 'error',
      title: '添加主机或目录失败',
      content: e,
      duration: 3000,
    })
  })
}
defineExpose({
  setData,
  openModel,
})
</script>

<template>
  <FixedOverlay v-model="state.visible" class="hosts-fixed-overlay" :close-on-click-overlay="false">
  <Modal
      v-model="state.visible"
      style="min-width: 60%;"
      :title="state.title"
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
      <FormItem field="folder_id" label="上级目录">
        <Select
            v-model="state.formModel.folder_id"
            placeholder="请选择上级目录"
        >
          <Option
              v-for="(item, index) in state.foldList"
              :key="index"
              :value="item.id"
              :name="item.label"
          />
        </Select>
      </FormItem>
      <FormItem field="label" label="标签">
        <Input v-model="state.formModel.label" placeholder="请设置标签名"/>
      </FormItem>
      <template v-if="!state.formModel.is_folder">
      <FormItem field="username" label="用户">
        <Input v-model="state.formModel.username" placeholder="请输入用户名"/>
      </FormItem>
        <FormItem field="address" label="地址">
          <Input v-model="state.formModel.address" placeholder="请输入用户名"/>
        </FormItem>
      <FormItem field="port" label="端口">
        <InputNumber v-model="state.formModel.port" :min="0" :max="65535"/>
      </FormItem>
      <FormItem v-if="state.useKey" field="key_id">
        <template #label>
          <Popover content="是否使用私钥?" trigger="hover" style="background-color: #7693f5; color: #fff">
            <Switch v-model="state.useKey">
              <template #checkedContent>是</template>
              <template #uncheckedContent>否</template>
            </Switch>
          </Popover>
        </template>
        <Row :gutter="8" style="width: 98%;">
          <Col :span="16" style="flex: 1;">
            <Select
                v-model="state.formModel.key_id"
                placeholder="请选择私钥"
            >
              <Option
                  v-for="(item, index) in state.keyList"
                  :key="index"
                  :value="item.id"
                  :name="item.label"
              >
                <Popover trigger="hover" style="background-color: transparent;" :position="['right']">
                  {{item.label}}
                  <template #content>
                    <Button
                        size="sm"
                        icon="icon-delete"
                        variant="solid"
                        color="danger"
                        @click="delKey(<number>item.id)"
                    >
                      删除
                    </Button>
                  </template>
                </Popover>
              </Option>
            </Select>
          </Col>
          <Col :span="6">
            <add_key :on-success="getKeys"/>
          </Col>
        </Row>
      </FormItem>
      <FormItem v-else field="password">
        <template #label>
          <Popover content="是否使用私钥?" trigger="hover" style="background-color: #7693f5; color: #fff">
            <Switch v-model="state.useKey">
              <template #checkedContent>是</template>
              <template #uncheckedContent>否</template>
            </Switch>
          </Popover>
        </template>
        <Input v-model="state.formModel.password" show-password placeholder="请输入ssh密码"/>
      </FormItem>
      </template>
      <FormOperation>
        <Row justify="end" style="width: 100%;">
          <Col :span="6">
            <Button @click="closeModel">取消</Button>
          </Col>
          <Col :span="6">
            <Button variant="solid" @click="addHost">提交</Button>
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

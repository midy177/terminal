<script setup lang="ts">
import {
  FixedOverlay,
  Modal,
  ModalFooter,
  Button,
  Breadcrumb,
  BreadcrumbItem,
  Row,
  Col,
  Input,
  Icon,
  Message
} from "vue-devui";
import {STableColumnsType, STableContextmenuPopupArg} from '@shene/table';
import { STable,STableProvider } from '@shene/table';
import {onMounted, reactive} from "vue";
import Add_host from "./add_host.vue";
import {main} from "../../../wailsjs/go/models";
import {GetFoldsAndHosts} from "../../../wailsjs/go/main/App";
const state = reactive({
  visible: false,
  tableData: <Array<main.HostEntry>>[],
})

const columns: STableColumnsType<main.HostEntry> = [
  {
    title: '名称',
    dataIndex: 'label',
    key: 'label',
    width: 120
  },
  {
    title: '主机',
    dataIndex: 'address',
    key: 'address',
    width: 120
  },
  {
    title: '端口',
    dataIndex: 'port',
    key: 'port',
    width: 80
  },
  {
    title: '用户名',
    key: 'username',
    dataIndex: 'username',
    width: 120
  }
]

function closeModel(){
  state.visible = false
}
function openModel() {
  state.visible = true
}
function handleContextMenuEdit(args: STableContextmenuPopupArg) {
  Message({
    message: 'Click edit of line ' + args.index,
    type: 'success'
  })
}
function handleContextMenuDelete(args: STableContextmenuPopupArg) {
  Message({
    message: 'Click delete of line ' + args.index,
    type: 'success'
  })
}

function GetList(id:number) {
  GetFoldsAndHosts(id).then((res:main.HostEntry[])=>{
    console.log(res)
  }).catch(e=>{
    console.log(e)
  })
}

onMounted(()=>{
  GetList(0)
})
</script>

<template>
  <i class="icon-go-tree" @click="openModel"></i>
  <FixedOverlay v-model="state.visible" class="hosts-fixed-overlay" :close-on-click-overlay="false">
  <Modal
      v-model="state.visible"
      style="width: 80%;"
      :show-close="false"
      :draggable="false"
      :show-overlay="false"
      :close-on-click-overlay="false"
  >
    <template #header>
        <Row type="flex" class="header-bar">
          <Col flex="2.5rem">
            <add_host/>
          </Col>
          <Col flex="auto">
            <Breadcrumb>
              <BreadcrumbItem>Homepage</BreadcrumbItem>
              <BreadcrumbItem>
                <span>DevUI</span>
              </BreadcrumbItem>
              <BreadcrumbItem>
                <span>Breadcrumb</span>
              </BreadcrumbItem>
          </Breadcrumb>
          </Col>
          <Col flex="6rem">
            <Input placeholder="请输入">
              <template #suffix>
                <Icon name="search" style="font-size: inherit;"></Icon>
              </template>
            </Input>
          </Col>
        </Row>
    </template>
    <template #default>
      <STableProvider size="small">
        <STable
            style="--s-bg-color-component: transport;"
            :columns="columns"
            :scroll="{ y: 300 }"
            :data-source="state.tableData"
            :pagination="true"
            :max-height="300"
        >
        <template #bodyCell="{ text, column, record }">
          <template v-if="column.key === 'name'">
            <a>{{ text }}</a>
          </template>
        </template>
          <template #contextmenuPopup="args">
            <ul class="popup">
              <li class="popup-item">打开</li>
              <li class="popup-item" @click="handleContextMenuEdit(args)">编辑</li>
              <li class="popup-item" style="color: #ed4014" @click="handleContextMenuDelete(args)">删除</li>
            </ul>
          </template>
      </STable>
      </STableProvider>
    </template>
    <template #footer>
      <ModalFooter style="text-align: right; padding-right: .4rem;bottom: 0;">
        <Button
            variant="solid"
            color="secondary"
            @click="closeModel"
        >
          关闭
        </Button>
      </ModalFooter>
    </template>
  </Modal>
  </FixedOverlay>
</template>

<style scoped lang="less">
.header-bar {
  padding: .5rem;
  align-items: center;
}
.hosts-fixed-overlay {
  align-items: center;
  justify-content: center;
  background-color: transparent;
}
.popup {
  border-radius: .3rem;
}
.popup-item {
  cursor: pointer;
  padding: .3rem 1rem .3rem .4rem;
  background-color: #3b3f41;
}
.popup-item:hover {
  background-color: #90f64c;
}
.popup-item.disabled {
  color: #00000040;
  cursor: not-allowed;
}
</style>

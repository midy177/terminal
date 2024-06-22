<script setup lang="ts">
import { Icon } from "vue-devui";
import { PropType, reactive, ref} from "vue";
import Add_host from "./add_host.vue";
import {logic} from "../../../wailsjs/go/models";
import {DelFoldOrHost, GetFoldsAndHosts} from '../../../wailsjs/go/logic/Logic';
import Update_host from "./update_host.vue";
import {
  Modal, Space, Table, Breadcrumb,
  BreadcrumbItem,
  TableProps, Row, Col, Button, Popover, Tooltip, message, notification,
} from "ant-design-vue";

const modifyHostRef = ref();
const props = defineProps({
  openSshTerminal: {
    type: Function as PropType<(id:number,label:string) => void>,
    required: true
  }
})
interface breadcrumbItem {
  id: number,
  name: string,
}
const initState = () => ({
  visible: false,
  tableData: <Array<logic.HostEntry>>[],
  currentDirId: 0,
  breadcrumbSource: <Array<breadcrumbItem>>[{
    id: 0,
    name: '根'
  }]
})

const state = reactive(initState())

const columns: TableProps['columns'] = [
  {
    title: '名称',
    dataIndex: 'label',
    key: 'label',
    width: 120,
    resizable: true,
    ellipsis: true,
  },
  {
    title: '用户名',
    key: 'username',
    dataIndex: 'username',
    width: 80,
    resizable: true,
    ellipsis: true,
  },
  {
    title: '主机',
    dataIndex: 'address',
    key: 'address',
    width: 100,
    resizable: true,
    ellipsis: true,
  },
  {
    title: '端口',
    dataIndex: 'port',
    key: 'port',
    width: 60,
    resizable: true,
    ellipsis: true,
  },
  {
    title: '操作',
    key: 'action',
    width: 60,
    resizable: true,
  }
]

function handleResizeColumn(w:any, col:any) {
  col.width = w;
}

function closeModel(){
  // reset reactive
  Object.assign(state, initState());
  state.visible = false
}
function openModel() {
  getList(state.currentDirId)
  state.visible = true
}

function reRender() {
  state.tableData = []
  getList(state.currentDirId)
}

function jumperFolder(index:number) {
  if (state.breadcrumbSource[index].id === state.currentDirId) return;
  state.breadcrumbSource.splice(index + 1);
  state.currentDirId = state.breadcrumbSource[index].id
  getList(state.currentDirId)
}
function handleOpenFolder(id: number,label: string){
  state.breadcrumbSource.push({
    id: id,
    name: label
  })
  state.currentDirId = id
  getList(state.currentDirId)
}
function handleConnect(record: logic.HostEntry) {
  if (props.openSshTerminal) props.openSshTerminal(record.id,record.label)
  closeModel()
}
function handleEdit(args: logic.HostEntry) {
  if (modifyHostRef.value) {
    modifyHostRef.value.openModel()
    modifyHostRef.value.setData(args)
  }
}
function handleDelete(args: logic.HostEntry) {
  DelFoldOrHost(args.id,args.is_folder).then(()=>{
    message.success('成功删除',1)
    reRender();
  }).catch(e=>{
    notification.error({
      message: '删除失败',
      description: e,
      duration: null
    });
  })
}

function getList(id:number) {
  GetFoldsAndHosts(id).then((res:logic.HostEntry[])=>{
    state.tableData = res
  }).catch(e=>{
    notification.error({
      message: '获取主机和目录列表失败',
      description: e,
      duration: null
    });
  })
}
</script>

<template>
  <Button type="text" size="small" @click="openModel">
    <template #icon>
      <Tooltip placement="bottom" title="SSH配置">
        <Icon name="icon-go-tree" color="#f2f3f5"/>
      </Tooltip>
    </template>
  </Button>
  <Modal
      v-model:open="state.visible"
      width="90%"
      centered
      :closable="false"
      :destroyOnClose="true"
      :maskClosable="false"
      :mask="false"
  >
    <template #title>
        <Row type="flex" class="header-bar">
          <Col flex="2.5rem">
            <add_host :folder_id="state.currentDirId" :callback="reRender"/>
          </Col>
          <Col flex="auto">
            <Breadcrumb>
              <BreadcrumbItem
                  v-for="(item,index) in state.breadcrumbSource"
                  :key="index"
              >
                <Button
                    type="link"
                    size="small"
                    @click="jumperFolder(index)"
                >
                  {{ item.name }}
                </Button>
              </BreadcrumbItem>
          </Breadcrumb>
          </Col>
        </Row>
    </template>
        <Table
            :rowKey="(record:logic.HostEntry) => record.label"
            :dataSource="state.tableData"
            :columns="columns"
            :pagination="{ pageSize: 10 ,showSizeChanger: true}"
            sticky
            :scroll="{ y: '44vh' }"
            @resizeColumn="handleResizeColumn"
            size="small"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'label'">
              <Button
                  v-if="record.is_folder"
                  type="link"
                  size="small"
                  @dblclick="handleOpenFolder(record.id,record.label)"
              >
                <template #icon>
                  <Icon  name="icon-open-folder-2" color="#3DCCA6" >
                    <template #suffix>
                      <span style="color: #f2f3f5;">{{ record.label }}</span>
                    </template>
                  </Icon>
                </template>
              </Button>
              <Button
                  v-else
                  type="link"
                  size="small"
                  @dblclick="handleConnect(record)"
              >
                <template #icon>
                  <Tooltip placement="bottom" title="链接ssh">
                  <Icon name="icon-console">
                    <template #suffix>
                      <span style="color: #f2f3f5;">{{ record.label }}</span>
                    </template>
                  </Icon>
                  </Tooltip>
                </template>
              </Button>
            </template>
            <template v-else-if="column.key === 'port'">
              {{record.is_folder ? '': record.port }}
            </template>
            <template v-else-if="column.key === 'action'">
              <Space :size="1">
                <Button
                    v-if="record.is_folder"
                    type="link"
                    size="small"
                    @click="handleOpenFolder(record.id,record.label)"
                >
                  <template #icon>
                    <Tooltip placement="bottom" title="打开文件夹">
                      <Icon name="icon-open-folder"></Icon>
                    </Tooltip>
                  </template>
                </Button>
                <Button
                    v-else
                    type="link"
                    size="small"
                    @click="handleConnect(record)"
                >
                  <template #icon>
                    <Tooltip placement="right" title="连接ssh">
                      <Icon name="icon-connect"></Icon>
                    </Tooltip>
                  </template>
                </Button>
                <Button
                    type="link"
                    size="small"
                    @click="handleEdit(record)"
                >
                  <template #icon>
                    <Tooltip placement="right" title="编辑">
                      <Icon name="icon-edit"></Icon>
                    </Tooltip>
                  </template>
                </Button>
                <Popover trigger="click" placement="topLeft">
                  <Button
                      type="link"
                      size="small"
                  >
                    <template #icon>
                      <Tooltip placement="right" title="删除">
                        <Icon name="delete"></Icon>
                      </Tooltip>
                    </template>
                  </Button>
                  <template #content>
                    <Button
                        type="link"
                        size="small"
                        danger
                        @click="handleDelete(record)"
                    >确认</Button>
                  </template>
                </Popover>
              </Space>
            </template>
          </template>
        </Table>
    <template #footer>
        <Button
            type="primary"
            ghost
            @click="closeModel"
        >
          关闭
        </Button>
    </template>
  </Modal>
  <update_host ref="modifyHostRef"/>
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
/deep/.ant-table {
  border-radius: .5rem;
}
/deep/.ant-table-body {
  min-height: 44vh !important;
}
/deep/.devui-icon__container {
  display: block;
}
</style>

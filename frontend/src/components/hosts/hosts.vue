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
  Popover,
  Icon,
  Message, NotificationService
} from "vue-devui";
import {onMounted, PropType, reactive, ref} from "vue";
import Add_host from "./add_host.vue";
import {logic} from "../../../wailsjs/go/models";
import {DelFoldOrHost, GetFoldsAndHosts} from '../../../wailsjs/go/logic/Logic';
import Update_host from "./update_host.vue";
import zhCN from "ant-design-vue/es/locale/zh_CN";
import {ConfigProvider, Space, Table, TableProps, theme} from "ant-design-vue";
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
    width: 100,
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
    Message({
      message: '成功删除',
      type: 'success'
    })
    reRender()
  }).catch(e=>{
    Message({
      message: e,
      type: 'error'
    })
  })
}

function getList(id:number) {
  GetFoldsAndHosts(id).then((res:logic.HostEntry[])=>{
    state.tableData = res
  }).catch(e=>{
    NotificationService.open({
      type: 'error',
      title: '获取主机和目录列表失败',
      content: e,
      duration: 3000,
    })
  })
}
</script>

<template>
  <Button
      icon="icon-go-tree"
      variant="text"
      title="Tree"
      @click="openModel"
  />
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
            <add_host :folder_id="state.currentDirId" :callback="reRender"/>
          </Col>
          <Col flex="auto">
            <Breadcrumb>
              <BreadcrumbItem
                  v-for="(item,index) in state.breadcrumbSource"
                  :key="index"
              >
                <Button
                    variant="text"
                    @click="jumperFolder(index)"
                >
                  {{ item.name }}
                </Button>
              </BreadcrumbItem>
          </Breadcrumb>
          </Col>
        </Row>
    </template>
    <template #default>
      <ConfigProvider
          :locale="zhCN"
          :theme="{
              algorithm: theme.darkAlgorithm,
              }"
      >
        <Table
            :rowKey="(record:logic.HostEntry) => record.label"
            :dataSource="state.tableData"
            :columns="columns"
            :pagination="{ pageSize: 10 ,showSizeChanger: true}"
            sticky
            :scroll="{ y: '44vh' }"
            @resizeColumn="handleResizeColumn"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'label'">
              <Icon v-if="record.is_folder" name="icon-open-folder-2" color="#3DCCA6" operable @dblclick="handleOpenFolder(record.id,record.label)">
                <template #suffix>
                  <span style="color: #f2f3f5;">{{ record.label }}</span>
                </template>
              </Icon>
              <Icon v-else name="icon-console" operable @dblclick="handleConnect(record)">
                <template #suffix>
                  <span style="color: #f2f3f5;">{{ record.label }}</span>
                </template>
              </Icon>
            </template>
            <template v-else-if="column.key === 'port'">
              {{record.is_folder ? '': record.port }}
            </template>
            <template v-else-if="column.key === 'action'">
              <Space :size="1">
                <Button
                    v-if="record.is_folder"
                    icon="icon-open-folder"
                    variant="text"
                    title="打开目录"
                    @click="handleOpenFolder(record.id,record.label)"
                />
                <Button
                    v-else
                    icon="icon-connect"
                    variant="text"
                    title="连接ssh"
                    @click="handleConnect(record)"
                />
                <Button
                    icon="icon-setting"
                    variant="text"
                    title="编辑"
                    @click="handleEdit(record)"
                />
                <Popover trigger="hover">
                  <Button
                      icon="delete"
                      variant="text"
                      title="删除"
                  />
                  <template #content>
                    <Button
                        variant="solid"
                        color="danger"
                        title="删除"
                        @click="handleDelete(record)"
                    >确认</Button>
                  </template>
                </Popover>
              </Space>
            </template>
          </template>
        </Table>
      </ConfigProvider>
      <update_host ref="modifyHostRef"/>
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
/deep/.ant-table {
  border-radius: .5rem;
}
/deep/.ant-table-body {
  min-height: 44vh !important;
}
</style>

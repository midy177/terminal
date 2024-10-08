<script setup lang="ts">
import { Icon } from "vue-devui";
import { PropType, reactive, ref} from "vue";
import Add_host from "./add_host.vue";
import {logic} from "../../../wailsjs/go/models";
import {DelFoldOrHost, GetFoldsAndHosts} from '../../../wailsjs/go/logic/Logic';
import Update_host from "./update_host.vue";
import {
  Modal,  Space, Table, Breadcrumb, Popconfirm,
  BreadcrumbItem,
  TableProps, Row, Col, Button, Popover, Tooltip, message, notification, SelectOption, Select,
} from 'ant-design-vue';
import { HomeFilled } from '@ant-design/icons-vue';
const modifyHostRef = ref();
const props = defineProps({
  openSshTerminal: {
    type: Function as PropType<(id:number,label:string) => void>,
    required: true
  },
  openSshTerminalWithJumper: {
    type: Function as PropType<(tid:number,jid: number,label:string) => void>,
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
    name: 'Home'
  }],
  jumperId: <number|null>null,
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
    width: 50,
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
  // Object.assign(state, initState());
  state.visible = false
  state.jumperId = null
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

function handleConnectWithJumper(record: logic.HostEntry) {
  if (state.jumperId === null) {
    message.warn('未选择跳板机')
    return
  }
  if (props.openSshTerminalWithJumper) props.openSshTerminalWithJumper(record.id,state.jumperId,record.label)
  closeModel()
}

function handleEdit(args: logic.HostEntry) {
  if (modifyHostRef.value) {
    modifyHostRef.value.setData(args)
    modifyHostRef.value.openModel()
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
      <Tooltip placement="bottom" title="SSH主机列表">
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
      :mask="true"
      :maskStyle="{borderRadius: '.5rem',backgroundColor: 'var(--d2h-dark-empty-placeholder-bg-color)'}"
      :footer="null"
  >
    <template #title>
        <Row type="flex" class="header-bar">
          <Col flex="2.5rem">
            <add_host :folder_id="state.currentDirId" :on-change="reRender"/>
          </Col>
          <Col flex="auto">
            <Breadcrumb>
              <BreadcrumbItem
                  v-for="(item,index) in state.breadcrumbSource"
                  :key="index"
              >
                <template v-if="item.id === 0">
                  <HomeFilled @click="jumperFolder(index)"/>
                </template>
                <template v-else>
                  {{ item.name }}
                </template>
              </BreadcrumbItem>
          </Breadcrumb>
          </Col>
          <Col>
            <Button
                danger
                block
                ghost
                size="small"
                @click="closeModel"
            >
              关闭
            </Button>
          </Col>
        </Row>
    </template>
        <Table
            :rowKey="(record:logic.HostEntry) => record.label"
            :dataSource="state.tableData"
            :columns="columns"
            :pagination="{ pageSize: 10 ,showSizeChanger: true}"
            sticky
            :scroll="{ y: '50vh' }"
            @resizeColumn="handleResizeColumn"
            size="middle"
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
                  @dblclick="handleConnect(record as logic.HostEntry)"
              >
                <template #icon>
<!--                  <Tooltip placement="bottom" title="双击链接ssh">-->
                  <Icon name="icon-console">
                    <template #suffix>
                      <span style="color: #f2f3f5;">{{ record.label }}</span>
                    </template>
                  </Icon>
<!--                  </Tooltip>-->
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
                <template v-else>
                  <Button
                      type="link"
                      size="small"
                      @click="handleConnect(record as logic.HostEntry)"
                  >
                    <template #icon>
                      <Tooltip placement="bottom" title="建立连接">
                        <Icon name="icon-connect"></Icon>
                      </Tooltip>
                    </template>
                  </Button>
                  <Popover :title=null trigger="click">
                    <template #content>
                      <Space :size="2">
                      <Select
                            v-model:value="<number>state.jumperId"
                            placeholder="请选择跳班机"
                            allowClear
                            style="width: auto;min-width: 20vw;max-width: 100vw;"
                        >
                          <template v-for="(item, index) in state.tableData" :key="index">
                            <SelectOption
                                v-if="!item.is_folder && item.id !== record.id"
                                :key="index"
                                :value="item.id"
                            >
                              {{item.label}}
                            </SelectOption>
                          </template>
<!--                        <template #suffixIcon>-->
<!--                          <Tooltip placement="bottom" title="跳板连接">-->
<!--                            <Button-->
<!--                                size="small"-->
<!--                                type="link"-->
<!--                                @click="handleConnectWithJumper(record as logic.HostEntry)"-->
<!--                            >-->
<!--                              连接-->
<!--                            </Button>-->
<!--                            <Icon-->
<!--                                style="cursor: pointer;"-->
<!--                                @click="handleConnectWithJumper(record as logic.HostEntry)"-->
<!--                                name="icon-go-pipeline"-->
<!--                            />-->
<!--                          </Tooltip>-->
<!--                        </template>-->
                        </Select>
                      <Tooltip placement="bottom" title="跳板连接">
                        <Button
                            @click="handleConnectWithJumper(record as logic.HostEntry)"
                        >
                          连接
                        </Button>
                      </Tooltip>
                      </Space>
                    </template>
                    <Button
                        type="link"
                        size="small"
                    >
                      <template #icon>
                        <Tooltip placement="bottom" title="跳板连接">
                          <Icon name="icon-go-pipeline"></Icon>
                        </Tooltip>
                      </template>
                    </Button>
<!--                     @confirm="handleConnectWithJumper(record as logic.HostEntry)"-->
                  </Popover>
                </template>
                <Button
                    type="link"
                    size="small"
                    @click="handleEdit(record as logic.HostEntry)"
                >
                  <template #icon>
                    <Tooltip placement="bottom" title="编辑">
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
                      <Tooltip placement="bottom" title="删除">
                        <Icon name="delete"></Icon>
                      </Tooltip>
                    </template>
                  </Button>
                  <template #content>
                    <Button
                        type="link"
                        size="small"
                        danger
                        @click="handleDelete(record as logic.HostEntry)"
                    >确认</Button>
                  </template>
                </Popover>
              </Space>
            </template>
          </template>
        </Table>
<!--    <template #footer>-->
<!--        <Button-->
<!--            type="primary"-->
<!--            ghost-->
<!--            @click="closeModel"-->
<!--        >-->
<!--          关闭-->
<!--        </Button>-->
<!--    </template>-->
  </Modal>
  <update_host ref="modifyHostRef" :on-change="reRender"/>
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
  overflow: hidden;
}
/deep/.ant-table-body {
  min-height: 50vh !important;
}
/deep/.devui-icon__container {
  display: block;
}
</style>

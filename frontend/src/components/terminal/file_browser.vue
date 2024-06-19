<script setup lang="ts">
import {
  Button,
  Col,
  FixedOverlay,
  Icon, Message,
  Modal,
  ModalFooter,
  NotificationService,
  Popover,
  Row
} from "vue-devui";
import zhCN from 'ant-design-vue/es/locale/zh_CN';
import {Table, ConfigProvider, theme, Space, TableProps} from "ant-design-vue";
import Update_host from "../hosts/update_host.vue";
import {logic} from "../../../wailsjs/go/models";
import {reactive, ref} from "vue";
import {
  SftpDelete,
  SftpDir,
  SftpDownload, SftpHomeDir,
  SftpUploadDirectory,
  SftpUploadMultipleFiles
} from "../../../wailsjs/go/logic/Logic";
import filter from "../hosts/filter.vue";
const props =defineProps({
  tid: {
    type: String,
    require: true,
  }
})
const initState = () => ({
  tid: '',
  visible: false,
  tableData: <Array<logic.FileInfo>>[],
  currentDir: ''
})
const state = reactive(initState())

const columns: TableProps['columns'] = [
  {
    title: '名称',
    dataIndex: 'name',
    key: 'name',
    width: 100,
    resizable: true,
    ellipsis: true,
  },
  {
    title: '大小',
    key: 'size',
    dataIndex: 'size',
    width: 50,
    resizable: true,
    ellipsis: true
  },
  {
    title: '权限',
    dataIndex: 'mode',
    key: 'mode',
    width: 60,
    resizable: true,
    ellipsis: true
  },
  {
    title: '修改时间',
    dataIndex: 'mod_time',
    key: 'mod_time',
    resizable: true,
    width: 100,
    ellipsis: true
  },
  {
    title: '操作',
    resizable: true,
    key: 'action',
    width: 60
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
  if (props.tid) {
    SftpHomeDir(props.tid).then((res: string)=>{
      state.currentDir = res;
      handleFoldList(state.currentDir);
    }).catch(e=>{
      NotificationService.open({
        type: 'error',
        title: '获取Home路径失败',
        content: e,
        duration: 3000,
      })
    })
    state.visible = true
  } else {
    Message({
      message: '当前没有打开ssh连接',
      type: 'warning',
      bordered: false,
    });
  }
}

function handleBack() {
  if (state.currentDir.includes('/')){
    state.currentDir = state.currentDir.substring(0, state.currentDir.lastIndexOf('/'));
    if (state.currentDir === '') state.currentDir = '/';
    handleFoldList(state.currentDir);
  }
}

function handleUploadFile(){
  if (props.tid) {
    SftpUploadMultipleFiles(props.tid,state.currentDir).then(()=>{
      NotificationService.open({
        type: 'success',
        title: '上传文件成功',
        content: '目的路径'+state.currentDir,
        duration: 3000,
      })
      handleFoldList(state.currentDir);
    }).catch(e=>{
      NotificationService.open({
        type: 'error',
        title: '上传文件失败，目的路径：'+state.currentDir,
        content: e,
        duration: 3000,
      })
    })
  }
}
function handleUploadFold(){
  if (props.tid) {
    SftpUploadDirectory(props.tid,state.currentDir).then(()=>{
      NotificationService.open({
        type: 'success',
        title: '上传文件夹成功',
        content: '目的路径'+state.currentDir,
        duration: 3000,
      })
      handleFoldList(state.currentDir);
    }).catch(e=>{
      NotificationService.open({
        type: 'error',
        title: '上传文件夹：'+state.currentDir,
        content: e,
        duration: 3000,
      })
    })
  }
}
function handleFoldList(dst: string) {
  if (props.tid) {
    SftpDir(props.tid,dst).then((res:Array<logic.FileInfo>)=>{
      state.currentDir = dst;
      state.tableData = res;
    }).catch(e=>{
      NotificationService.open({
        type: 'error',
        title: '获取目录失败：'+dst,
        content: e,
        duration: 3000,
      })
    })
  }
}
function handleDownload(dst: string) {
  if (props.tid) {
    SftpDownload(props.tid, dst).then(()=>{
      NotificationService.open({
        type: 'success',
        title: '下载成功',
        content: dst,
        duration: 3000,
      })
    }).catch(e=>{
      NotificationService.open({
        type: 'error',
        title: '下载失败：'+dst,
        content: e,
        duration: 3000,
      })
    })
  }
}
function handleDelete(dst: string){
  if (props.tid) {
    SftpDelete(props.tid,dst).then(()=>{
      NotificationService.open({
        type: 'success',
        title: '删除成功',
        content: dst,
        duration: 3000,
      })
      handleFoldList(state.currentDir);
    }).catch(e=>{
      NotificationService.open({
        type: 'error',
        title: '删除失败：'+dst,
        content: e,
        duration: 3000,
      })
    })
  }
}

function formatTimestamp(timestamp: number): string {
  const date = new Date(timestamp * 1000);

  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}
defineExpose({
  closeModel,
  openModel,
})
</script>

<template>
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
        <Row :gutter="10" class="header-bar">
          <Col :span="12">
            <Popover :content="state.currentDir" trigger="hover" style="background-color: #7693f5; color: #fff">
              <Icon name="icon-folder" operable>
                <template #suffix>
                  当前路径
                </template>
              </Icon>
            </Popover>
          </Col>
          <Col :span="4">
            <Button icon="icon-go-back" @click="handleBack">上一级</Button>
          </Col>
          <Col :span="4">
            <Button icon="icon-upload" @click="handleUploadFile">多文件</Button>
          </Col>
          <Col :span="4">
            <Button icon="icon-upload" @click="handleUploadFold">文件夹</Button>
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
              :rowKey="(record:logic.FileInfo) => record.name"
              :dataSource="state.tableData"
              :columns="columns"
              :pagination="{ pageSize: 10 ,showSizeChanger: true}"
              sticky
              :scroll="{ y: '44vh' }"
              @resizeColumn="handleResizeColumn"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'name'">
                <Icon v-if="record.is_dir" name="icon-open-folder-2" color="#3DCCA6" operable @dblclick="handleFoldList(record.full_path)">
                  <template #suffix>
                    <span style="color: #f2f3f5;">{{ record.name }}</span>
                  </template>
                </Icon>
                <Icon v-else name="icon-file" operable @dblclick="handleDownload(record.full_path)">
                  <template #suffix>
                    <span style="color: #f2f3f5;">{{ record.name }}</span>
                  </template>
                </Icon>
              </template>
              <template v-else-if="column.key === 'size'">
                {{ record.is_dir ? '目录' : record.size }}
              </template>
              <template v-else-if="column.key === 'mod_time'">
                {{ formatTimestamp(record.mod_time) }}
              </template>
              <template v-else-if="column.key === 'action'">
                <Space :size="1">
                  <Button
                      icon="icon-download"
                      variant="text"
                      title="Connect"
                      @click="handleDownload(record.full_path)"
                  />
                  <Popover trigger="hover">
                    <Button
                        icon="delete"
                        variant="text"
                        title="Delete"
                    />
                    <template #content>
                      <Button
                          variant="solid"
                          color="danger"
                          title="Delete"
                          @click="handleDelete(record.full_path)"
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

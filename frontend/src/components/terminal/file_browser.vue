<script setup lang="ts">
import {
  Icon
} from "vue-devui";
import {
  Modal, Space, Table,
  TableProps, Row, Col, Button, Popover, Tooltip, message, notification,
} from "ant-design-vue";
import Update_host from "../hosts/update_host.vue";
import {logic} from "../../../wailsjs/go/models";
import {reactive} from "vue";
import {
  CloseSftpClient, OpenLink,
  SftpDelete,
  SftpDir,
  SftpDownload, SftpHomeDir,
  SftpUploadDirectory,
  SftpUploadMultipleFiles
} from "../../../wailsjs/go/logic/Logic";
import {ColumnType} from "ant-design-vue/es/table/interface";
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
  currentDir: '',
  loading: true
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
    width: 80,
    ellipsis: true
  },
  {
    title: '操作',
    key: 'action',
    width: 40
  }
]

function handleResizeColumn(w:number, col:ColumnType<any>) {
  col.width = w;
}

function closeModel(){
  if (props.tid) {
    CloseSftpClient(props.tid).then().catch(e=>{
      notification.error({
        message: '关闭sftp客户端失败',
        description: e,
        duration: null
      });
    });
  }
  // reset reactive
  Object.assign(state, initState());
  state.visible = false
}
function openModel() {
  if (props.tid) {
    state.visible = true;
    state.loading = true;
    SftpHomeDir(props.tid).then((res: string)=>{
      state.currentDir = res;
      handleFoldList(state.currentDir);
    }).catch(e=>{
      notification.error({
        message: '获取Home路径失败',
        description: e,
        duration: null
      });
    })
  } else {
    message.warning('当前没有打开ssh连接',1);
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
      notification.success({
        message: '文件上传成功',
        description: '目标路径'+state.currentDir,
        duration: 3
      });
      handleFoldList(state.currentDir);
    }).catch(e=>{
      notification.error({
        message: '文件上传失败',
        description: '目标路径：'+state.currentDir+'错误信息：'+ e,
        duration: null
      });
    })
  }
}
function handleUploadFold(){
  if (props.tid) {
    SftpUploadDirectory(props.tid,state.currentDir).then(()=>{
      notification.success({
        message: '文件夹上传成功',
        description: '目标路径'+state.currentDir,
        duration: 3
      });
      handleFoldList(state.currentDir);
    }).catch(e=>{
      notification.error({
        message: '文件夹上传失败',
        description: '目标路径：'+state.currentDir+'错误信息：'+ e,
        duration: null
      });
    })
  }
}
function handleFoldList(dst: string) {
  if (props.tid) {
    state.loading = true
    SftpDir(props.tid,dst).then((res:Array<logic.FileInfo>)=>{
      state.currentDir = dst;
      state.tableData = res;
    }).catch(e=>{
      notification.error({
        message: '目录获取失败',
        description: '目标路径：' + dst + '错误信息：'+ e,
        duration: null
      });
    }).finally(()=>{
      state.loading = false
    })
  }
}
function handleDownload(dst: string) {
  if (props.tid) {
    SftpDownload(props.tid, dst).then(()=>{
      notification.success({
        message: '下载成功',
        description: '远端路径信息：' + dst,
        duration: 3
      });
    }).catch(e=>{
      notification.error({
        message: '下载失败',
        description: '远端路径信息：' + dst + ' 错误信息：'+ e,
        duration: null
      });
    })
  }
}
function handleDelete(dst: string){
  if (props.tid) {
    SftpDelete(props.tid,dst).then(()=>{
      notification.success({
        message: '删除成功',
        description: '远端路径信息：' + dst,
        duration: 3
      });
      handleFoldList(state.currentDir);
    }).catch(e=>{
      notification.error({
        message: '删除失败',
        description: '远端路径信息：' + dst,
        duration: null
      });
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
        <Row :gutter="10" class="header-bar">
          <Col :span="8">
            <Popover :content="state.currentDir" trigger="hover" style="background-color: #7693f5; color: #fff">
              <Button
                  type="link"
                  size="small"
              >
                <template #icon>
                  <Icon  name="icon-folder" color="#3DCCA6" >
                    <template #suffix>
                      <span style="color: #f2f3f5;">当前路径</span>
                    </template>
                  </Icon>
                </template>
              </Button>
            </Popover>
          </Col>
          <Col :span="4">
            <Button
                type="link"
                block
                size="small"
                @click="handleBack"
                :disabled="state.loading"
            >
              <template #icon>
                <Icon  name="icon-go-back" color="#3DCCA6" >
                  <template #suffix>
                    <span style="color: #f2f3f5;">上一级</span>
                  </template>
                </Icon>
              </template>
            </Button>
          </Col>
          <Col :span="4">
            <Button
                type="link"
                block
                size="small"
                @click="handleUploadFile"
                :disabled="state.loading"
                title="多文件上传"
            >
              <template #icon>
                <Icon name="icon-upload" color="#3DCCA6">
                  <template #suffix>
                    <span style="color: #f2f3f5;">多文件</span>
                  </template>
                </Icon>
              </template>
            </Button>
          </Col>
          <Col :span="4">
            <Button
                type="link"
                block
                size="small"
                @click="handleUploadFold"
                :disabled="state.loading"
                title="文件夹上传"
            >
              <template #icon>
                <Icon  name="icon-upload" color="#3DCCA6" >
                  <template #suffix>
                    <span style="color: #f2f3f5;">文件夹</span>
                  </template>
                </Icon>
              </template>
            </Button>
          </Col>
          <Col>
            <Button
                type="link"
                block
                size="small"
                @click="OpenLink('https://trzsz.github.io/cn/')"
                title="终端支持trzsz,点击打开安装和使用Trzsz教程"
            >Trzsz</Button>
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
      <template #default>
          <Table
              :rowKey="(record:logic.FileInfo) => record.name"
              :dataSource="state.tableData"
              :columns="columns"
              :pagination="{ pageSize: 10 ,showSizeChanger: true}"
              sticky
              :scroll="{ y: '50vh' }"
              @resizeColumn="handleResizeColumn"
              size="middle"
              :loading="state.loading"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'name'">
                <Button
                    v-if="record.is_dir"
                    type="link"
                    size="small"
                    @dblclick="handleFoldList(record.full_path)"
                    title="双击进入文件夹"
                >
                  <template #icon>
                    <Icon  name="icon-open-folder-2" color="#3DCCA6" >
                      <template #suffix>
                        <span style="color: #f2f3f5;">{{ record.name }}</span>
                      </template>
                    </Icon>
                  </template>
                </Button>
                <Button
                    v-else
                    type="link"
                    size="small"
                    @dblclick="handleDownload(record.full_path)"
                    title="双击下载文件"
                >
                  <template #icon>
                    <Icon name="icon-file">
                      <template #suffix>
                        <span style="color: #f2f3f5;">{{ record.name }}</span>
                      </template>
                    </Icon>
                  </template>
                </Button>
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
                      type="text"
                      ghost
                      size="small"
                      @click="handleDownload(record.full_path)"
                  >
                    <template #icon>
                      <Tooltip placement="bottom" title="下载">
                        <Icon name="icon-download" color="#f2f3f5"/>
                      </Tooltip>
                    </template>
                  </Button>

                  <Popover trigger="click">
                    <Button type="text" ghost danger size="small">
                      <template #icon>
                        <Tooltip placement="bottom" title="删除">
                          <Icon name="delete" color="#aa1111"/>
                        </Tooltip>
                      </template>
                    </Button>
                    <template #content>
                      <Button type="text" ghost danger size="small" @click="handleDelete(record.full_path)">
                        确认
                      </Button>
                    </template>
                  </Popover>
                  </Space>
              </template>
            </template>
          </Table>
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
  overflow: hidden;
}
/deep/.ant-table-body {
  min-height: 50vh !important;
}
/deep/.devui-icon__container {
  display: block;
}
</style>

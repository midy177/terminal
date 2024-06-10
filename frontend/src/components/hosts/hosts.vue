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
const state = reactive({
  visible: false,
  tableData: <Array<DataType>>[],
})
interface DataType {
  key: string
  name: string
  age: number
  sex: string
  address: string
  tags: string[]
}
const columns: STableColumnsType<DataType> = [
  {
    title: '姓名',
    dataIndex: 'name',
    key: 'name',
    width: 120
  },
  {
    title: '年龄',
    dataIndex: 'age',
    key: 'age',
    width: 100
  },
  {
    title: '性别',
    dataIndex: 'sex',
    key: 'sex',
    width: 100
  },
  {
    title: '职业',
    key: 'tags',
    dataIndex: 'tags',
    width: 220
  },
  {
    title: '地址',
    dataIndex: 'address',
    key: 'address',
    ellipsis: true
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
onMounted(()=>{
  for (let i = 0; i < 20; i++) {
    state.tableData.push({
      key: i.toString(),
      name: ['张三', '李四', '王五', '马六'][i % 4],
      age: [18, 30, 26, 45][i % 4],
      sex: ['男', '女'][i % 2],
      address: ['北京', '上海', '天津', '重庆'][i % 4] + '市某某区某某大街520号',
      tags: [['前端', '后端'], ['后端'], ['前端', '产品', '项目'], ['测试']][i % 4]
    })
  }
})
</script>

<template>
  <i class="icon-go-tree" @click="openModel"></i>
<!--  <Icon class="iconfont icon-tree"></Icon>-->
  <FixedOverlay v-model="state.visible" class="hosts-fixed-overlay">
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
            <Button
                icon="add"
                variant="solid"
                title="Add"
            />
          </Col>
          <Col flex="auto">
            <Breadcrumb>
              <BreadcrumbItem to="{ path: '/' }">Homepage</BreadcrumbItem>
              <BreadcrumbItem>
                <a href="/">DevUI</a>
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

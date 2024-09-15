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
      <Row type="flex" class="header-bar">
        <Col flex="auto">
          <Upload
              accept=".cast"
              :customRequest="customRequest"
              :showUploadList="false"
          >
            <Button>
              <upload-outlined></upload-outlined>
              选择录屏文件
            </Button>
          </Upload>
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
      <div ref="playerContainer"></div>
  </Modal>
</template>

<script setup>
import {nextTick, onUnmounted, reactive, ref} from 'vue';
import * as AsciinemaPlayer from 'asciinema-player';
import 'asciinema-player/dist/bundle/asciinema-player.css';
import {Button, Col, message, Modal, Row, Upload} from "ant-design-vue";
import { UploadOutlined } from '@ant-design/icons-vue';

const playerContainer = ref(null);
const initState = () => ({
  visible: false,
  player: null,
  objectUrl: null,
  playerOptions: {
    cols: 80,
    rows: 20,
  },
});

const state = reactive(initState())

function openModel() {
  state.visible = true;
  nextTick(() => {
    state.player = AsciinemaPlayer.create({data: ''}, playerContainer.value, state.playerOptions);
  })
}

function closeModel() {
  try {
    if (state.objectUrl) {
      URL.revokeObjectURL(state.objectUrl);
    }
    if (state.player) {
      state.player.dispose();
    }
  } catch (error) {
    console.error('释放播放器资源失败:', error);
  } finally {
    state.visible = false;
    state.player = undefined;
  }
}

defineExpose({
    openModel,
  })

const customRequest = ({ file, onSuccess, onError }) => {
  if (!file.name.endsWith('.cast')) {
    message.error('只能上传 .cast 文件!');
    onError(new Error('文件类型不正确'));
    return;
  }

  if (file.size / 1024 / 1024 > 2) {
    message.error('文件必须小于 2MB!');
    onError(new Error('文件太大'));
    return;
  }
  // 如果之前有创建的 URL,先撤销它
  if (state.objectUrl) {
    URL.revokeObjectURL(state.objectUrl);
  }

  state.objectUrl = URL.createObjectURL(file);

  const src = {
    driver: 'recording',
    url: state.objectUrl,
    parser: 'asciicast'
  };

  try {
    if (state.player) {
      state.player.dispose();
      // state.player.src = src;
    }
    state.player = AsciinemaPlayer.create(src, playerContainer.value, state.playerOptions);
    onSuccess();
  } catch (error) {
    onError(error);
    message.error('播放器创建失败'+error.message);
  }
};

onUnmounted(
  () => {
    closeModel();
  },
)

</script>

<style scoped>
.ascii-player {
  width: 100%;
  max-width: 80vw;
  max-height: 70vh;
  margin: 0 auto;
}
</style>

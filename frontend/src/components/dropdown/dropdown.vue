<script setup lang="ts">
import { Dropdown,Menu,MenuItem} from 'ant-design-vue';

</script>

<template>
  <Dropdown overlayClassName="dropdown-wrap" placement="bottomRight">
      <i class="iconfont icon-add" @click.prevent></i>
    <template #overlay>
      <Menu>
        <MenuItem>
          <a>1st menu item</a>
        </MenuItem>
        <MenuItem>
          <a>2nd menu item</a>
        </MenuItem>
        <MenuItem>
          <a>3rd menu item</a>
        </MenuItem>
      </Menu>
    </template>
  </Dropdown>
</template>

<style scoped lang="less">
.dropdown-wrap {
  background-color: #1A1B1E;
}
</style>

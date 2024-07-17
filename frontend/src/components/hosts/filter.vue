<script setup lang="ts">
import { PropType, ref, watch } from 'vue'
import { Input } from 'ant-design-vue';
const props = defineProps({
  value: String,
  onChange: Function as PropType<(value: string) => void>
})
const emit = defineEmits(['change'])
const innerValue = ref(props.value)

const onChange = (event: Event) => {
  emit('change', (event.target as HTMLInputElement).value)
}

watch(
    () => props.value,
    newValue => {
      if (newValue !== innerValue.value) {
        innerValue.value = newValue
      }
    }
)
</script>

<template>
  <Input v-bind="$attrs" v-model="innerValue" @change="onChange" />
</template>

<script setup lang="ts">
// 金额输入组件 - 支持高精度数字输入，失焦时自动格式化
import { ref, watch } from 'vue'
import Decimal from 'decimal.js'

const props = defineProps<{
  modelValue: string // 金额值（字符串类型，避免精度丢失）
  currency?: string // 货币代码，显示在输入框前缀
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

// 显示值，用于输入过程中的临时值
const displayValue = ref(props.modelValue)

// 监听外部值变化，同步到显示值
watch(
  () => props.modelValue,
  (val) => {
    displayValue.value = val
  }
)

// 输入时更新显示值（不立即触发v-model更新，避免输入体验差）
function handleInput(val: string) {
  displayValue.value = val
}

// 失焦时格式化金额 - 使用Decimal.js确保精度，无效输入重置为0
function handleBlur() {
  try {
    const d = new Decimal(displayValue.value || '0')
    const formatted = d.toString()
    emit('update:modelValue', formatted)
    displayValue.value = formatted
  } catch {
    emit('update:modelValue', '0')
    displayValue.value = '0'
  }
}
</script>

<template>
  <el-input
    :model-value="displayValue"
    @input="handleInput"
    @blur="handleBlur"
    type="text"
    clearable
  >
    <!-- 有货币代码时，在输入框前显示货币前缀 -->
    <template v-if="currency" #prepend>
      {{ currency }}
    </template>
  </el-input>
</template>

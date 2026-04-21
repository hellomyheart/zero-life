<script setup lang="ts">
import { ref, watch } from 'vue'
import Decimal from 'decimal.js'

const props = defineProps<{
  modelValue: string
  currency?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const displayValue = ref(props.modelValue)

watch(
  () => props.modelValue,
  (val) => {
    displayValue.value = val
  }
)

function handleInput(val: string) {
  displayValue.value = val
}

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
    <template v-if="currency" #prepend>
      {{ currency }}
    </template>
  </el-input>
</template>

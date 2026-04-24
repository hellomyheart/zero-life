<script setup lang="ts">
// 货币选择组件 - 下拉选择已启用的货币
import type { Currency } from '@/types/currency'
import { computed } from 'vue'

const props = defineProps<{
  modelValue: string // 当前选中的货币代码
  currencies: Currency[] // 可选货币列表
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

// 只显示已启用的货币
const enabledCurrencies = computed(() => props.currencies.filter((c) => c.is_enabled))

// 选择货币时触发更新
function handleChange(val: string) {
  emit('update:modelValue', val)
}
</script>

<template>
  <!-- filterable支持搜索过滤 -->
  <el-select
    :model-value="modelValue"
    @update:model-value="handleChange"
    :placeholder="$t('common.selectPlaceholder')"
    filterable
  >
    <el-option
      v-for="currency in enabledCurrencies"
      :key="currency.id"
      :label="`${currency.code} - ${currency.name}`"
      :value="currency.code"
    />
  </el-select>
</template>

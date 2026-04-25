<script setup lang="ts">
// 货币选择组件 - 下拉选择已启用的货币
// 支持两种模式：
// 1. 按 ID 选择（mode="id"）：v-model 绑定 currency_id (number)，用于创建账户等场景
// 2. 按代码选择（mode="code"）：v-model 绑定 currency code (string)，用于筛选等场景
import type { Currency } from '@/types/currency'
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  modelValue: string | number   // 当前选中的货币ID或代码
  currencies: Currency[]        // 可选货币列表
  mode?: 'id' | 'code'         // 选择模式，默认按代码选择
  disabled?: boolean            // 是否禁用
}>(), {
  mode: 'code',
  disabled: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: string | number]
}>()

// 只显示已启用的货币
const enabledCurrencies = computed(() => props.currencies.filter((c) => c.is_enabled))

// 选择货币时触发更新，根据模式返回 ID 或代码
function handleChange(val: string | number) {
  emit('update:modelValue', val)
}
</script>

<template>
  <!-- filterable支持搜索过滤 -->
  <el-select
    :model-value="modelValue"
    @update:model-value="handleChange"
    :placeholder="$t('common.selectPlaceholder')"
    :disabled="disabled"
    filterable
  >
    <el-option
      v-for="currency in enabledCurrencies"
      :key="currency.id"
      :label="`${currency.code} - ${currency.name}`"
      :value="mode === 'id' ? currency.id : currency.code"
    />
  </el-select>
</template>

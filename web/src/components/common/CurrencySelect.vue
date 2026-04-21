<script setup lang="ts">
import type { Currency } from '@/types/currency'
import { computed } from 'vue'

const props = defineProps<{
  modelValue: string
  currencies: Currency[]
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const enabledCurrencies = computed(() => props.currencies.filter((c) => c.is_enabled))

function handleChange(val: string) {
  emit('update:modelValue', val)
}
</script>

<template>
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

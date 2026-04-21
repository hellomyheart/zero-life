<script setup lang="ts">
import { ref, watch } from 'vue'
import dayjs from 'dayjs'

const props = defineProps<{
  startDate: string
  endDate: string
}>()

const emit = defineEmits<{
  'update:startDate': [value: string]
  'update:endDate': [value: string]
}>()

const dateRange = ref<[string, string]>([
  props.startDate || dayjs().startOf('month').format('YYYY-MM-DD'),
  props.endDate || dayjs().endOf('month').format('YYYY-MM-DD'),
])

watch(
  () => [props.startDate, props.endDate],
  ([start, end]) => {
    if (start && end) {
      dateRange.value = [start, end]
    }
  }
)

function handleChange(val: [string, string] | null) {
  if (val) {
    emit('update:startDate', val[0])
    emit('update:endDate', val[1])
  }
}
</script>

<template>
  <el-date-picker
    v-model="dateRange"
    type="daterange"
    range-separator="-"
    :start-placeholder="$t('common.startDate')"
    :end-placeholder="$t('common.endDate')"
    value-format="YYYY-MM-DD"
    @change="handleChange"
  />
</template>

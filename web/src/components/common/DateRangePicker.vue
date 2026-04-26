<script setup lang="ts">
// 日期范围选择组件 - 选择开始和结束日期，默认当月
import { ref, watch } from 'vue'
import dayjs from 'dayjs'

const props = defineProps<{
  startDate: string // 开始日期
  endDate: string // 结束日期
}>()

const emit = defineEmits<{
  'update:startDate': [value: string] // 更新开始日期
  'update:endDate': [value: string] // 更新结束日期
}>()

// 日期范围，默认为当月第一天到最后一天
const dateRange = ref<[string, string]>([
  props.startDate || dayjs().startOf('month').format('YYYY-MM-DD'),
  props.endDate || dayjs().endOf('month').format('YYYY-MM-DD'),
])

// 监听外部日期变化，同步到内部
watch(
  () => [props.startDate, props.endDate],
  ([start, end]) => {
    if (start && end) {
      dateRange.value = [start, end]
    }
  }
)

// 选择日期范围后，分别触发开始和结束日期的更新
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
    style="width: 100%"
    @change="handleChange"
  />
</template>

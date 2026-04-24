<script setup lang="ts">
/**
 * 趋势报表页面
 * 功能：
 * - 展示收入、支出和净额随时间的变化趋势（折线图）
 * - 显示每个时间段的详细收支数据表格
 * - 支持自定义日期范围查询
 * 
 * 数据来源：后端 /reports/trend 接口
 * 使用组件：DateRangePicker（日期选择器）、LineChart（折线图）
 */
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { trend } from '@/api/report'
import { formatAmount } from '@/utils/format'
import dayjs from 'dayjs'
import type { TrendResp, ReportReq } from '@/types/report'
import DateRangePicker from '@/components/common/DateRangePicker.vue'
import LineChart from '@/components/charts/LineChart.vue'

const { t } = useI18n()

// 趋势报表数据
const data = ref<TrendResp | null>(null)
// 加载状态
const loading = ref(false)

// 查询参数：默认显示最近 6 个月的趋势数据
const params = reactive<ReportReq>({
  start_date: dayjs().subtract(6, 'month').startOf('month').format('YYYY-MM-DD'),
  end_date: dayjs().endOf('month').format('YYYY-MM-DD'),
})

/**
 * 获取趋势报表数据
 * 调用后端 API 获取指定日期范围内按时间分组的收支趋势数据
 */
async function fetchData() {
  loading.value = true
  try {
    data.value = await trend(params) as unknown as TrendResp
  } catch (error) {
    console.error('Failed to fetch trend data:', error)
  } finally {
    loading.value = false
  }
}

// 组件挂载时自动加载数据
onMounted(fetchData)
</script>

<template>
  <div class="report-page">
    <h2>{{ t('report.trend') }}</h2>
    
    <!-- 日期范围选择器 -->
    <el-card style="margin-bottom: 16px">
      <DateRangePicker
        v-model:start-date="params.start_date"
        v-model:end-date="params.end_date"
      />
      <el-button type="primary" style="margin-left: 12px" @click="fetchData" :loading="loading">{{ t('common.search') }}</el-button>
    </el-card>

    <template v-if="data">
      <!-- 收支趋势折线图 -->
      <el-card>
        <LineChart
          :data="data.items.map((item) => ({ date: item.date, income: Number(item.income), expense: Number(item.expense), net: Number(item.net) }))"
          x-field="date"
          :y-fields="[{ field: 'income', name: t('report.income') }, { field: 'expense', name: t('report.expense') }, { field: 'net', name: t('report.netIncome') }]"
          :title="t('report.trend')"
        />
      </el-card>

      <!-- 收支趋势明细表 -->
      <el-card style="margin-top: 20px">
        <el-table :data="data.items" stripe>
          <el-table-column prop="date" :label="t('common.date')" />
          <el-table-column prop="income" :label="t('report.income')">
            <template #default="{ row }">{{ formatAmount(row.income) }}</template>
          </el-table-column>
          <el-table-column prop="expense" :label="t('report.expense')">
            <template #default="{ row }">{{ formatAmount(row.expense) }}</template>
          </el-table-column>
          <el-table-column prop="net" :label="t('report.netIncome')">
            <template #default="{ row }">{{ formatAmount(row.net) }}</template>
          </el-table-column>
        </el-table>
      </el-card>
    </template>
  </div>
</template>

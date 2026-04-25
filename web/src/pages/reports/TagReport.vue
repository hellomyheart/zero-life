<script setup lang="ts">
// 标签报表 - 按标签展示收支占比饼图和明细表
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { tag } from '@/api/report'
import { formatAmount, formatPercent } from '@/utils/format'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import type { CategoryReportResp, ReportReq } from '@/types/report' // TODO: CategoryReportResp is used as a stand-in; no TagReportResp type exists yet
import DateRangePicker from '@/components/common/DateRangePicker.vue'
import PieChart from '@/components/charts/PieChart.vue'

const { t } = useI18n()

const data = ref<CategoryReportResp | null>(null)
const loading = ref(false)

const params = ref<ReportReq>({
  start_date: dayjs().startOf('month').format('YYYY-MM-DD'),
  end_date: dayjs().endOf('month').format('YYYY-MM-DD'),
})

async function fetchReport() {
  loading.value = true
  try {
    data.value = await tag(params.value) as unknown as CategoryReportResp
  } catch {
    ElMessage.error(t('common.fetchError') || 'Failed to load data')
  } finally {
    loading.value = false
  }
}

onMounted(fetchReport)
</script>

<template>
  <div class="tag-report-page">
    <h2>{{ t('report.tag') }}</h2>
    <el-card style="margin-bottom: 16px">
      <DateRangePicker
        v-model:start-date="params.start_date"
        v-model:end-date="params.end_date"
      />
      <el-button type="primary" style="margin-left: 12px" @click="fetchReport" :loading="loading">{{ t('common.search') }}</el-button>
    </el-card>

    <template v-if="data">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-card>
            <template #header>{{ t('report.incomeByTag') }}</template>
            <PieChart
              :data="data.income_categories.map((c) => ({ name: c.category_name, value: Number(c.amount) }))"
              name-field="name"
              value-field="value"
            />
            <el-table :data="data.income_categories" size="small" style="margin-top: 16px">
              <el-table-column prop="category_name" :label="t('tag.name')" />
              <el-table-column prop="amount" :label="t('transaction.amount')">
                <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
              </el-table-column>
              <el-table-column prop="percentage" :label="t('budget.usageRate')" width="120">
                <template #default="{ row }">{{ formatPercent(row.percentage) }}</template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card>
            <template #header>{{ t('report.expenseByTag') }}</template>
            <PieChart
              :data="data.expense_categories.map((c) => ({ name: c.category_name, value: Number(c.amount) }))"
              name-field="name"
              value-field="value"
            />
            <el-table :data="data.expense_categories" size="small" style="margin-top: 16px">
              <el-table-column prop="category_name" :label="t('tag.name')" />
              <el-table-column prop="amount" :label="t('transaction.amount')">
                <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
              </el-table-column>
              <el-table-column prop="percentage" :label="t('budget.usageRate')" width="120">
                <template #default="{ row }">{{ formatPercent(row.percentage) }}</template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>
      </el-row>
    </template>
  </div>
</template>

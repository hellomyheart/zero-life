<script setup lang="ts">
// 分类报表 - 各分类展示收支占比饼图和明细表
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { category } from '@/api/report'
import { formatAmount, formatPercent } from '@/utils/format'
import dayjs from 'dayjs'
import type { CategoryReportResp, ReportReq } from '@/types/report'
import DateRangePicker from '@/components/common/DateRangePicker.vue'
import PieChart from '@/components/charts/PieChart.vue'

const { t } = useI18n()

const data = ref<CategoryReportResp | null>(null)
const loading = ref(false)

const params = reactive<ReportReq>({
  start_date: dayjs().startOf('month').format('YYYY-MM-DD'),
  end_date: dayjs().endOf('month').format('YYYY-MM-DD'),
})

async function fetchData() {
  loading.value = true
  try {
    data.value = await category(params) as unknown as CategoryReportResp
  } catch {
    // handle error
  } finally {
    loading.value = false
  }
}

onMounted(fetchData)
</script>

<template>
  <div class="report-page">
    <h2>{{ t('report.category') }}</h2>
    <el-card style="margin-bottom: 16px">
      <DateRangePicker
        v-model:start-date="params.start_date"
        v-model:end-date="params.end_date"
      />
      <el-button type="primary" style="margin-left: 12px" @click="fetchData" :loading="loading">{{ t('common.search') }}</el-button>
    </el-card>

    <template v-if="data">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-card>
            <template #header>Income by Category</template>
            <PieChart
              :data="data.income_categories.map((c) => ({ name: c.category_name, value: Number(c.amount) }))"
              name-field="name"
              value-field="value"
            />
            <el-table :data="data.income_categories" size="small" style="margin-top: 16px">
              <el-table-column prop="category_name" label="Category" />
              <el-table-column prop="amount" label="Amount">
                <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
              </el-table-column>
              <el-table-column prop="percentage" label="Percentage">
                <template #default="{ row }">{{ formatPercent(row.percentage) }}</template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card>
            <template #header>Expense by Category</template>
            <PieChart
              :data="data.expense_categories.map((c) => ({ name: c.category_name, value: Number(c.amount) }))"
              name-field="name"
              value-field="value"
            />
            <el-table :data="data.expense_categories" size="small" style="margin-top: 16px">
              <el-table-column prop="category_name" label="Category" />
              <el-table-column prop="amount" label="Amount">
                <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
              </el-table-column>
              <el-table-column prop="percentage" label="Percentage">
                <template #default="{ row }">{{ formatPercent(row.percentage) }}</template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>
      </el-row>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { budget } from '@/api/report'
import { formatAmount } from '@/utils/format'
import dayjs from 'dayjs'
import type { BudgetReportResp, ReportReq } from '@/types/report'
import DateRangePicker from '@/components/common/DateRangePicker.vue'
import BarChart from '@/components/charts/BarChart.vue'

const { t } = useI18n()

const data = ref<BudgetReportResp | null>(null)
const loading = ref(false)

const params = reactive<ReportReq>({
  start_date: dayjs().startOf('month').format('YYYY-MM-DD'),
  end_date: dayjs().endOf('month').format('YYYY-MM-DD'),
})

async function fetchData() {
  loading.value = true
  try {
    data.value = await budget(params) as unknown as BudgetReportResp
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
    <h2>{{ t('report.budget') }}</h2>
    <el-card style="margin-bottom: 16px">
      <DateRangePicker
        v-model:start-date="params.start_date"
        v-model:end-date="params.end_date"
      />
      <el-button type="primary" style="margin-left: 12px" @click="fetchData" :loading="loading">{{ t('common.search') }}</el-button>
    </el-card>

    <template v-if="data">
      <el-card>
        <BarChart
          :data="data.budgets.map((b) => ({ name: b.budget_name, amount: Number(b.amount), spent: Number(b.spent) }))"
          x-field="name"
          :y-fields="[{ field: 'amount', name: 'Budget' }, { field: 'spent', name: 'Spent' }]"
          :title="t('report.budget')"
        />
      </el-card>

      <el-card style="margin-top: 20px">
        <el-table :data="data.budgets" stripe>
          <el-table-column prop="budget_name" :label="t('budget.name')" />
          <el-table-column prop="amount" :label="t('budget.amount')">
            <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
          </el-table-column>
          <el-table-column prop="spent" :label="t('budget.spent')">
            <template #default="{ row }">{{ formatAmount(row.spent) }}</template>
          </el-table-column>
          <el-table-column :label="t('budget.usageRate')">
            <template #default="{ row }">
              <el-progress :percentage="Math.round(row.usage_rate * 100)" />
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </template>
  </div>
</template>

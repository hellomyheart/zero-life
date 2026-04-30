<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { budget } from '@/api/report'
import { formatAmount } from '@/utils/format'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import type { BudgetReportResp, ReportReq } from '@/types/report'
import BarChart from '@/components/charts/BarChart.vue'

const { t } = useI18n()

const data = ref<BudgetReportResp | null>(null)
const loading = ref(false)

const params: ReportReq = {
  start_date: dayjs().startOf('month').format('YYYY-MM-DD'),
  end_date: dayjs().endOf('month').format('YYYY-MM-DD'),
}

async function fetchData() {
  loading.value = true
  try {
    data.value = await budget(params) as unknown as BudgetReportResp
  } catch {
    ElMessage.error(t('common.operationFailed') || 'Failed to load data')
  } finally {
    loading.value = false
  }
}

onMounted(fetchData)
</script>

<template>
  <div class="report-page">
    <h2>{{ t('report.budget') }}</h2>

    <template v-if="data">
      <el-card>
        <BarChart
          :data="data.items.map((b) => ({ name: b.budget_name, amount: Number(b.amount), spent: Number(b.spent) }))"
          x-field="name"
          :y-fields="[{ field: 'amount', name: t('budget.amount') }, { field: 'spent', name: t('budget.spent') }]"
          :title="t('report.budget')"
        />
      </el-card>

      <el-card style="margin-top: 20px">
        <el-table :data="data.items" stripe>
          <el-table-column prop="budget_name" :label="t('budget.name')" />
          <el-table-column prop="amount" :label="t('budget.amount')">
            <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
          </el-table-column>
          <el-table-column prop="spent" :label="t('budget.spent')">
            <template #default="{ row }">{{ formatAmount(row.spent) }}</template>
          </el-table-column>
          <el-table-column :label="t('budget.usageRate')">
            <template #default="{ row }">
              <el-progress :percentage="Math.min(Math.round(row.usage_rate * 100), 100)" />
              <span v-if="row.usage_rate > 1" class="overspent-label">{{ (row.usage_rate * 100).toFixed(1) }}%</span>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </template>
  </div>
</template>

<style scoped>
.overspent-label {
  color: var(--el-color-danger);
  font-size: 12px;
  margin-left: 8px;
}
</style>

<script setup lang="ts">
// 仪表盘页面 - 展示月度收支预算账单消费和资产概况
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { getDashboard } from '@/api/dashboard'
import { formatAmount, formatDate } from '@/utils/format'
import type { DashboardResp } from '@/types/dashboard'

const { t } = useI18n()

const dashboard = ref<DashboardResp | null>(null)
const loading = ref(true)

onMounted(async () => {
  try {
    dashboard.value = await getDashboard() as unknown as DashboardResp
  } catch {
    ElMessage.error(t('common.fetchError') || 'Failed to load dashboard')
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="dashboard-page" v-loading="loading">
    <template v-if="dashboard">
      <el-row :gutter="20">
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>{{ t('dashboard.income') }}</template>
            <div class="amount income">{{ formatAmount(dashboard.month_income) }}</div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>{{ t('dashboard.expense') }}</template>
            <div class="amount expense">{{ formatAmount(dashboard.month_expense) }}</div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>{{ t('dashboard.netIncome') }}</template>
            <div class="amount">{{ formatAmount(dashboard.net_income) }}</div>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px">
        <el-col :span="12">
          <el-card shadow="hover">
            <template #header>{{ t('dashboard.totalBalance') }}</template>
            <div class="amount large">{{ formatAmount(dashboard.total_balance) }}</div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card shadow="hover">
            <template #header>{{ t('dashboard.budgetAlerts') }}</template>
            <el-table :data="dashboard.budget_alerts" size="small" v-if="dashboard.budget_alerts?.length">
              <el-table-column prop="budget_name" :label="t('budget.name')" />
              <el-table-column prop="usage_rate" :label="t('budget.usageRate')">
                <template #default="{ row }">
                  <el-progress :percentage="Math.round(row.usage_rate * 100)" :status="row.status === 'exceeded' ? 'exception' : row.status === 'warning' ? 'warning' : undefined" />
                </template>
              </el-table-column>
            </el-table>
            <el-empty v-else :description="t('common.noData')" :image-size="60" />
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px">
        <el-col :span="12">
          <el-card shadow="hover">
            <template #header>{{ t('dashboard.billReminders') }}</template>
            <el-table :data="dashboard.bill_reminders" size="small" v-if="dashboard.bill_reminders?.length">
              <el-table-column prop="bill_name" :label="t('bill.name')" />
              <el-table-column prop="amount" :label="t('bill.amount')">
                <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
              </el-table-column>
              <el-table-column prop="next_due" :label="t('bill.nextDueDate')">
                <template #default="{ row }">{{ formatDate(row.next_due) }}</template>
              </el-table-column>
            </el-table>
            <el-empty v-else :description="t('common.noData')" :image-size="60" />
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card shadow="hover">
            <template #header>{{ t('dashboard.recentTransactions') }}</template>
            <el-table :data="dashboard.recent_txns" size="small" v-if="dashboard.recent_txns?.length">
              <el-table-column prop="date" :label="t('transaction.date')">
                <template #default="{ row }">{{ formatDate(row.date) }}</template>
              </el-table-column>
              <el-table-column prop="description" :label="t('transaction.description')" />
              <el-table-column prop="amount" :label="t('transaction.amount')">
                <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
              </el-table-column>
            </el-table>
            <el-empty v-else :description="t('common.noData')" :image-size="60" />
          </el-card>
        </el-col>
      </el-row>
    </template>
  </div>
</template>

<style scoped>
.amount {
  font-size: 24px;
  font-weight: 600;
  text-align: center;
  padding: 8px 0;
}

.amount.large {
  font-size: 32px;
}

.amount.income {
  color: #67c23a;
}

.amount.expense {
  color: #f56c6c;
}
</style>

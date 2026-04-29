<script setup lang="ts">
/**
 * 仪表盘页面
 * 展示月度收支、资产总余额、预算预警、账单提醒和最近交易
 */
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { getDashboard } from '@/api/dashboard'
import { formatAmount, formatDate } from '@/utils/format'
import type { DashboardResp } from '@/types/dashboard'
import dayjs from 'dayjs'

const { t } = useI18n()

const dashboard = ref<DashboardResp | null>(null)
const loading = ref(true)

onMounted(async () => {
  try {
    dashboard.value = await getDashboard() as unknown as DashboardResp
  } catch {
    ElMessage.error(t('common.fetchError'))
  } finally {
    loading.value = false
  }
})

/**
 * 获取交易类型标签类型
 * deposit=收入(success), withdrawal=支出(danger), transfer=转账(info)
 */
function getTypeTagType(type: string) {
  switch (type) {
    case 'deposit': return 'success'
    case 'withdrawal': return 'danger'
    case 'transfer': return 'info'
    default: return ''
  }
}

function getTypeLabel(type: string) {
  switch (type) {
    case 'deposit': return t('transaction.deposit')
    case 'withdrawal': return t('transaction.withdrawal')
    case 'transfer': return t('transaction.transfer')
    default: return type
  }
}

/**
 * 格式化交易金额，根据类型添加正负号
 * deposit 显示 +金额（绿色），withdrawal 显示 -金额（红色），transfer 显示金额
 */
function formatTxnAmount(type: string, amount: string) {
  const val = formatAmount(amount)
  switch (type) {
    case 'deposit': return `+${val}`
    case 'withdrawal': return `-${val}`
    default: return val
  }
}

function getAmountColor(type: string) {
  switch (type) {
    case 'deposit': return '#67c23a'
    case 'withdrawal': return '#f56c6c'
    default: return 'var(--el-text-color-primary)'
  }
}
</script>

<template>
  <div class="dashboard-page" v-loading="loading">
    <template v-if="dashboard">
      <!-- 月度收支概览 -->
      <el-row :gutter="20">
        <el-col :xs="24" :sm="8">
          <el-card shadow="hover">
            <template #header>{{ t('dashboard.income') }}</template>
            <div class="amount income">{{ formatAmount(dashboard.month_income) }}</div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="8">
          <el-card shadow="hover">
            <template #header>{{ t('dashboard.expense') }}</template>
            <div class="amount expense">{{ formatAmount(dashboard.month_expense) }}</div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="8">
          <el-card shadow="hover">
            <template #header>{{ t('dashboard.netIncome') }}</template>
            <div class="amount" :class="{ income: parseFloat(dashboard.net_income) >= 0, expense: parseFloat(dashboard.net_income) < 0 }">
              {{ formatAmount(dashboard.net_income) }}
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 净资产/总资产/总负债 + 预算预警 -->
      <el-row :gutter="20" style="margin-top: 20px">
        <el-col :xs="24" :sm="8">
          <el-card shadow="hover">
            <template #header>{{ t('dashboard.totalBalance') }}</template>
            <div class="amount large">{{ formatAmount(dashboard.total_balance) }}</div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="8">
          <el-card shadow="hover">
            <template #header>{{ t('dashboard.totalAssets') }}</template>
            <div class="amount large">{{ formatAmount(dashboard.total_assets) }}</div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="8">
          <el-card shadow="hover">
            <template #header>{{ t('dashboard.totalLiabilities') }}</template>
            <div class="amount large">{{ formatAmount(dashboard.total_liabilities) }}</div>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px">
        <el-col :xs="24">
          <el-card shadow="hover">
            <template #header>{{ t('dashboard.budgetAlerts') }}</template>
            <el-table :data="dashboard.budget_alerts" size="small" v-if="dashboard.budget_alerts?.length">
              <el-table-column prop="budget_name" :label="t('budget.name')" />
              <el-table-column :label="t('budget.spent') + '/' + t('budget.amount')" width="140">
                <template #default="{ row }">{{ formatAmount(row.spent) }} / {{ formatAmount(row.amount) }}</template>
              </el-table-column>
              <el-table-column prop="usage_rate" :label="t('budget.usageRate')">
                <template #default="{ row }">
                  <el-progress :percentage="Math.round(row.usage_rate * 100)" :stroke-width="18" :text-inside="true" :status="row.status === 'overspent' ? 'exception' : row.status === 'warning' ? 'warning' : undefined" />
                </template>
              </el-table-column>
            </el-table>
            <el-empty v-else :description="t('common.noData')" :image-size="60" />
          </el-card>
        </el-col>
      </el-row>

      <!-- 循环交易提醒 + 最近交易 -->
      <el-row :gutter="20" style="margin-top: 20px">
        <el-col :xs="24" :sm="12">
          <el-card shadow="hover">
            <template #header>{{ t('dashboard.recurringReminders') }}</template>
            <el-table :data="dashboard.recurring_reminders" size="small" v-if="dashboard.recurring_reminders?.length">
              <el-table-column prop="name" :label="t('recurringTransaction.description')" />
              <el-table-column prop="amount" :label="t('recurringTransaction.amount')">
                <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
              </el-table-column>
              <el-table-column prop="next_due" :label="t('recurringTransaction.nextOccurrence')">
                <template #default="{ row }">
                  <span :class="{ 'overdue-text': dayjs(row.next_due).isBefore(dayjs(), 'day') }">
                    {{ formatDate(row.next_due) }}
                  </span>
                </template>
              </el-table-column>
            </el-table>
            <el-empty v-else :description="t('common.noData')" :image-size="60" />
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12">
          <el-card shadow="hover">
            <template #header>{{ t('dashboard.recentTransactions') }}</template>
            <el-table :data="dashboard.recent_txns" size="small" v-if="dashboard.recent_txns?.length">
              <el-table-column :label="t('transaction.type')" width="80">
                <template #default="{ row }">
                  <el-tag :type="getTypeTagType(row.type)" size="small">{{ getTypeLabel(row.type) }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="description" :label="t('transaction.description')" show-overflow-tooltip />
              <el-table-column :label="t('transaction.amount')" width="120" align="right">
                <template #default="{ row }">
                  <span :style="{ color: getAmountColor(row.type), fontVariantNumeric: 'tabular-nums' }">
                    {{ formatTxnAmount(row.type, row.amount) }}
                  </span>
                </template>
              </el-table-column>
              <el-table-column prop="date" :label="t('transaction.date')" width="100">
                <template #default="{ row }">{{ formatDate(row.date) }}</template>
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
  font-variant-numeric: tabular-nums;
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

.overdue-text {
  color: #f56c6c;
  font-weight: 600;
}
</style>
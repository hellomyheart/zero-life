<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { getBudget, getHistory } from '@/api/budget'
import { list } from '@/api/transaction'
import { formatAmount, formatPercent, formatDate } from '@/utils/format'
import { ElMessage } from 'element-plus'
import type { Budget, BudgetHistory } from '@/types/budget'
import type { Transaction } from '@/types/transaction'
import { TransactionType } from '@/types/transaction'
import LineChart from '@/components/charts/LineChart.vue'
import Pagination from '@/components/common/Pagination.vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const budget = ref<Budget | null>(null)
const history = ref<BudgetHistory[]>([])
const loading = ref(true)

const transactions = ref<Transaction[]>([])
const txnTotal = ref(0)
const txnLoading = ref(false)
const txnPage = ref(1)
const txnPageSize = ref(20)

const periodRange = computed(() => {
  if (!budget.value) return { start: '', end: '' }
  const now = new Date()
  if (budget.value.period === 'monthly') {
    const start = new Date(now.getFullYear(), now.getMonth(), 1)
    const end = new Date(now.getFullYear(), now.getMonth() + 1, 0)
    return {
      start: `${start.getFullYear()}-${String(start.getMonth() + 1).padStart(2, '0')}-${String(start.getDate()).padStart(2, '0')}`,
      end: `${end.getFullYear()}-${String(end.getMonth() + 1).padStart(2, '0')}-${String(end.getDate()).padStart(2, '0')}`,
    }
  }
  return {
    start: `${now.getFullYear()}-01-01`,
    end: `${now.getFullYear()}-12-31`,
  }
})

const categoryIds = computed(() => {
  if (!budget.value) return ''
  return budget.value.categories.map(c => c.id).join(',')
})

function formatTransactionAmount(row: Transaction): string {
  const num = Number(row.amount)
  if (isNaN(num)) return row.amount
  const absFormatted = Math.abs(num).toLocaleString('en-US', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
  if (row.type === TransactionType.Deposit) return `+${absFormatted}`
  if (row.type === TransactionType.Withdrawal) return `-${absFormatted}`
  return absFormatted
}

function amountClass(row: Transaction): string {
  if (row.type === TransactionType.Deposit) return 'amount-deposit'
  if (row.type === TransactionType.Withdrawal) return 'amount-withdrawal'
  return 'amount-transfer'
}

function typeTagType(row: Transaction): '' | 'success' | 'danger' | 'info' {
  if (row.type === TransactionType.Deposit) return 'success'
  if (row.type === TransactionType.Withdrawal) return 'danger'
  return 'info'
}

function typeLabel(row: Transaction): string {
  if (row.type === TransactionType.Deposit) return t('transaction.deposit')
  if (row.type === TransactionType.Withdrawal) return t('transaction.withdrawal')
  return t('transaction.transfer')
}

async function fetchTransactions() {
  if (!budget.value) return
  txnLoading.value = true
  try {
    const params: Record<string, unknown> = {
      type: TransactionType.Withdrawal,
      start_date: periodRange.value.start,
      end_date: periodRange.value.end,
      category_ids: categoryIds.value,
      page: txnPage.value,
      page_size: txnPageSize.value,
    }
    const res = await list(params)
    const data = res as unknown as { items: Transaction[]; total: number }
    transactions.value = data.items || []
    txnTotal.value = data.total || 0
  } catch {
    ElMessage.error(t('common.fetchError'))
  } finally {
    txnLoading.value = false
  }
}

function handleTxnPageChange(page: number) {
  txnPage.value = page
  fetchTransactions()
}

function handleTxnSizeChange(size: number) {
  txnPageSize.value = size
  txnPage.value = 1
  fetchTransactions()
}

function handleEditTxn(id: number) {
  router.push(`/transactions/${id}/edit`)
}

onMounted(async () => {
  const id = Number(route.params.id)
  try {
    budget.value = await getBudget(id) as unknown as Budget
    history.value = await getHistory(id) as unknown as BudgetHistory[]
    fetchTransactions()
  } catch {
    ElMessage.error(t('common.fetchError'))
  } finally {
    loading.value = false
  }
})

function goBack() {
  router.push('/budgets')
}
</script>

<template>
  <div class="budget-detail-page" v-loading="loading">
    <template v-if="budget">
      <div class="page-header">
        <h2>{{ budget.name }}</h2>
        <el-button @click="goBack">{{ t('common.back') }}</el-button>
      </div>

      <el-row :gutter="16">
        <el-col :xs="24" :sm="8">
          <el-card shadow="hover">
            <template #header>{{ t('budget.amount') }}</template>
            <div class="amount">{{ formatAmount(budget.amount) }}</div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="8">
          <el-card shadow="hover">
            <template #header>{{ t('budget.spent') }}</template>
            <div class="amount expense">{{ formatAmount(budget.spent) }}</div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="8">
          <el-card shadow="hover">
            <template #header>{{ t('budget.remaining') }}</template>
            <div class="amount remaining">{{ formatAmount(budget.remaining) }}</div>
          </el-card>
        </el-col>
      </el-row>

      <el-card style="margin-top: 16px">
        <template #header>{{ t('budget.usageRate') }}</template>
        <el-progress
          :percentage="Math.min(Math.round(budget.usage_rate * 100), 100)"
          :stroke-width="20"
          :status="budget.status === 'overspent' ? 'exception' : budget.status === 'warning' ? 'warning' : undefined"
        />
        <div class="usage-detail">
          {{ formatPercent(budget.usage_rate) }}
        </div>
      </el-card>

      <el-card style="margin-top: 16px">
        <template #header>{{ t('budget.categories') }}</template>
        <el-tag v-for="cat in budget.categories" :key="cat.id" style="margin: 4px">{{ cat.name }}</el-tag>
        <span v-if="!budget.categories?.length" class="no-data">{{ t('common.noData') }}</span>
      </el-card>

      <el-card v-if="history.length" style="margin-top: 16px">
        <template #header>{{ t('budget.history') }}</template>
        <LineChart
          :data="history.map((h) => ({ period: h.period_start?.substring(0, 10) || '', amount: Number(h.amount), spent: Number(h.spent) }))"
          x-field="period"
          :y-fields="[{ field: 'amount', name: t('budget.amount') }, { field: 'spent', name: t('budget.spent') }]"
        />
      </el-card>

      <el-card style="margin-top: 16px">
        <template #header>{{ t('transaction.title') }}</template>
        <el-table :data="transactions" v-loading="txnLoading" stripe>
          <el-table-column prop="date" :label="t('transaction.date')" width="170">
            <template #default="{ row }">{{ formatDate(row.date) }}</template>
          </el-table-column>
          <el-table-column :label="t('transaction.type')" width="90">
            <template #default="{ row }">
              <el-tag :type="typeTagType(row)" size="small" effect="dark">{{ typeLabel(row) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="description" :label="t('transaction.description')" min-width="160" show-overflow-tooltip />
          <el-table-column :label="t('transaction.amount')" width="140" align="right">
            <template #default="{ row }">
              <span :class="amountClass(row)" class="amount-text">{{ formatTransactionAmount(row) }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('transaction.sourceAccount')" width="130" show-overflow-tooltip class-name="hidden-md-and-down-col">
            <template #default="{ row }">{{ row.source?.name || '' }}</template>
          </el-table-column>
          <el-table-column :label="t('transaction.category')" width="120" show-overflow-tooltip class-name="hidden-md-and-down-col">
            <template #default="{ row }">{{ row.category?.name || '' }}</template>
          </el-table-column>
          <el-table-column :label="t('transaction.tags')" width="150" class-name="hidden-md-and-down-col">
            <template #default="{ row }">
              <el-tag
                v-for="tag in (row.tags || []).slice(0, 3)"
                :key="tag.id"
                size="small"
                :color="tag.color"
                style="color: #fff; margin: 2px"
              >{{ tag.name }}</el-tag>
              <span v-if="(row.tags || []).length > 3" class="more-tags">+{{ row.tags.length - 3 }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.edit')" width="80" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="handleEditTxn(row.id)">{{ t('common.edit') }}</el-button>
            </template>
          </el-table-column>
        </el-table>
        <Pagination
          :total="txnTotal"
          :page="txnPage"
          :page-size="txnPageSize"
          @update:page="handleTxnPageChange"
          @update:page-size="handleTxnSizeChange"
        />
      </el-card>
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

.amount.expense {
  color: var(--app-amount-withdrawal);
}

.amount.remaining {
  color: var(--app-amount-deposit);
}

.amount-text {
  font-weight: 600;
  font-variant-numeric: tabular-nums;
}

.amount-deposit {
  color: var(--app-amount-deposit);
}

.amount-withdrawal {
  color: var(--app-amount-withdrawal);
}

.amount-transfer {
  color: var(--app-amount-transfer);
}

.usage-detail {
  margin-top: 8px;
  text-align: center;
  color: var(--app-text-secondary);
  font-size: 14px;
}

.no-data {
  color: var(--app-text-secondary);
  font-size: 14px;
}

.more-tags {
  font-size: 12px;
  color: var(--app-text-secondary);
  margin-left: 4px;
}

@media (max-width: 767px) {
  .amount {
    font-size: 20px;
  }
}

@media (max-width: 1023px) {
  :deep(.hidden-md-and-down-col) {
    display: none;
  }
}
</style>
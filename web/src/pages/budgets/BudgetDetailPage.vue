<script setup lang="ts">
/**
 * 预算详情页面
 * 功能：
 * - 展示预算基本信息（金额、已用、剩余、使用率）
 * - 展示关联分类标签
 * - 展示历史趋势图表
 * - 响应式布局
 */
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { getBudget, getHistory } from '@/api/budget'
import { formatAmount, formatPercent } from '@/utils/format'
import { ElMessage } from 'element-plus'
import type { Budget, BudgetHistory } from '@/types/budget'
import LineChart from '@/components/charts/LineChart.vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const budget = ref<Budget | null>(null)
const history = ref<BudgetHistory[]>([])
const loading = ref(true)

onMounted(async () => {
  const id = Number(route.params.id)
  try {
    budget.value = await getBudget(id) as unknown as Budget
    history.value = await getHistory(id) as unknown as BudgetHistory[]
  } catch {
    ElMessage.error(t('common.fetchError') || 'Failed to load data')
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

@media (max-width: 767px) {
  .amount {
    font-size: 20px;
  }
}
</style>
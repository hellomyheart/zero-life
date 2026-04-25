<script setup lang="ts">
// 预算详情页面 - 展示预算信息和历史消费图表
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { getBudget, getHistory } from '@/api/budget'
import { formatAmount, formatPercent } from '@/utils/format'
import { ElMessage } from 'element-plus'
import type { Budget, BudgetHistory } from '@/types/budget'
import LineChart from '@/components/charts/LineChart.vue'

const { t } = useI18n()
const route = useRoute()

const budget = ref<Budget | null>(null)
const history = ref<BudgetHistory[]>([])
const loading = ref(true)

onMounted(async () => {
  const id = route.params.id as string
  try {
    budget.value = await getBudget(id) as unknown as Budget
    history.value = await getHistory(id) as unknown as BudgetHistory[]
  } catch {
    ElMessage.error(t('common.fetchError') || 'Failed to load data')
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="budget-detail-page" v-loading="loading">
    <template v-if="budget">
      <h2>{{ budget.name }}</h2>
      <el-row :gutter="20">
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>{{ t('budget.amount') }}</template>
            <div class="amount">{{ formatAmount(budget.amount) }}</div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>{{ t('budget.spent') }}</template>
            <div class="amount expense">{{ formatAmount(budget.spent) }}</div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>{{ t('budget.usageRate') }}</template>
            <div class="amount">{{ formatPercent(Number(budget.amount) ? Number(budget.spent) / Number(budget.amount) : 0) }}</div>
          </el-card>
        </el-col>
      </el-row>

      <el-card style="margin-top: 20px">
        <template #header>{{ t('budget.usageRate') }}</template>
        <el-progress
          :percentage="Number(budget.amount) ? Math.round((Number(budget.spent) / Number(budget.amount)) * 100) : 0"
          :stroke-width="20"
          :status="budget.status === 'exceeded' ? 'exception' : budget.status === 'warning' ? 'warning' : undefined"
        />
      </el-card>

      <el-card style="margin-top: 20px">
        <template #header>{{ t('budget.categories') }}</template>
        <el-tag v-for="name in budget.category_names" :key="name" style="margin: 4px">{{ name }}</el-tag>
      </el-card>

      <el-card v-if="history.length" style="margin-top: 20px">
        <template #header>{{ t('budget.history') }}</template>
        <LineChart
          :data="history.map((h) => ({ period: h.period, amount: Number(h.amount), spent: Number(h.spent) }))"
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
  color: #f56c6c;
}
</style>

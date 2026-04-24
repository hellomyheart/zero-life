<script setup lang="ts">
/**
 * 收支报表页面
 * 功能：
 * - 展示指定时间段内的总收入、总支出和净收入
 * - 按账户分组显示收入明细和支出明细
 * - 支持自定义日期范围查询
 * 
 * 数据来源：后端 /reports/income-expense 接口
 * 使用组件：DateRangePicker（日期选择器）
 */
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { incomeExpense } from '@/api/report'
import { formatAmount } from '@/utils/format'
import dayjs from 'dayjs'
import type { IncomeExpenseResp, ReportReq } from '@/types/report'
import DateRangePicker from '@/components/common/DateRangePicker.vue'

const { t } = useI18n()

// 收支报表数据
const data = ref<IncomeExpenseResp | null>(null)
// 加载状态
const loading = ref(false)

// 查询参数：默认显示当前月份的收支数据
const params = reactive<ReportReq>({
  start_date: dayjs().startOf('month').format('YYYY-MM-DD'),
  end_date: dayjs().endOf('month').format('YYYY-MM-DD'),
})

/**
 * 获取收支报表数据
 * 调用后端 API 获取指定日期范围内的收支汇总和按账户分组的明细
 */
async function fetchData() {
  loading.value = true
  try {
    data.value = await incomeExpense(params) as unknown as IncomeExpenseResp
  } catch (error) {
    console.error('Failed to fetch income expense data:', error)
  } finally {
    loading.value = false
  }
}

// 组件挂载时自动加载数据
onMounted(fetchData)
</script>

<template>
  <div class="report-page">
    <h2>{{ t('report.incomeExpense') }}</h2>
    
    <!-- 日期范围选择器 -->
    <el-card style="margin-bottom: 16px">
      <DateRangePicker
        v-model:start-date="params.start_date"
        v-model:end-date="params.end_date"
      />
      <el-button type="primary" style="margin-left: 12px" @click="fetchData" :loading="loading">{{ t('common.search') }}</el-button>
    </el-card>

    <template v-if="data">
      <!-- 收支概览卡片 -->
      <el-row :gutter="20">
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>{{ t('report.income') }}</template>
            <div class="amount income">{{ formatAmount(data.total_income) }}</div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>{{ t('report.expense') }}</template>
            <div class="amount expense">{{ formatAmount(data.total_expense) }}</div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>{{ t('report.netIncome') }}</template>
            <div class="amount">{{ formatAmount(data.net_income) }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 按账户分组的收支明细 -->
      <el-row :gutter="20" style="margin-top: 20px">
        <el-col :span="12">
          <el-card>
            <template #header>{{ t('report.incomeByAccount') }}</template>
            <el-table :data="data.income_by_account" size="small">
              <el-table-column prop="account_name" :label="t('account.name')" />
              <el-table-column prop="amount" :label="t('transaction.amount')">
                <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card>
            <template #header>{{ t('report.expenseByAccount') }}</template>
            <el-table :data="data.expense_by_account" size="small">
              <el-table-column prop="account_name" :label="t('account.name')" />
              <el-table-column prop="amount" :label="t('transaction.amount')">
                <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
              </el-table-column>
            </el-table>
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

.amount.income { color: #67c23a; }
.amount.expense { color: #f56c6c; }
</style>

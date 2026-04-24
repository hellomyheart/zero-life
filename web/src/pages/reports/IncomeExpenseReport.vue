<script setup lang="ts">
// 收支报表页面 - 展示月度收入支出和按账户分类的收支明细
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { incomeExpense } from '@/api/report'
import { formatAmount } from '@/utils/format'
import dayjs from 'dayjs'
import type { IncomeExpenseResp, ReportReq } from '@/types/report'
import DateRangePicker from '@/components/common/DateRangePicker.vue'

const { t } = useI18n()

const data = ref<IncomeExpenseResp | null>(null)
const loading = ref(false)

const params = reactive<ReportReq>({
  start_date: dayjs().startOf('month').format('YYYY-MM-DD'),
  end_date: dayjs().endOf('month').format('YYYY-MM-DD'),
})

async function fetchData() {
  loading.value = true
  try {
    data.value = await incomeExpense(params) as unknown as IncomeExpenseResp
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
    <h2>{{ t('report.incomeExpense') }}</h2>
    <el-card style="margin-bottom: 16px">
      <DateRangePicker
        v-model:start-date="params.start_date"
        v-model:end-date="params.end_date"
      />
      <el-button type="primary" style="margin-left: 12px" @click="fetchData" :loading="loading">{{ t('common.search') }}</el-button>
    </el-card>

    <template v-if="data">
      <el-row :gutter="20">
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>Income</template>
            <div class="amount income">{{ formatAmount(data.total_income) }}</div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>Expense</template>
            <div class="amount expense">{{ formatAmount(data.total_expense) }}</div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>Net Income</template>
            <div class="amount">{{ formatAmount(data.net_income) }}</div>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px">
        <el-col :span="12">
          <el-card>
            <template #header>Income by Account</template>
            <el-table :data="data.income_by_account" size="small">
              <el-table-column prop="account_name" label="Account" />
              <el-table-column prop="amount" label="Amount">
                <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card>
            <template #header>Expense by Account</template>
            <el-table :data="data.expense_by_account" size="small">
              <el-table-column prop="account_name" label="Account" />
              <el-table-column prop="amount" label="Amount">
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

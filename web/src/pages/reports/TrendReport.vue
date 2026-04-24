<script setup lang="ts">
// 趋势报表 - 展示月度收入支出和净值变化趋势
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { trend } from '@/api/report'
import { formatAmount } from '@/utils/format'
import dayjs from 'dayjs'
import type { TrendResp, ReportReq } from '@/types/report'
import DateRangePicker from '@/components/common/DateRangePicker.vue'
import LineChart from '@/components/charts/LineChart.vue'

const { t } = useI18n()

const data = ref<TrendResp | null>(null)
const loading = ref(false)

const params = reactive<ReportReq>({
  start_date: dayjs().subtract(6, 'month').startOf('month').format('YYYY-MM-DD'),
  end_date: dayjs().endOf('month').format('YYYY-MM-DD'),
})

async function fetchData() {
  loading.value = true
  try {
    data.value = await trend(params) as unknown as TrendResp
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
    <h2>{{ t('report.trend') }}</h2>
    <el-card style="margin-bottom: 16px">
      <DateRangePicker
        v-model:start-date="params.start_date"
        v-model:end-date="params.end_date"
      />
      <el-button type="primary" style="margin-left: 12px" @click="fetchData" :loading="loading">{{ t('common.search') }}</el-button>
    </el-card>

    <template v-if="data">
      <el-card>
        <LineChart
          :data="data.items.map((item) => ({ date: item.date, income: Number(item.income), expense: Number(item.expense), net: Number(item.net) }))"
          x-field="date"
          :y-fields="[{ field: 'income', name: 'Income' }, { field: 'expense', name: 'Expense' }, { field: 'net', name: 'Net' }]"
          :title="t('report.trend')"
        />
      </el-card>

      <el-card style="margin-top: 20px">
        <el-table :data="data.items" stripe>
          <el-table-column prop="date" label="Date" />
          <el-table-column prop="income" label="Income">
            <template #default="{ row }">{{ formatAmount(row.income) }}</template>
          </el-table-column>
          <el-table-column prop="expense" label="Expense">
            <template #default="{ row }">{{ formatAmount(row.expense) }}</template>
          </el-table-column>
          <el-table-column prop="net" label="Net">
            <template #default="{ row }">{{ formatAmount(row.net) }}</template>
          </el-table-column>
        </el-table>
      </el-card>
    </template>
  </div>
</template>

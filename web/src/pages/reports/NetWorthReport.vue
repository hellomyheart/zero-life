<script setup lang="ts">
// 净值报表 - 展示资产和债务及净值变化趋势
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { netWorth } from '@/api/report'
import { formatAmount } from '@/utils/format'
import dayjs from 'dayjs'
import type { NetWorthResp, ReportReq } from '@/types/report'
import DateRangePicker from '@/components/common/DateRangePicker.vue'
import LineChart from '@/components/charts/LineChart.vue'

const { t } = useI18n()

const data = ref<NetWorthResp | null>(null)
const loading = ref(false)

const params = reactive<ReportReq>({
  start_date: dayjs().subtract(6, 'month').startOf('month').format('YYYY-MM-DD'),
  end_date: dayjs().endOf('month').format('YYYY-MM-DD'),
})

async function fetchData() {
  loading.value = true
  try {
    data.value = await netWorth(params) as unknown as NetWorthResp
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
    <h2>{{ t('report.netWorth') }}</h2>
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
            <template #header>Total Assets</template>
            <div class="amount income">{{ formatAmount(data.total_assets) }}</div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>Total Liabilities</template>
            <div class="amount expense">{{ formatAmount(data.total_liabilities) }}</div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>Net Worth</template>
            <div class="amount">{{ formatAmount(data.net_worth) }}</div>
          </el-card>
        </el-col>
      </el-row>

      <el-card style="margin-top: 20px">
        <LineChart
          :data="data.items.map((item) => ({ date: item.date, netWorth: Number(item.net_worth) }))"
          x-field="date"
          :y-fields="[{ field: 'netWorth', name: 'Net Worth' }]"
          :title="t('report.netWorth')"
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

.amount.income { color: #67c23a; }
.amount.expense { color: #f56c6c; }
</style>

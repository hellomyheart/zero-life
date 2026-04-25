<script setup lang="ts">
/**
 * 净值报表页面
 * 功能：
 * - 展示用户的总资产、总负债和净值
 * - 按日期范围显示净值变化趋势图表
 * - 支持自定义日期范围查询
 * 
 * 数据来源：后端 /reports/net-worth 接口
 * 使用组件：DateRangePicker（日期选择器）、LineChart（折线图）
 */
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { netWorth } from '@/api/report'
import { formatAmount } from '@/utils/format'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import type { NetWorthResp, ReportReq } from '@/types/report'
import DateRangePicker from '@/components/common/DateRangePicker.vue'
import LineChart from '@/components/charts/LineChart.vue'

const { t } = useI18n()

// 净值报表数据
const data = ref<NetWorthResp | null>(null)
// 加载状态
const loading = ref(false)

// 查询参数：默认显示最近 6 个月的净值趋势
const params = reactive<ReportReq>({
  start_date: dayjs().subtract(6, 'month').startOf('month').format('YYYY-MM-DD'),
  end_date: dayjs().endOf('month').format('YYYY-MM-DD'),
})

/**
 * 获取净值报表数据
 * 调用后端 API 获取指定日期范围内的净值数据
 */
async function fetchData() {
  loading.value = true
  try {
    data.value = await netWorth(params) as unknown as NetWorthResp
  } catch {
    ElMessage.error(t('common.operationFailed') || 'Failed to load data')
  } finally {
    loading.value = false
  }
}

// 组件挂载时自动加载数据
onMounted(fetchData)
</script>

<template>
  <div class="report-page">
    <h2>{{ t('report.netWorth') }}</h2>
    
    <!-- 日期范围选择器 -->
    <el-card style="margin-bottom: 16px">
      <DateRangePicker
        v-model:start-date="params.start_date"
        v-model:end-date="params.end_date"
      />
      <el-button type="primary" style="margin-left: 12px" @click="fetchData" :loading="loading">{{ t('common.search') }}</el-button>
    </el-card>

    <template v-if="data">
      <!-- 净值概览卡片 -->
      <el-row :gutter="20">
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>{{ t('report.totalAssets') }}</template>
            <div class="amount income">{{ formatAmount(data.total_assets) }}</div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>{{ t('report.totalLiabilities') }}</template>
            <div class="amount expense">{{ formatAmount(data.total_liabilities) }}</div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>{{ t('report.netWorth') }}</template>
            <div class="amount">{{ formatAmount(data.net_worth) }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 净值变化趋势图 -->
      <el-card style="margin-top: 20px">
        <LineChart
          :data="data.trend.map((item) => ({ date: item.date, netWorth: Number(item.net_worth) }))"
          x-field="date"
          :y-fields="[{ field: 'netWorth', name: t('report.netWorth') }]"
          :title="t('report.netWorthTrend')"
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

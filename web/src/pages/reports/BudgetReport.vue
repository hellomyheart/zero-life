<script setup lang="ts">
/**
 * 预算报表页面
 * 功能：
 * - 展示各预算的金额、已花费金额和使用率
 * - 使用柱状图对比预算与实际支出
 * - 支持自定义日期范围查询
 * 
 * 数据来源：后端 /reports/budget 接口
 * 使用组件：DateRangePicker（日期选择器）、BarChart（柱状图）
 */
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { budget } from '@/api/report'
import { formatAmount } from '@/utils/format'
import dayjs from 'dayjs'
import type { BudgetReportResp, ReportReq } from '@/types/report'
import DateRangePicker from '@/components/common/DateRangePicker.vue'
import BarChart from '@/components/charts/BarChart.vue'

const { t } = useI18n()

// 预算报表数据
const data = ref<BudgetReportResp | null>(null)
// 加载状态
const loading = ref(false)

// 查询参数：默认显示当前月份的预算数据
const params = reactive<ReportReq>({
  start_date: dayjs().startOf('month').format('YYYY-MM-DD'),
  end_date: dayjs().endOf('month').format('YYYY-MM-DD'),
})

/**
 * 获取预算报表数据
 * 调用后端 API 获取指定日期范围内各预算的使用情况
 */
async function fetchData() {
  loading.value = true
  try {
    data.value = await budget(params) as unknown as BudgetReportResp
  } catch (error) {
    console.error('Failed to fetch budget data:', error)
  } finally {
    loading.value = false
  }
}

// 组件挂载时自动加载数据
onMounted(fetchData)
</script>

<template>
  <div class="report-page">
    <h2>{{ t('report.budget') }}</h2>
    
    <!-- 日期范围选择器 -->
    <el-card style="margin-bottom: 16px">
      <DateRangePicker
        v-model:start-date="params.start_date"
        v-model:end-date="params.end_date"
      />
      <el-button type="primary" style="margin-left: 12px" @click="fetchData" :loading="loading">{{ t('common.search') }}</el-button>
    </el-card>

    <template v-if="data">
      <!-- 预算对比柱状图 -->
      <el-card>
        <BarChart
          :data="data.budgets.map((b) => ({ name: b.budget_name, amount: Number(b.amount), spent: Number(b.spent) }))"
          x-field="name"
          :y-fields="[{ field: 'amount', name: t('budget.amount') }, { field: 'spent', name: t('budget.spent') }]"
          :title="t('report.budget')"
        />
      </el-card>

      <!-- 预算使用明细表 -->
      <el-card style="margin-top: 20px">
        <el-table :data="data.budgets" stripe>
          <el-table-column prop="budget_name" :label="t('budget.name')" />
          <el-table-column prop="amount" :label="t('budget.amount')">
            <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
          </el-table-column>
          <el-table-column prop="spent" :label="t('budget.spent')">
            <template #default="{ row }">{{ formatAmount(row.spent) }}</template>
          </el-table-column>
          <el-table-column :label="t('budget.usageRate')">
            <template #default="{ row }">
              <el-progress :percentage="Math.round(row.usage_rate * 100)" />
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </template>
  </div>
</template>

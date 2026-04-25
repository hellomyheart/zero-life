<script setup lang="ts">
/**
 * 分类报表页面
 * 功能：
 * - 按分类展示收入和支出的占比（饼图）
 * - 显示每个分类的具体金额和百分比
 * - 支持自定义日期范围查询
 * 
 * 数据来源：后端 /reports/category 接口
 * 使用组件：DateRangePicker（日期选择器）、PieChart（饼图）
 */
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { category } from '@/api/report'
import { formatAmount, formatPercent } from '@/utils/format'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import type { CategoryReportResp, ReportReq } from '@/types/report'
import DateRangePicker from '@/components/common/DateRangePicker.vue'
import PieChart from '@/components/charts/PieChart.vue'

const { t } = useI18n()

// 分类报表数据
const data = ref<CategoryReportResp | null>(null)
// 加载状态
const loading = ref(false)

// 查询参数：默认显示当前月份的分类数据
const params = reactive<ReportReq>({
  start_date: dayjs().startOf('month').format('YYYY-MM-DD'),
  end_date: dayjs().endOf('month').format('YYYY-MM-DD'),
})

/**
 * 获取分类报表数据
 * 调用后端 API 获取指定日期范围内按分类汇总的收支数据
 */
async function fetchData() {
  loading.value = true
  try {
    data.value = await category(params) as unknown as CategoryReportResp
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
    <h2>{{ t('report.category') }}</h2>
    
    <!-- 日期范围选择器 -->
    <el-card style="margin-bottom: 16px">
      <DateRangePicker
        v-model:start-date="params.start_date"
        v-model:end-date="params.end_date"
      />
      <el-button type="primary" style="margin-left: 12px" @click="fetchData" :loading="loading">{{ t('common.search') }}</el-button>
    </el-card>

    <template v-if="data">
      <!-- 收入分类和支出分类并排显示 -->
      <el-row :gutter="20">
        <el-col :span="12">
          <el-card>
            <template #header>{{ t('report.incomeByCategory') }}</template>
            <!-- 收入分类饼图 -->
            <PieChart
              :data="data.income_categories.map((c) => ({ name: c.category_name, value: Number(c.amount) }))"
              name-field="name"
              value-field="value"
            />
            <!-- 收入分类明细表 -->
            <el-table :data="data.income_categories" size="small" style="margin-top: 16px">
              <el-table-column prop="category_name" :label="t('category.name')" />
              <el-table-column prop="amount" :label="t('transaction.amount')">
                <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
              </el-table-column>
              <el-table-column prop="percentage" :label="t('report.percentage')">
                <template #default="{ row }">{{ formatPercent(row.percentage) }}</template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card>
            <template #header>{{ t('report.expenseByCategory') }}</template>
            <!-- 支出分类饼图 -->
            <PieChart
              :data="data.expense_categories.map((c) => ({ name: c.category_name, value: Number(c.amount) }))"
              name-field="name"
              value-field="value"
            />
            <!-- 支出分类明细表 -->
            <el-table :data="data.expense_categories" size="small" style="margin-top: 16px">
              <el-table-column prop="category_name" :label="t('category.name')" />
              <el-table-column prop="amount" :label="t('transaction.amount')">
                <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
              </el-table-column>
              <el-table-column prop="percentage" :label="t('report.percentage')">
                <template #default="{ row }">{{ formatPercent(row.percentage) }}</template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>
      </el-row>
    </template>
  </div>
</template>

<!-- 标签报表页面 - 按标签维度展示收支统计，包含数据表格和饼图可视化 -->
<script setup lang="ts">
// 导入 Vue 响应式和生命周期 API
import { ref, onMounted } from 'vue'
// 导入国际化钩子函数
import { useI18n } from 'vue-i18n'
// 导入标签报表 API
import { tag } from '@/api/report'
// 导入金额格式化工具
import { formatAmount } from '@/utils/format'
// 导入 Element Plus 消息提示组件
import { ElMessage } from 'element-plus'
// 导入 dayjs 日期处理库
import dayjs from 'dayjs'
// 导入报表相关类型定义
import type { TagReportResp, ReportReq } from '@/types/report'
// 导入日期范围选择器组件
import DateRangePicker from '@/components/common/DateRangePicker.vue'
// 导入饼图组件
import PieChart from '@/components/charts/PieChart.vue'

// 国际化翻译函数
const { t } = useI18n()

// 标签报表数据
const data = ref<TagReportResp | null>(null)
// 加载状态
const loading = ref(false)

// 查询参数，默认为当月的时间范围
const params = ref<ReportReq>({
  start_date: dayjs().startOf('month').format('YYYY-MM-DD'),
  end_date: dayjs().endOf('month').format('YYYY-MM-DD'),
})

// 获取标签报表数据
async function fetchReport() {
  loading.value = true
  try {
    data.value = await tag(params.value) as unknown as TagReportResp
  } catch {
    ElMessage.error(t('common.operationFailed') || 'Failed to load data')
  } finally {
    loading.value = false
  }
}

// 页面挂载时自动加载报表数据
onMounted(fetchReport)
</script>

<template>
  <div class="tag-report-page">
    <!-- 页面标题 -->
    <h2>{{ t('report.tag') }}</h2>
    <!-- 日期范围筛选区域 -->
    <el-card style="margin-bottom: 16px">
      <DateRangePicker
        v-model:start-date="params.start_date"
        v-model:end-date="params.end_date"
      />
      <el-button type="primary" style="margin-left: 12px" @click="fetchReport" :loading="loading">{{ t('common.search') }}</el-button>
    </el-card>

    <!-- 报表数据展示区域 -->
    <template v-if="data">
      <!-- 标签收支数据表格 -->
      <el-card>
        <el-table :data="data.items" stripe>
          <el-table-column prop="tag_name" :label="t('tag.name')" />
          <el-table-column prop="income" :label="t('report.income')">
            <template #default="{ row }">{{ formatAmount(row.income) }}</template>
          </el-table-column>
          <el-table-column prop="expense" :label="t('report.expense')">
            <template #default="{ row }">{{ formatAmount(row.expense) }}</template>
          </el-table-column>
        </el-table>
      </el-card>

      <!-- 饼图可视化区域 -->
      <el-row :gutter="20" style="margin-top: 20px">
        <!-- 按标签的收入饼图 -->
        <el-col :span="12">
          <el-card>
            <template #header>{{ t('report.incomeByTag') }}</template>
            <PieChart
              :data="data.items.map((item) => ({ name: item.tag_name, value: Number(item.income) })).filter((d) => d.value > 0)"
              name-field="name"
              value-field="value"
            />
          </el-card>
        </el-col>
        <!-- 按标签的支出饼图 -->
        <el-col :span="12">
          <el-card>
            <template #header>{{ t('report.expenseByTag') }}</template>
            <PieChart
              :data="data.items.map((item) => ({ name: item.tag_name, value: Number(item.expense) })).filter((d) => d.value > 0)"
              name-field="name"
              value-field="value"
            />
          </el-card>
        </el-col>
      </el-row>
    </template>
  </div>
</template>

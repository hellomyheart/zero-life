<script setup lang="ts">
// 数据导出页面 - 支持导出为CSV和JSON格式
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { exportTransactions, exportAccounts, exportBudgets, exportCategories, exportTags } from '@/api/export'
import { ElMessage } from 'element-plus'
import type { ExportReq } from '@/types/export'

const { t } = useI18n()

const dateRange = ref<[string, string] | null>(null)
const format = ref<'csv' | 'json'>('csv')

const exportItems = [
  { key: 'transactions', fn: exportTransactions },
  { key: 'accounts', fn: exportAccounts },
  { key: 'budgets', fn: exportBudgets },
  { key: 'categories', fn: exportCategories },
  { key: 'tags', fn: exportTags },
]

async function handleExport(fn: (params: ExportReq) => Promise<unknown>, _key: string) {
  try {
    const params: ExportReq = { format: format.value }
    if (dateRange.value) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    await fn(params)
    ElMessage.success(t('common.success'))
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}
</script>

<template>
  <div class="export-page">
    <div class="page-header">
      <h2>{{ t('export.title') }}</h2>
    </div>

    <el-form inline style="margin-bottom: 20px;">
      <el-form-item :label="t('common.startDate')">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          value-format="YYYY-MM-DD"
          :start-placeholder="t('common.startDate')"
          :end-placeholder="t('common.endDate')"
        />
      </el-form-item>
      <el-form-item :label="t('export.format')">
        <el-select v-model="format">
          <el-option label="CSV" value="csv" />
          <el-option label="JSON" value="json" />
        </el-select>
      </el-form-item>
    </el-form>

    <el-row :gutter="16">
      <el-col :span="8" v-for="item in exportItems" :key="item.key">
        <el-card shadow="hover" style="margin-bottom: 16px;">
          <template #header>{{ t(`export.${item.key}`) }}</template>
          <el-button type="primary" @click="handleExport(item.fn, item.key)">{{ t('common.export') }}</el-button>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.page-header h2 {
  margin: 0;
}
</style>

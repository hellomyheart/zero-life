<script setup lang="ts">
// 预算列表页面 - 展示预算已用和剩余使用率等
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { list, remove } from '@/api/budget'
import { formatAmount } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useCategoryStore } from '@/stores/category'
import type { Budget } from '@/types/budget'
import { BudgetPeriod } from '@/types/budget'

const { t } = useI18n()
const router = useRouter()
// 分类状态管理 - 用于获取分类下拉选项
const categoryStore = useCategoryStore()

const budgets = ref<Budget[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref<string | null>(null)

const form = ref({
  name: '',
  amount: '0',
  period: 'monthly' as BudgetPeriod,
  category_ids: [] as string[],
  start_date: '',
})

async function fetchBudgets() {
  loading.value = true
  try {
    budgets.value = await list() as unknown as Budget[]
  } catch {
    // handle error
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  dialogTitle.value = t('budget.create')
  editingId.value = null
  form.value = { name: '', amount: '0', period: BudgetPeriod.Monthly, category_ids: [], start_date: '' }
  dialogVisible.value = true
}

function handleEdit(budget: Budget) {
  dialogTitle.value = t('budget.edit')
  editingId.value = budget.id
  form.value = {
    name: budget.name,
    amount: budget.amount,
    period: budget.period,
    category_ids: budget.category_ids,
    start_date: budget.start_date,
  }
  dialogVisible.value = true
}

function handleDetail(id: string) {
  router.push(`/budgets/${id}`)
}

async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm(t('budget.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchBudgets()
  } catch {
    // cancelled or error
  }
}

async function handleSubmit() {
  if (!form.value.name) {
    ElMessage.warning(t('common.required'))
    return
  }
  try {
    const { create: createBudget, update: updateBudget } = await import('@/api/budget')
    if (editingId.value) {
      await updateBudget(editingId.value, form.value)
    } else {
      await createBudget(form.value)
    }
    ElMessage.success(t('common.success'))
    dialogVisible.value = false
    await fetchBudgets()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

function getStatusType(status: string) {
  if (status === 'exceeded') return 'danger'
  if (status === 'warning') return 'warning'
  return 'success'
}

onMounted(async () => {
  // 打开页面时加载分类列表，供下拉选择使用
  await categoryStore.fetchCategories()
  await fetchBudgets()
})
</script>

<template>
  <div class="budget-list-page">
    <div class="page-header">
      <h2>{{ t('budget.title') }}</h2>
      <el-button type="primary" @click="handleCreate">{{ t('budget.create') }}</el-button>
    </div>

    <el-table :data="budgets" v-loading="loading" stripe>
      <el-table-column prop="name" :label="t('budget.name')" />
      <el-table-column prop="amount" :label="t('budget.amount')" width="150">
        <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
      </el-table-column>
      <el-table-column prop="spent" :label="t('budget.spent')" width="150">
        <template #default="{ row }">{{ formatAmount(row.spent) }}</template>
      </el-table-column>
      <el-table-column :label="t('budget.usageRate')" width="200">
        <template #default="{ row }">
          <el-progress :percentage="Math.round(row.spent / row.amount * 100)" :status="getStatusType(row.status) === 'danger' ? 'exception' : getStatusType(row.status) === 'warning' ? 'warning' : undefined" />
        </template>
      </el-table-column>
      <el-table-column :label="t('budget.status')" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)" size="small">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.edit')" width="200" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleDetail(row.id)">{{ t('common.confirm') }}</el-button>
          <el-button link type="primary" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
          <el-button link type="danger" @click="handleDelete(row.id)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item :label="t('budget.name')">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="t('budget.amount')">
          <el-input v-model="form.amount" />
        </el-form-item>
        <el-form-item :label="t('budget.period')">
          <el-select v-model="form.period">
            <el-option label="Monthly" value="monthly" />
            <el-option label="Quarterly" value="quarterly" />
            <el-option label="Yearly" value="yearly" />
          </el-select>
        </el-form-item>
        <!-- 分类多选下拉框 - 选择预算关联的多个分类 -->
        <el-form-item :label="t('transaction.category')">
          <el-select v-model="form.category_ids" multiple :placeholder="t('common.selectPlaceholder')" filterable clearable>
            <el-option v-for="cat in categoryStore.categories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('budget.startDate')">
          <el-date-picker v-model="form.start_date" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
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

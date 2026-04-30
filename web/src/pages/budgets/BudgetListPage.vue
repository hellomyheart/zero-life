<script setup lang="ts">
/**
 * 预算列表页面
 * 功能：
 * - 展示预算列表，含已用/剩余/使用率/状态
 * - 支持创建、编辑、删除预算
 * - 分类选择支持树形结构
 * - 响应式布局
 */
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { list, create, update, remove } from '@/api/budget'
import { formatAmount } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useCategoryStore } from '@/stores/category'
import type { Budget } from '@/types/budget'
import { BudgetPeriod } from '@/types/budget'
import type { Category } from '@/types/category'

const { t } = useI18n()
const router = useRouter()
const categoryStore = useCategoryStore()

const budgets = ref<Budget[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref<number | null>(null)

const form = ref({
  name: '',
  amount: '0',
  period: BudgetPeriod.Monthly,
  category_ids: [] as number[],
})

const categoryTreeData = computed(() => {
  function transform(categories: Category[]): { value: number; label: string; children?: { value: number; label: string }[] }[] {
    return categories.map(cat => {
      const node: { value: number; label: string; children?: { value: number; label: string }[] } = {
        value: cat.id,
        label: cat.name,
      }
      if (cat.children?.length) {
        node.children = transform(cat.children)
      }
      return node
    })
  }
  return transform(categoryStore.categories ?? [])
})

async function fetchBudgets() {
  loading.value = true
  try {
    budgets.value = await list() as unknown as Budget[]
  } catch {
    ElMessage.error(t('common.fetchError') || 'Failed to load data')
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  dialogTitle.value = t('budget.create')
  editingId.value = null
  form.value = { name: '', amount: '', period: BudgetPeriod.Monthly, category_ids: [] }
  dialogVisible.value = true
}

function handleEdit(budget: Budget) {
  dialogTitle.value = t('budget.edit')
  editingId.value = budget.id
  form.value = {
    name: budget.name,
    amount: budget.amount,
    period: budget.period as BudgetPeriod,
    category_ids: budget.categories?.map(c => c.id) || [],
  }
  dialogVisible.value = true
}

function handleDetail(id: number) {
  router.push(`/budgets/${id}`)
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm(t('budget.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchBudgets()
  } catch (err) {
    if (err !== 'cancel' && (err as { message?: string }).message !== 'cancel') ElMessage.error(t('common.fetchError') || 'Failed to load data')
  }
}

async function handleSubmit() {
  if (!form.value.name || !form.value.amount || form.value.category_ids.length === 0) {
    ElMessage.warning(t('common.required'))
    return
  }
  const numAmount = Number(form.value.amount)
  if (isNaN(numAmount) || numAmount <= 0) {
    ElMessage.warning(t('budget.amountInvalid') || 'Amount must be greater than 0')
    return
  }
  try {
    if (editingId.value) {
      await update(editingId.value, {
        name: form.value.name,
        amount: form.value.amount,
        period: form.value.period,
        category_ids: form.value.category_ids,
      })
    } else {
      await create(form.value)
    }
    ElMessage.success(t('common.success'))
    dialogVisible.value = false
    await fetchBudgets()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

/**
 * 根据预算状态返回 el-tag 的 type
 * 后端返回：normal / warning / overspent
 */
function getStatusType(status: string): '' | 'success' | 'warning' | 'danger' {
  if (status === 'overspent') return 'danger'
  if (status === 'warning') return 'warning'
  return 'success'
}

function getStatusLabel(status: string): string {
  if (status === 'overspent') return t('budget.overspent')
  if (status === 'warning') return t('budget.warning')
  return t('budget.normal')
}

onMounted(async () => {
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
      <el-table-column prop="name" :label="t('budget.name')" min-width="140" show-overflow-tooltip />
      <el-table-column :label="t('budget.amount')" width="130">
        <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
      </el-table-column>
      <el-table-column :label="t('budget.spent')" width="130">
        <template #default="{ row }">{{ formatAmount(row.spent) }}</template>
      </el-table-column>
      <el-table-column :label="t('budget.remaining')" width="130">
        <template #default="{ row }">{{ formatAmount(row.remaining) }}</template>
      </el-table-column>
      <el-table-column :label="t('budget.usageRate')" width="180">
        <template #default="{ row }">
          <el-progress
            :percentage="Math.min(Math.round(row.usage_rate * 100), 100)"
            :status="row.status === 'overspent' ? 'exception' : row.status === 'warning' ? 'warning' : undefined"
          />
          <span v-if="row.usage_rate > 1" class="overspent-label">{{ Math.round(row.usage_rate * 100) }}%</span>
        </template>
      </el-table-column>
      <el-table-column :label="t('budget.period')" width="80">
        <template #default="{ row }">
          {{ t(`budget.${row.period}`) }}
        </template>
      </el-table-column>
      <el-table-column :label="t('budget.status')" width="90">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.edit')" width="200" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleDetail(row.id)">{{ t('budget.detail') }}</el-button>
          <el-button link type="primary" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
          <el-button link type="danger" @click="handleDelete(row.id)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item :label="t('budget.name')">
          <el-input v-model="form.name" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="t('budget.amount')">
          <el-input v-model="form.amount" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="t('budget.period')">
          <el-select v-model="form.period" style="width: 100%">
            <el-option :label="t('budget.daily')" :value="BudgetPeriod.Daily" />
            <el-option :label="t('budget.weekly')" :value="BudgetPeriod.Weekly" />
            <el-option :label="t('budget.monthly')" :value="BudgetPeriod.Monthly" />
            <el-option :label="t('budget.quarterly')" :value="BudgetPeriod.Quarterly" />
            <el-option :label="t('budget.yearly')" :value="BudgetPeriod.Yearly" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('budget.categories')">
          <el-tree-select
            v-model="form.category_ids"
            :data="categoryTreeData"
            :placeholder="t('common.selectPlaceholder')"
            check-strictly
            multiple
            filterable
            clearable
            collapse-tags
            collapse-tags-tooltip
            style="width: 100%"
            :render-after-expand="false"
          />
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
.overspent-label {
  font-size: 12px;
  color: var(--app-amount-withdrawal);
  font-weight: 600;
}
</style>

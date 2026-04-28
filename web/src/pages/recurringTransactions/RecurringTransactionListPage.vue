<script setup lang="ts">
/**
 * 循环交易管理页面
 * 字段名与后端 JSON tag 完全对应：
 * - description（非 title）
 * - source_id（非 source_account_id）
 * - destination_id（非 destination_account_id）
 * - repeat_every（非 repeat_interval）
 * - is_active（非 active）
 * - next_occurrence（非 next_date）
 */
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, create, update, remove, processDue } from '@/api/recurringTransaction'
import { formatAmount, formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAccountStore } from '@/stores/account'
import { useCategoryStore } from '@/stores/category'
import type { RecurringTransaction, CreateRecurringTransactionReq, UpdateRecurringTransactionReq } from '@/types/recurringTransaction'
import { RecurrenceType } from '@/types/recurringTransaction'
import Pagination from '@/components/common/Pagination.vue'

const { t } = useI18n()
const accountStore = useAccountStore()
const categoryStore = useCategoryStore()

const items = ref<RecurringTransaction[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref<number | null>(null)
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

const recurrenceTypeOptions = [
  { value: RecurrenceType.Daily, label: t('recurringTransaction.daily') },
  { value: RecurrenceType.Weekly, label: t('recurringTransaction.weekly') },
  { value: RecurrenceType.Monthly, label: t('recurringTransaction.monthly') },
  { value: RecurrenceType.Yearly, label: t('recurringTransaction.yearly') },
]

const form = ref<CreateRecurringTransactionReq>({
  description: '',
  amount: '0',
  source_id: 0,
  destination_id: null,
  category_id: null,
  notes: '',
  recurrence_type: RecurrenceType.Monthly,
  repeat_every: 1,
  start_date: '',
})

async function fetchList() {
  loading.value = true
  try {
    const res = await list({ page: pagination.page, page_size: pagination.page_size }) as unknown as { items: RecurringTransaction[], total: number }
    items.value = res.items || []
    pagination.total = res.total || 0
  } catch {
    ElMessage.error(t('common.fetchError'))
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  pagination.page = page
  fetchList()
}

function handleSizeChange(size: number) {
  pagination.page_size = size
  pagination.page = 1
  fetchList()
}

function handleCreate() {
  dialogTitle.value = t('recurringTransaction.create')
  editingId.value = null
  form.value = {
    description: '',
    amount: '0',
    source_id: 0,
    destination_id: null,
    category_id: null,
    notes: '',
    recurrence_type: RecurrenceType.Monthly,
    repeat_every: 1,
    start_date: '',
    reminder_days: 0,
  }
  dialogVisible.value = true
}

function handleEdit(row: RecurringTransaction) {
  dialogTitle.value = t('recurringTransaction.edit')
  editingId.value = row.id
  form.value = {
    description: row.description,
    amount: row.amount,
    source_id: row.source_id,
    destination_id: row.destination_id,
    category_id: row.category_id,
    notes: row.notes,
    recurrence_type: row.recurrence_type,
    repeat_every: row.repeat_every,
    start_date: row.start_date ? row.start_date.substring(0, 10) : '',
    end_date: row.end_date ? row.end_date.substring(0, 10) : undefined,
    reminder_days: row.reminder_days ?? 0,
  }
  dialogVisible.value = true
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm(t('recurringTransaction.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchList()
  } catch {
    ElMessage.error(t('common.failed'))
  }
}

async function handleSubmit() {
  if (!form.value.description) {
    ElMessage.warning(t('common.required'))
    return
  }
  try {
    if (editingId.value) {
      await update(editingId.value, form.value as UpdateRecurringTransactionReq)
    } else {
      await create(form.value)
    }
    ElMessage.success(t('common.success'))
    dialogVisible.value = false
    await fetchList()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

async function handleProcessDue() {
  try {
    await processDue()
    ElMessage.success(t('common.success'))
    await fetchList()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

onMounted(async () => {
  await Promise.all([accountStore.fetchAccounts(), categoryStore.fetchCategories()])
  await fetchList()
})
</script>

<template>
  <div class="recurring-transaction-list-page">
    <div class="page-header">
      <h2>{{ t('recurringTransaction.title') }}</h2>
      <div>
        <el-button type="success" @click="handleProcessDue">{{ t('recurringTransaction.processDue') }}</el-button>
        <el-button type="primary" @click="handleCreate">{{ t('recurringTransaction.create') }}</el-button>
      </div>
    </div>

    <el-table :data="items" v-loading="loading" stripe>
      <!-- description 字段名与后端 RecurringTransactionResp JSON tag 对应（非 title） -->
      <el-table-column prop="description" :label="t('recurringTransaction.description')" />
      <el-table-column prop="amount" :label="t('recurringTransaction.amount')" width="130">
        <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
      </el-table-column>
      <el-table-column prop="recurrence_type" :label="t('recurringTransaction.recurrenceType')" width="100">
        <template #default="{ row }">
          {{ recurrenceTypeOptions.find(r => r.value === row.recurrence_type)?.label || row.recurrence_type }}
        </template>
      </el-table-column>
      <!-- next_occurrence 字段名与后端 JSON tag 对应（非 next_date） -->
      <el-table-column prop="next_occurrence" :label="t('recurringTransaction.nextDate')" width="120">
        <template #default="{ row }">{{ formatDate(row.next_occurrence) }}</template>
      </el-table-column>
      <!-- is_active 字段名与后端 JSON tag 对应（非 active） -->
      <el-table-column :label="t('recurringTransaction.active')" width="80">
        <template #default="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'info'" size="small">{{ row.is_active ? t('rule.enabled') : t('rule.disabled') }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.edit')" width="160" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
          <el-button link type="danger" @click="handleDelete(row.id)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <Pagination
      :total="pagination.total"
      :page="pagination.page"
      :page-size="pagination.page_size"
      @update:page="handlePageChange"
      @update:page-size="handleSizeChange"
    />

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="550px">
      <el-form :model="form" label-width="110px">
        <el-form-item :label="t('recurringTransaction.description')">
          <el-input v-model="form.description" />
        </el-form-item>
        <el-form-item :label="t('recurringTransaction.amount')">
          <el-input v-model="form.amount" />
        </el-form-item>
        <!-- source_id 字段名与后端 CreateRecurringTransactionReq JSON tag 对应（非 source_account_id） -->
        <el-form-item :label="t('transaction.sourceAccount')">
          <el-select v-model="form.source_id" :placeholder="t('common.selectPlaceholder')" filterable clearable style="width: 100%">
            <el-option v-for="acc in accountStore.accounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <!-- destination_id 字段名与后端 JSON tag 对应（非 destination_account_id） -->
        <el-form-item :label="t('transaction.destinationAccount')">
          <el-select v-model="form.destination_id" :placeholder="t('common.selectPlaceholder')" filterable clearable style="width: 100%">
            <el-option v-for="acc in accountStore.accounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('transaction.category')">
          <el-select v-model="form.category_id" :placeholder="t('common.selectPlaceholder')" filterable clearable style="width: 100%">
            <el-option v-for="cat in categoryStore.categories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('recurringTransaction.recurrenceType')">
          <el-select v-model="form.recurrence_type" style="width: 100%">
            <el-option v-for="opt in recurrenceTypeOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <!-- repeat_every 字段名与后端 JSON tag 对应（非 repeat_interval） -->
        <el-form-item :label="t('recurringTransaction.repeatEvery')">
          <el-input-number v-model="form.repeat_every" :min="1" />
        </el-form-item>
        <el-form-item :label="t('recurringTransaction.startDate')">
          <el-date-picker v-model="form.start_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="t('recurringTransaction.endDate')">
          <el-date-picker v-model="form.end_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="t('recurringTransaction.reminderDays')">
          <el-input-number v-model="form.reminder_days" :min="0" :max="365" />
        </el-form-item>
        <el-form-item :label="t('recurringTransaction.notes')">
          <el-input v-model="form.notes" type="textarea" />
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
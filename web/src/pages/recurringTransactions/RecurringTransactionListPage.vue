<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, create, update, remove, processDue } from '@/api/recurringTransaction'
import { formatAmount, formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { RecurringTransaction, CreateRecurringTransactionReq, UpdateRecurringTransactionReq } from '@/types/recurringTransaction'
import { RecurrenceType } from '@/types/recurringTransaction'

const { t } = useI18n()

const items = ref<RecurringTransaction[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref<string | null>(null)

const recurrenceTypeOptions = [
  { value: RecurrenceType.Daily, label: t('recurringTransaction.daily') },
  { value: RecurrenceType.Weekly, label: t('recurringTransaction.weekly') },
  { value: RecurrenceType.Monthly, label: t('recurringTransaction.monthly') },
  { value: RecurrenceType.Yearly, label: t('recurringTransaction.yearly') },
]

const form = ref<CreateRecurringTransactionReq>({
  title: '',
  type: 'withdrawal',
  amount: '0',
  source_account_id: '',
  recurrence_type: RecurrenceType.Monthly,
  repeat_interval: 1,
  start_date: '',
  description: '',
})

async function fetchList() {
  loading.value = true
  try {
    const res = await list({}) as unknown as { items: RecurringTransaction[] }
    items.value = res.items || []
  } catch {
    // handle error
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  dialogTitle.value = t('recurringTransaction.create')
  editingId.value = null
  form.value = {
    title: '',
    type: 'withdrawal',
    amount: '0',
    source_account_id: '',
    recurrence_type: RecurrenceType.Monthly,
    repeat_interval: 1,
    start_date: '',
    description: '',
  }
  dialogVisible.value = true
}

function handleEdit(row: RecurringTransaction) {
  dialogTitle.value = t('recurringTransaction.edit')
  editingId.value = row.id
  form.value = {
    title: row.title,
    type: row.type,
    amount: row.amount,
    source_account_id: row.source_account_id,
    recurrence_type: row.recurrence_type,
    repeat_interval: row.repeat_interval,
    start_date: row.start_date,
    end_date: row.end_date,
    description: row.description,
  }
  dialogVisible.value = true
}

async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm(t('recurringTransaction.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchList()
  } catch {
    // cancelled or error
  }
}

async function handleSubmit() {
  if (!form.value.title) {
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

onMounted(fetchList)
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
      <el-table-column prop="title" :label="t('recurringTransaction.title')" />
      <el-table-column prop="type" :label="t('recurringTransaction.type')" width="100" />
      <el-table-column prop="amount" :label="t('recurringTransaction.amount')" width="130">
        <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
      </el-table-column>
      <el-table-column prop="recurrence_type" :label="t('recurringTransaction.recurrenceType')" width="100">
        <template #default="{ row }">
          {{ recurrenceTypeOptions.find(r => r.value === row.recurrence_type)?.label || row.recurrence_type }}
        </template>
      </el-table-column>
      <el-table-column prop="next_date" :label="t('recurringTransaction.nextDate')" width="120">
        <template #default="{ row }">{{ formatDate(row.next_date) }}</template>
      </el-table-column>
      <el-table-column :label="t('recurringTransaction.active')" width="80">
        <template #default="{ row }">
          <el-tag :type="row.active ? 'success' : 'info'" size="small">{{ row.active ? t('rule.enabled') : t('rule.disabled') }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.edit')" width="160" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
          <el-button link type="danger" @click="handleDelete(row.id)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="550px">
      <el-form :model="form" label-width="110px">
        <el-form-item :label="t('recurringTransaction.title')">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item :label="t('recurringTransaction.type')">
          <el-select v-model="form.type">
            <el-option label="Withdrawal" value="withdrawal" />
            <el-option label="Deposit" value="deposit" />
            <el-option label="Transfer" value="transfer" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('recurringTransaction.amount')">
          <el-input v-model="form.amount" />
        </el-form-item>
        <el-form-item :label="t('recurringTransaction.recurrenceType')">
          <el-select v-model="form.recurrence_type">
            <el-option v-for="opt in recurrenceTypeOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('recurringTransaction.repeatInterval')">
          <el-input-number v-model="form.repeat_interval" :min="1" />
        </el-form-item>
        <el-form-item :label="t('recurringTransaction.startDate')">
          <el-date-picker v-model="form.start_date" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item :label="t('recurringTransaction.endDate')">
          <el-date-picker v-model="form.end_date" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item :label="t('transaction.description')">
          <el-input v-model="form.description" type="textarea" />
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

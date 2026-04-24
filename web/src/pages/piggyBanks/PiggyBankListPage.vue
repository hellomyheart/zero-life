<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, create, update, remove, addAmount, removeAmount } from '@/api/piggyBank'
import { formatAmount, formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { PiggyBank, CreatePiggyBankReq, UpdatePiggyBankReq, AddAmountReq } from '@/types/piggyBank'

const { t } = useI18n()

const piggyBanks = ref<PiggyBank[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref<string | null>(null)
const amountDialogVisible = ref(false)
const amountDialogTitle = ref('')
const amountDialogType = ref<'add' | 'remove'>('add')
const amountDialogId = ref<string>('')
const amountForm = ref<AddAmountReq>({ amount: '0', note: '' })

const form = ref<CreatePiggyBankReq>({
  name: '',
  account_id: '',
  target_amount: '0',
  target_date: '',
  notes: '',
})

async function fetchPiggyBanks() {
  loading.value = true
  try {
    const res = await list({}) as unknown as { items: PiggyBank[] }
    piggyBanks.value = res.items || []
  } catch {
    // handle error
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  dialogTitle.value = t('piggyBank.create')
  editingId.value = null
  form.value = { name: '', account_id: '', target_amount: '0', target_date: '', notes: '' }
  dialogVisible.value = true
}

function handleEdit(row: PiggyBank) {
  dialogTitle.value = t('piggyBank.edit')
  editingId.value = row.id
  form.value = {
    name: row.name,
    account_id: row.account_id,
    target_amount: row.target_amount,
    target_date: row.target_date,
    notes: row.notes,
  }
  dialogVisible.value = true
}

async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm(t('piggyBank.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchPiggyBanks()
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
    if (editingId.value) {
      await update(editingId.value, form.value as UpdatePiggyBankReq)
    } else {
      await create(form.value)
    }
    ElMessage.success(t('common.success'))
    dialogVisible.value = false
    await fetchPiggyBanks()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

function handleAddMoney(row: PiggyBank) {
  amountDialogTitle.value = t('piggyBank.addMoney')
  amountDialogType.value = 'add'
  amountDialogId.value = row.id
  amountForm.value = { amount: '0', note: '' }
  amountDialogVisible.value = true
}

function handleRemoveMoney(row: PiggyBank) {
  amountDialogTitle.value = t('piggyBank.removeMoney')
  amountDialogType.value = 'remove'
  amountDialogId.value = row.id
  amountForm.value = { amount: '0', note: '' }
  amountDialogVisible.value = true
}

async function handleAmountSubmit() {
  try {
    if (amountDialogType.value === 'add') {
      await addAmount(amountDialogId.value, amountForm.value)
    } else {
      await removeAmount(amountDialogId.value, amountForm.value)
    }
    ElMessage.success(t('common.success'))
    amountDialogVisible.value = false
    await fetchPiggyBanks()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

onMounted(fetchPiggyBanks)
</script>

<template>
  <div class="piggy-bank-list-page">
    <div class="page-header">
      <h2>{{ t('piggyBank.title') }}</h2>
      <el-button type="primary" @click="handleCreate">{{ t('piggyBank.create') }}</el-button>
    </div>

    <el-table :data="piggyBanks" v-loading="loading" stripe>
      <el-table-column prop="name" :label="t('piggyBank.name')" />
      <el-table-column prop="target_amount" :label="t('piggyBank.targetAmount')" width="140">
        <template #default="{ row }">{{ formatAmount(row.target_amount) }}</template>
      </el-table-column>
      <el-table-column prop="current_amount" :label="t('piggyBank.currentAmount')" width="140">
        <template #default="{ row }">{{ formatAmount(row.current_amount) }}</template>
      </el-table-column>
      <el-table-column :label="t('piggyBank.percentage')" width="180">
        <template #default="{ row }">
          <el-progress :percentage="row.percentage || 0" :stroke-width="18" />
        </template>
      </el-table-column>
      <el-table-column prop="target_date" :label="t('piggyBank.targetDate')" width="120">
        <template #default="{ row }">{{ row.target_date ? formatDate(row.target_date) : '-' }}</template>
      </el-table-column>
      <el-table-column :label="t('common.edit')" width="220" fixed="right">
        <template #default="{ row }">
          <el-button link type="success" @click="handleAddMoney(row)">{{ t('piggyBank.addMoney') }}</el-button>
          <el-button link type="warning" @click="handleRemoveMoney(row)">{{ t('piggyBank.removeMoney') }}</el-button>
          <el-button link type="primary" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
          <el-button link type="danger" @click="handleDelete(row.id)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item :label="t('piggyBank.name')">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="t('piggyBank.targetAmount')">
          <el-input v-model="form.target_amount" />
        </el-form-item>
        <el-form-item :label="t('piggyBank.targetDate')">
          <el-date-picker v-model="form.target_date" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item :label="t('piggyBank.notes')">
          <el-input v-model="form.notes" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="amountDialogVisible" :title="amountDialogTitle" width="400px">
      <el-form :model="amountForm" label-width="80px">
        <el-form-item :label="t('piggyBank.amount')">
          <el-input v-model="amountForm.amount" />
        </el-form-item>
        <el-form-item :label="t('piggyBank.notes')">
          <el-input v-model="amountForm.note" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="amountDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleAmountSubmit">{{ t('common.save') }}</el-button>
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

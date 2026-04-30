<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { list, create, update, remove, addAmount, removeAmount } from '@/api/piggyBank'
import { formatAmount, formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAccountStore } from '@/stores/account'
import type { PiggyBank, CreatePiggyBankReq, UpdatePiggyBankReq, AddAmountReq, RemoveAmountReq } from '@/types/piggyBank'
import { AccountType } from '@/types/account'
import Decimal from 'decimal.js'

const { t } = useI18n()
const router = useRouter()
const accountStore = useAccountStore()

const piggyBanks = ref<PiggyBank[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref<number | null>(null)
const amountDialogVisible = ref(false)
const amountDialogTitle = ref('')
const amountDialogType = ref<'add' | 'remove'>('add')
const amountDialogId = ref<number>(0)
const amountDialogMaxAmount = ref('0')
const amountForm = ref<AddAmountReq>({ amount: '0', note: '' })

const assetAccounts = computed(() =>
  accountStore.accounts.filter(a => a.type === AccountType.Asset)
)

const form = ref<CreatePiggyBankReq>({
  name: '',
  account_id: 0,
  target_amount: '0',
  target_date: '',
  notes: '',
})

async function fetchPiggyBanks() {
  loading.value = true
  try {
    const res = await list()
    piggyBanks.value = res as unknown as PiggyBank[]
  } catch {
    ElMessage.error(t('common.fetchError'))
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  dialogTitle.value = t('piggyBank.create')
  editingId.value = null
  form.value = { name: '', account_id: 0, target_amount: '0', target_date: '', notes: '' }
  dialogVisible.value = true
}

function handleEdit(row: PiggyBank) {
  router.push(`/piggy-banks/${row.id}`)
}

function handleRowClick(row: PiggyBank) {
  router.push(`/piggy-banks/${row.id}`)
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm(t('piggyBank.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchPiggyBanks()
  } catch (err) {
    if ((err as string) !== 'cancel') ElMessage.error(t('common.failed'))
  }
}

async function handleSubmit() {
  if (!form.value.name) {
    ElMessage.warning(t('common.required'))
    return
  }
  if (!form.value.account_id) {
    ElMessage.warning(t('piggyBank.accountRequired'))
    return
  }
  const amount = new Decimal(form.value.target_amount)
  if (!amount.isPositive()) {
    ElMessage.warning(t('piggyBank.amountPositive'))
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
  amountDialogMaxAmount.value = row.available_deposit || '0'
  amountForm.value = { amount: '0', note: '' }
  amountDialogVisible.value = true
}

function handleRemoveMoney(row: PiggyBank) {
  amountDialogTitle.value = t('piggyBank.removeMoney')
  amountDialogType.value = 'remove'
  amountDialogId.value = row.id
  amountDialogMaxAmount.value = row.current_amount
  amountForm.value = { amount: '0', note: '' }
  amountDialogVisible.value = true
}

async function handleAmountSubmit() {
  const amount = new Decimal(amountForm.value.amount)
  if (!amount.isPositive()) {
    ElMessage.warning(t('piggyBank.amountPositive'))
    return
  }
  const max = new Decimal(amountDialogMaxAmount.value)
  if (amount.greaterThan(max)) {
    ElMessage.warning(amountDialogType.value === 'add' ? t('piggyBank.overDeposit') : t('piggyBank.overWithdraw'))
    return
  }

  try {
    if (amountDialogType.value === 'add') {
      await addAmount(amountDialogId.value, amountForm.value)
    } else {
      await removeAmount(amountDialogId.value, amountForm.value as RemoveAmountReq)
    }
    ElMessage.success(t('common.success'))
    amountDialogVisible.value = false
    await fetchPiggyBanks()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

function progressStatus(percentage: number): '' | 'success' | 'warning' | 'exception' {
  if (percentage >= 100) return 'success'
  if (percentage >= 80) return 'warning'
  return ''
}

onMounted(async () => {
  await accountStore.fetchAccounts()
  await fetchPiggyBanks()
})
</script>

<template>
  <div class="piggy-bank-list-page">
    <div class="page-header">
      <h2>{{ t('piggyBank.title') }}</h2>
      <el-button type="primary" @click="handleCreate">{{ t('piggyBank.create') }}</el-button>
    </div>

    <el-table :data="piggyBanks" v-loading="loading" stripe @row-click="handleRowClick" class="clickable-table">
      <el-table-column prop="name" :label="t('piggyBank.name')" min-width="120" />
      <el-table-column :label="t('piggyBank.account')" width="130" show-overflow-tooltip>
        <template #default="{ row }">{{ row.account?.name || '-' }}</template>
      </el-table-column>
      <el-table-column prop="target_amount" :label="t('piggyBank.targetAmount')" width="140" align="right">
        <template #default="{ row }">{{ formatAmount(row.target_amount) }}</template>
      </el-table-column>
      <el-table-column prop="current_amount" :label="t('piggyBank.currentAmount')" width="140" align="right">
        <template #default="{ row }">{{ formatAmount(row.current_amount) }}</template>
      </el-table-column>
      <el-table-column :label="t('piggyBank.percentage')" width="180">
        <template #default="{ row }">
          <el-progress :percentage="Math.round(row.percentage || 0)" :stroke-width="18" :status="progressStatus(row.percentage)" />
        </template>
      </el-table-column>
      <el-table-column prop="target_date" :label="t('piggyBank.targetDate')" width="120">
        <template #default="{ row }">{{ row.target_date ? formatDate(row.target_date, 'YYYY-MM-DD') : '-' }}</template>
      </el-table-column>
      <el-table-column :label="t('common.edit')" width="220" fixed="right">
        <template #default="{ row }">
          <el-button link type="success" @click.stop="handleAddMoney(row)">{{ t('piggyBank.addMoney') }}</el-button>
          <el-button link type="warning" @click.stop="handleRemoveMoney(row)">{{ t('piggyBank.removeMoney') }}</el-button>
          <el-button link type="primary" @click.stop="handleEdit(row)">{{ t('common.edit') }}</el-button>
          <el-button link type="danger" @click.stop="handleDelete(row.id)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item :label="t('piggyBank.name')">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="t('piggyBank.account')">
          <el-select v-model="form.account_id" :placeholder="t('common.selectPlaceholder')" filterable>
            <el-option v-for="acc in assetAccounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('piggyBank.targetAmount')">
          <el-input v-model="form.target_amount" type="number" step="0.01" min="0" />
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
          <el-input v-model="amountForm.amount" type="number" step="0.01" min="0" />
        </el-form-item>
        <el-form-item :label="t('piggyBank.maxAllowed')">
          <span class="hint-text">{{ formatAmount(amountDialogMaxAmount) }}</span>
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

.clickable-table :deep(.el-table__row) {
  cursor: pointer;
}

.hint-text {
  color: var(--app-text-secondary);
  font-size: 13px;
}
</style>

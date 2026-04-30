<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { getPiggyBank, getEvents, addAmount, removeAmount, update, remove as removePiggyBank, resetHistory } from '@/api/piggyBank'
import { formatAmount, formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAccountStore } from '@/stores/account'
import type { PiggyBank, PiggyEvent, UpdatePiggyBankReq, AddAmountReq, RemoveAmountReq } from '@/types/piggyBank'
import Decimal from 'decimal.js'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const accountStore = useAccountStore()

const piggyBank = ref<PiggyBank | null>(null)
const events = ref<PiggyEvent[]>([])
const loading = ref(true)
const eventsLoading = ref(false)

const amountDialogVisible = ref(false)
const amountDialogTitle = ref('')
const amountDialogType = ref<'add' | 'remove'>('add')
const amountForm = ref<AddAmountReq>({ amount: '0', note: '' })

const editDialogVisible = ref(false)
const editForm = ref<UpdatePiggyBankReq>({})

const maxAddAmount = computed(() => {
  if (!piggyBank.value) return '0'
  return piggyBank.value.available_deposit || '0'
})

function formatEventAmount(event: PiggyEvent): string {
  const num = new Decimal(event.amount)
  if (num.isPositive()) {
    return `+${num.abs().toFixed(2)}`
  }
  return `-${num.abs().toFixed(2)}`
}

function eventAmountClass(event: PiggyEvent): string {
  return new Decimal(event.amount).isPositive() ? 'amount-deposit' : 'amount-withdrawal'
}

async function fetchPiggyBank() {
  const id = Number(route.params.id)
  try {
    piggyBank.value = await getPiggyBank(id) as unknown as PiggyBank
  } catch {
    ElMessage.error(t('common.fetchError'))
  } finally {
    loading.value = false
  }
}

async function fetchEvents() {
  const id = Number(route.params.id)
  eventsLoading.value = true
  try {
    events.value = await getEvents(id) as unknown as PiggyEvent[]
  } catch {
    ElMessage.error(t('common.fetchError'))
  } finally {
    eventsLoading.value = false
  }
}

function handleAddMoney() {
  amountDialogTitle.value = t('piggyBank.addMoney')
  amountDialogType.value = 'add'
  amountForm.value = { amount: '0', note: '' }
  amountDialogVisible.value = true
}

function handleRemoveMoney() {
  amountDialogTitle.value = t('piggyBank.removeMoney')
  amountDialogType.value = 'remove'
  amountForm.value = { amount: '0', note: '' }
  amountDialogVisible.value = true
}

async function handleAmountSubmit() {
  const id = Number(route.params.id)
  const amount = new Decimal(amountForm.value.amount)
  if (!amount.isPositive()) {
    ElMessage.warning(t('piggyBank.amountPositive'))
    return
  }
  if (amountDialogType.value === 'add') {
    const max = new Decimal(maxAddAmount.value)
    if (amount.greaterThan(max)) {
      ElMessage.warning(t('piggyBank.overDeposit'))
      return
    }
  } else {
    const current = new Decimal(piggyBank.value?.current_amount || '0')
    if (amount.greaterThan(current)) {
      ElMessage.warning(t('piggyBank.overWithdraw'))
      return
    }
  }

  try {
    if (amountDialogType.value === 'add') {
      piggyBank.value = await addAmount(id, amountForm.value) as unknown as PiggyBank
    } else {
      piggyBank.value = await removeAmount(id, amountForm.value as RemoveAmountReq) as unknown as PiggyBank
    }
    ElMessage.success(t('common.success'))
    amountDialogVisible.value = false
    await fetchEvents()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

function handleEdit() {
  if (!piggyBank.value) return
  editForm.value = {
    name: piggyBank.value.name,
    target_amount: piggyBank.value.target_amount,
    target_date: piggyBank.value.target_date,
    notes: piggyBank.value.notes,
  }
  editDialogVisible.value = true
}

async function handleEditSubmit() {
  const id = Number(route.params.id)
  try {
    piggyBank.value = await update(id, editForm.value) as unknown as PiggyBank
    ElMessage.success(t('common.success'))
    editDialogVisible.value = false
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

async function handleDelete() {
  const id = Number(route.params.id)
  try {
    await ElMessageBox.confirm(t('piggyBank.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await removePiggyBank(id)
    ElMessage.success(t('common.success'))
    router.push('/piggy-banks')
  } catch (err) {
    if ((err as string) !== 'cancel') ElMessage.error(t('common.failed'))
  }
}

async function handleResetHistory() {
  const id = Number(route.params.id)
  try {
    await ElMessageBox.confirm(t('piggyBank.resetConfirm'), t('common.confirm'), { type: 'warning' })
    piggyBank.value = await resetHistory(id) as unknown as PiggyBank
    ElMessage.success(t('common.success'))
    await fetchEvents()
  } catch (err) {
    if ((err as string) !== 'cancel') ElMessage.error(t('common.failed'))
  }
}

function goBack() {
  router.push('/piggy-banks')
}

onMounted(async () => {
  await accountStore.fetchAccounts()
  await fetchPiggyBank()
  await fetchEvents()
})
</script>

<template>
  <div class="piggy-bank-detail-page" v-loading="loading">
    <template v-if="piggyBank">
      <div class="page-header">
        <h2>{{ piggyBank.name }}</h2>
        <div class="header-actions">
          <el-button type="success" @click="handleAddMoney">{{ t('piggyBank.addMoney') }}</el-button>
          <el-button type="warning" @click="handleRemoveMoney">{{ t('piggyBank.removeMoney') }}</el-button>
          <el-button @click="handleEdit">{{ t('common.edit') }}</el-button>
          <el-button type="danger" @click="handleDelete">{{ t('common.delete') }}</el-button>
          <el-button @click="goBack">{{ t('common.back') }}</el-button>
        </div>
      </div>

      <el-row :gutter="16">
        <el-col :xs="24" :sm="8">
          <el-card shadow="hover">
            <template #header>{{ t('piggyBank.targetAmount') }}</template>
            <div class="amount">{{ formatAmount(piggyBank.target_amount) }}</div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="8">
          <el-card shadow="hover">
            <template #header>{{ t('piggyBank.currentAmount') }}</template>
            <div class="amount deposit">{{ formatAmount(piggyBank.current_amount) }}</div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="8">
          <el-card shadow="hover">
            <template #header>{{ t('piggyBank.remaining') }}</template>
            <div class="amount remaining">{{ formatAmount(new Decimal(piggyBank.target_amount).minus(new Decimal(piggyBank.current_amount)).toString()) }}</div>
          </el-card>
        </el-col>
      </el-row>

      <el-card style="margin-top: 16px">
        <template #header>{{ t('piggyBank.percentage') }}</template>
        <el-progress
          :percentage="Math.round(piggyBank.percentage)"
          :stroke-width="20"
          :status="piggyBank.percentage >= 100 ? 'success' : undefined"
        />
        <div class="usage-detail">{{ piggyBank.percentage.toFixed(1) }}%</div>
      </el-card>

      <el-card style="margin-top: 16px">
        <template #header>{{ t('piggyBank.info') }}</template>
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="t('piggyBank.account')">
            {{ piggyBank.account?.name || piggyBank.account_id }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('piggyBank.targetDate')">
            {{ piggyBank.target_date ? formatDate(piggyBank.target_date, 'YYYY-MM-DD') : '-' }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('piggyBank.notes')" :span="2">
            {{ piggyBank.notes || '-' }}
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <el-card style="margin-top: 16px">
        <template #header>
          <div class="card-header-with-action">
            <span>{{ t('piggyBank.events') }}</span>
            <el-button type="danger" size="small" @click="handleResetHistory">{{ t('piggyBank.resetHistory') }}</el-button>
          </div>
        </template>
        <el-table :data="events" v-loading="eventsLoading" stripe>
          <el-table-column prop="created_at" :label="t('common.createdAt')" width="180">
            <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
          </el-table-column>
          <el-table-column :label="t('piggyBank.amount')" width="160" align="right">
            <template #default="{ row }">
              <span :class="eventAmountClass(row)" class="amount-text">{{ formatEventAmount(row) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="note" :label="t('piggyBank.notes')" min-width="200" show-overflow-tooltip>
            <template #default="{ row }">{{ row.note || '-' }}</template>
          </el-table-column>
        </el-table>
        <div v-if="!events.length && !eventsLoading" class="no-data">{{ t('common.noData') }}</div>
      </el-card>
    </template>
  </div>

  <el-dialog v-model="amountDialogVisible" :title="amountDialogTitle" width="400px">
    <el-form :model="amountForm" label-width="80px">
      <el-form-item :label="t('piggyBank.amount')">
        <el-input v-model="amountForm.amount" type="number" step="0.01" min="0" />
      </el-form-item>
      <el-form-item v-if="amountDialogType === 'add'" :label="t('piggyBank.maxAdd')">
        <span class="hint-text">{{ maxAddAmount }}</span>
      </el-form-item>
      <el-form-item v-else :label="t('piggyBank.maxRemove')">
        <span class="hint-text">{{ piggyBank?.current_amount || '0' }}</span>
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

  <el-dialog v-model="editDialogVisible" :title="t('piggyBank.edit')" width="500px">
    <el-form :model="editForm" label-width="100px">
      <el-form-item :label="t('piggyBank.name')">
        <el-input v-model="editForm.name" />
      </el-form-item>
      <el-form-item :label="t('piggyBank.targetAmount')">
        <el-input v-model="editForm.target_amount" type="number" step="0.01" min="0" />
      </el-form-item>
      <el-form-item :label="t('piggyBank.targetDate')">
        <el-date-picker v-model="editForm.target_date" type="date" value-format="YYYY-MM-DD" />
      </el-form-item>
      <el-form-item :label="t('piggyBank.notes')">
        <el-input v-model="editForm.notes" type="textarea" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="editDialogVisible = false">{{ t('common.cancel') }}</el-button>
      <el-button type="primary" @click="handleEditSubmit">{{ t('common.save') }}</el-button>
    </template>
  </el-dialog>
</template>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 8px;
}

.page-header h2 {
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.amount {
  font-size: 24px;
  font-weight: 600;
  text-align: center;
  padding: 8px 0;
}

.amount.deposit {
  color: var(--app-amount-deposit);
}

.amount.remaining {
  color: var(--app-amount-transfer);
}

.amount-text {
  font-weight: 600;
  font-variant-numeric: tabular-nums;
}

.amount-deposit {
  color: var(--app-amount-deposit);
}

.amount-withdrawal {
  color: var(--app-amount-withdrawal);
}

.usage-detail {
  margin-top: 8px;
  text-align: center;
  color: var(--app-text-secondary);
  font-size: 14px;
}

.card-header-with-action {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.no-data {
  color: var(--app-text-secondary);
  text-align: center;
  padding: 20px;
  font-size: 14px;
}

.hint-text {
  color: var(--app-text-secondary);
  font-size: 13px;
}

@media (max-width: 767px) {
  .amount {
    font-size: 20px;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>

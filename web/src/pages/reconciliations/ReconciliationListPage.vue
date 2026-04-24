<script setup lang="ts">
// 对账管理页面 - 核对账户余额与实际余额
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, create, update, remove } from '@/api/reconciliation'
import { formatAmount, formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Reconciliation, CreateReconciliationReq, UpdateReconciliationReq } from '@/types/reconciliation'

const { t } = useI18n()

const items = ref<Reconciliation[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref<string | null>(null)

const form = ref<CreateReconciliationReq>({
  account_id: '',
  start_date: '',
  end_date: '',
  start_balance: '0',
  end_balance: '0',
  notes: '',
})

async function fetchList() {
  loading.value = true
  try {
    const res = await list({}) as unknown as { items: Reconciliation[] }
    items.value = res.items || []
  } catch {
    // handle error
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  dialogTitle.value = t('reconciliation.create')
  editingId.value = null
  form.value = { account_id: '', start_date: '', end_date: '', start_balance: '0', end_balance: '0', notes: '' }
  dialogVisible.value = true
}

function handleEdit(row: Reconciliation) {
  dialogTitle.value = t('reconciliation.edit')
  editingId.value = row.id
  form.value = {
    account_id: row.account_id,
    start_date: row.start_date,
    end_date: row.end_date,
    start_balance: row.start_balance,
    end_balance: row.end_balance,
    notes: row.notes,
  }
  dialogVisible.value = true
}

async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm(t('reconciliation.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchList()
  } catch {
    // cancelled or error
  }
}

async function handleSubmit() {
  try {
    if (editingId.value) {
      await update(editingId.value, { end_balance: form.value.end_balance, notes: form.value.notes } as UpdateReconciliationReq)
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

onMounted(fetchList)
</script>

<template>
  <div class="reconciliation-list-page">
    <div class="page-header">
      <h2>{{ t('reconciliation.title') }}</h2>
      <el-button type="primary" @click="handleCreate">{{ t('reconciliation.create') }}</el-button>
    </div>

    <el-table :data="items" v-loading="loading" stripe>
      <el-table-column prop="start_date" :label="t('common.startDate')" width="120">
        <template #default="{ row }">{{ formatDate(row.start_date) }}</template>
      </el-table-column>
      <el-table-column prop="end_date" :label="t('common.endDate')" width="120">
        <template #default="{ row }">{{ formatDate(row.end_date) }}</template>
      </el-table-column>
      <el-table-column prop="start_balance" :label="t('reconciliation.startBalance')" width="130">
        <template #default="{ row }">{{ formatAmount(row.start_balance) }}</template>
      </el-table-column>
      <el-table-column prop="end_balance" :label="t('reconciliation.endBalance')" width="130">
        <template #default="{ row }">{{ formatAmount(row.end_balance) }}</template>
      </el-table-column>
      <el-table-column prop="difference" :label="t('reconciliation.difference')" width="130">
        <template #default="{ row }">{{ formatAmount(row.difference) }}</template>
      </el-table-column>
      <el-table-column prop="status" :label="t('reconciliation.status')" width="100" />
      <el-table-column :label="t('common.edit')" width="160" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
          <el-button link type="danger" @click="handleDelete(row.id)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item :label="t('common.startDate')">
          <el-date-picker v-model="form.start_date" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item :label="t('common.endDate')">
          <el-date-picker v-model="form.end_date" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item :label="t('reconciliation.startBalance')">
          <el-input v-model="form.start_balance" />
        </el-form-item>
        <el-form-item :label="t('reconciliation.endBalance')">
          <el-input v-model="form.end_balance" />
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

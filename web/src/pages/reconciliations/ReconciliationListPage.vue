<script setup lang="ts">
/**
 * 对账管理页面
 * 字段名与后端 JSON tag 完全对应：
 * - starting_balance（非 start_balance）
 * - ending_balance（非 end_balance）
 * - 后端 CreateReconciliationReq 不接受 notes 字段
 */
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, create, update, remove } from '@/api/reconciliation'
import { formatAmount, formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAccountStore } from '@/stores/account'
import type { Reconciliation, CreateReconciliationReq, UpdateReconciliationReq } from '@/types/reconciliation'
import Pagination from '@/components/common/Pagination.vue'

const { t } = useI18n()
const accountStore = useAccountStore()

const items = ref<Reconciliation[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref<number | null>(null)
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

const form = ref<CreateReconciliationReq>({
  account_id: 0,
  start_date: '',
  end_date: '',
  starting_balance: '0',
  ending_balance: '0',
})

async function fetchList() {
  loading.value = true
  try {
    const res = await list({ page: pagination.page, page_size: pagination.page_size }) as unknown as { items: Reconciliation[], total: number }
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
  dialogTitle.value = t('reconciliation.create')
  editingId.value = null
  form.value = { account_id: 0, start_date: '', end_date: '', starting_balance: '0', ending_balance: '0' }
  dialogVisible.value = true
}

function handleEdit(row: Reconciliation) {
  dialogTitle.value = t('reconciliation.edit')
  editingId.value = row.id
  form.value = {
    account_id: row.account_id,
    start_date: row.start_date ? row.start_date.substring(0, 10) : '',
    end_date: row.end_date ? row.end_date.substring(0, 10) : '',
    starting_balance: row.starting_balance,
    ending_balance: row.ending_balance,
  }
  dialogVisible.value = true
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm(t('reconciliation.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchList()
  } catch {
    ElMessage.error(t('common.failed'))
  }
}

async function handleSubmit() {
  try {
    if (editingId.value) {
      await update(editingId.value, { ending_balance: form.value.ending_balance } as UpdateReconciliationReq)
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

onMounted(async () => {
  await accountStore.fetchAccounts()
  await fetchList()
})
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
      <!-- starting_balance 字段名与后端 ReconciliationResp JSON tag 对应（非 start_balance） -->
      <el-table-column prop="starting_balance" :label="t('reconciliation.startBalance')" width="130">
        <template #default="{ row }">{{ formatAmount(row.starting_balance) }}</template>
      </el-table-column>
      <!-- ending_balance 字段名与后端 JSON tag 对应（非 end_balance） -->
      <el-table-column prop="ending_balance" :label="t('reconciliation.endBalance')" width="130">
        <template #default="{ row }">{{ formatAmount(row.ending_balance) }}</template>
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

    <Pagination
      :total="pagination.total"
      :page="pagination.page"
      :page-size="pagination.page_size"
      @update:page="handlePageChange"
      @update:page-size="handleSizeChange"
    />

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item :label="t('transaction.sourceAccount')">
          <el-select v-model="form.account_id" :placeholder="t('common.selectPlaceholder')" filterable clearable style="width: 100%">
            <el-option v-for="acc in accountStore.accounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('common.startDate')">
          <el-date-picker v-model="form.start_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="t('common.endDate')">
          <el-date-picker v-model="form.end_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="t('reconciliation.startBalance')">
          <el-input v-model="form.starting_balance" />
        </el-form-item>
        <el-form-item :label="t('reconciliation.endBalance')">
          <el-input v-model="form.ending_balance" />
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
<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { list, remove } from '@/api/transaction'
import { formatAmount, formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Transaction } from '@/types/transaction'
import { TransactionType } from '@/types/transaction'
import Pagination from '@/components/common/Pagination.vue'
import DateRangePicker from '@/components/common/DateRangePicker.vue'

const { t } = useI18n()
const router = useRouter()

const transactions = ref<Transaction[]>([])
const total = ref(0)
const loading = ref(false)

const filter = reactive({
  type: '' as string,
  start_date: '',
  end_date: '',
  source_account_id: '',
  category_id: '',
  keyword: '',
  page: 1,
  page_size: 20,
})

async function fetchTransactions() {
  loading.value = true
  try {
    const params: Record<string, unknown> = { page: filter.page, page_size: filter.page_size }
    if (filter.type) params.type = filter.type
    if (filter.start_date) params.start_date = filter.start_date
    if (filter.end_date) params.end_date = filter.end_date
    if (filter.source_account_id) params.source_account_id = filter.source_account_id
    if (filter.category_id) params.category_id = filter.category_id
    if (filter.keyword) params.keyword = filter.keyword

    const res = await list(params as typeof filter)
    const data = res as unknown as { items: Transaction[]; total: number }
    transactions.value = data.items || []
    total.value = data.total || 0
  } catch {
    // handle error
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  router.push('/transactions/create')
}

function handleEdit(id: string) {
  router.push(`/transactions/${id}/edit`)
}

async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm(t('transaction.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchTransactions()
  } catch {
    // cancelled or error
  }
}

function handlePageChange(page: number) {
  filter.page = page
  fetchTransactions()
}

function handleSizeChange(size: number) {
  filter.page_size = size
  filter.page = 1
  fetchTransactions()
}

function handleFilter() {
  filter.page = 1
  fetchTransactions()
}

function resetFilter() {
  filter.type = ''
  filter.start_date = ''
  filter.end_date = ''
  filter.source_account_id = ''
  filter.category_id = ''
  filter.keyword = ''
  filter.page = 1
  fetchTransactions()
}

onMounted(fetchTransactions)
</script>

<template>
  <div class="transaction-list-page">
    <div class="page-header">
      <h2>{{ t('transaction.title') }}</h2>
      <el-button type="primary" @click="handleCreate">{{ t('transaction.create') }}</el-button>
    </div>

    <el-card class="filter-card">
      <el-form inline>
        <el-form-item :label="t('transaction.type')">
          <el-select v-model="filter.type" clearable :placeholder="t('common.selectPlaceholder')">
            <el-option :label="t('transaction.deposit')" :value="TransactionType.Deposit" />
            <el-option :label="t('transaction.withdrawal')" :value="TransactionType.Withdrawal" />
            <el-option :label="t('transaction.transfer')" :value="TransactionType.Transfer" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('common.startDate')">
          <DateRangePicker
            v-model:start-date="filter.start_date"
            v-model:end-date="filter.end_date"
          />
        </el-form-item>
        <el-form-item :label="t('transaction.search')">
          <el-input v-model="filter.keyword" clearable :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleFilter">{{ t('transaction.filter') }}</el-button>
          <el-button @click="resetFilter">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-table :data="transactions" v-loading="loading" stripe style="margin-top: 16px">
      <el-table-column prop="date" :label="t('transaction.date')" width="120">
        <template #default="{ row }">{{ formatDate(row.date) }}</template>
      </el-table-column>
      <el-table-column prop="description" :label="t('transaction.description')" />
      <el-table-column prop="amount" :label="t('transaction.amount')" width="150">
        <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
      </el-table-column>
      <el-table-column prop="source_account_name" :label="t('transaction.sourceAccount')" width="150" />
      <el-table-column prop="category_name" :label="t('transaction.category')" width="150" />
      <el-table-column :label="t('common.edit')" width="160" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleEdit(row.id)">{{ t('common.edit') }}</el-button>
          <el-button link type="danger" @click="handleDelete(row.id)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <Pagination
      :total="total"
      :page="filter.page"
      :page-size="filter.page_size"
      @update:page="handlePageChange"
      @update:page-size="handleSizeChange"
    />
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

.filter-card {
  margin-bottom: 0;
}
</style>

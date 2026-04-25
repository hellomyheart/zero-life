<script setup lang="ts">
/**
 * 交易列表页面
 * 功能：
 * - 展示交易列表，支持分页
 * - 支持按类型、日期范围、账户、分类、关键词筛选
 * - 支持创建、编辑、删除交易
 * - 支持查看交易详情
 */
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { list, remove } from '@/api/transaction'
import { formatAmount, formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAccountStore } from '@/stores/account'
import { useCategoryStore } from '@/stores/category'
import type { Transaction } from '@/types/transaction'
import { TransactionType } from '@/types/transaction'
import Pagination from '@/components/common/Pagination.vue'
import DateRangePicker from '@/components/common/DateRangePicker.vue'

const { t } = useI18n()
const router = useRouter()
const accountStore = useAccountStore()
const categoryStore = useCategoryStore()

// 响应式数据
const transactions = ref<Transaction[]>([]) // 交易列表数据
const total = ref(0) // 总记录数
const loading = ref(false) // 加载状态

// 筛选条件
const filter = reactive({
  type: undefined as TransactionType | undefined, // 交易类型
  start_date: '', // 开始日期
  end_date: '', // 结束日期
  source_account_id: '', // 源账户ID
  category_id: '', // 分类ID
  keyword: '', // 搜索关键词
  page: 1, // 当前页码
  page_size: 20, // 每页记录数
})

/**
 * 获取交易列表
 * 根据筛选条件从后端获取交易数据
 */
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

    const res = await list(params as unknown as typeof filter)
    const data = res as unknown as { items: Transaction[]; total: number }
    transactions.value = data.items || []
    total.value = data.total || 0
  } catch {
    ElMessage.error(t('common.fetchError') || 'Failed to load data')
  } finally {
    loading.value = false
  }
}

/**
 * 跳转到创建交易页面
 */
function handleCreate() {
  router.push('/transactions/create')
}

/**
 * 跳转到编辑交易页面
 * @param id 交易ID
 */
function handleEdit(id: string) {
  router.push(`/transactions/${id}/edit`)
}

/**
 * 删除交易
 * 弹出确认框，确认后删除交易记录
 * @param id 交易ID
 */
async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm(t('transaction.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchTransactions()
  } catch (err) {
    if ((err as string) !== 'cancel') ElMessage.error(t('common.fetchError') || 'Failed to load data')
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
  filter.type = undefined
  filter.start_date = ''
  filter.end_date = ''
  filter.source_account_id = ''
  filter.category_id = ''
  filter.keyword = ''
  filter.page = 1
  fetchTransactions()
}

onMounted(async () => {
  await Promise.all([accountStore.fetchAccounts(), categoryStore.fetchCategories()])
  fetchTransactions()
})
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
        <el-form-item :label="t('transaction.sourceAccount')">
          <el-select v-model="filter.source_account_id" clearable :placeholder="t('common.selectPlaceholder')" filterable>
            <el-option v-for="acc in accountStore.accounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('transaction.category')">
          <el-select v-model="filter.category_id" clearable :placeholder="t('common.selectPlaceholder')" filterable>
            <el-option v-for="cat in categoryStore.categories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
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

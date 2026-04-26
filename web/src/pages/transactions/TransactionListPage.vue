<script setup lang="ts">
/**
 * 交易列表页面
 * 功能：
 * - 展示交易列表，支持分页
 * - 支持按类型、日期范围、账户、分类、关键词筛选
 * - 支持创建、编辑、删除交易
 * - 金额按类型显示正负号和颜色：存款(+)绿色，取款(-)红色，转账灰色
 * - 类型标签用不同颜色区分：存款(绿)，取款(红)，转账(蓝)
 * - 响应式：筛选区自适应网格，手机端隐藏部分列
 */
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { list, remove } from '@/api/transaction'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAccountStore } from '@/stores/account'
import { useCategoryStore } from '@/stores/category'
import { useTagStore } from '@/stores/tag'
import type { Category } from '@/types/category'
import type { Tag } from '@/types/tag'
import type { Transaction } from '@/types/transaction'
import { TransactionType } from '@/types/transaction'
import Pagination from '@/components/common/Pagination.vue'
import DateRangePicker from '@/components/common/DateRangePicker.vue'

const { t } = useI18n()
const router = useRouter()
const accountStore = useAccountStore()
const categoryStore = useCategoryStore()
const tagStore = useTagStore()

const transactions = ref<Transaction[]>([])
const total = ref(0)
const loading = ref(false)

const filter = reactive({
  type: undefined as TransactionType | undefined,
  start_date: '',
  end_date: '',
  account_id: undefined as number | undefined,
  category_ids: [] as number[],
  tag_ids: [] as number[],
  keyword: '',
  page: 1,
  page_size: 20,
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
  return transform(categoryStore.categories)
})

const tagTreeData = computed(() => {
  function transform(tags: Tag[]): { value: number; label: string; children?: { value: number; label: string }[] }[] {
    return tags.map(tag => {
      const node: { value: number; label: string; children?: { value: number; label: string }[] } = {
        value: tag.id,
        label: tag.name,
      }
      if (tag.children?.length) {
        node.children = transform(tag.children)
      }
      return node
    })
  }
  return transform(tagStore.tags)
})

function formatTransactionAmount(row: Transaction): string {
  const num = Number(row.amount)
  if (isNaN(num)) return row.amount
  const absFormatted = Math.abs(num).toLocaleString('en-US', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
  if (row.type === TransactionType.Deposit) return `+${absFormatted}`
  if (row.type === TransactionType.Withdrawal) return `-${absFormatted}`
  return absFormatted
}

function amountClass(row: Transaction): string {
  if (row.type === TransactionType.Deposit) return 'amount-deposit'
  if (row.type === TransactionType.Withdrawal) return 'amount-withdrawal'
  return 'amount-transfer'
}

function typeTagType(row: Transaction): '' | 'success' | 'danger' | 'info' {
  if (row.type === TransactionType.Deposit) return 'success'
  if (row.type === TransactionType.Withdrawal) return 'danger'
  return 'info'
}

function typeLabel(row: Transaction): string {
  if (row.type === TransactionType.Deposit) return t('transaction.deposit')
  if (row.type === TransactionType.Withdrawal) return t('transaction.withdrawal')
  return t('transaction.transfer')
}

async function fetchTransactions() {
  loading.value = true
  try {
    const params: Record<string, unknown> = { page: filter.page, page_size: filter.page_size }
    if (filter.type) params.type = filter.type
    if (filter.start_date) params.start_date = filter.start_date
    if (filter.end_date) params.end_date = filter.end_date
    if (filter.account_id) params.account_id = filter.account_id
    if (filter.category_ids.length > 0) params.category_ids = filter.category_ids.join(',')
    if (filter.tag_ids.length > 0) params.tag_ids = filter.tag_ids.join(',')
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

function handleCreate() {
  router.push('/transactions/create')
}

function handleEdit(id: number) {
  router.push(`/transactions/${id}/edit`)
}

function handleClone(row: Transaction) {
  const data = {
    type: row.type,
    description: row.description,
    amount: row.amount,
    source_id: row.source_id,
    destination_id: row.destination_id ?? '',
    category_id: row.category_id ?? '',
    notes: row.notes || '',
    tags: row.tags?.map(t => t.id).join(',') ?? '',
  }
  const query = Object.fromEntries(
    Object.entries(data).filter(([, v]) => v !== '' && v !== undefined)
  )
  router.push({ path: '/transactions/create', query })
}

async function handleDelete(id: number) {
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
  filter.account_id = undefined
  filter.category_ids = []
  filter.tag_ids = []
  filter.keyword = ''
  filter.page = 1
  fetchTransactions()
}

onMounted(async () => {
  await Promise.all([accountStore.fetchAccounts(), categoryStore.fetchCategories(), tagStore.fetchTags()])
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
      <div class="filter-grid">
        <el-form-item :label="t('transaction.type')">
          <el-select v-model="filter.type" clearable :placeholder="t('common.selectPlaceholder')" style="width: 100%">
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
          <el-select v-model="filter.account_id" clearable :placeholder="t('common.selectPlaceholder')" filterable style="width: 100%">
            <el-option v-for="acc in accountStore.accounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('transaction.category')">
          <el-tree-select
            v-model="filter.category_ids"
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
        <el-form-item :label="t('transaction.tags')">
          <el-tree-select
            v-model="filter.tag_ids"
            :data="tagTreeData"
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
        <el-form-item :label="t('transaction.search')">
          <el-input v-model="filter.keyword" clearable :placeholder="t('common.inputPlaceholder')" style="width: 100%" />
        </el-form-item>
      </div>
      <div class="filter-actions">
        <el-button type="primary" @click="handleFilter">{{ t('transaction.filter') }}</el-button>
        <el-button @click="resetFilter">{{ t('common.reset') }}</el-button>
      </div>
    </el-card>

    <el-table :data="transactions" v-loading="loading" stripe style="margin-top: 16px">
      <el-table-column prop="date" :label="t('transaction.date')" width="170">
        <template #default="{ row }">{{ formatDate(row.date) }}</template>
      </el-table-column>
      <el-table-column :label="t('transaction.type')" width="90">
        <template #default="{ row }">
          <el-tag :type="typeTagType(row)" size="small" effect="dark">{{ typeLabel(row) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="description" :label="t('transaction.description')" min-width="160" show-overflow-tooltip />
      <el-table-column :label="t('transaction.amount')" width="140" align="right">
        <template #default="{ row }">
          <span :class="amountClass(row)" class="amount-text">{{ formatTransactionAmount(row) }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="t('transaction.sourceAccount')" width="130" show-overflow-tooltip class-name="hidden-md-and-down-col">
        <template #default="{ row }">{{ row.source?.name || '' }}</template>
      </el-table-column>
      <el-table-column :label="t('transaction.destinationAccount')" width="130" show-overflow-tooltip class-name="hidden-md-and-down-col">
        <template #default="{ row }">{{ row.destination?.name || '' }}</template>
      </el-table-column>
      <el-table-column :label="t('transaction.category')" width="120" show-overflow-tooltip class-name="hidden-md-and-down-col">
        <template #default="{ row }">{{ row.category?.name || '' }}</template>
      </el-table-column>
      <el-table-column :label="t('transaction.tags')" width="150" class-name="hidden-md-and-down-col">
        <template #default="{ row }">
          <el-tag
            v-for="tag in (row.tags || []).slice(0, 3)"
            :key="tag.id"
            size="small"
            :color="tag.color"
            style="color: #fff; margin: 2px"
          >{{ tag.name }}</el-tag>
          <span v-if="(row.tags || []).length > 3" class="more-tags">+{{ row.tags.length - 3 }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.edit')" width="200" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleEdit(row.id)">{{ t('common.edit') }}</el-button>
          <el-button link type="primary" @click="handleClone(row)">{{ t('common.clone') }}</el-button>
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
.filter-card {
  margin-bottom: 0;
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

.amount-transfer {
  color: var(--app-amount-transfer);
}

.more-tags {
  font-size: 12px;
  color: var(--app-text-secondary);
  margin-left: 4px;
}

/* 平板以下隐藏部分列 */
@media (max-width: 1023px) {
  :deep(.hidden-md-and-down-col) {
    display: none;
  }
}
</style>
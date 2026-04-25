<script setup lang="ts">
/**
 * 循环交易管理页面
 * 功能：
 * - 创建循环交易模板，定义定期执行的交易
 * - 设置循环频率（日/周/月/年）和间隔
 * - 手动触发执行到期的循环交易
 * - 查看和管理所有循环交易
 * 
 * 业务流程：
 * 1. 创建循环交易模板，设置金额、账户、分类、频率等
 * 2. 系统根据频率自动计算下次执行日期
 * 3. 点击"执行到期交易"按钮，手动触发所有到期的循环交易
 * 4. 循环交易执行后创建实际交易记录，更新账户余额
 * 
 * 数据来源：后端 /api/v1/recurring-transactions 接口
 * 使用 Store：accountStore（账户列表）、categoryStore（分类列表）
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
// 账户状态管理 - 用于获取账户下拉选项
const accountStore = useAccountStore()
// 分类状态管理 - 用于获取分类下拉选项
const categoryStore = useCategoryStore()

// 循环交易列表
const items = ref<RecurringTransaction[]>([])
// 加载状态
const loading = ref(false)
// 对话框显示状态
const dialogVisible = ref(false)
// 对话框标题
const dialogTitle = ref('')
// 当前编辑的循环交易 ID
const editingId = ref<number | null>(null)
/** 分页参数 - page: 当前页码, page_size: 每页数量, total: 总记录数 */
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

// 循环频率选项
const recurrenceTypeOptions = [
  { value: RecurrenceType.Daily, label: t('recurringTransaction.daily') },
  { value: RecurrenceType.Weekly, label: t('recurringTransaction.weekly') },
  { value: RecurrenceType.Monthly, label: t('recurringTransaction.monthly') },
  { value: RecurrenceType.Yearly, label: t('recurringTransaction.yearly') },
]

// 表单数据
const form = ref<CreateRecurringTransactionReq>({
  title: '',
  type: 'withdrawal',
  amount: '0',
  source_account_id: 0,
  destination_account_id: null,
  category_id: null,
  recurrence_type: RecurrenceType.Monthly,
  repeat_interval: 1,
  start_date: '',
  description: '',
})

/**
 * 获取循环交易列表
 * 传入分页参数，从后端获取当前页的数据和总记录数
 */
async function fetchList() {
  loading.value = true
  try {
    const res = await list({ page: pagination.page, page_size: pagination.page_size }) as unknown as { items: RecurringTransaction[], total: number }
    items.value = res.items || []
    pagination.total = res.total || 0
  } catch {
    ElMessage.error(t('common.operationFailed') || 'Failed to load data')
  } finally {
    loading.value = false
  }
}

/**
 * 页码变化处理函数
 * @param page 新的页码
 */
function handlePageChange(page: number) {
  pagination.page = page
  fetchList()
}

/**
 * 每页数量变化处理函数
 * @param size 新的每页数量
 */
function handleSizeChange(size: number) {
  pagination.page_size = size
  pagination.page = 1
  fetchList()
}

/**
 * 打开创建对话框
 * 初始化表单为空值
 */
function handleCreate() {
  dialogTitle.value = t('recurringTransaction.create')
  editingId.value = null
  form.value = {
    title: '',
    type: 'withdrawal',
    amount: '0',
    source_account_id: 0,
    destination_account_id: null,
    category_id: null,
    recurrence_type: RecurrenceType.Monthly,
    repeat_interval: 1,
    start_date: '',
    description: '',
  }
  dialogVisible.value = true
}

/**
 * 打开编辑对话框
 * @param row 选中的循环交易记录
 */
function handleEdit(row: RecurringTransaction) {
  dialogTitle.value = t('recurringTransaction.edit')
  editingId.value = row.id
  form.value = {
    title: row.title,
    type: row.type,
    amount: row.amount,
    source_account_id: row.source_account_id,
    destination_account_id: row.destination_account_id,
    category_id: row.category_id,
    recurrence_type: row.recurrence_type,
    repeat_interval: row.repeat_interval,
    start_date: row.start_date,
    end_date: row.end_date,
    description: row.description,
  }
  dialogVisible.value = true
}

/**
 * 删除循环交易
 * @param id 循环交易 ID
 */
async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm(t('recurringTransaction.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchList()
  } catch {
    ElMessage.error(t('common.operationFailed') || 'Failed to delete')
  }
}

/**
 * 提交表单数据
 * 根据 editingId 判断是创建还是更新操作
 */
async function handleSubmit() {
  if (!form.value.title) {
    ElMessage.warning(t('common.required'))
    return
  }
  try {
    if (editingId.value) {
      // 更新操作
      await update(editingId.value, form.value as UpdateRecurringTransactionReq)
    } else {
      // 创建操作
      await create(form.value)
    }
    ElMessage.success(t('common.success'))
    dialogVisible.value = false
    await fetchList()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

/**
 * 执行所有到期的循环交易
 * 调用后端 API 处理所有到期日期在今天或之前的循环交易
 */
async function handleProcessDue() {
  try {
    await processDue()
    ElMessage.success(t('common.success'))
    await fetchList()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

// 组件挂载时加载账户、分类列表和循环交易列表
onMounted(async () => {
  // 打开页面时同时加载账户和分类列表，供下拉选择使用
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

    <Pagination
      :total="pagination.total"
      :page="pagination.page"
      :page-size="pagination.page_size"
      @update:page="handlePageChange"
      @update:page-size="handleSizeChange"
    />

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
        <!-- 来源账户选择下拉框 - 选择交易扣款的账户 -->
        <el-form-item :label="t('transaction.sourceAccount')">
          <el-select v-model="form.source_account_id" :placeholder="t('common.selectPlaceholder')" filterable clearable>
            <el-option v-for="acc in accountStore.accounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <!-- 目标账户选择下拉框 - 转账时选择收款账户 -->
        <el-form-item v-if="form.type === 'transfer'" :label="t('transaction.destinationAccount')">
          <el-select v-model="form.destination_account_id" :placeholder="t('common.selectPlaceholder')" filterable clearable>
            <el-option v-for="acc in accountStore.accounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <!-- 分类选择下拉框 - 选择交易所属的分类 -->
        <el-form-item :label="t('transaction.category')">
          <el-select v-model="form.category_id" :placeholder="t('common.selectPlaceholder')" filterable clearable>
            <el-option v-for="cat in categoryStore.categories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
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

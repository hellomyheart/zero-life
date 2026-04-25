<script setup lang="ts">
/**
 * 账单列表页面
 * 功能：
 * - 展示账单列表，支持分页
 * - 显示账单重复频率、到期日和逾期状态
 * - 支持创建、编辑、删除账单
 * - 关联账户和分类选择
 */
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, create, update, remove } from '@/api/bill'
import { formatAmount, formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAccountStore } from '@/stores/account'
import { useCategoryStore } from '@/stores/category'
import type { Bill, CreateBillReq } from '@/types/bill'
import { RepeatRule } from '@/types/bill'
import Pagination from '@/components/common/Pagination.vue'

const { t } = useI18n()
const accountStore = useAccountStore()
const categoryStore = useCategoryStore()

/** 账单列表数据 */
const bills = ref<Bill[]>([])
/** 加载状态 */
const loading = ref(false)
/** 对话框显示状态 */
const dialogVisible = ref(false)
/** 对话框标题 */
const dialogTitle = ref('')
/** 当前编辑的账单ID，null表示新建模式 */
const editingId = ref<number | null>(null)

/** 分页参数 */
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0,
})

/** 表单数据 - 创建/编辑账单时使用 */
const form = ref<CreateBillReq>({
  name: '',
  amount: '0',
  account_id: 0,
  category_id: null,
  repeat_rule: RepeatRule.Monthly,
  next_due_date: '',
  description: '',
})

/** 重复规则选项 */
const repeatRuleOptions = [
  { value: RepeatRule.Daily, label: t('bill.daily') },
  { value: RepeatRule.Weekly, label: t('bill.weekly') },
  { value: RepeatRule.Monthly, label: t('bill.monthly') },
  { value: RepeatRule.Quarterly, label: t('bill.quarterly') },
  { value: RepeatRule.Yearly, label: t('bill.yearly') },
]

/**
 * 获取账单列表
 * 从后端API获取账单数据
 */
async function fetchBills() {
  loading.value = true
  try {
    const res = await list()
    const data = res as unknown as { items: Bill[]; total: number }
    bills.value = data.items || (res as unknown as Bill[])
    pagination.total = data.total || 0
  } catch {
    ElMessage.error(t('common.fetchError') || 'Failed to load data')
  } finally {
    loading.value = false
  }
}

/** 打开创建账单对话框 */
function handleCreate() {
  dialogTitle.value = t('bill.create')
  editingId.value = null
  form.value = { name: '', amount: '0', account_id: 0, category_id: null, repeat_rule: RepeatRule.Monthly, next_due_date: '', description: '' }
  dialogVisible.value = true
}

/**
 * 打开编辑账单对话框
 * @param bill 要编辑的账单数据
 */
function handleEdit(bill: Bill) {
  dialogTitle.value = t('bill.edit')
  editingId.value = bill.id
  form.value = {
    name: bill.name,
    amount: bill.amount,
    account_id: bill.account_id,
    category_id: bill.category_id,
    repeat_rule: bill.repeat_rule,
    next_due_date: bill.next_due_date,
    description: bill.description,
  }
  dialogVisible.value = true
}

/**
 * 删除账单
 * @param id 账单ID
 */
async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm(t('bill.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchBills()
  } catch (err) {
    if ((err as string) !== 'cancel') ElMessage.error(t('common.fetchError') || 'Failed to load data')
  }
}

/**
 * 提交表单
 * 根据editingId判断是创建还是编辑操作
 */
async function handleSubmit() {
  if (!form.value.name) {
    ElMessage.warning(t('common.required'))
    return
  }
  try {
    if (editingId.value) {
      await update(editingId.value, form.value)
    } else {
      await create(form.value)
    }
    ElMessage.success(t('common.success'))
    dialogVisible.value = false
    await fetchBills()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

/** 页码变化 */
function handlePageChange(page: number) {
  pagination.page = page
  fetchBills()
}

/** 每页数量变化 */
function handleSizeChange(size: number) {
  pagination.page_size = size
  pagination.page = 1
  fetchBills()
}

onMounted(async () => {
  await Promise.all([accountStore.fetchAccounts(), categoryStore.fetchCategories()])
  await fetchBills()
})
</script>

<template>
  <div class="bill-list-page">
    <div class="page-header">
      <h2>{{ t('bill.title') }}</h2>
      <el-button type="primary" @click="handleCreate">{{ t('bill.create') }}</el-button>
    </div>

    <el-table :data="bills" v-loading="loading" stripe>
      <el-table-column prop="name" :label="t('bill.name')" />
      <el-table-column prop="amount" :label="t('bill.amount')" width="150">
        <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
      </el-table-column>
      <el-table-column prop="repeat_rule" :label="t('bill.repeatRule')" width="120">
        <template #default="{ row }">
          {{ repeatRuleOptions.find(r => r.value === row.repeat_rule)?.label || row.repeat_rule }}
        </template>
      </el-table-column>
      <el-table-column prop="next_due_date" :label="t('bill.nextDueDate')" width="150">
        <template #default="{ row }">{{ formatDate(row.next_due_date) }}</template>
      </el-table-column>
      <el-table-column :label="t('bill.overdue')" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.is_overdue" type="danger" size="small">{{ t('bill.overdue') }}</el-tag>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item :label="t('bill.name')">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="t('bill.amount')">
          <el-input v-model="form.amount" />
        </el-form-item>
        <el-form-item :label="t('bill.repeatRule')">
          <el-select v-model="form.repeat_rule">
            <el-option v-for="opt in repeatRuleOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('bill.nextDueDate')">
          <el-date-picker v-model="form.next_due_date" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
        <!-- 账户选择下拉框 - 选择账单关联的扣款账户 -->
        <el-form-item :label="t('transaction.sourceAccount')">
          <el-select v-model="form.account_id" :placeholder="t('common.selectPlaceholder')" filterable clearable>
            <el-option v-for="acc in accountStore.accounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <!-- 分类选择下拉框 - 选择账单所属的分类 -->
        <el-form-item :label="t('transaction.category')">
          <el-select v-model="form.category_id" :placeholder="t('common.selectPlaceholder')" filterable clearable>
            <el-option v-for="cat in categoryStore.categories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
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

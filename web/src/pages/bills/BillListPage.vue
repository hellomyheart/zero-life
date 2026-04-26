<script setup lang="ts">
/**
 * 账单列表页面
 * 功能：
 * - 展示账单列表，支持分页
 * - 显示账单重复频率、到期日和逾期状态
 * - 支持创建、编辑、删除账单
 * - 关联账户和分类选择
 * 
 * 字段名与后端 JSON tag 完全对应：
 * - source_id（非 account_id）
 * - next_due（非 next_due_date）
 * - notes（非 description）
 */
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, create, update, remove } from '@/api/bill'
import { formatAmount, formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAccountStore } from '@/stores/account'
import { useCategoryStore } from '@/stores/category'
import type { Bill, CreateBillReq, UpdateBillReq } from '@/types/bill'
import { RepeatRule } from '@/types/bill'
import Pagination from '@/components/common/Pagination.vue'
import dayjs from 'dayjs'

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

/**
 * 表单数据 - 创建/编辑账单时使用
 * 字段名与后端 CreateBillReq / UpdateBillReq 的 JSON tag 完全对应
 */
const form = ref<CreateBillReq & { source_id: number | null }>({
  name: '',
  amount: '0',
  repeat_rule: RepeatRule.Monthly,
  next_due: '',
  source_id: null,
  category_id: null,
  notes: '',
})

/** 重复规则选项（后端仅支持 daily/weekly/monthly/yearly） */
const repeatRuleOptions = [
  { value: RepeatRule.Daily, label: t('bill.daily') },
  { value: RepeatRule.Weekly, label: t('bill.weekly') },
  { value: RepeatRule.Monthly, label: t('bill.monthly') },
  { value: RepeatRule.Yearly, label: t('bill.yearly') },
]

/**
 * 判断账单是否逾期
 * 后端 BillResp 不返回 is_overdue 字段，前端根据 next_due 与当前日期比较计算
 * @param bill 账单数据
 * @returns 是否逾期
 */
function isOverdue(bill: Bill): boolean {
  if (!bill.next_due) return false
  return dayjs(bill.next_due).isBefore(dayjs(), 'day')
}

/**
 * 获取账单列表
 * 从后端API获取账单数据
 */
async function fetchBills() {
  loading.value = true
  try {
    const res = await list()
    const data = res as unknown as Bill[]
    bills.value = data
    pagination.total = data.length
  } catch {
    ElMessage.error(t('common.fetchError'))
  } finally {
    loading.value = false
  }
}

/** 打开创建账单对话框 */
function handleCreate() {
  dialogTitle.value = t('bill.create')
  editingId.value = null
  form.value = {
    name: '',
    amount: '0',
    repeat_rule: RepeatRule.Monthly,
    next_due: '',
    source_id: null,
    category_id: null,
    notes: '',
  }
  dialogVisible.value = true
}

/**
 * 打开编辑账单对话框
 * 将后端返回的 Bill 数据映射到表单字段
 * @param bill 要编辑的账单数据
 */
function handleEdit(bill: Bill) {
  dialogTitle.value = t('bill.edit')
  editingId.value = bill.id
  form.value = {
    name: bill.name,
    amount: bill.amount,
    repeat_rule: bill.repeat_rule,
    next_due: bill.next_due ? dayjs(bill.next_due).format('YYYY-MM-DD') : '',
    source_id: bill.source_id ?? null,
    category_id: bill.category_id ?? null,
    notes: bill.notes || '',
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
    if ((err as string) !== 'cancel') ElMessage.error(t('common.failed'))
  }
}

/**
 * 提交表单
 * 根据editingId判断是创建还是编辑操作
 * 提交前将 source_id 为 0 或 null 的字段设为 null（后端 *uint64 期望 null 表示未选择）
 */
async function handleSubmit() {
  if (!form.value.name) {
    ElMessage.warning(t('common.required'))
    return
  }
  try {
    // 构建请求数据，确保 source_id 和 category_id 的 null 值正确传递
    const payload: CreateBillReq | UpdateBillReq = {
      name: form.value.name,
      amount: form.value.amount,
      repeat_rule: form.value.repeat_rule,
      next_due: form.value.next_due,
      source_id: form.value.source_id || null,
      category_id: form.value.category_id || null,
      notes: form.value.notes,
    }

    if (editingId.value) {
      await update(editingId.value, payload as UpdateBillReq)
    } else {
      await create(payload as CreateBillReq)
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
      <!-- next_due 字段名与后端 BillResp JSON tag 对应 -->
      <el-table-column prop="next_due" :label="t('bill.nextDueDate')" width="150">
        <template #default="{ row }">{{ formatDate(row.next_due) }}</template>
      </el-table-column>
      <!-- 逾期状态由前端根据 next_due 与当前日期比较计算 -->
      <el-table-column :label="t('bill.overdue')" width="100">
        <template #default="{ row }">
          <el-tag v-if="isOverdue(row)" type="danger" size="small">{{ t('bill.overdue') }}</el-tag>
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
          <el-select v-model="form.repeat_rule" style="width: 100%">
            <el-option v-for="opt in repeatRuleOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <!-- next_due 字段名与后端 CreateBillReq JSON tag 对应 -->
        <el-form-item :label="t('bill.nextDueDate')">
          <el-date-picker v-model="form.next_due" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <!-- source_id 字段名与后端 CreateBillReq JSON tag 对应（非 account_id） -->
        <el-form-item :label="t('transaction.sourceAccount')">
          <el-select v-model="form.source_id" :placeholder="t('common.selectPlaceholder')" filterable clearable style="width: 100%">
            <el-option v-for="acc in accountStore.accounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('transaction.category')">
          <el-select v-model="form.category_id" :placeholder="t('common.selectPlaceholder')" filterable clearable style="width: 100%">
            <el-option v-for="cat in categoryStore.categories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
        </el-form-item>
        <!-- notes 字段名与后端 CreateBillReq JSON tag 对应（非 description） -->
        <el-form-item :label="t('bill.notes')">
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

<script setup lang="ts">
// 存钱罐列表页面 - 展示储蓄目标和进度支持存取款
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, create, update, remove, addAmount, removeAmount } from '@/api/piggyBank'
import { formatAmount, formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAccountStore } from '@/stores/account'
import type { PiggyBank, CreatePiggyBankReq, UpdatePiggyBankReq, AddAmountReq } from '@/types/piggyBank'
import Pagination from '@/components/common/Pagination.vue'

const { t } = useI18n()
// 账户状态管理 - 用于获取账户下拉选项
const accountStore = useAccountStore()

/** 存钱罐列表数据 */
const piggyBanks = ref<PiggyBank[]>([])
/** 加载状态 */
const loading = ref(false)
/** 对话框显示状态 */
const dialogVisible = ref(false)
/** 对话框标题 */
const dialogTitle = ref('')
/** 当前编辑的存钱罐 ID，null表示新建 */
const editingId = ref<string | null>(null)
/** 存取款对话框显示状态 */
const amountDialogVisible = ref(false)
/** 存取款对话框标题 */
const amountDialogTitle = ref('')
/** 存取款对话框类型：add=存入, remove=取出 */
const amountDialogType = ref<'add' | 'remove'>('add')
/** 当前操作的存钱罐 ID */
const amountDialogId = ref<string>('')
/** 存取款表单数据 */
const amountForm = ref<AddAmountReq>({ amount: '0', note: '' })
/** 分页参数 - page: 当前页码, page_size: 每页数量, total: 总记录数 */
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

/** 存钱罐创建/编辑表单数据 */
const form = ref<CreatePiggyBankReq>({
  name: '',
  account_id: '',
  target_amount: '0',
  target_date: '',
  notes: '',
})

/**
 * 获取存钱罐列表
 * 传入分页参数，从后端获取当前页的数据和总记录数
 */
async function fetchPiggyBanks() {
  loading.value = true
  try {
    const res = await list({ page: pagination.page, page_size: pagination.page_size }) as unknown as { items: PiggyBank[], total: number }
    piggyBanks.value = res.items || []
    pagination.total = res.total || 0
  } catch {
    // handle error
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
  fetchPiggyBanks()
}

/**
 * 每页数量变化处理函数
 * @param size 新的每页数量
 */
function handleSizeChange(size: number) {
  pagination.page_size = size
  pagination.page = 1
  fetchPiggyBanks()
}

/**
 * 打开创建存钱罐对话框
 * 初始化表单为空值
 */
function handleCreate() {
  dialogTitle.value = t('piggyBank.create')
  editingId.value = null
  form.value = { name: '', account_id: '', target_amount: '0', target_date: '', notes: '' }
  dialogVisible.value = true
}

/**
 * 打开编辑存钱罐对话框
 * @param row 选中的存钱罐数据
 */
function handleEdit(row: PiggyBank) {
  dialogTitle.value = t('piggyBank.edit')
  editingId.value = row.id
  form.value = {
    name: row.name,
    account_id: row.account_id,
    target_amount: row.target_amount,
    target_date: row.target_date,
    notes: row.notes,
  }
  dialogVisible.value = true
}

/**
 * 删除存钱罐
 * 弹出确认框后调用API删除
 * @param id 存钱罐ID
 */
async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm(t('piggyBank.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchPiggyBanks()
  } catch {
    // cancelled or error
  }
}

/**
 * 提交创建/编辑表单
 * 根据editingId判断是创建还是更新操作
 */
async function handleSubmit() {
  if (!form.value.name) {
    ElMessage.warning(t('common.required'))
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

/**
 * 打开存入金额对话框
 * @param row 选中的存钱罐数据
 */
function handleAddMoney(row: PiggyBank) {
  amountDialogTitle.value = t('piggyBank.addMoney')
  amountDialogType.value = 'add'
  amountDialogId.value = row.id
  amountForm.value = { amount: '0', note: '' }
  amountDialogVisible.value = true
}

/**
 * 打开取出金额对话框
 * @param row 选中的存钱罐数据
 */
function handleRemoveMoney(row: PiggyBank) {
  amountDialogTitle.value = t('piggyBank.removeMoney')
  amountDialogType.value = 'remove'
  amountDialogId.value = row.id
  amountForm.value = { amount: '0', note: '' }
  amountDialogVisible.value = true
}

/**
 * 提交存取款操作
 * 根据amountDialogType判断是存入还是取出，调用对应API
 */
async function handleAmountSubmit() {
  try {
    if (amountDialogType.value === 'add') {
      await addAmount(amountDialogId.value, amountForm.value)
    } else {
      await removeAmount(amountDialogId.value, amountForm.value)
    }
    ElMessage.success(t('common.success'))
    amountDialogVisible.value = false
    await fetchPiggyBanks()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

onMounted(async () => {
  // 打开页面时加载账户列表，供下拉选择使用
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

    <el-table :data="piggyBanks" v-loading="loading" stripe>
      <el-table-column prop="name" :label="t('piggyBank.name')" />
      <el-table-column prop="target_amount" :label="t('piggyBank.targetAmount')" width="140">
        <template #default="{ row }">{{ formatAmount(row.target_amount) }}</template>
      </el-table-column>
      <el-table-column prop="current_amount" :label="t('piggyBank.currentAmount')" width="140">
        <template #default="{ row }">{{ formatAmount(row.current_amount) }}</template>
      </el-table-column>
      <el-table-column :label="t('piggyBank.percentage')" width="180">
        <template #default="{ row }">
          <el-progress :percentage="row.percentage || 0" :stroke-width="18" />
        </template>
      </el-table-column>
      <el-table-column prop="target_date" :label="t('piggyBank.targetDate')" width="120">
        <template #default="{ row }">{{ row.target_date ? formatDate(row.target_date) : '-' }}</template>
      </el-table-column>
      <el-table-column :label="t('common.edit')" width="220" fixed="right">
        <template #default="{ row }">
          <el-button link type="success" @click="handleAddMoney(row)">{{ t('piggyBank.addMoney') }}</el-button>
          <el-button link type="warning" @click="handleRemoveMoney(row)">{{ t('piggyBank.removeMoney') }}</el-button>
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
        <el-form-item :label="t('piggyBank.name')">
          <el-input v-model="form.name" />
        </el-form-item>
        <!-- 账户选择下拉框 - 选择存钱罐关联的账户 -->
        <el-form-item :label="t('transaction.sourceAccount')">
          <el-select v-model="form.account_id" :placeholder="t('common.selectPlaceholder')" filterable clearable>
            <el-option v-for="acc in accountStore.accounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('piggyBank.targetAmount')">
          <el-input v-model="form.target_amount" />
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
          <el-input v-model="amountForm.amount" />
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
</style>

<script setup lang="ts">
/**
 * 对账管理页面
 * 功能：
 * - 创建对账记录，核对账户余额与实际余额
 * - 查看对账历史，包括期初余额、期末余额和差额
 * - 编辑和删除对账记录
 * 
 * 业务流程：
 * 1. 选择要对账的账户
 * 2. 输入对账周期（开始日期和结束日期）
 * 3. 输入期初余额和期末余额
 * 4. 系统自动计算差额（期末余额 - 账面余额）
 * 5. 保存对账记录，状态为 open
 * 
 * 数据来源：后端 /api/v1/reconciliations 接口
 * 使用 Store：accountStore（获取账户列表）
 */
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, create, update, remove } from '@/api/reconciliation'
import { formatAmount, formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAccountStore } from '@/stores/account'
import type { Reconciliation, CreateReconciliationReq, UpdateReconciliationReq } from '@/types/reconciliation'

const { t } = useI18n()
// 账户状态管理 - 用于获取账户下拉选项
const accountStore = useAccountStore()

// 对账记录列表
const items = ref<Reconciliation[]>([])
// 加载状态
const loading = ref(false)
// 对话框显示状态
const dialogVisible = ref(false)
// 对话框标题
const dialogTitle = ref('')
// 当前编辑的对账记录 ID
const editingId = ref<string | null>(null)

// 表单数据
const form = ref<CreateReconciliationReq>({
  account_id: '',
  start_date: '',
  end_date: '',
  start_balance: '0',
  end_balance: '0',
  notes: '',
})

/**
 * 获取对账记录列表
 * 从后端 API 获取所有对账记录
 */
async function fetchList() {
  loading.value = true
  try {
    const res = await list({}) as unknown as { items: Reconciliation[] }
    items.value = res.items || []
  } catch (error) {
    console.error('Failed to fetch reconciliation list:', error)
  } finally {
    loading.value = false
  }
}

/**
 * 打开创建对话框
 * 初始化表单为空值
 */
function handleCreate() {
  dialogTitle.value = t('reconciliation.create')
  editingId.value = null
  form.value = { account_id: '', start_date: '', end_date: '', start_balance: '0', end_balance: '0', notes: '' }
  dialogVisible.value = true
}

/**
 * 打开编辑对话框
 * @param row 选中的对账记录
 */
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

/**
 * 删除对账记录
 * @param id 对账记录 ID
 */
async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm(t('reconciliation.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchList()
  } catch (error) {
    console.error('Failed to delete reconciliation:', error)
  }
}

/**
 * 提交表单数据
 * 根据 editingId 判断是创建还是更新操作
 */
async function handleSubmit() {
  try {
    if (editingId.value) {
      // 更新操作：只允许更新期末余额和备注
      await update(editingId.value, { end_balance: form.value.end_balance, notes: form.value.notes } as UpdateReconciliationReq)
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

// 组件挂载时加载账户列表和对账记录
onMounted(async () => {
  // 打开页面时加载账户列表，供下拉选择使用
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
        <!-- 账户选择下拉框 - 选择要对账的账户 -->
        <el-form-item :label="t('transaction.sourceAccount')">
          <el-select v-model="form.account_id" :placeholder="t('common.selectPlaceholder')" filterable clearable>
            <el-option v-for="acc in accountStore.accounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
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

<script setup lang="ts">
// 账单列表页面 - 展示账单重复频率和到期日
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, create, update, remove } from '@/api/bill'
import { formatAmount, formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Bill, CreateBillReq } from '@/types/bill'
import { RepeatRule } from '@/types/bill'

const { t } = useI18n()

const bills = ref<Bill[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref<string | null>(null)

const form = ref<CreateBillReq>({
  name: '',
  amount: '0',
  account_id: '',
  category_id: '',
  repeat_rule: RepeatRule.Monthly,
  next_due_date: '',
  description: '',
})

const repeatRuleOptions = [
  { value: RepeatRule.Daily, label: t('bill.daily') },
  { value: RepeatRule.Weekly, label: t('bill.weekly') },
  { value: RepeatRule.Monthly, label: t('bill.monthly') },
  { value: RepeatRule.Quarterly, label: t('bill.quarterly') },
  { value: RepeatRule.Yearly, label: t('bill.yearly') },
]

async function fetchBills() {
  loading.value = true
  try {
    bills.value = await list() as unknown as Bill[]
  } catch {
    // handle error
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  dialogTitle.value = t('bill.create')
  editingId.value = null
  form.value = { name: '', amount: '0', account_id: '', category_id: '', repeat_rule: RepeatRule.Monthly, next_due_date: '', description: '' }
  dialogVisible.value = true
}

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

async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm(t('bill.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchBills()
  } catch {
    // cancelled or error
  }
}

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

onMounted(fetchBills)
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

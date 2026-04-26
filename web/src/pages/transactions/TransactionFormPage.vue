<script setup lang="ts">
// 交易编辑页面 - 用于创建和编辑交易，支持拆分交易
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { create, update, getTransaction } from '@/api/transaction'
import { useAccountStore } from '@/stores/account'
import { useCategoryStore } from '@/stores/category'
import { useTagStore } from '@/stores/tag'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import type { CreateTransactionReq, Transaction } from '@/types/transaction'
import type { Category } from '@/types/category'
import type { Tag } from '@/types/tag'
import { TransactionType } from '@/types/transaction'
import AmountInput from '@/components/common/AmountInput.vue'

import { AccountType } from '@/types/account'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const accountStore = useAccountStore()
const categoryStore = useCategoryStore()
const tagStore = useTagStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const isEdit = computed(() => !!route.params.id)
const enableSplits = ref(false)

// 获取当前本地时间，格式为 YYYY-MM-DD HH:mm（带设备时区）
function getLocalDateTimeStr(): string {
  const now = new Date()
  const y = now.getFullYear()
  const m = String(now.getMonth() + 1).padStart(2, '0')
  const d = String(now.getDate()).padStart(2, '0')
  const h = String(now.getHours()).padStart(2, '0')
  const min = String(now.getMinutes()).padStart(2, '0')
  return `${y}-${m}-${d} ${h}:${min}`
}

function formatDateForPicker(dateStr: string): string {
  if (!dateStr) return getLocalDateTimeStr()
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return dateStr
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const h = String(d.getHours()).padStart(2, '0')
  const min = String(d.getMinutes()).padStart(2, '0')
  return `${y}-${m}-${day} ${h}:${min}`
}

const form = reactive<CreateTransactionReq>({
  type: TransactionType.Withdrawal,
  date: getLocalDateTimeStr(),
  description: '',
  amount: '0',
  source_id: 0,
  destination_id: undefined,
  category_id: undefined,
  notes: '',
  tags: [],
  splits: [],
})

const rules: FormRules = {
  type: [{ required: true, message: t('common.required'), trigger: 'change' }],
  date: [{ required: true, message: t('common.required'), trigger: 'change' }],
  description: [{ required: true, message: t('common.required'), trigger: 'blur' }],
  amount: [{ required: true, message: t('common.required'), trigger: 'blur' }],
  source_id: [
    { required: true, message: t('common.required'), trigger: 'change' },
    {
      validator: (_rule: unknown, value: number | undefined, callback: (err?: Error) => void) => {
        if (value && form.destination_id && value === form.destination_id) {
          callback(new Error(t('transaction.sameAccountError')))
        } else {
          callback()
        }
      },
      trigger: 'change',
    },
  ],
  destination_id: [
    { required: true, message: t('common.required'), trigger: 'change' },
    {
      validator: (_rule: unknown, value: number | undefined, callback: (err?: Error) => void) => {
        if (value && form.source_id && value === form.source_id) {
          callback(new Error(t('transaction.sameAccountError')))
        } else {
          callback()
        }
      },
      trigger: 'change',
    },
  ],
}

const transactionTypeOptions = [
  { value: TransactionType.Deposit, label: t('transaction.deposit') },
  { value: TransactionType.Withdrawal, label: t('transaction.withdrawal') },
  { value: TransactionType.Transfer, label: t('transaction.transfer') },
]

// 根据交易类型筛选账户选项
const assetAccounts = computed(() => accountStore.accounts.filter((a: { type: string }) => a.type === AccountType.Asset))
const expenseAccounts = computed(() => accountStore.accounts.filter((a: { type: string }) => a.type === AccountType.Expense))
const revenueAccounts = computed(() => accountStore.accounts.filter((a: { type: string }) => a.type === AccountType.Revenue))

// 取款：source=资产账户，destination=支出账户
// 存款：source=收入账户，destination=资产账户
// 转账：source=资产账户，destination=资产账户
const sourceAccountOptions = computed(() => {
  if (form.type === TransactionType.Deposit) return revenueAccounts.value
  return assetAccounts.value
})
const destinationAccountOptions = computed(() => {
  if (form.type === TransactionType.Withdrawal) return expenseAccounts.value
  return assetAccounts.value
})
const showDestination = computed(() => true)
const sourceAccountLabel = computed(() => {
  if (form.type === TransactionType.Deposit) return t('transaction.sourceAccount') + '（' + t('transaction.revenueAccount') + '）'
  return t('transaction.sourceAccount') + '（' + t('transaction.assetAccount') + '）'
})
const destinationAccountLabel = computed(() => {
  if (form.type === TransactionType.Withdrawal) return t('transaction.destinationAccount') + '（' + t('transaction.expenseAccount') + '）'
  return t('transaction.destinationAccount') + '（' + t('transaction.assetAccount') + '）'
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

function handleTypeChange() {
  form.source_id = 0
  form.destination_id = undefined
}

function addSplit() {
  if (!form.splits) form.splits = []
  form.splits.push({ amount: '0', category_id: undefined, tags: [], description: '', notes: '' })
}

function removeSplit(index: number) {
  form.splits?.splice(index, 1)
}

onMounted(async () => {
  await Promise.all([
    accountStore.fetchAccounts(),
    categoryStore.fetchCategories(),
    tagStore.fetchTags(),
  ])

  if (isEdit.value) {
    try {
      const tx = (await getTransaction(Number(route.params.id))) as unknown as Transaction
      form.type = tx.type
      form.date = formatDateForPicker(tx.date)
      form.description = tx.description
      form.amount = tx.amount
      form.source_id = tx.source_id
      form.destination_id = tx.destination_id ?? undefined
      form.category_id = tx.category_id ?? undefined
      form.notes = tx.notes || ''
      form.tags = tx.tags?.map(tag => tag.id) ?? []
      form.splits = tx.splits?.map(s => ({
        amount: s.amount,
        category_id: s.category_id ?? undefined,
        tags: s.tags?.map(t => t.id) ?? [],
        description: '',
        notes: s.notes || '',
      })) ?? []
      enableSplits.value = (tx.splits?.length ?? 0) > 0
    } catch {
      ElMessage.error(t('common.failed'))
    }
  }
})

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const data = { ...form }
    if (!enableSplits.value) {
      data.splits = undefined
    }
    if (isEdit.value) {
      await update(Number(route.params.id), data)
    } else {
      await create(data)
    }
    ElMessage.success(t('common.success'))
    router.push('/transactions')
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="transaction-form-page">
    <h2>{{ isEdit ? t('transaction.edit') : t('transaction.create') }}</h2>
    <el-card>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item :label="t('transaction.type')" prop="type">
          <el-radio-group v-model="form.type" @change="handleTypeChange">
            <el-radio v-for="opt in transactionTypeOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="t('transaction.date')" prop="date">
          <el-date-picker v-model="form.date" type="datetime" value-format="YYYY-MM-DD HH:mm" />
        </el-form-item>
        <el-form-item :label="t('transaction.description')" prop="description">
          <el-input v-model="form.description" :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('transaction.amount')" prop="amount">
          <AmountInput v-model="form.amount" />
        </el-form-item>
        <el-form-item :label="sourceAccountLabel" prop="source_id">
          <el-select v-model="form.source_id" :placeholder="t('common.selectPlaceholder')" filterable>
            <el-option v-for="acc in sourceAccountOptions" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="showDestination" :label="destinationAccountLabel" prop="destination_id">
          <el-select v-model="form.destination_id" :placeholder="t('common.selectPlaceholder')" filterable>
            <el-option v-for="acc in destinationAccountOptions" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('transaction.category')">
          <el-tree-select
            v-model="form.category_id"
            :data="categoryTreeData"
            :placeholder="t('common.selectPlaceholder')"
            check-strictly
            filterable
            clearable
            :render-after-expand="false"
          />
        </el-form-item>
        <el-form-item :label="t('transaction.tags')">
          <el-tree-select
            v-model="form.tags"
            :data="tagTreeData"
            :placeholder="t('common.selectPlaceholder')"
            check-strictly
            multiple
            filterable
            clearable
            :render-after-expand="false"
          />
        </el-form-item>
        <el-form-item :label="t('transaction.notes')">
          <el-input v-model="form.notes" type="textarea" :rows="2" :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>

        <el-form-item :label="t('transaction.enableSplits')">
          <el-switch v-model="enableSplits" />
        </el-form-item>

        <template v-if="enableSplits">
          <el-divider>{{ t('transaction.splits') }}</el-divider>
          <div v-for="(split, index) in form.splits" :key="index" class="split-item">
            <el-row :gutter="10">
              <el-col :span="6">
                <el-form-item :label="t('transaction.amount')">
                  <AmountInput v-model="split.amount" />
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item :label="t('transaction.category')">
                  <el-tree-select
                    v-model="split.category_id"
                    :data="categoryTreeData"
                    :placeholder="t('common.selectPlaceholder')"
                    check-strictly
                    filterable
                    clearable
                    :render-after-expand="false"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item :label="t('transaction.tags')">
                  <el-tree-select
                    v-model="split.tags"
                    :data="tagTreeData"
                    :placeholder="t('common.selectPlaceholder')"
                    check-strictly
                    multiple
                    filterable
                    clearable
                    :render-after-expand="false"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item :label="t('common.delete')">
                  <el-button type="danger" link @click="removeSplit(index)">{{ t('common.delete') }}</el-button>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="10">
              <el-col :span="12">
                <el-form-item :label="t('transaction.description')">
                  <el-input v-model="split.description" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="t('transaction.notes')">
                  <el-input v-model="split.notes" />
                </el-form-item>
              </el-col>
            </el-row>
          </div>
          <el-button type="primary" link @click="addSplit">+ {{ t('transaction.splits') }}</el-button>
        </template>

        <el-form-item style="margin-top: 20px">
          <el-button type="primary" :loading="loading" @click="handleSubmit">{{ t('common.save') }}</el-button>
          <el-button @click="router.push('/transactions')">{{ t('common.cancel') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.split-item {
  padding: 8px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}
</style>
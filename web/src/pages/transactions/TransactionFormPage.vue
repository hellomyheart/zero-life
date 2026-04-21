<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { create, update, getTransaction } from '@/api/transaction'
import { useAccountStore } from '@/stores/account'
import { useCategoryStore } from '@/stores/category'
import { useTagStore } from '@/stores/tag'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import type { CreateTransactionReq, CreateSplitReq } from '@/types/transaction'
import { TransactionType } from '@/types/transaction'
import AmountInput from '@/components/common/AmountInput.vue'

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

const form = reactive<CreateTransactionReq>({
  type: TransactionType.Withdrawal,
  date: new Date().toISOString().slice(0, 10),
  description: '',
  amount: '0',
  source_account_id: '',
  destination_account_id: '',
  category_id: '',
  tag_ids: [],
  splits: [],
})

const rules: FormRules = {
  type: [{ required: true, message: t('common.required'), trigger: 'change' }],
  date: [{ required: true, message: t('common.required'), trigger: 'change' }],
  amount: [{ required: true, message: t('common.required'), trigger: 'blur' }],
  source_account_id: [{ required: true, message: t('common.required'), trigger: 'change' }],
}

const transactionTypeOptions = [
  { value: TransactionType.Deposit, label: t('transaction.deposit') },
  { value: TransactionType.Withdrawal, label: t('transaction.withdrawal') },
  { value: TransactionType.Transfer, label: t('transaction.transfer') },
]

function addSplit() {
  if (!form.splits) form.splits = []
  form.splits.push({ amount: '0', category_id: '', tag_ids: [], description: '' })
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
      const tx = await getTransaction(route.params.id as string) as unknown as CreateTransactionReq & { id: string; splits: CreateSplitReq[] }
      form.type = tx.type
      form.date = tx.date
      form.description = tx.description
      form.amount = tx.amount
      form.source_account_id = tx.source_account_id
      form.destination_account_id = tx.destination_account_id
      form.category_id = tx.category_id
      form.tag_ids = tx.tag_ids
      form.splits = tx.splits
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
      await update(route.params.id as string, data)
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
          <el-radio-group v-model="form.type">
            <el-radio v-for="opt in transactionTypeOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="t('transaction.date')" prop="date">
          <el-date-picker v-model="form.date" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item :label="t('transaction.description')">
          <el-input v-model="form.description" :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('transaction.amount')" prop="amount">
          <AmountInput v-model="form.amount" />
        </el-form-item>
        <el-form-item :label="t('transaction.sourceAccount')" prop="source_account_id">
          <el-select v-model="form.source_account_id" :placeholder="t('common.selectPlaceholder')" filterable>
            <el-option v-for="acc in accountStore.accounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="form.type === TransactionType.Transfer" :label="t('transaction.destinationAccount')">
          <el-select v-model="form.destination_account_id" :placeholder="t('common.selectPlaceholder')" filterable>
            <el-option v-for="acc in accountStore.accounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('transaction.category')">
          <el-select v-model="form.category_id" :placeholder="t('common.selectPlaceholder')" filterable clearable>
            <el-option v-for="cat in categoryStore.categories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('transaction.tags')">
          <el-select v-model="form.tag_ids" multiple :placeholder="t('common.selectPlaceholder')" filterable clearable>
            <el-option v-for="tag in tagStore.tags" :key="tag.id" :label="tag.name" :value="tag.id" />
          </el-select>
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
                  <el-select v-model="split.category_id" :placeholder="t('common.selectPlaceholder')" filterable>
                    <el-option v-for="cat in categoryStore.categories" :key="cat.id" :label="cat.name" :value="cat.id" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item :label="t('transaction.description')">
                  <el-input v-model="split.description" />
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-button type="danger" link @click="removeSplit(index)">{{ t('common.delete') }}</el-button>
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

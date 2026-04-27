<script setup lang="ts">
// 账户表单页面 - 用于创建和编辑账户
// 编辑模式下：账户类型、货币、初始余额不可修改（只读），仅名称、备注、虚拟属性可编辑
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { create, update, getAccount } from '@/api/account'
import { useCurrencyStore } from '@/stores/currency'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import type { CreateAccountReq, UpdateAccountReq, Account } from '@/types/account'
import { AccountType } from '@/types/account'
import AmountInput from '@/components/common/AmountInput.vue'
import CurrencySelect from '@/components/common/CurrencySelect.vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const currencyStore = useCurrencyStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const isEdit = computed(() => !!route.params.id)

// 表单数据，创建和编辑共用
const form = reactive({
  name: '',
  account_number: '',
  type: AccountType.Asset as string,
  currency_id: 0 as number,
  initial_balance: '0',
  notes: '',
  is_virtual: false,
})

const rules: FormRules = {
  name: [{ required: true, message: t('common.required'), trigger: 'blur' }],
  account_number: [{ required: true, message: t('common.required'), trigger: 'blur' }],
  type: [{ required: true, message: t('common.required'), trigger: 'change' }],
  currency_id: [{ required: true, message: t('common.required'), trigger: 'change' }],
  initial_balance: [{ required: true, message: t('common.required'), trigger: 'blur' }],
}

const selectedCurrencyCode = computed(() => {
  const c = (currencyStore.currencies ?? []).find((c: { id: number }) => c.id === form.currency_id)
  return c?.code || 'CNY'
})

const accountTypeOptions = [
  { value: AccountType.Asset, label: t('account.asset') },
  { value: AccountType.Expense, label: t('account.expense') },
  { value: AccountType.Revenue, label: t('account.income') },
  { value: AccountType.Liability, label: t('account.liability') },
]

onMounted(async () => {
  await currencyStore.fetchCurrencies()
  // 新建模式：设置默认货币ID作为初始值，避免下拉框显示0
  if (!isEdit.value && currencyStore.defaultCurrency) {
    form.currency_id = currencyStore.defaultCurrency.id
  }
  if (isEdit.value) {
    try {
      const account = (await getAccount(Number(route.params.id))) as unknown as Account
      form.name = account.name
      form.account_number = account.account_number
      form.type = account.type
      form.currency_id = account.currency_id
      form.initial_balance = account.initial_balance
      form.notes = account.notes || ''
      form.is_virtual = account.is_virtual
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
    if (isEdit.value) {
      // 编辑模式：只提交可修改的字段（名称、备注、虚拟属性）
      const req: UpdateAccountReq = {
        name: form.name,
        notes: form.notes,
        is_virtual: form.is_virtual,
      }
      await update(Number(route.params.id), req)
    } else {
      // 创建模式：提交所有字段
      const req: CreateAccountReq = {
        name: form.name,
        account_number: form.account_number,
        type: form.type as AccountType,
        currency_id: form.currency_id,
        initial_balance: form.initial_balance,
        notes: form.notes,
        is_virtual: form.is_virtual,
      }
      await create(req)
    }
    ElMessage.success(t('common.success'))
    router.push('/accounts')
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="account-form-page">
    <h2>{{ isEdit ? t('account.edit') : t('account.create') }}</h2>
    <el-card>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item :label="t('account.name')" prop="name">
          <el-input v-model="form.name" :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('account.accountNumber')" prop="account_number">
          <el-input v-model="form.account_number" :placeholder="t('common.inputPlaceholder')" :disabled="isEdit" />
        </el-form-item>
        <el-form-item :label="t('account.type')" prop="type">
          <el-select v-model="form.type" :placeholder="t('common.selectPlaceholder')" :disabled="isEdit">
            <el-option v-for="opt in accountTypeOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('account.currency')" prop="currency_id">
          <CurrencySelect v-model="form.currency_id" :currencies="currencyStore.currencies" mode="id" :disabled="isEdit" />
        </el-form-item>
        <el-form-item :label="t('account.initialBalance')" prop="initial_balance">
          <AmountInput v-model="form.initial_balance" :currency="selectedCurrencyCode" :disabled="isEdit" />
        </el-form-item>
        <el-form-item :label="t('account.isVirtual')">
          <el-switch v-model="form.is_virtual" />
        </el-form-item>
        <el-form-item :label="t('account.notes')">
          <el-input v-model="form.notes" type="textarea" :rows="3" :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">{{ t('common.save') }}</el-button>
          <el-button @click="router.push('/accounts')">{{ t('common.cancel') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>
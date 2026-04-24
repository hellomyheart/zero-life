<script setup lang="ts">
// 账户表单页面 - 用于创建和编辑账户
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { create, update, getAccount } from '@/api/account'
import { useCurrencyStore } from '@/stores/currency'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import type { CreateAccountReq } from '@/types/account'
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

const form = reactive<CreateAccountReq & { is_virtual?: boolean }>({
  name: '',
  type: AccountType.Asset,
  currency: 'CNY',
  initial_balance: '0',
  is_virtual: false,
})

const rules: FormRules = {
  name: [{ required: true, message: t('common.required'), trigger: 'blur' }],
  type: [{ required: true, message: t('common.required'), trigger: 'change' }],
  currency: [{ required: true, message: t('common.required'), trigger: 'change' }],
  initial_balance: [{ required: true, message: t('common.required'), trigger: 'blur' }],
}

const accountTypeOptions = [
  { value: AccountType.Asset, label: t('account.asset') },
  { value: AccountType.Expense, label: t('account.expense') },
  { value: AccountType.Income, label: t('account.income') },
  { value: AccountType.Liability, label: t('account.liability') },
]

onMounted(async () => {
  await currencyStore.fetchCurrencies()
  if (isEdit.value) {
    try {
      const account = await getAccount(route.params.id as string) as unknown as CreateAccountReq & { is_virtual?: boolean; id: string }
      form.name = account.name
      form.type = account.type
      form.currency = account.currency
      form.initial_balance = account.initial_balance
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
      await update(route.params.id as string, form)
    } else {
      await create(form)
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
        <el-form-item :label="t('account.type')" prop="type">
          <el-select v-model="form.type" :placeholder="t('common.selectPlaceholder')">
            <el-option v-for="opt in accountTypeOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('account.currency')" prop="currency">
          <CurrencySelect v-model="form.currency" :currencies="currencyStore.currencies" />
        </el-form-item>
        <el-form-item :label="t('account.initialBalance')" prop="initial_balance">
          <AmountInput v-model="form.initial_balance" :currency="form.currency" />
        </el-form-item>
        <el-form-item :label="t('account.isVirtual')">
          <el-switch v-model="form.is_virtual" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">{{ t('common.save') }}</el-button>
          <el-button @click="router.push('/accounts')">{{ t('common.cancel') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

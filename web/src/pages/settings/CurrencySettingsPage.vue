<script setup lang="ts">
// 货币管理页面 - 管理货币启用状态默认货币和汇率
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, updateStatus, setDefault, getExchangeRates, setExchangeRate } from '@/api/currency'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Currency, ExchangeRate, SetExchangeRateReq } from '@/types/currency'

const { t } = useI18n()

const currencies = ref<Currency[]>([])
const exchangeRates = ref<ExchangeRate[]>([])
const loading = ref(false)
const rateDialogVisible = ref(false)

const rateForm = ref<SetExchangeRateReq>({
  from_currency_id: 0,
  to_currency_id: 0,
  rate: '1',
})

async function fetchData() {
  loading.value = true
  try {
    currencies.value = await list() as unknown as Currency[]
    exchangeRates.value = await getExchangeRates() as unknown as ExchangeRate[]
  } catch {
    ElMessage.error(t('common.fetchError') || 'Failed to load data')
  } finally {
    loading.value = false
  }
}

async function handleToggleEnabled(currency: Currency) {
  try {
    await updateStatus(currency.id, { is_enabled: !currency.is_enabled })
    ElMessage.success(t('common.success'))
    await fetchData()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

async function handleSetDefault(currency: Currency) {
  try {
    await ElMessageBox.confirm(
      `Set ${currency.code} as default currency?`,
      t('common.confirm'),
      { type: 'info' }
    )
    await setDefault(currency.id)
    ElMessage.success(t('common.success'))
    await fetchData()
  } catch (err) {
    if ((err as string) !== 'cancel') ElMessage.error(t('common.fetchError') || 'Failed to load data')
  }
}

function handleAddRate() {
  rateForm.value = { from_currency_id: 0, to_currency_id: 0, rate: '1' }
  rateDialogVisible.value = true
}

async function handleRateSubmit() {
  try {
    await setExchangeRate(rateForm.value)
    ElMessage.success(t('common.success'))
    rateDialogVisible.value = false
    await fetchData()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

onMounted(fetchData)
</script>

<template>
  <div class="currency-settings-page">
    <h2>{{ t('currency.title') }}</h2>

    <el-card style="margin-bottom: 20px">
      <template #header>Currencies</template>
      <el-table :data="currencies" v-loading="loading" stripe>
        <el-table-column prop="code" :label="t('currency.code')" width="100" />
        <el-table-column prop="name" :label="t('currency.name')" />
        <el-table-column prop="symbol" :label="t('currency.symbol')" width="80" />
        <el-table-column :label="t('currency.enabled')" width="100">
          <template #default="{ row }">
            <el-switch :model-value="row.is_enabled" @change="handleToggleEnabled(row)" />
          </template>
        </el-table-column>
        <el-table-column :label="t('currency.default')" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.is_default" type="success" size="small">{{ t('currency.default') }}</el-tag>
            <el-button v-else link type="primary" size="small" @click="handleSetDefault(row)">{{ t('currency.setDefault') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>{{ t('currency.exchangeRate') }}</span>
          <el-button type="primary" size="small" @click="handleAddRate">+ {{ t('currency.exchangeRate') }}</el-button>
        </div>
      </template>
      <el-table :data="exchangeRates" stripe>
        <el-table-column :label="t('currency.sourceCurrency')">
          <template #default="{ row }">{{ currencies.find(c => c.id === row.from_currency_id)?.code || row.from_currency_id }}</template>
        </el-table-column>
        <el-table-column :label="t('currency.targetCurrency')">
          <template #default="{ row }">{{ currencies.find(c => c.id === row.to_currency_id)?.code || row.to_currency_id }}</template>
        </el-table-column>
        <el-table-column prop="rate" :label="t('currency.rate')" />
        <el-table-column prop="updated_at" label="Updated At" />
      </el-table>
    </el-card>

    <el-dialog v-model="rateDialogVisible" :title="t('currency.exchangeRate')" width="400px">
      <el-form :model="rateForm" label-width="120px">
        <el-form-item :label="t('currency.sourceCurrency')">
          <el-select v-model="rateForm.from_currency_id" filterable>
            <el-option v-for="c in currencies" :key="c.id" :label="c.code" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('currency.targetCurrency')">
          <el-select v-model="rateForm.to_currency_id" filterable>
            <el-option v-for="c in currencies" :key="c.id" :label="c.code" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('currency.rate')">
          <el-input v-model="rateForm.rate" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rateDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleRateSubmit">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

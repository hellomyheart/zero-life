<!-- 汇率管理页面 - 展示汇率列表，支持设置汇率 -->
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, set as setRate } from '@/api/exchange-rate'
import { ElMessage } from 'element-plus'
import { useCurrencyStore } from '@/stores/currency'
import type { ExchangeRate, SetExchangeRateReq } from '@/types/currency'
import type { Currency } from '@/types/currency'

const { t } = useI18n()
const currencyStore = useCurrencyStore()

const enabledCurrencies = computed(() => currencyStore.currencies.filter((c: Currency) => c.is_enabled))

const exchangeRates = ref<ExchangeRate[]>([])
const loading = ref(false)
const dialogVisible = ref(false)

const form = ref<SetExchangeRateReq>({
  from_currency_id: 0,
  to_currency_id: 0,
  rate: '1',
})

async function fetchExchangeRates() {
  loading.value = true
  try {
    exchangeRates.value = await list() as unknown as ExchangeRate[]
  } catch {
    ElMessage.error(t('common.fetchError') || 'Failed to load data')
  } finally {
    loading.value = false
  }
}

function handleSetRate() {
  form.value = { from_currency_id: 0, to_currency_id: 0, rate: '1' }
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!form.value.from_currency_id || !form.value.to_currency_id || !form.value.rate) {
    ElMessage.warning(t('common.required'))
    return
  }
  try {
    await setRate(form.value)
    ElMessage.success(t('common.success'))
    dialogVisible.value = false
    await fetchExchangeRates()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

onMounted(async () => {
  await currencyStore.fetchCurrencies()
  fetchExchangeRates()
})
</script>

<template>
  <div class="exchange-rate-list-page">
    <el-card>
      <template #header>
        <div class="header">
          <h2>{{ t('exchangeRate.title') }}</h2>
          <el-button type="primary" @click="handleSetRate">
            {{ t('common.create') }}
          </el-button>
        </div>
      </template>

      <el-table :data="exchangeRates" v-loading="loading" stripe>
        <el-table-column prop="source_currency" :label="t('exchangeRate.fromCurrency')" />
        <el-table-column prop="target_currency" :label="t('exchangeRate.toCurrency')" />
        <el-table-column prop="rate" :label="t('exchangeRate.rate')" />
        <el-table-column prop="updated_at" label="Updated At" />
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="t('exchangeRate.title')" width="500px">
      <el-form :model="form" label-width="120px">
        <el-form-item :label="t('exchangeRate.fromCurrency')">
          <el-select v-model="form.from_currency_id" :placeholder="t('common.selectPlaceholder')" filterable style="width: 100%">
            <el-option v-for="c in enabledCurrencies" :key="c.id" :label="`${c.code} - ${c.name}`" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('exchangeRate.toCurrency')">
          <el-select v-model="form.to_currency_id" :placeholder="t('common.selectPlaceholder')" filterable style="width: 100%">
            <el-option v-for="c in enabledCurrencies" :key="c.id" :label="`${c.code} - ${c.name}`" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('exchangeRate.rate')">
          <el-input v-model="form.rate" />
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
.exchange-rate-list-page {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header h2 {
  margin: 0;
}
</style>

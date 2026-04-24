// 货币状态管理 - 管理货币列表和默认货币
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Currency } from '@/types/currency'
import { list as listCurrencies } from '@/api/currency'

export const useCurrencyStore = defineStore('currency', () => {
  const currencies = ref<Currency[]>([]) // 货币列表
  const defaultCurrency = ref<Currency | null>(null) // 默认货币

  // fetchCurrencies 获取货币列表并找出默认货币
  async function fetchCurrencies() {
    currencies.value = await listCurrencies() as unknown as Currency[]
    defaultCurrency.value = currencies.value.find((c) => c.is_default) || null
  }

  return {
    currencies,
    defaultCurrency,
    fetchCurrencies,
  }
})

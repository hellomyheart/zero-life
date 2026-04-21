import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Currency } from '@/types/currency'
import { list as listCurrencies } from '@/api/currency'

export const useCurrencyStore = defineStore('currency', () => {
  const currencies = ref<Currency[]>([])
  const defaultCurrency = ref<Currency | null>(null)

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

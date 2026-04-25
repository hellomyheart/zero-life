/**
 * 账户状态管理（Pinia Store）
 *
 * 组件功能：
 * - 缓存货币列表数据，避免重复请求
 * - 自动识别默认货币，供金额格式化等场景使用
 *
 * 数据流：
 * - currencies: 从后端API获取所有货币信息
 * - defaultCurrency: 从currencies中自动筛选is_default为true的货币
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Currency } from '@/types/currency'
import { list as listCurrencies } from '@/api/currency'

export const useCurrencyStore = defineStore('currency', () => {
  /** 货币列表数据，包含是否启用和是否默认等信息 */
  const currencies = ref<Currency[]>([])
  /** 默认货币，用于新建账户时的默认选项和系统展示 */
  const defaultCurrency = ref<Currency | null>(null)

  /**
   * 获取货币列表并找出默认货币
   * 从后端API获取所有货币数据，自动筛选出默认货币
   */
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

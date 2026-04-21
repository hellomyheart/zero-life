# 状态管理与 API 层

## 1. API 请求层

### Axios 实例配置

```typescript
// api/index.ts
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.accessToken) {
      config.headers.Authorization = `Bearer ${authStore.accessToken}`
    }

    // 注入当前用户组
    if (authStore.currentGroupID) {
      config.headers['X-User-Group-ID'] = authStore.currentGroupID
    }

    return config
  },
  (error) => Promise.reject(error)
)

// 响应拦截器
let isRefreshing = false
let pendingRequests: Array<(token: string) => void> = []

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config

    // 401: Token 过期，尝试刷新
    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve) => {
          pendingRequests.push((token: string) => {
            originalRequest.headers.Authorization = `Bearer ${token}`
            resolve(api(originalRequest))
          })
        })
      }

      originalRequest._retry = true
      isRefreshing = true

      try {
        const authStore = useAuthStore()
        const newToken = await authStore.refreshToken()

        pendingRequests.forEach((cb) => cb(newToken))
        pendingRequests = []

        originalRequest.headers.Authorization = `Bearer ${newToken}`
        return api(originalRequest)
      } catch (refreshError) {
        const authStore = useAuthStore()
        authStore.logout()
        router.push({ name: 'Login' })
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }

    return Promise.reject(error)
  }
)

export default api
```

### API 模块示例

```typescript
// api/accounts.ts
import api from './index'
import type { Account, AccountListParams, AccountForm, PaginatedResponse } from '@/types'

export const accountApi = {
  list(params: AccountListParams): Promise<PaginatedResponse<Account>> {
    return api.get('/accounts', { params })
  },

  get(id: number): Promise<{ data: Account }> {
    return api.get(`/accounts/${id}`)
  },

  create(data: AccountForm): Promise<{ data: Account }> {
    return api.post('/accounts', data)
  },

  update(id: number, data: Partial<AccountForm>): Promise<{ data: Account }> {
    return api.put(`/accounts/${id}`, data)
  },

  delete(id: number): Promise<void> {
    return api.delete(`/accounts/${id}`)
  },

  listTransactions(id: number, params?: any): Promise<PaginatedResponse<Transaction>> {
    return api.get(`/accounts/${id}/transactions`, { params })
  },

  listPiggyBanks(id: number): Promise<{ data: PiggyBank[] }> {
    return api.get(`/accounts/${id}/piggy-banks`)
  },
}
```

```typescript
// api/transactions.ts
import api from './index'
import type { TransactionGroup, TransactionForm, TransactionListParams, PaginatedResponse } from '@/types'

export const transactionApi = {
  list(params: TransactionListParams): Promise<PaginatedResponse<TransactionGroup>> {
    return api.get('/transactions', { params })
  },

  get(id: number): Promise<{ data: TransactionGroup }> {
    return api.get(`/transactions/${id}`)
  },

  create(data: TransactionForm): Promise<{ data: TransactionGroup }> {
    return api.post('/transactions', data)
  },

  update(id: number, data: TransactionForm): Promise<{ data: TransactionGroup }> {
    return api.put(`/transactions/${id}`, data)
  },

  delete(id: number): Promise<void> {
    return api.delete(`/transactions/${id}`)
  },
}
```

```typescript
// api/autocomplete.ts
import api from './index'

export const autocompleteApi = {
  accounts(query: string, type?: string): Promise<{ data: AutocompleteItem[] }> {
    return api.get('/autocomplete/accounts', { params: { query, type, limit: 10 } })
  },
  categories(query: string): Promise<{ data: AutocompleteItem[] }> {
    return api.get('/autocomplete/categories', { params: { query, limit: 10 } })
  },
  tags(query: string): Promise<{ data: AutocompleteItem[] }> {
    return api.get('/autocomplete/tags', { params: { query, limit: 10 } })
  },
  budgets(query: string): Promise<{ data: AutocompleteItem[] }> {
    return api.get('/autocomplete/budgets', { params: { query, limit: 10 } })
  },
  bills(query: string): Promise<{ data: AutocompleteItem[] }> {
    return api.get('/autocomplete/bills', { params: { query, limit: 10 } })
  },
  currencies(query: string): Promise<{ data: AutocompleteItem[] }> {
    return api.get('/autocomplete/currencies', { params: { query, limit: 10 } })
  },
  piggyBanks(query: string): Promise<{ data: AutocompleteItem[] }> {
    return api.get('/autocomplete/piggy-banks', { params: { query, limit: 10 } })
  },
  rules(query: string): Promise<{ data: AutocompleteItem[] }> {
    return api.get('/autocomplete/rules', { params: { query, limit: 10 } })
  },
  ruleGroups(query: string): Promise<{ data: AutocompleteItem[] }> {
    return api.get('/autocomplete/rule-groups', { params: { query, limit: 10 } })
  },
  recurrences(query: string): Promise<{ data: AutocompleteItem[] }> {
    return api.get('/autocomplete/recurrences', { params: { query, limit: 10 } })
  },
  webhooks(query: string): Promise<{ data: AutocompleteItem[] }> {
    return api.get('/autocomplete/webhooks', { params: { query, limit: 10 } })
  },
  transactionTypes(query: string): Promise<{ data: AutocompleteItem[] }> {
    return api.get('/autocomplete/transaction-types', { params: { query, limit: 10 } })
  },
}
```

## 2. Pinia 状态管理

### Auth Store

```typescript
// stores/auth.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/auth'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref<string | null>(localStorage.getItem('access_token'))
  const refreshToken = ref<string | null>(localStorage.getItem('refresh_token'))
  const user = ref<User | null>(null)
  const currentGroupID = ref<number | null>(null)
  const roles = ref<string[]>([])

  const isAuthenticated = computed(() => !!accessToken.value)

  async function login(email: string, password: string) {
    const { data } = await authApi.login({ email, password })
    accessToken.value = data.access_token
    refreshToken.value = data.refresh_token
    localStorage.setItem('access_token', data.access_token)
    localStorage.setItem('refresh_token', data.refresh_token)
    await fetchUser()
  }

  async function fetchUser() {
    const { data } = await authApi.me()
    user.value = data
    currentGroupID.value = data.default_group_id
    roles.value = data.roles || []
  }

  async function refreshAccessToken(): Promise<string> {
    const { data } = await authApi.refresh(refreshToken.value!)
    accessToken.value = data.access_token
    localStorage.setItem('access_token', data.access_token)
    return data.access_token
  }

  function hasRole(role: string): boolean {
    // OWNER 和 FULL 角色覆盖所有权限
    if (roles.value.includes('OWNER') || roles.value.includes('FULL')) {
      return true
    }
    return roles.value.includes(role)
  }

  function logout() {
    accessToken.value = null
    refreshToken.value = null
    user.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    router.push({ name: 'Login' })
  }

  return {
    accessToken, refreshToken, user, currentGroupID, roles,
    isAuthenticated,
    login, fetchUser, refreshAccessToken, hasRole, logout,
  }
})
```

### Account Store

```typescript
// stores/account.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { accountApi } from '@/api/accounts'

export const useAccountStore = defineStore('account', () => {
  const accounts = ref<Account[]>([])
  const assetAccounts = ref<Account[]>([])
  const expenseAccounts = ref<Account[]>([])
  const revenueAccounts = ref<Account[]>([])
  const loading = ref(false)

  async function fetchAccounts(params?: AccountListParams) {
    loading.value = true
    try {
      const { data, meta } = await accountApi.list(params)
      accounts.value = data
      return meta
    } finally {
      loading.value = false
    }
  }

  async function fetchAssetAccounts() {
    const { data } = await accountApi.list({ type: 'asset', per_page: 500 })
    assetAccounts.value = data
  }

  async function fetchExpenseAccounts() {
    const { data } = await accountApi.list({ type: 'expense', per_page: 500 })
    expenseAccounts.value = data
  }

  async function fetchRevenueAccounts() {
    const { data } = await accountApi.list({ type: 'revenue', per_page: 500 })
    revenueAccounts.value = data
  }

  return {
    accounts, assetAccounts, expenseAccounts, revenueAccounts, loading,
    fetchAccounts, fetchAssetAccounts, fetchExpenseAccounts, fetchRevenueAccounts,
  }
})
```

### Currency Store

```typescript
// stores/currency.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { currencyApi } from '@/api/currencies'

export const useCurrencyStore = defineStore('currency', () => {
  const currencies = ref<TransactionCurrency[]>([])
  const defaultCurrency = computed(() => currencies.value.find(c => c.is_default))

  async function fetchCurrencies() {
    const { data } = await currencyApi.list()
    currencies.value = data
  }

  function getCurrency(id: number): TransactionCurrency | undefined {
    return currencies.value.find(c => c.id === id)
  }

  function formatAmount(amount: string | number, currencyId?: number): string {
    const currency = currencyId ? getCurrency(currencyId) : defaultCurrency.value
    const num = new Decimal(amount)
    return `${currency?.symbol || ''} ${num.toFixed(currency?.decimal_places ?? 2)}`
  }

  return {
    currencies, defaultCurrency,
    fetchCurrencies, getCurrency, formatAmount,
  }
})
```

### UI Store

```typescript
// stores/ui.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUIStore = defineStore('ui', () => {
  const sidebarCollapsed = ref(false)
  const dateRange = ref<{ start: string; end: string }>({
    start: dayjs().startOf('month').format('YYYY-MM-DD'),
    end: dayjs().endOf('month').format('YYYY-MM-DD'),
  })
  const viewRange = ref<string>('1M') // 1D/1W/1M/3M/6M/1Y

  function setDateRange(start: string, end: string) {
    dateRange.value = { start, end }
  }

  function setViewRange(range: string) {
    viewRange.value = range
    // 根据视图范围计算日期
    const now = dayjs()
    switch (range) {
      case '1D':
        dateRange.value = { start: now.format('YYYY-MM-DD'), end: now.format('YYYY-MM-DD') }
        break
      case '1W':
        dateRange.value = {
          start: now.startOf('week').format('YYYY-MM-DD'),
          end: now.endOf('week').format('YYYY-MM-DD'),
        }
        break
      case '1M':
        dateRange.value = {
          start: now.startOf('month').format('YYYY-MM-DD'),
          end: now.endOf('month').format('YYYY-MM-DD'),
        }
        break
      // ... 其他范围
    }
  }

  return {
    sidebarCollapsed, dateRange, viewRange,
    setDateRange, setViewRange,
  }
})
```

## 3. TypeScript 类型定义

### 通用类型

```typescript
// types/common.ts
export interface PaginatedResponse<T> {
  data: T[]
  meta: {
    page: number
    per_page: number
    total: number
    total_pages: number
  }
  links?: {
    self: string
    first: string
    last: string
    next?: string
    prev?: string
  }
}

export interface ApiError {
  error: {
    code: string
    message: string
    details?: Array<{ field: string; message: string }>
  }
}

export interface AutocompleteItem {
  id: number
  name: string
  type?: string
}
```

### 账户类型

```typescript
// types/account.ts
export interface Account {
  id: number
  user_id: number
  user_group_id: number
  name: string
  account_type: AccountType
  account_role?: AccountRole
  currency_id: number
  object_group_id?: number
  virtual_balance: string
  current_balance?: string  // API 计算返回
  is_active: boolean
  is_virtual: boolean
  order: number
  created_at: string
  updated_at: string
  notes?: Note
  location?: Location
  meta?: AccountMeta[]
}

export type AccountType =
  | 'asset' | 'default' | 'cash' | 'credit_card' | 'saving' | 'shared_asset'
  | 'expense' | 'revenue'
  | 'debt' | 'loan' | 'mortgage'
  | 'initial_balance' | 'reconciliation' | 'import' | 'liability_credit'

export type AccountRole = 'defaultAsset' | 'sharedAsset' | 'savingAsset' | 'ccAsset' | 'cashWalletAsset'

export interface AccountForm {
  name: string
  type: AccountType
  account_role?: AccountRole
  currency_id: number
  virtual_balance?: string
  opening_balance?: string
  opening_balance_date?: string
  is_active?: boolean
  is_virtual?: boolean
  interest?: string
  interest_period?: string
  notes?: string
  object_group_id?: number
}
```

### 交易类型

```typescript
// types/transaction.ts
export interface TransactionGroup {
  id: number
  user_id: number
  user_group_id: number
  title?: string
  journals: TransactionJournal[]
  created_at: string
  updated_at: string
}

export interface TransactionJournal {
  id: number
  group_id: number
  transaction_type_id: TransactionType
  description: string
  date: string
  currency_id: number
  bill_id?: number
  category_id?: number
  budget_id?: number
  completed: boolean
  transactions: Transaction[]
  tags: Tag[]
  notes?: Note
}

export interface Transaction {
  id: number
  journal_id: number
  account_id: number
  description?: string
  amount: string
  native_amount: string
  foreign_amount?: string
  native_foreign_amount?: string
  currency_id: number
  foreign_currency_id?: number
  identifier: number
  reconciled: boolean
}

export type TransactionType = 'withdrawal' | 'deposit' | 'transfer' | 'opening_balance' | 'reconciliation' | 'liability_credit'

export interface TransactionForm {
  type: TransactionType
  date: string
  description: string
  title?: string
  transactions: TransactionItemForm[]
}

export interface TransactionItemForm {
  source_id: number
  destination_id: number
  amount: string
  currency_id: number
  foreign_amount?: string
  foreign_currency_id?: number
  category_id?: number
  budget_id?: number
  bill_id?: number
  tags?: string[]
  description?: string
  notes?: string
  reconciled?: boolean
}
```

## 4. 组合式函数 (Composables)

### usePagination

```typescript
// composables/usePagination.ts
export function usePagination<T>(fetchFn: (params: any) => Promise<PaginatedResponse<T>>) {
  const items = ref<T[]>([]) as Ref<T[]>
  const loading = ref(false)
  const page = ref(1)
  const perPage = ref(50)
  const total = ref(0)
  const totalPages = ref(0)

  async function fetch(params: Record<string, any> = {}) {
    loading.value = true
    try {
      const response = await fetchFn({
        page: page.value,
        per_page: perPage.value,
        ...params,
      })
      items.value = response.data
      total.value = response.meta.total
      totalPages.value = response.meta.total_pages
    } finally {
      loading.value = false
    }
  }

  function goToPage(p: number) {
    page.value = p
  }

  return { items, loading, page, perPage, total, totalPages, fetch, goToPage }
}
```

### useDateRange

```typescript
// composables/useDateRange.ts
export function useDateRange() {
  const uiStore = useUIStore()

  const start = computed(() => uiStore.dateRange.start)
  const end = computed(() => uiStore.dateRange.end)

  const presetRanges = [
    { label: '今日', value: '1D' },
    { label: '本周', value: '1W' },
    { label: '本月', value: '1M' },
    { label: '本季', value: '3M' },
    { label: '半年', value: '6M' },
    { label: '本年', value: '1Y' },
    { label: '最近7天', value: 'last7' },
    { label: '最近30天', value: 'last30' },
    { label: '最近90天', value: 'last90' },
    { label: '最近365天', value: 'last365' },
    { label: '月初至今', value: 'MTD' },
    { label: '季初至今', value: 'QTD' },
    { label: '年初至今', value: 'YTD' },
  ]

  function setRange(range: string) {
    uiStore.setViewRange(range)
  }

  function setCustomRange(s: string, e: string) {
    uiStore.setDateRange(s, e)
  }

  return { start, end, presetRanges, setRange, setCustomRange }
}
```

### useAutocomplete

```typescript
// composables/useAutocomplete.ts
export function useAutocomplete(type: 'accounts' | 'categories' | 'tags' | 'budgets' | 'bills' | 'currencies' | 'piggyBanks' | 'rules' | 'ruleGroups' | 'recurrences' | 'webhooks' | 'transactionTypes') {
  const items = ref<AutocompleteItem[]>([])
  const loading = ref(false)
  let debounceTimer: ReturnType<typeof setTimeout>

  async function search(query: string) {
    clearTimeout(debounceTimer)
    if (!query || query.length < 1) {
      items.value = []
      return
    }

    debounceTimer = setTimeout(async () => {
      loading.value = true
      try {
        const { data } = await autocompleteApi[type](query)
        items.value = data
      } finally {
        loading.value = false
      }
    }, 300)
  }

  return { items, loading, search }
}
```

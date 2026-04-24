// 账户状态管理 - 管理账户列表和当前选中账户
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Account } from '@/types/account'
import { list as listAccounts } from '@/api/account'

export const useAccountStore = defineStore('account', () => {
  const accounts = ref<Account[]>([]) // 账户列表
  const currentAccount = ref<Account | null>(null) // 当前选中的账户

  // fetchAccounts 获取账户列表
  async function fetchAccounts() {
    const res = await listAccounts({})
    accounts.value = (res as unknown as { items: Account[] }).items || (res as unknown as Account[])
  }

  // setCurrentAccount 设置当前选中账户
  function setCurrentAccount(account: Account | null) {
    currentAccount.value = account
  }

  return {
    accounts,
    currentAccount,
    fetchAccounts,
    setCurrentAccount,
  }
})

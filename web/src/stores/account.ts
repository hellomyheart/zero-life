import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Account } from '@/types/account'
import { list as listAccounts } from '@/api/account'

export const useAccountStore = defineStore('account', () => {
  const accounts = ref<Account[]>([])
  const currentAccount = ref<Account | null>(null)

  async function fetchAccounts() {
    const res = await listAccounts({})
    accounts.value = (res as unknown as { items: Account[] }).items || (res as unknown as Account[])
  }

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

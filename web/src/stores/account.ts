/**
 * 账户状态管理（Pinia Store）
 *
 * 组件功能：
 * - 缓存账户列表数据，避免重复请求
 * - 管理当前选中账户（用于表单等场景）
 *
 * 数据流：
 * - accounts: 从后端API获取，供各页面下拉选择使用
 * - currentAccount: 用于编辑等场景的当前选中账户
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Account } from '@/types/account'
import { list as listAccounts } from '@/api/account'

export const useAccountStore = defineStore('account', () => {
  /** 账户列表数据，供下拉选择等场景使用 */
  const accounts = ref<Account[]>([])
  /** 当前选中的账户，用于编辑等场景 */
  const currentAccount = ref<Account | null>(null)

  /**
   * 获取账户列表
   * 从后端API获取所有账户数据并缓存到store中
   */
  async function fetchAccounts() {
    const res = await listAccounts({})
    accounts.value = (res as unknown as { items: Account[] }).items || (res as unknown as Account[])
  }

  /**
   * 设置当前选中账户
   * @param account - 选中的账户对象，null表示取消选中
   */
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

<script setup lang="ts">
// 账户列表页面 - 按类型标签展示账户
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { list, remove } from '@/api/account'
import { formatAmount } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Account } from '@/types/account'
import { AccountType } from '@/types/account'

const { t } = useI18n()
const router = useRouter()

const accounts = ref<Account[]>([])
const loading = ref(false)
const activeTab = ref(AccountType.Asset)

const accountTypes = [
  { key: AccountType.Asset, label: t('account.asset') },
  { key: AccountType.Expense, label: t('account.expense') },
  { key: AccountType.Income, label: t('account.income') },
  { key: AccountType.Liability, label: t('account.liability') },
]

const filteredAccounts = ref<Account[]>([])

async function fetchAccounts() {
  loading.value = true
  try {
    const res = await list({})
    accounts.value = (res as unknown as { items: Account[] }).items || (res as unknown as Account[])
    filterAccounts()
  } catch {
    // handle error
  } finally {
    loading.value = false
  }
}

function filterAccounts() {
  filteredAccounts.value = accounts.value.filter((a) => a.type === activeTab.value)
}

function handleTabChange() {
  filterAccounts()
}

function handleCreate() {
  router.push('/accounts/create')
}

function handleEdit(id: string) {
  router.push(`/accounts/${id}/edit`)
}

async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm(t('account.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchAccounts()
  } catch {
    // cancelled or error
  }
}

onMounted(fetchAccounts)
</script>

<template>
  <div class="account-list-page">
    <div class="page-header">
      <h2>{{ t('account.title') }}</h2>
      <el-button type="primary" @click="handleCreate">{{ t('account.create') }}</el-button>
    </div>

    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane v-for="typeItem in accountTypes" :key="typeItem.key" :label="typeItem.label" :name="typeItem.key" />
    </el-tabs>

    <el-table :data="filteredAccounts" v-loading="loading" stripe>
      <el-table-column prop="name" :label="t('account.name')" />
      <el-table-column prop="balance" :label="t('account.balance')">
        <template #default="{ row }">{{ formatAmount(row.balance, row.currency) }}</template>
      </el-table-column>
      <el-table-column prop="currency" :label="t('account.currency')" width="100" />
      <el-table-column :label="t('account.isVirtual')" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.is_virtual" type="info" size="small">Virtual</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.edit')" width="160" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleEdit(row.id)">{{ t('common.edit') }}</el-button>
          <el-button link type="danger" @click="handleDelete(row.id)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.page-header h2 {
  margin: 0;
}
</style>

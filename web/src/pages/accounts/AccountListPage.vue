<script setup lang="ts">
/**
 * 账户列表页面
 * 功能：
 * - 按类型标签展示账户（资产/支出/收入/负债）
 * - 支持分页浏览
 * - 支持创建、编辑、删除账户
 * - 支持格式化金额显示
 */
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { list, remove } from '@/api/account'
import { formatAmount } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Account } from '@/types/account'
import { AccountType } from '@/types/account'
import Pagination from '@/components/common/Pagination.vue'

const { t } = useI18n()
const router = useRouter()

/** 所有账户数据 */
const accounts = ref<Account[]>([])
/** 加载状态 */
const loading = ref(false)
/** 当前选中的账户类型标签 */
const activeTab = ref(AccountType.Asset)

/** 分页参数 */
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0,
})

/** 账户类型选项卡配置 */
const accountTypes = [
  { key: AccountType.Asset, label: t('account.asset') },
  { key: AccountType.Expense, label: t('account.expense') },
  { key: AccountType.Income, label: t('account.income') },
  { key: AccountType.Liability, label: t('account.liability') },
]

/** 当前类型标签下的账户列表（前端过滤） */
const filteredAccounts = ref<Account[]>([])

/**
 * 获取账户列表
 * 从后端API获取所有账户数据，然后按当前类型标签过滤
 */
async function fetchAccounts() {
  loading.value = true
  try {
    const res = await list({ page: pagination.page, page_size: pagination.page_size })
    const data = res as unknown as { items: Account[]; total: number }
    accounts.value = data.items || []
    pagination.total = data.total || 0
    filterAccounts()
  } catch {
    ElMessage.error(t('common.fetchError') || 'Failed to load data')
  } finally {
    loading.value = false
  }
}

/**
 * 按当前类型标签过滤账户
 * 在前端根据activeTab过滤，因为后端返回所有类型
 */
function filterAccounts() {
  filteredAccounts.value = accounts.value.filter((a) => a.type === activeTab.value)
}

/** 切换类型标签时重新过滤 */
function handleTabChange() {
  filterAccounts()
}

/** 跳转到创建账户页面 */
function handleCreate() {
  router.push('/accounts/create')
}

/**
 * 跳转到编辑账户页面
 * @param id 账户ID
 */
function handleEdit(id: string) {
  router.push(`/accounts/${id}/edit`)
}

/**
 * 删除账户
 * 弹出确认框，确认后调用API删除
 * @param id 账户ID
 */
async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm(t('account.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchAccounts()
  } catch (err) {
    if ((err as string) !== 'cancel') ElMessage.error(t('common.fetchError') || 'Failed to load data')
  }
}

/** 页码变化时重新获取数据 */
function handlePageChange(page: number) {
  pagination.page = page
  fetchAccounts()
}

/** 每页数量变化时重置到第一页并重新获取 */
function handleSizeChange(size: number) {
  pagination.page_size = size
  pagination.page = 1
  fetchAccounts()
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
      <el-table-column prop="current_balance" :label="t('account.balance')">
        <template #default="{ row }">{{ formatAmount(row.current_balance, row.currency?.code) }}</template>
      </el-table-column>
      <el-table-column :label="t('account.currency')" width="100">
        <template #default="{ row }">{{ row.currency?.code || '' }}</template>
      </el-table-column>
      <el-table-column :label="t('account.isVirtual')" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.is_virtual" type="info" size="small">{{ t('account.isVirtual') }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.edit')" width="160" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleEdit(row.id)">{{ t('common.edit') }}</el-button>
          <el-button link type="danger" @click="handleDelete(row.id)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <Pagination
      :total="pagination.total"
      :page="pagination.page"
      :page-size="pagination.page_size"
      @update:page="handlePageChange"
      @update:page-size="handleSizeChange"
    />
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

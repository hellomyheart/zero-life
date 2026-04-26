<script setup lang="ts">
/**
 * Webhook列表页面
 * 字段名与后端 JSON tag 完全对应：
 * - is_active（非 active）
 * - trigger 值使用后端定义的事件名（如 transaction.created）
 */
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, create, update, remove } from '@/api/webhook'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Webhook, CreateWebhookReq, UpdateWebhookReq } from '@/types/webhook'
import { WebhookTrigger } from '@/types/webhook'
import Pagination from '@/components/common/Pagination.vue'

const { t } = useI18n()

const webhooks = ref<Webhook[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref<number | null>(null)
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

/** 触发事件选项 - 值与后端 oneof 验证完全对应 */
const triggerOptions = [
  { value: WebhookTrigger.TransactionCreated, label: t('webhook.transactionCreated') },
  { value: WebhookTrigger.TransactionUpdated, label: t('webhook.transactionUpdated') },
  { value: WebhookTrigger.TransactionDeleted, label: t('webhook.transactionDeleted') },
  { value: WebhookTrigger.BillPaid, label: t('webhook.billPaid') },
  { value: WebhookTrigger.BudgetCreated, label: t('webhook.budgetCreated') },
  { value: WebhookTrigger.BudgetUpdated, label: t('webhook.budgetUpdated') },
  { value: WebhookTrigger.BudgetDeleted, label: t('webhook.budgetDeleted') },
]

const form = ref<CreateWebhookReq & { is_active?: boolean }>({
  name: '',
  url: '',
  trigger: WebhookTrigger.TransactionCreated,
  is_active: true,
})

async function fetchWebhooks() {
  loading.value = true
  try {
    const res = await list()
    const data = res as unknown as Webhook[]
    webhooks.value = data
    pagination.total = data.length
  } catch {
    ElMessage.error(t('common.fetchError'))
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  pagination.page = page
  fetchWebhooks()
}

function handleSizeChange(size: number) {
  pagination.page_size = size
  pagination.page = 1
  fetchWebhooks()
}

function handleCreate() {
  dialogTitle.value = t('webhook.create')
  editingId.value = null
  form.value = { name: '', url: '', trigger: WebhookTrigger.TransactionCreated, is_active: true }
  dialogVisible.value = true
}

function handleEdit(row: Webhook) {
  dialogTitle.value = t('webhook.edit')
  editingId.value = row.id
  form.value = { name: row.name, url: row.url, trigger: row.trigger, is_active: row.is_active }
  dialogVisible.value = true
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm(t('webhook.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchWebhooks()
  } catch (err) {
    if ((err as string) !== 'cancel') ElMessage.error(t('common.failed'))
  }
}

async function handleSubmit() {
  if (!form.value.name || !form.value.url) {
    ElMessage.warning(t('common.required'))
    return
  }
  try {
    if (editingId.value) {
      const updateData: UpdateWebhookReq = {
        name: form.value.name,
        url: form.value.url,
        trigger: form.value.trigger,
        is_active: form.value.is_active,
      }
      await update(editingId.value, updateData)
    } else {
      await create(form.value)
    }
    ElMessage.success(t('common.success'))
    dialogVisible.value = false
    await fetchWebhooks()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

onMounted(fetchWebhooks)
</script>

<template>
  <div class="webhook-list-page">
    <div class="page-header">
      <h2>{{ t('webhook.title') }}</h2>
      <el-button type="primary" @click="handleCreate">{{ t('webhook.create') }}</el-button>
    </div>

    <el-table :data="webhooks" v-loading="loading" stripe>
      <el-table-column prop="name" :label="t('webhook.name')" />
      <el-table-column prop="url" :label="t('webhook.url')" min-width="200" />
      <el-table-column prop="trigger" :label="t('webhook.trigger')" width="180">
        <template #default="{ row }">
          {{ triggerOptions.find(o => o.value === row.trigger)?.label || row.trigger }}
        </template>
      </el-table-column>
      <!-- is_active 字段名与后端 WebhookResp JSON tag 对应（非 active） -->
      <el-table-column :label="t('webhook.active')" width="80">
        <template #default="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'info'" size="small">{{ row.is_active ? t('rule.enabled') : t('rule.disabled') }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="last_delivered_at" :label="t('webhook.lastDelivered')" width="160">
        <template #default="{ row }">{{ row.last_delivered_at ? formatDate(row.last_delivered_at) : '-' }}</template>
      </el-table-column>
      <el-table-column :label="t('common.edit')" width="160" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item :label="t('webhook.name')">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="t('webhook.url')">
          <el-input v-model="form.url" />
        </el-form-item>
        <el-form-item :label="t('webhook.trigger')">
          <el-select v-model="form.trigger" style="width: 100%">
            <el-option v-for="opt in triggerOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('webhook.active')">
          <el-switch v-model="form.is_active" />
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
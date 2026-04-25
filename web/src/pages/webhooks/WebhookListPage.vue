<script setup lang="ts">
// Webhook列表页面 - 配置事件通知发送到URL和触发条件
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, create, update, remove } from '@/api/webhook'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Webhook, CreateWebhookReq, UpdateWebhookReq } from '@/types/webhook'
import Pagination from '@/components/common/Pagination.vue'

const { t } = useI18n()

/** Webhook列表数据 */
const webhooks = ref<Webhook[]>([])
/** 加载状态 */
const loading = ref(false)
/** 对话框显示状态 */
const dialogVisible = ref(false)
/** 对话框标题 */
const dialogTitle = ref('')
/** 当前编辑的Webhook ID，null表示新建 */
const editingId = ref<string | null>(null)
/** 分页参数 - page: 当前页码, page_size: 每页数量, total: 总记录数 */
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

const triggerOptions = [
  { value: 'STORE_TRANSACTION', label: 'Store Transaction' },
  { value: 'UPDATE_TRANSACTION', label: 'Update Transaction' },
  { value: 'DESTROY_TRANSACTION', label: 'Destroy Transaction' },
]

const form = ref<CreateWebhookReq>({ name: '', url: '', trigger: 'STORE_TRANSACTION' })

/**
 * 获取Webhook列表
 * 传入分页参数，从后端获取当前页的数据和总记录数
 */
async function fetchWebhooks() {
  loading.value = true
  try {
    const res = await list({ page: pagination.page, page_size: pagination.page_size }) as unknown as { items: Webhook[], total: number }
    webhooks.value = res.items || []
    pagination.total = res.total || 0
  } catch {
    // handle error
  } finally {
    loading.value = false
  }
}

/**
 * 页码变化处理函数
 * @param page 新的页码
 */
function handlePageChange(page: number) {
  pagination.page = page
  fetchWebhooks()
}

/**
 * 每页数量变化处理函数
 * @param size 新的每页数量
 */
function handleSizeChange(size: number) {
  pagination.page_size = size
  pagination.page = 1
  fetchWebhooks()
}

function handleCreate() {
  dialogTitle.value = t('webhook.create')
  editingId.value = null
  form.value = { name: '', url: '', trigger: 'STORE_TRANSACTION' }
  dialogVisible.value = true
}

function handleEdit(row: Webhook) {
  dialogTitle.value = t('webhook.edit')
  editingId.value = row.id
  form.value = { name: row.name, url: row.url, trigger: row.trigger }
  dialogVisible.value = true
}

async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm(t('webhook.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchWebhooks()
  } catch {
    // cancelled or error
  }
}

async function handleSubmit() {
  if (!form.value.name || !form.value.url) {
    ElMessage.warning(t('common.required'))
    return
  }
  try {
    if (editingId.value) {
      await update(editingId.value, form.value as UpdateWebhookReq)
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
      <el-table-column prop="trigger" :label="t('webhook.trigger')" width="180" />
      <el-table-column :label="t('webhook.active')" width="80">
        <template #default="{ row }">
          <el-tag :type="row.active ? 'success' : 'info'" size="small">{{ row.active ? t('rule.enabled') : t('rule.disabled') }}</el-tag>
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
          <el-select v-model="form.trigger">
            <el-option v-for="opt in triggerOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
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

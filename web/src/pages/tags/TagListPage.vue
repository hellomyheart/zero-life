<script setup lang="ts">
/**
 * 标签列表页面
 * 功能：
 * - 展示标签列表，支持分页
 * - 显示标签名称、颜色和关联交易数量
 * - 支持创建、编辑、删除标签
 * - 颜色选择器支持自定义颜色
 */
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, create, update, remove } from '@/api/tag'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Tag, CreateTagReq } from '@/types/tag'
import Pagination from '@/components/common/Pagination.vue'

const { t } = useI18n()

/** 标签列表数据 */
const tags = ref<Tag[]>([])
/** 加载状态 */
const loading = ref(false)
/** 对话框显示状态 */
const dialogVisible = ref(false)
/** 对话框标题 */
const dialogTitle = ref('')
/** 当前编辑的标签ID，null表示新建模式 */
const editingId = ref<string | null>(null)

/** 分页参数 */
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0,
})

/** 表单数据 - 创建/编辑标签时使用 */
const form = ref<CreateTagReq>({
  name: '',
  color: '#409EFF',
})

/**
 * 获取标签列表
 * 从后端API获取标签数据
 */
async function fetchTags() {
  loading.value = true
  try {
    const res = await list()
    const data = res as unknown as { items: Tag[]; total: number }
    tags.value = data.items || (res as unknown as Tag[])
    pagination.total = data.total || 0
  } catch {
    ElMessage.error(t('common.fetchError') || 'Failed to load data')
  } finally {
    loading.value = false
  }
}

/** 打开创建标签对话框 */
function handleCreate() {
  dialogTitle.value = t('tag.create')
  editingId.value = null
  form.value = { name: '', color: '#409EFF' }
  dialogVisible.value = true
}

/**
 * 打开编辑标签对话框
 * @param tag 要编辑的标签数据
 */
function handleEdit(tag: Tag) {
  dialogTitle.value = t('tag.edit')
  editingId.value = tag.id
  form.value = { name: tag.name, color: tag.color }
  dialogVisible.value = true
}

/**
 * 删除标签
 * 弹出确认框，确认后调用API删除
 * @param id 标签ID
 */
async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm(t('tag.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchTags()
  } catch (err) {
    if ((err as string) !== 'cancel') ElMessage.error(t('common.fetchError') || 'Failed to load data')
  }
}

/**
 * 提交表单
 * 根据editingId判断是创建还是编辑操作
 */
async function handleSubmit() {
  if (!form.value.name) {
    ElMessage.warning(t('common.required'))
    return
  }
  try {
    if (editingId.value) {
      await update(editingId.value, form.value)
    } else {
      await create(form.value)
    }
    ElMessage.success(t('common.success'))
    dialogVisible.value = false
    await fetchTags()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

/** 页码变化 */
function handlePageChange(page: number) {
  pagination.page = page
  fetchTags()
}

/** 每页数量变化 */
function handleSizeChange(size: number) {
  pagination.page_size = size
  pagination.page = 1
  fetchTags()
}

onMounted(fetchTags)
</script>

<template>
  <div class="tag-list-page">
    <div class="page-header">
      <h2>{{ t('tag.title') }}</h2>
      <el-button type="primary" @click="handleCreate">{{ t('tag.create') }}</el-button>
    </div>

    <el-table :data="tags" v-loading="loading" stripe>
      <el-table-column prop="name" :label="t('tag.name')" />
      <el-table-column :label="t('tag.color')" width="100">
        <template #default="{ row }">
          <el-color-picker v-model="row.color" disabled size="small" />
        </template>
      </el-table-column>
      <el-table-column prop="transaction_count" :label="t('tag.transactionCount')" width="150" />
      <el-table-column :label="t('common.edit')" width="160">
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="400px">
      <el-form :model="form" label-width="80px">
        <el-form-item :label="t('tag.name')">
          <el-input v-model="form.name" :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('tag.color')">
          <el-color-picker v-model="form.color" />
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

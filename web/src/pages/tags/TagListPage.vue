<script setup lang="ts">
/**
 * 标签列表页面
 * 功能：
 * - 展示标签树形列表（支持两级结构）
 * - 显示标签名称、颜色和关联交易数量
 * - 支持创建、编辑、删除标签
 * - 颜色选择器支持预设调色板 + 自定义颜色
 * - 创建/编辑时可选择父标签
 */
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, create, update, remove } from '@/api/tag'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Tag, CreateTagReq, UpdateTagReq } from '@/types/tag'
import ColorPicker from '@/components/common/ColorPicker.vue'

const { t } = useI18n()

/** 标签树形列表数据 */
const tags = ref<Tag[]>([])
/** 加载状态 */
const loading = ref(false)
/** 对话框显示状态 */
const dialogVisible = ref(false)
/** 对话框标题 */
const dialogTitle = ref('')
/** 当前编辑的标签ID，null表示新建模式 */
const editingId = ref<number | null>(null)

/** 表单数据 - 创建/编辑标签时使用 */
const form = ref<CreateTagReq & { parent_id: number | null }>({
  name: '',
  color: '#409EFF',
  parent_id: null,
})

/**
 * 将树形标签数据扁平化，用于父标签选择器
 * 只取顶级标签作为可选父标签（两级限制）
 */
const topLevelTags = computed(() => {
  return flattenTree(tags.value).filter(t => t.parent_id === null)
})

/**
 * 将树形结构扁平化为列表
 * @param tree 树形标签数据
 */
function flattenTree(tree: Tag[]): Tag[] {
  const result: Tag[] = []
  for (const item of tree) {
    result.push(item)
    if (item.children && item.children.length > 0) {
      result.push(...flattenTree(item.children))
    }
  }
  return result
}

/**
 * 获取标签列表
 * 从后端API获取标签数据（树形结构）
 */
async function fetchTags() {
  loading.value = true
  try {
    const res = await list()
    tags.value = (res as unknown as Tag[]) || []
  } catch {
    ElMessage.error(t('common.fetchError') || 'Failed to load data')
  } finally {
    loading.value = false
  }
}

/** 打开创建标签对话框 */
function handleCreate(parentId: number | null = null) {
  dialogTitle.value = t('tag.create')
  editingId.value = null
  form.value = { name: '', color: '#409EFF', parent_id: parentId }
  dialogVisible.value = true
}

/**
 * 打开编辑标签对话框
 * @param tag 要编辑的标签数据
 */
function handleEdit(tag: Tag) {
  dialogTitle.value = t('tag.edit')
  editingId.value = tag.id
  form.value = { name: tag.name, color: tag.color, parent_id: tag.parent_id }
  dialogVisible.value = true
}

/**
 * 删除标签
 * 弹出确认框，确认后调用API删除
 * 删除父标签将同时删除所有子标签
 * @param tag 要删除的标签
 */
async function handleDelete(tag: Tag) {
  const hasChildren = tag.children && tag.children.length > 0
  const confirmMsg = hasChildren ? t('tag.deleteWithChildren') : t('tag.deleteConfirm')
  try {
    await ElMessageBox.confirm(confirmMsg, t('common.confirm'), { type: 'warning' })
    await remove(tag.id)
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
      const data: UpdateTagReq = {
        name: form.value.name,
        color: form.value.color,
        parent_id: form.value.parent_id,
      }
      await update(editingId.value, data)
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

onMounted(fetchTags)
</script>

<template>
  <div class="tag-list-page">
    <div class="page-header">
      <h2>{{ t('tag.title') }}</h2>
      <el-button type="primary" @click="handleCreate()">{{ t('tag.create') }}</el-button>
    </div>

    <el-table
      :data="tags"
      v-loading="loading"
      stripe
      row-key="id"
      default-expand-all
      :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
    >
      <el-table-column prop="name" :label="t('tag.name')" />
      <el-table-column :label="t('tag.color')" width="100">
        <template #default="{ row }">
          <el-color-picker v-model="row.color" disabled size="small" />
        </template>
      </el-table-column>
      <el-table-column prop="transaction_count" :label="t('tag.transactionCount')" width="150" />
      <el-table-column :label="t('common.edit')" width="200">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleCreate(row.id)">{{ t('tag.create') }}</el-button>
          <el-button link type="primary" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
          <el-button link type="danger" @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="480px">
      <el-form :model="form" label-width="100px">
        <el-form-item :label="t('tag.name')">
          <el-input v-model="form.name" :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('tag.parentTag')">
          <el-select
            v-model="form.parent_id"
            clearable
            :placeholder="t('tag.selectParent')"
            style="width: 100%"
          >
            <el-option :label="t('tag.topLevel')" :value="null" />
            <el-option
              v-for="tag in topLevelTags"
              :key="tag.id"
              :label="tag.name"
              :value="tag.id"
              :disabled="tag.id === editingId"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('tag.color')">
          <ColorPicker v-model="form.color" />
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

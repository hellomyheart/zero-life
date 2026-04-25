<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, create, update, remove } from '@/api/category'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ElIcon } from 'element-plus'
import * as Icons from '@element-plus/icons-vue'
import type { Category, CreateCategoryReq } from '@/types/category'
import IconSelect from '@/components/common/IconSelect.vue'

const { t } = useI18n()

const categories = ref<Category[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref<number | null>(null)

const form = ref<CreateCategoryReq>({
  name: '',
  parent_id: null,
  icon: '',
  notes: '',
})

const treeProps = {
  children: 'children',
  label: 'name',
}

async function fetchCategories() {
  loading.value = true
  try {
    categories.value = await list() as unknown as Category[]
  } catch {
    ElMessage.error(t('common.fetchError') || 'Failed to load data')
  } finally {
    loading.value = false
  }
}

function handleCreate(parentId: number | null = null) {
  dialogTitle.value = t('category.create')
  editingId.value = null
  form.value = { name: '', parent_id: parentId, icon: '', notes: '' }
  dialogVisible.value = true
}

function handleEdit(data: Category) {
  dialogTitle.value = t('category.edit')
  editingId.value = data.id
  form.value = { name: data.name, parent_id: data.parent_id, icon: data.icon, notes: data.notes }
  dialogVisible.value = true
}

async function handleDelete(data: Category) {
  try {
    await ElMessageBox.confirm(t('category.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(data.id)
    ElMessage.success(t('common.success'))
    await fetchCategories()
  } catch (err) {
    if ((err as string) !== 'cancel') ElMessage.error(t('common.fetchError') || 'Failed to load data')
  }
}

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
    await fetchCategories()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

onMounted(fetchCategories)
</script>

<template>
  <div class="category-list-page">
    <div class="page-header">
      <h2>{{ t('category.title') }}</h2>
      <el-button type="primary" @click="handleCreate()">{{ t('category.create') }}</el-button>
    </div>

    <el-tree
      :data="categories"
      :props="treeProps"
      node-key="id"
      default-expand-all
      v-loading="loading"
    >
      <template #default="{ data }">
        <div class="tree-node">
          <span class="node-label">
            <el-icon v-if="data.icon && (Icons as any)[data.icon]" :size="16">
              <component :is="(Icons as any)[data.icon]" />
            </el-icon>
            <span>{{ data.name }}</span>
            <span v-if="data.notes" class="node-notes">{{ data.notes }}</span>
          </span>
          <span class="tree-actions">
            <el-button link type="primary" size="small" @click="handleCreate(data.id)">{{ t('category.create') }}</el-button>
            <el-button link type="primary" size="small" @click="handleEdit(data)">{{ t('common.edit') }}</el-button>
            <el-button link type="danger" size="small" @click="handleDelete(data)">{{ t('common.delete') }}</el-button>
          </span>
        </div>
      </template>
    </el-tree>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="400px">
      <el-form :model="form" label-width="80px">
        <el-form-item :label="t('category.name')">
          <el-input v-model="form.name" :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('category.icon')">
          <IconSelect :model-value="form.icon ?? ''" @update:model-value="form.icon = $event" />
        </el-form-item>
        <el-form-item :label="t('category.notes')">
          <el-input v-model="form.notes" type="textarea" :rows="2" :placeholder="t('category.notesPlaceholder')" />
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

.tree-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex: 1;
  padding-right: 8px;
}

.tree-actions {
  display: none;
}

.tree-node:hover .tree-actions {
  display: inline-flex;
  gap: 4px;
}

.node-label {
  display: flex;
  align-items: center;
  gap: 6px;
}

.node-notes {
  color: #999;
  font-size: 12px;
  margin-left: 4px;
}
</style>

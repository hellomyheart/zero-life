<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, create, update, remove } from '@/api/tag'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Tag, CreateTagReq } from '@/types/tag'

const { t } = useI18n()

const tags = ref<Tag[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref<string | null>(null)

const form = ref<CreateTagReq>({
  name: '',
  color: '#409EFF',
})

async function fetchTags() {
  loading.value = true
  try {
    tags.value = await list() as unknown as Tag[]
  } catch {
    // handle error
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  dialogTitle.value = t('tag.create')
  editingId.value = null
  form.value = { name: '', color: '#409EFF' }
  dialogVisible.value = true
}

function handleEdit(tag: Tag) {
  dialogTitle.value = t('tag.edit')
  editingId.value = tag.id
  form.value = { name: tag.name, color: tag.color }
  dialogVisible.value = true
}

async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm(t('tag.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchTags()
  } catch {
    // cancelled or error
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

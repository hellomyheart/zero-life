<script setup lang="ts">
// 附件管理页面 - 上传下载和删除附件文件
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, upload, remove } from '@/api/attachment'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Attachment } from '@/types/attachment'

const { t } = useI18n()

const attachments = ref<Attachment[]>([])
const loading = ref(false)
const uploadVisible = ref(false)
const fileList = ref<File[]>([])

async function fetchAttachments() {
  loading.value = true
  try {
    const res = await list({}) as unknown as { items: Attachment[] }
    attachments.value = res.items || []
  } catch {
    // handle error
  } finally {
    loading.value = false
  }
}

async function handleUpload() {
  if (!fileList.value.length) {
    ElMessage.warning(t('common.required'))
    return
  }
  try {
    const formData = new FormData()
    formData.append('file', fileList.value[0])
    await upload(formData)
    ElMessage.success(t('common.success'))
    uploadVisible.value = false
    fileList.value = []
    await fetchAttachments()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm(t('attachment.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchAttachments()
  } catch {
    // cancelled or error
  }
}

function handleFileChange(file: File) {
  fileList.value = [file]
}

onMounted(fetchAttachments)
</script>

<template>
  <div class="attachment-list-page">
    <div class="page-header">
      <h2>{{ t('attachment.title') }}</h2>
      <el-button type="primary" @click="uploadVisible = true">{{ t('attachment.upload') }}</el-button>
    </div>

    <el-table :data="attachments" v-loading="loading" stripe>
      <el-table-column prop="filename" :label="t('attachment.filename')" />
      <el-table-column prop="mime" :label="t('attachment.mime')" width="150" />
      <el-table-column prop="size" :label="t('attachment.size')" width="100">
        <template #default="{ row }">{{ (row.size / 1024).toFixed(1) }} KB</template>
      </el-table-column>
      <el-table-column prop="created_at" :label="t('attachment.createdAt')" width="160">
        <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
      </el-table-column>
      <el-table-column :label="t('common.edit')" width="200" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="">{{ t('attachment.download') }}</el-button>
          <el-button link type="danger" @click="handleDelete(row.id)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="uploadVisible" :title="t('attachment.upload')" width="400px">
      <el-upload :auto-upload="false" :limit="1" :on-change="handleFileChange as any">
        <el-button type="primary">{{ t('attachment.selectFile') }}</el-button>
      </el-upload>
      <template #footer>
        <el-button @click="uploadVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleUpload">{{ t('attachment.upload') }}</el-button>
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

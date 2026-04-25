<script setup lang="ts">
// 附件管理页面 - 上传下载和删除附件文件
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, upload, download, remove } from '@/api/attachment'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Attachment } from '@/types/attachment'
import Pagination from '@/components/common/Pagination.vue'

const { t } = useI18n()

/** 附件列表数据 */
const attachments = ref<Attachment[]>([])
/** 加载状态 */
const loading = ref(false)
/** 上传对话框显示状态 */
const uploadVisible = ref(false)
/** 上传文件列表 */
const fileList = ref<File[]>([])
/** 分页参数 - page: 当前页码, page_size: 每页数量, total: 总记录数 */
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

/**
 * 获取附件列表
 * 传入分页参数，从后端获取当前页的数据和总记录数
 */
async function fetchAttachments() {
  loading.value = true
  try {
    const res = await list({ page: pagination.page, page_size: pagination.page_size }) as unknown as { items: Attachment[], total: number }
    attachments.value = res.items || []
    pagination.total = res.total || 0
  } catch {
    ElMessage.error(t('common.fetchError') || 'Failed to load data')
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
  fetchAttachments()
}

/**
 * 每页数量变化处理函数
 * @param size 新的每页数量
 */
function handleSizeChange(size: number) {
  pagination.page_size = size
  pagination.page = 1
  fetchAttachments()
}

/**
 * 上传附件
 * 将选中的文件通过FormData上传到后端，上传成功后关闭对话框并刷新列表
 */
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

/**
 * 删除附件
 * 弹出确认框后调用API删除指定附件
 * @param id 附件ID
 */
async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm(t('attachment.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchAttachments()
  } catch (err) {
    if ((err as string) !== 'cancel') ElMessage.error(t('common.fetchError') || 'Failed to load data')
  }
}

// 下载附件 - 调用API获取文件并触发浏览器下载
async function handleDownload(row: Attachment) {
  try {
    const res = await download(row.id) as unknown as Blob
    const blob = res instanceof Blob ? res : new Blob([res as BlobPart])
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = row.filename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

/**
 * 文件选择变化处理
 * 只保留最新选择的文件（单文件上传）
 * @param file 用户选择的文件
 */
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
        <!-- 文件大小从字节转换为KB显示，保留1位小数 -->
        <template #default="{ row }">{{ (row.size / 1024).toFixed(1) }} KB</template>
      </el-table-column>
      <el-table-column prop="created_at" :label="t('attachment.createdAt')" width="160">
        <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
      </el-table-column>
      <el-table-column :label="t('common.edit')" width="200" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleDownload(row)">{{ t('attachment.download') }}</el-button>
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

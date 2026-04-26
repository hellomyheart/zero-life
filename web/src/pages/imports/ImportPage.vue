<script setup lang="ts">
// 数据导入页面 - 4步向导上传映射预览确认导入
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { upload, parse, execute } from '@/api/import'
import { ElMessage } from 'element-plus'
import type { ImportPreviewResp } from '@/types/import'

const { t } = useI18n()

const step = ref(1)
const fileId = ref('')
const loading = ref(false)
const parsedData = ref<Record<string, unknown>[]>([])
const columnMapping = reactive<Record<string, string>>({})
const availableColumns = ref<string[]>([])
const targetFields = ['date', 'description', 'amount', 'source_account', 'category', 'tags']

async function handleUpload(options: { file: File }) {
  loading.value = true
  try {
    const res = await upload(options.file) as { file_id: string }
    fileId.value = res.file_id
    step.value = 2
    ElMessage.success(t('common.success'))
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  } finally {
    loading.value = false
  }
}

async function handleParse() {
  loading.value = true
  try {
    const res = await parse({ file_id: fileId.value, mapping: columnMapping }) as unknown as ImportPreviewResp
    parsedData.value = res.rows?.map(r => r.data) || []
    step.value = 3
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  } finally {
    loading.value = false
  }
}

async function handleExecute() {
  loading.value = true
  try {
    await execute({ file_id: fileId.value, mapping: columnMapping })
    ElMessage.success(t('common.success'))
    step.value = 4
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="import-page">
    <h2>{{ t('import.title') }}</h2>

    <el-steps :active="step - 1" finish-status="success" style="margin-bottom: 24px">
      <el-step :title="t('import.upload')" />
      <el-step :title="t('import.columnMapping')" />
      <el-step :title="t('import.preview')" />
      <el-step :title="t('import.execute')" />
    </el-steps>

    <el-card v-if="step === 1">
      <el-upload
        drag
        accept=".csv"
        :http-request="handleUpload"
        :show-file-list="false"
      >
        <el-icon class="el-icon--upload"><Upload /></el-icon>
        <div class="el-upload__text">{{ t('import.fileType') }}</div>
      </el-upload>
    </el-card>

    <el-card v-if="step === 2">
      <h3>{{ t('import.columnMapping') }}</h3>
      <el-form label-width="150px">
        <el-form-item v-for="field in targetFields" :key="field" :label="field">
          <el-select v-model="columnMapping[field]" :placeholder="t('common.selectPlaceholder')" clearable>
            <el-option v-for="col in availableColumns" :key="col" :label="col" :value="col" />
          </el-select>
        </el-form-item>
      </el-form>
      <el-button type="primary" @click="handleParse" :loading="loading">{{ t('import.parse') }}</el-button>
      <el-button @click="step = 1">{{ t('common.previous') }}</el-button>
    </el-card>

    <el-card v-if="step === 3">
      <h3>{{ t('import.preview') }}</h3>
      <el-table :data="parsedData.slice(0, 20)" stripe max-height="400">
        <el-table-column v-for="col in availableColumns" :key="col" :prop="col" :label="col" />
      </el-table>
      <div style="margin-top: 16px">
        <el-button type="primary" @click="handleExecute" :loading="loading">{{ t('import.execute') }}</el-button>
        <el-button @click="step = 2">{{ t('common.previous') }}</el-button>
      </div>
    </el-card>

    <el-card v-if="step === 4">
      <el-result icon="success" :title="t('import.completed')" />
      <div style="text-align: center; margin-top: 16px">
        <el-button @click="step = 1; fileId = ''; parsedData = []; Object.keys(columnMapping).forEach(k => delete columnMapping[k])">{{ t('import.uploadAnother') }}</el-button>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts">
import { Upload } from '@element-plus/icons-vue'
export default {
  components: { Upload },
}
</script>

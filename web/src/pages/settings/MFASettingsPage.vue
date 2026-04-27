<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { status, setup, enable, disable, regenerateBackupCodes } from '@/api/mfa'
import { ElMessage } from 'element-plus'
import type { MFAStatus, MFASetupResp, BackupCodesResp } from '@/types/mfa'

const { t } = useI18n()

const mfaStatus = ref<MFAStatus>({ enabled: false })
const qrCodeUrl = ref('')
const secret = ref('')
const code = ref('')
const loading = ref(false)
const backupCodes = ref<string[]>([])
const showBackupCodes = ref(false)

async function fetchStatus() {
  try {
    mfaStatus.value = await status() as unknown as MFAStatus
  } catch {
    ElMessage.error(t('common.fetchError') || 'Failed to load data')
  }
}

async function handleSetup() {
  loading.value = true
  try {
    const res = await setup() as unknown as MFASetupResp
    qrCodeUrl.value = res.qr_code
    secret.value = res.secret
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  } finally {
    loading.value = false
  }
}

async function handleEnable() {
  if (!code.value) {
    ElMessage.warning(t('common.required'))
    return
  }
  try {
    const res = await enable({ code: code.value }) as unknown as BackupCodesResp
    ElMessage.success(t('common.success'))
    code.value = ''
    qrCodeUrl.value = ''
    backupCodes.value = res.codes
    showBackupCodes.value = true
    await fetchStatus()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

async function handleDisable() {
  if (!code.value) {
    ElMessage.warning(t('common.required'))
    return
  }
  try {
    await disable({ code: code.value })
    ElMessage.success(t('common.success'))
    code.value = ''
    backupCodes.value = []
    showBackupCodes.value = false
    await fetchStatus()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

async function handleRegenerateBackupCodes() {
  loading.value = true
  try {
    const res = await regenerateBackupCodes() as unknown as BackupCodesResp
    backupCodes.value = res.codes
    showBackupCodes.value = true
    ElMessage.success(t('common.success'))
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  } finally {
    loading.value = false
  }
}

onMounted(fetchStatus)
</script>

<template>
  <div class="mfa-settings-page">
    <h2>{{ t('mfa.title') }}</h2>

    <el-card style="margin-top: 16px;">
      <el-descriptions :column="1">
        <el-descriptions-item :label="t('mfa.status')">
          <el-tag :type="mfaStatus.enabled ? 'success' : 'info'">{{ mfaStatus.enabled ? t('mfa.enabled') : t('mfa.disabled') }}</el-tag>
        </el-descriptions-item>
      </el-descriptions>

      <div style="margin-top: 20px;">
        <template v-if="!mfaStatus.enabled">
          <el-button type="primary" @click="handleSetup" :loading="loading">{{ t('mfa.setup') }}</el-button>
          <div v-if="qrCodeUrl" style="margin-top: 16px;">
            <p>{{ t('mfa.scanQR') }}</p>
            <img
              :src="`https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(qrCodeUrl)}`"
              alt="QR Code"
              style="margin: 8px 0;"
            />
            <p>{{ t('mfa.secret') }}: <code>{{ secret }}</code></p>
            <el-input v-model="code" :placeholder="t('mfa.enterCode')" style="width: 200px; margin-top: 8px;" />
            <el-button type="success" @click="handleEnable" style="margin-left: 8px;">{{ t('mfa.enable') }}</el-button>
          </div>
        </template>
        <template v-else>
          <el-input v-model="code" :placeholder="t('mfa.enterCode')" style="width: 200px;" />
          <el-button type="danger" @click="handleDisable" style="margin-left: 8px;">{{ t('mfa.disable') }}</el-button>

          <el-divider />

          <el-button type="warning" @click="handleRegenerateBackupCodes" :loading="loading">{{ t('mfa.regenerateCodes') }}</el-button>
        </template>
      </div>
    </el-card>

    <el-card v-if="showBackupCodes && backupCodes.length" style="margin-top: 16px;">
      <template #header>
        <span>{{ t('mfa.backupCodes') }}</span>
      </template>
      <el-alert :title="t('mfa.backupCodesWarning')" type="warning" :closable="false" show-icon style="margin-bottom: 16px;" />
      <div class="backup-codes-grid">
        <code v-for="c in backupCodes" :key="c">{{ c }}</code>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
h2 {
  margin: 0 0 16px 0;
}

.backup-codes-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

.backup-codes-grid code {
  display: block;
  padding: 6px 12px;
  background: var(--el-fill-color-light);
  border-radius: 4px;
  text-align: center;
  font-size: 16px;
  letter-spacing: 1px;
}
</style>

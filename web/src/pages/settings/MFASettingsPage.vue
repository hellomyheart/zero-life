<script setup lang="ts">
// 多因素认证页面 - 设置启用和禁用TOTP认证码验证
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { status, setup, enable, disable } from '@/api/mfa'
import { ElMessage } from 'element-plus'
import type { MFAStatus, MFASetupResp } from '@/types/mfa'

const { t } = useI18n()

const mfaStatus = ref<MFAStatus>({ enabled: false })
const qrCodeUrl = ref('')
const secret = ref('')
const code = ref('')
const loading = ref(false)

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
    await enable({ code: code.value })
    ElMessage.success(t('common.success'))
    code.value = ''
    qrCodeUrl.value = ''
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
    await fetchStatus()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
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
          <el-tag :type="mfaStatus.enabled ? 'success' : 'info'">{{ mfaStatus.enabled ? t('rule.enabled') : t('rule.disabled') }}</el-tag>
        </el-descriptions-item>
      </el-descriptions>

      <div style="margin-top: 20px;">
        <template v-if="!mfaStatus.enabled">
          <el-button type="primary" @click="handleSetup" :loading="loading">{{ t('mfa.setup') }}</el-button>
          <div v-if="qrCodeUrl" style="margin-top: 16px;">
            <p>{{ t('mfa.scanQR') }}</p>
            <!-- 使用QR码生成服务将URL渲染为二维码图片 -->
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
        </template>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
h2 {
  margin: 0 0 16px 0;
}
</style>

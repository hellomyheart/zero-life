<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { forgotPassword, resetPassword } from '@/api/auth'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'

const { t } = useI18n()

const step = ref(1)
const loading = ref(false)
const emailFormRef = ref<FormInstance>()
const resetFormRef = ref<FormInstance>()

const emailForm = reactive({
  email: '',
})

const resetForm = reactive({
  token: '',
  new_password: '',
})

const emailRules: FormRules = {
  email: [
    { required: true, message: t('common.required'), trigger: 'blur' },
    { type: 'email', message: 'Please enter a valid email', trigger: 'blur' },
  ],
}

const resetRules: FormRules = {
  token: [{ required: true, message: t('common.required'), trigger: 'blur' }],
  new_password: [
    { required: true, message: t('common.required'), trigger: 'blur' },
    { min: 8, message: 'Password must be at least 8 characters', trigger: 'blur' },
  ],
}

async function handleSendEmail() {
  const valid = await emailFormRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await forgotPassword({ email: emailForm.email })
    ElMessage.success('Reset email sent')
    step.value = 2
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  } finally {
    loading.value = false
  }
}

async function handleReset() {
  const valid = await resetFormRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await resetPassword({ token: resetForm.token, new_password: resetForm.new_password })
    ElMessage.success('Password reset successfully')
    step.value = 3
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="reset-page">
    <el-card class="reset-card">
      <template #header>
        <h2>{{ t('auth.resetPassword') }}</h2>
      </template>

      <template v-if="step === 1">
        <el-form ref="emailFormRef" :model="emailForm" :rules="emailRules" label-position="top" @submit.prevent="handleSendEmail">
          <el-form-item :label="t('auth.email')" prop="email">
            <el-input v-model="emailForm.email" type="email" :placeholder="t('common.inputPlaceholder')" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="loading" style="width: 100%" native-type="submit">
              {{ t('auth.sendResetEmail') }}
            </el-button>
          </el-form-item>
        </el-form>
      </template>

      <template v-else-if="step === 2">
        <el-form ref="resetFormRef" :model="resetForm" :rules="resetRules" label-position="top" @submit.prevent="handleReset">
          <el-form-item :label="t('auth.resetToken')" prop="token">
            <el-input v-model="resetForm.token" :placeholder="t('common.inputPlaceholder')" />
          </el-form-item>
          <el-form-item :label="t('auth.newPassword')" prop="new_password">
            <el-input v-model="resetForm.new_password" type="password" show-password :placeholder="t('common.inputPlaceholder')" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="loading" style="width: 100%" native-type="submit">
              {{ t('auth.resetPassword') }}
            </el-button>
          </el-form-item>
        </el-form>
      </template>

      <template v-else>
        <el-result icon="success" title="Password Reset Successfully" sub-title="You can now login with your new password">
          <template #extra>
            <router-link to="/login">
              <el-button type="primary">{{ t('auth.goLogin') }}</el-button>
            </router-link>
          </template>
        </el-result>
      </template>

      <div v-if="step < 3" class="reset-links">
        <router-link to="/login">{{ t('auth.goLogin') }}</router-link>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.reset-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #f5f7fa;
}

.reset-card {
  width: 400px;
}

.reset-card h2 {
  text-align: center;
  margin: 0;
}

.reset-links {
  text-align: center;
  margin-top: 8px;
}

.reset-links a {
  font-size: 14px;
  color: var(--el-color-primary);
  text-decoration: none;
}
</style>

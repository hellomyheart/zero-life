<!-- 登录页面 - 用户通过邮箱和密码登录系统，启用MFA时需二次验证 -->
<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import type { LoginReq } from '@/types/auth'

const { t } = useI18n()
const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

const formRef = ref<FormInstance>()
const loading = ref(false)
const mfaCode = ref('')
const mfaLoading = ref(false)

const form = reactive<LoginReq>({
  email: '',
  password: '',
})

const rules: FormRules<LoginReq> = {
  email: [
    { required: true, message: t('common.required'), trigger: 'blur' },
    { type: 'email', message: 'Please enter a valid email', trigger: 'blur' },
  ],
  password: [
    { required: true, message: t('common.required'), trigger: 'blur' },
  ],
}

async function handleLogin() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await authStore.login(form)
    if (authStore.mfaRequired) {
      return
    }
    ElMessage.success(t('auth.loginSuccess'))
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  } finally {
    loading.value = false
  }
}

async function handleMFAVerify() {
  if (!mfaCode.value || mfaCode.value.length !== 6) {
    ElMessage.warning('Please enter 6-digit code')
    return
  }

  mfaLoading.value = true
  try {
    await authStore.verifyMFALogin(mfaCode.value)
    ElMessage.success(t('auth.loginSuccess'))
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  } finally {
    mfaLoading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <el-card class="login-card">
      <template #header>
        <h2>{{ t('auth.login') }}</h2>
      </template>

      <!-- MFA二次验证 -->
      <div v-if="authStore.mfaRequired">
        <p class="mfa-hint">Please enter the 6-digit code from your authenticator app</p>
        <el-form label-position="top" @submit.prevent="handleMFAVerify">
          <el-form-item label="MFA Code">
            <el-input v-model="mfaCode" maxlength="6" placeholder="000000" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="mfaLoading" style="width: 100%" native-type="submit">
              Verify
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 普通登录表单 -->
      <el-form v-else ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="handleLogin">
        <el-form-item :label="t('auth.email')" prop="email">
          <el-input v-model="form.email" type="email" :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('auth.password')" prop="password">
          <el-input v-model="form.password" type="password" show-password :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" style="width: 100%" native-type="submit">
            {{ t('auth.login') }}
          </el-button>
        </el-form-item>
      </el-form>

      <div v-if="!authStore.mfaRequired" class="login-links">
        <router-link to="/reset-password">{{ t('auth.forgotPassword') }}</router-link>
        <router-link to="/register">{{ t('auth.noAccount') }} {{ t('auth.goRegister') }}</router-link>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: var(--app-bg);
  padding: 20px;
}

.login-card {
  width: 400px;
}

.login-card h2 {
  text-align: center;
  margin: 0;
}

.login-links {
  display: flex;
  justify-content: space-between;
  margin-top: 8px;
}

.login-links a {
  font-size: 14px;
  color: var(--el-color-primary);
  text-decoration: none;
}

.mfa-hint {
  text-align: center;
  color: var(--el-text-color-secondary);
  margin-bottom: 16px;
}
</style>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import type { RegisterReq } from '@/types/auth'

const { t } = useI18n()
const authStore = useAuthStore()
const router = useRouter()

const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive<RegisterReq>({
  email: '',
  password: '',
  nickname: '',
})

const rules: FormRules<RegisterReq> = {
  email: [
    { required: true, message: t('common.required'), trigger: 'blur' },
    { type: 'email', message: 'Please enter a valid email', trigger: 'blur' },
  ],
  password: [
    { required: true, message: t('common.required'), trigger: 'blur' },
    { min: 8, message: 'Password must be at least 8 characters', trigger: 'blur' },
  ],
  nickname: [
    { required: true, message: t('common.required'), trigger: 'blur' },
  ],
}

async function handleRegister() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await authStore.register(form)
    ElMessage.success(t('auth.registerSuccess'))
    router.push('/')
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="register-page">
    <el-card class="register-card">
      <template #header>
        <h2>{{ t('auth.register') }}</h2>
      </template>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="handleRegister">
        <el-form-item :label="t('auth.email')" prop="email">
          <el-input v-model="form.email" type="email" :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('auth.password')" prop="password">
          <el-input v-model="form.password" type="password" show-password :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('auth.nickname')" prop="nickname">
          <el-input v-model="form.nickname" :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" style="width: 100%" native-type="submit">
            {{ t('auth.register') }}
          </el-button>
        </el-form-item>
      </el-form>
      <div class="register-links">
        <router-link to="/login">{{ t('auth.hasAccount') }} {{ t('auth.goLogin') }}</router-link>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.register-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #f5f7fa;
}

.register-card {
  width: 400px;
}

.register-card h2 {
  text-align: center;
  margin: 0;
}

.register-links {
  text-align: center;
  margin-top: 8px;
}

.register-links a {
  font-size: 14px;
  color: var(--el-color-primary);
  text-decoration: none;
}
</style>

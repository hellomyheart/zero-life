<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import type { UpdateProfileReq, ChangePasswordReq } from '@/types/auth'

const { t } = useI18n()
const authStore = useAuthStore()

const profileFormRef = ref<FormInstance>()
const passwordFormRef = ref<FormInstance>()
const profileLoading = ref(false)
const passwordLoading = ref(false)

const profileForm = reactive<UpdateProfileReq>({
  nickname: '',
  language: '',
  timezone: '',
})

const passwordForm = reactive<ChangePasswordReq & { confirm_password: string }>({
  old_password: '',
  new_password: '',
  confirm_password: '',
})

const profileRules: FormRules = {
  nickname: [{ required: true, message: t('common.required'), trigger: 'blur' }],
}

const passwordRules: FormRules = {
  old_password: [{ required: true, message: t('common.required'), trigger: 'blur' }],
  new_password: [
    { required: true, message: t('common.required'), trigger: 'blur' },
    { min: 8, message: 'Password must be at least 8 characters', trigger: 'blur' },
  ],
  confirm_password: [
    { required: true, message: t('common.required'), trigger: 'blur' },
    {
      validator: (_rule, value, callback) => {
        if (value !== passwordForm.new_password) {
          callback(new Error('Passwords do not match'))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
}

onMounted(async () => {
  await authStore.loadProfile()
  if (authStore.user) {
    profileForm.nickname = authStore.user.nickname
    profileForm.language = authStore.user.language
    profileForm.timezone = authStore.user.timezone
  }
})

async function handleProfileSubmit() {
  const valid = await profileFormRef.value?.validate().catch(() => false)
  if (!valid) return

  profileLoading.value = true
  try {
    await authStore.updateProfile(profileForm)
    ElMessage.success(t('common.success'))
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  } finally {
    profileLoading.value = false
  }
}

async function handlePasswordSubmit() {
  const valid = await passwordFormRef.value?.validate().catch(() => false)
  if (!valid) return

  passwordLoading.value = true
  try {
    await authStore.changePassword({
      old_password: passwordForm.old_password,
      new_password: passwordForm.new_password,
    })
    ElMessage.success(t('common.success'))
    passwordForm.old_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  } finally {
    passwordLoading.value = false
  }
}
</script>

<template>
  <div class="profile-page">
    <h2>{{ t('profile.title') }}</h2>

    <el-card style="margin-bottom: 20px">
      <template #header>{{ t('profile.title') }}</template>
      <el-form ref="profileFormRef" :model="profileForm" :rules="profileRules" label-width="120px">
        <el-form-item :label="t('profile.nickname')" prop="nickname">
          <el-input v-model="profileForm.nickname" />
        </el-form-item>
        <el-form-item :label="t('profile.language')">
          <el-select v-model="profileForm.language">
            <el-option label="中文" value="zh-CN" />
            <el-option label="English" value="en-US" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('profile.timezone')">
          <el-input v-model="profileForm.timezone" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="profileLoading" @click="handleProfileSubmit">{{ t('common.save') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card>
      <template #header>{{ t('profile.changePassword') }}</template>
      <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-width="120px">
        <el-form-item :label="t('profile.oldPassword')" prop="old_password">
          <el-input v-model="passwordForm.old_password" type="password" show-password />
        </el-form-item>
        <el-form-item :label="t('profile.newPassword')" prop="new_password">
          <el-input v-model="passwordForm.new_password" type="password" show-password />
        </el-form-item>
        <el-form-item :label="t('profile.confirmNewPassword')" prop="confirm_password">
          <el-input v-model="passwordForm.confirm_password" type="password" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="passwordLoading" @click="handlePasswordSubmit">{{ t('common.save') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

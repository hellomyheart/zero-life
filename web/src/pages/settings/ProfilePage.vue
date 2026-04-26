<script setup lang="ts">
/**
 * 用户资料页面
 * 功能：
 * - 显示和编辑用户基本信息
 * - 修改密码
 */
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { getProfile, updateProfile, changePassword } from '@/api/profile'
import { ElMessage } from 'element-plus'
import type { FormRules } from 'element-plus'
import type { ProfileResp } from '@/types/auth'

const { t } = useI18n()

const user = ref<ProfileResp>({
  id: 0,
  email: '',
  nickname: '',
  language: '',
  timezone: '',
  created_at: '',
})

const form = ref({
  nickname: '',
  language: '',
  timezone: '',
})

// 密码表单
const passwordForm = ref({
  old_password: '',
  new_password: '',
  confirm_password: '',
})

const loading = ref(false)

// 表单引用 - 用于调用validate方法
const profileFormRef = ref()
const passwordFormRef = ref()

// 个人资料表单验证规则 - 姓名必填、邮箱格式校验
const profileRules: FormRules = {
  nickname: [{ required: true, message: t('profile.nameRequired'), trigger: 'blur' }],
  language: [{ required: true, message: t('common.required'), trigger: 'blur' }],
  timezone: [{ required: true, message: t('common.required'), trigger: 'blur' }],
}

// 密码表单验证规则 - 新密码最小长度6位
const passwordRules: FormRules = {
  old_password: [{ required: true, message: t('profile.oldPasswordRequired'), trigger: 'blur' }],
  new_password: [
    { required: true, message: t('profile.newPasswordRequired'), trigger: 'blur' },
    { min: 6, message: t('profile.passwordMinLength'), trigger: 'blur' },
  ],
  confirm_password: [{ required: true, message: t('profile.confirmPasswordRequired'), trigger: 'blur' }],
}

/**
 * 获取用户资料
 */
async function fetchProfile() {
  try {
    const res = await getProfile()
    user.value = res as unknown as User
    form.value.nickname = user.value.nickname
    form.value.language = user.value.language
    form.value.timezone = user.value.timezone
  } catch {
    ElMessage.error(t('common.fetchError') || 'Failed to load data')
  }
}

/**
 * 更新用户资料
 */
async function handleUpdateProfile() {
  // 先进行表单验证，通过后才提交
  try {
    await profileFormRef.value?.validate()
  } catch {
    return
  }
  loading.value = true
  try {
    await updateProfile(form.value)
    ElMessage.success(t('common.success'))
    await fetchProfile()
  } catch {
    ElMessage.error(t('common.fetchError') || 'Failed to load data')
  } finally {
    loading.value = false
  }
}

/**
 * 修改密码
 */
async function handleChangePassword() {
  // 先进行表单验证
  try {
    await passwordFormRef.value?.validate()
  } catch {
    return
  }
  if (passwordForm.value.new_password !== passwordForm.value.confirm_password) {
    ElMessage.error(t('profile.passwordMismatch'))
    return
  }

  loading.value = true
  try {
    await changePassword({
      old_password: passwordForm.value.old_password,
      new_password: passwordForm.value.new_password,
    })
    ElMessage.success(t('common.success'))
    passwordForm.value = {
      old_password: '',
      new_password: '',
      confirm_password: '',
    }
  } catch {
    ElMessage.error(t('common.fetchError') || 'Failed to load data')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchProfile()
})
</script>

<template>
  <div class="profile-page">
    <el-card>
      <template #header>
        <h2>{{ t('profile.title') }}</h2>
      </template>

      <el-form ref="profileFormRef" :model="form" :rules="profileRules" label-width="120px">
        <el-form-item :label="t('profile.name')" prop="nickname">
          <el-input v-model="form.nickname" />
        </el-form-item>
        <el-form-item :label="t('profile.language')" prop="language">
          <el-select v-model="form.language" style="width: 100%">
            <el-option label="简体中文" value="zh-CN" />
            <el-option label="English" value="en-US" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('profile.timezone')" prop="timezone">
          <el-input v-model="form.timezone" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleUpdateProfile" :loading="loading">
            {{ t('common.save') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card style="margin-top: 20px">
      <template #header>
        <h3>{{ t('profile.changePassword') }}</h3>
      </template>

      <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-width="120px">
        <el-form-item :label="t('profile.oldPassword')" prop="old_password">
          <el-input v-model="passwordForm.old_password" type="password" />
        </el-form-item>
        <el-form-item :label="t('profile.newPassword')" prop="new_password">
          <el-input v-model="passwordForm.new_password" type="password" />
        </el-form-item>
        <el-form-item :label="t('profile.confirmPassword')" prop="confirm_password">
          <el-input v-model="passwordForm.confirm_password" type="password" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleChangePassword" :loading="loading">
            {{ t('common.save') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.profile-page {
  padding: 20px;
}
</style>

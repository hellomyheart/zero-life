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
import type { User } from '@/types/user'

const { t } = useI18n()

// 用户信息
const user = ref<User>({
  id: 0,
  email: '',
  name: '',
  created_at: '',
  updated_at: '',
})

// 表单数据
const form = ref({
  name: '',
  email: '',
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
  name: [{ required: true, message: t('profile.nameRequired'), trigger: 'blur' }],
  email: [
    { required: true, message: t('profile.emailRequired'), trigger: 'blur' },
    { type: 'email', message: t('profile.emailInvalid'), trigger: 'blur' },
  ],
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
    form.value.name = user.value.name
    form.value.email = user.value.email
  } catch {
    // handle error
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
    // handle error
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
    // handle error
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
        <el-form-item :label="t('profile.name')" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="t('profile.email')" prop="email">
          <el-input v-model="form.email" />
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

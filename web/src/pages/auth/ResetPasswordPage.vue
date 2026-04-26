<!-- 重置密码页面 - 分三步完成密码重置：1.输入邮箱发送重置邮件 2.输入令牌和新密码 3.重置成功 -->
<script setup lang="ts">
// 导入 Vue 响应式 API
import { ref, reactive } from 'vue'
// 导入国际化钩子函数
import { useI18n } from 'vue-i18n'
// 导入认证相关 API（忘记密码、重置密码）
import { forgotPassword, resetPassword } from '@/api/auth'
// 导入 Element Plus 消息提示组件
import { ElMessage } from 'element-plus'
// 导入表单类型定义
import type { FormInstance, FormRules } from 'element-plus'

// 国际化翻译函数
const { t } = useI18n()

// 当前步骤：1=输入邮箱, 2=输入令牌和新密码, 3=重置成功
const step = ref(1)
// 加载状态，防止重复提交
const loading = ref(false)
// 邮箱表单引用
const emailFormRef = ref<FormInstance>()
// 重置密码表单引用
const resetFormRef = ref<FormInstance>()

// 第一步：邮箱表单数据
const emailForm = reactive({
  email: '',
})

// 第二步：重置密码表单数据（token 从邮件中获取）
const resetForm = reactive({
  token: '',
  new_password: '',
})

// 邮箱表单验证规则
const emailRules: FormRules = {
  email: [
    { required: true, message: t('common.required'), trigger: 'blur' },
    { type: 'email', message: 'Please enter a valid email', trigger: 'blur' },
  ],
}

// 重置密码表单验证规则
const resetRules: FormRules = {
  // 重置令牌：必�?  token: [{ required: true, message: t('common.required'), trigger: 'blur' }],
  // 新密码：必填 + 最�?�?  new_password: [
    { required: true, message: t('common.required'), trigger: 'blur' },
    { min: 8, message: 'Password must be at least 8 characters', trigger: 'blur' },
  ],
}

// 处理发送重置邮件（第一步）
async function handleSendEmail() {
  // 验证邮箱表单
  const valid = await emailFormRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    // 调用忘记密码 API，发送重置邮件
    await forgotPassword({ email: emailForm.email })
    ElMessage.success('Reset email sent')
    // 进入第二步：输入令牌和新密码
    step.value = 2
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  } finally {
    loading.value = false
  }
}

// 处理重置密码提交（第二步）
async function handleReset() {
  // 验证重置密码表单
  const valid = await resetFormRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    // 调用重置密码 API，提交令牌和新密�?    await resetPassword({ token: resetForm.token, password: resetForm.new_password })
    ElMessage.success('Password reset successfully')
    // 进入第三步：显示成功提示
    step.value = 3
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <!-- 重置密码页面容器，垂直水平居�?-->
  <div class="reset-page">
    <el-card class="reset-card">
      <!-- 卡片标题 -->
      <template #header>
        <h2>{{ t('auth.resetPassword') }}</h2>
      </template>

      <!-- 第一步：输入邮箱发送重置邮�?-->
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

      <!-- 第二步：输入令牌和新密码 -->
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

      <!-- 第三步：重置成功提示 -->
      <template v-else>
        <el-result icon="success" title="Password Reset Successfully" sub-title="You can now login with your new password">
          <template #extra>
            <router-link to="/login">
              <el-button type="primary">{{ t('auth.goLogin') }}</el-button>
            </router-link>
          </template>
        </el-result>
      </template>

      <!-- 前两步显示返回登录链�?-->
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
  background: var(--app-bg); padding: 20px;
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

<!-- 注册页面 - 新用户通过邮箱、密码和昵称注册账号 -->
<script setup lang="ts">
// 导入 Vue 响应式 API
import { ref, reactive } from 'vue'
// 导入国际化钩子函数
import { useI18n } from 'vue-i18n'
// 导入认证状态管理
import { useAuthStore } from '@/stores/auth'
// 导入路由钩子
import { useRouter } from 'vue-router'
// 导入 Element Plus 消息提示组件
import { ElMessage } from 'element-plus'
// 导入表单类型定义
import type { FormInstance, FormRules } from 'element-plus'
// 导入注册请求参数类型
import type { RegisterReq } from '@/types/auth'

// 国际化翻译函数
const { t } = useI18n()
// 认证状态管理实例，用于调用注册方法
const authStore = useAuthStore()
// 路由实例，用于注册成功后跳转
const router = useRouter()

// 表单引用，用于调用表单验证方法
const formRef = ref<FormInstance>()
// 注册加载状态，防止重复提交
const loading = ref(false)

// 注册表单数据，包含确认密码字段
const form = reactive<RegisterReq & { confirmPassword: string }>({
  email: '',
  password: '',
  nickname: '',
  confirmPassword: '',
})

// 表单验证规则
const rules: FormRules<typeof form> = {
  // 邮箱：必填 + 格式校验
  email: [
    { required: true, message: t('common.required'), trigger: 'blur' },
    { type: 'email', message: 'Please enter a valid email', trigger: 'blur' },
  ],
  // 密码：必填 + 最少8位
  password: [
    { required: true, message: t('common.required'), trigger: 'blur' },
    { min: 8, message: 'Password must be at least 8 characters', trigger: 'blur' },
  ],
  // 确认密码：必填 + 必须与密码一致
  confirmPassword: [
    { required: true, message: t('common.required'), trigger: 'blur' },
    {
      // 自定义验证器：检查两次输入的密码是否一致
      validator: (_rule: unknown, value: string, callback: (error?: Error) => void) => {
        if (value !== form.password) {
          callback(new Error('Passwords do not match'))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
  // 昵称：必填
  nickname: [
    { required: true, message: t('common.required'), trigger: 'blur' },
  ],
}

// 处理注册提交
async function handleRegister() {
  // 先进行表单验证，验证失败则不继续
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    // 调用认证 store 的注册方法
    await authStore.register(form)
    ElMessage.success(t('auth.registerSuccess'))
    // 注册成功后跳转首页
    router.push('/')
  } catch (err) {
    // 注册失败显示错误信息
    ElMessage.error((err as Error).message || t('common.failed'))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <!-- 注册页面容器，垂直水平居中 -->
  <div class="register-page">
    <el-card class="register-card">
      <!-- 卡片标题 -->
      <template #header>
        <h2>{{ t('auth.register') }}</h2>
      </template>
      <!-- 注册表单，提交时调用 handleRegister -->
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="handleRegister">
        <!-- 邮箱输入框 -->
        <el-form-item :label="t('auth.email')" prop="email">
          <el-input v-model="form.email" type="email" :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>
        <!-- 密码输入框，支持显示/隐藏密码 -->
        <el-form-item :label="t('auth.password')" prop="password">
          <el-input v-model="form.password" type="password" show-password :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>
        <!-- 确认密码输入框 -->
        <el-form-item :label="t('auth.confirmPassword')" prop="confirmPassword">
          <el-input v-model="form.confirmPassword" type="password" show-password :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>
        <!-- 昵称输入框 -->
        <el-form-item :label="t('auth.nickname')" prop="nickname">
          <el-input v-model="form.nickname" :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>
        <!-- 注册按钮 -->
        <el-form-item>
          <el-button type="primary" :loading="loading" style="width: 100%" native-type="submit">
            {{ t('auth.register') }}
          </el-button>
        </el-form-item>
      </el-form>
      <!-- 底部链接：已有账号去登录 -->
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

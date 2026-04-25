<!-- 登录页面 - 用户通过邮箱和密码登录系统 -->
<script setup lang="ts">
// 导入 Vue 响应式 API
import { ref, reactive } from 'vue'
// 导入国际化钩子函数
import { useI18n } from 'vue-i18n'
// 导入认证状态管理
import { useAuthStore } from '@/stores/auth'
// 导入路由钩子
import { useRouter, useRoute } from 'vue-router'
// 导入 Element Plus 消息提示组件
import { ElMessage } from 'element-plus'
// 导入表单类型定义
import type { FormInstance, FormRules } from 'element-plus'
// 导入登录请求参数类型
import type { LoginReq } from '@/types/auth'

// 国际化翻译函数
const { t } = useI18n()
// 认证状态管理实例，用于调用登录方法
const authStore = useAuthStore()
// 路由实例，用于登录后跳转
const router = useRouter()
// 当前路由信息，用于获取重定向地址
const route = useRoute()

// 表单引用，用于调用表单验证方法
const formRef = ref<FormInstance>()
// 登录加载状态，防止重复提交
const loading = ref(false)

// 登录表单数据
const form = reactive<LoginReq>({
  email: '',
  password: '',
})

// 表单验证规则
const rules: FormRules<LoginReq> = {
  // 邮箱：必填 + 格式校验
  email: [
    { required: true, message: t('common.required'), trigger: 'blur' },
    { type: 'email', message: 'Please enter a valid email', trigger: 'blur' },
  ],
  // 密码：必填
  password: [
    { required: true, message: t('common.required'), trigger: 'blur' },
  ],
}

// 处理登录提交
async function handleLogin() {
  // 先进行表单验证，验证失败则不继续
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    // 调用认证 store 的登录方法
    await authStore.login(form)
    ElMessage.success(t('auth.loginSuccess'))
    // 登录成功后跳转：优先跳转到之前被拦截的页面，否则跳转首页
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (err) {
    // 登录失败显示错误信息
    ElMessage.error((err as Error).message || t('common.failed'))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <!-- 登录页面容器，垂直水平居中 -->
  <div class="login-page">
    <el-card class="login-card">
      <!-- 卡片标题 -->
      <template #header>
        <h2>{{ t('auth.login') }}</h2>
      </template>
      <!-- 登录表单，提交时调用 handleLogin -->
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="handleLogin">
        <!-- 邮箱输入框 -->
        <el-form-item :label="t('auth.email')" prop="email">
          <el-input v-model="form.email" type="email" :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>
        <!-- 密码输入框，支持显示/隐藏密码 -->
        <el-form-item :label="t('auth.password')" prop="password">
          <el-input v-model="form.password" type="password" show-password :placeholder="t('common.inputPlaceholder')" />
        </el-form-item>
        <!-- 登录按钮 -->
        <el-form-item>
          <el-button type="primary" :loading="loading" style="width: 100%" native-type="submit">
            {{ t('auth.login') }}
          </el-button>
        </el-form-item>
      </el-form>
      <!-- 底部链接：忘记密码 / 去注册 -->
      <div class="login-links">
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
  background: #f5f7fa;
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
</style>

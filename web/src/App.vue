<script setup lang="ts">
// 根组件 - 根据认证状态和路由决定显示布局还是裸页面
// 使用 el-config-provider 动态切换 Element Plus 的语言
import { useAuthStore } from '@/stores/auth'
import AppLayout from '@/components/layout/AppLayout.vue'
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { elementLocales } from '@/utils/element-locale'

const authStore = useAuthStore()
const route = useRoute()
const { locale } = useI18n()

// 根据 i18n 当前语言获取 Element Plus 对应的 locale
const elementLocale = computed(() => elementLocales[locale.value] || elementLocales['zh-CN'])

// 判断当前是否为认证页面（登录/注册/重置密码）
const isAuthPage = computed(() => {
  return ['/login', '/register', '/reset-password'].includes(route.path)
})
</script>

<template>
  <!-- el-config-provider 动态切换 Element Plus 组件的语言（分页、日期选择等） -->
  <el-config-provider :locale="elementLocale">
    <!-- 已登录且不在认证页面时，显示应用布局（侧边栏+顶栏+内容区） -->
    <AppLayout v-if="authStore.isAuthenticated && !isAuthPage" />
    <!-- 未登录或在认证页面时，直接显示路由页面 -->
    <router-view v-else />
  </el-config-provider>
</template>

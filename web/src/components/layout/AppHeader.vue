<script setup lang="ts">
// 顶栏组件 - 包含侧边栏折叠按钮、语言切换、用户菜单
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import { useRouter } from 'vue-router'
import { Fold, Expand } from '@element-plus/icons-vue'

const { t, locale } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()
const router = useRouter()

// 切换语言并保存到本地存储
function handleLanguageChange(lang: string) {
  locale.value = lang
  localStorage.setItem('locale', lang)
}

// 退出登录
function handleLogout() {
  authStore.logout()
}

// 跳转到个人设置页
function goToProfile() {
  router.push('/settings/profile')
}
</script>

<template>
  <div class="header-container">
    <div class="header-left">
      <!-- 侧边栏折叠/展开切换按钮 -->
      <el-icon class="collapse-btn" @click="appStore.toggleSidebar">
        <Fold v-if="!appStore.sidebarCollapsed" />
        <Expand v-else />
      </el-icon>
      <span class="app-title">{{ t('app.title') }}</span>
    </div>
    <div class="header-right">
      <!-- 语言切换下拉菜单 -->
      <el-dropdown trigger="click" @command="handleLanguageChange">
        <span class="dropdown-link">
          {{ locale === 'zh-CN' ? '中文' : 'English' }}
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="zh-CN">中文</el-dropdown-item>
            <el-dropdown-item command="en-US">English</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      <!-- 用户信息下拉菜单 -->
      <el-dropdown trigger="click">
        <span class="dropdown-link">
          <el-icon><User /></el-icon>
          <!-- 优先显示昵称，其次显示邮箱 -->
          {{ authStore.user?.nickname || authStore.user?.email || '' }}
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="goToProfile">
              {{ t('setting.profile') }}
            </el-dropdown-item>
            <el-dropdown-item @click="handleLogout">
              {{ t('auth.logout') }}
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script lang="ts">
import { User } from '@element-plus/icons-vue'
export default {
  components: { User },
}
</script>

<style scoped>
.header-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  padding: 0 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.collapse-btn {
  font-size: 20px;
  cursor: pointer;
  color: var(--el-text-color-primary);
}

.app-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.dropdown-link {
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  color: var(--el-text-color-primary);
  font-size: 14px;
}

.dropdown-link:hover {
  color: var(--el-color-primary);
}
</style>

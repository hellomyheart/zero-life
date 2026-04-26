<script setup lang="ts">
/**
 * 顶栏组件 - 包含侧边栏折叠按钮、主题切换、语言切换、用户菜单
 * 响应式行为：
 * - 桌面端：显示折叠按钮 + 标题 + 主题切换 + 语言 + 用户菜单
 * - 移动端：显示汉堡菜单按钮 + 标题 + 主题切换 + 用户菜单
 */
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import type { ThemeMode } from '@/stores/app'
import { useRouter } from 'vue-router'
import { computed, onMounted, onUnmounted, ref } from 'vue'

const { t, locale } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()
const router = useRouter()

/** 是否为移动端 */
const isMobile = ref(false)

function checkMobile() {
  isMobile.value = window.innerWidth < 1024
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

/** 主题模式选项 */
const themeOptions = computed(() => [
  { label: t('theme.light'), value: 'light' as ThemeMode },
  { label: t('theme.dark'), value: 'dark' as ThemeMode },
  { label: t('theme.system'), value: 'system' as ThemeMode },
])

/** 当前主题图标 */
const themeIcon = computed(() => {
  if (appStore.theme === 'dark') return 'Moon'
  if (appStore.theme === 'light') return 'Sunny'
  return 'Monitor'
})

/** 切换语言并保存到本地存储 */
function handleLanguageChange(lang: string) {
  locale.value = lang
  localStorage.setItem('locale', lang)
}

/** 退出登录 */
function handleLogout() {
  authStore.logout()
}

/** 跳转到个人设置页 */
function goToProfile() {
  router.push('/settings/profile')
}
</script>

<template>
  <div class="header-container">
    <div class="header-left">
      <!-- 桌面端：折叠按钮 -->
      <el-icon v-if="!isMobile" class="collapse-btn" @click="appStore.toggleSidebar">
        <Fold v-if="!appStore.sidebarCollapsed" />
        <Expand v-else />
      </el-icon>
      <!-- 移动端：汉堡菜单按钮 -->
      <el-icon v-else class="collapse-btn" @click="appStore.toggleMobileMenu">
        <Expand />
      </el-icon>
      <span class="app-title">{{ t('app.title') }}</span>
    </div>
    <div class="header-right">
      <!-- 主题切换 -->
      <el-dropdown trigger="click" @command="(cmd: ThemeMode) => appStore.setTheme(cmd)">
        <span class="dropdown-link theme-link">
          <el-icon><component :is="themeIcon" /></el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item
              v-for="opt in themeOptions"
              :key="opt.value"
              :command="opt.value"
              :class="{ 'is-active': appStore.theme === opt.value }"
            >
              {{ opt.label }}
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>

      <!-- 语言切换 -->
      <el-dropdown trigger="click" @command="handleLanguageChange" class="hidden-sm-and-down">
        <span class="dropdown-link">
          {{ locale === 'zh-CN' ? '中文' : 'EN' }}
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
          <span class="hidden-sm-and-down">{{ authStore.user?.nickname || authStore.user?.email || '' }}</span>
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
import { User, Fold, Expand, Sunny, Moon, Monitor } from '@element-plus/icons-vue'
export default {
  components: { User, Fold, Expand, Sunny, Moon, Monitor },
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
  color: var(--app-text-primary);
}

.app-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--app-text-primary);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.dropdown-link {
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  color: var(--app-text-secondary);
  font-size: 14px;
}

.dropdown-link:hover {
  color: var(--el-color-primary);
}

.theme-link {
  font-size: 18px;
}

@media (max-width: 767px) {
  .header-container {
    padding: 0 12px;
  }

  .app-title {
    font-size: 16px;
  }

  .header-right {
    gap: 10px;
  }
}
</style>
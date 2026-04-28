<script setup lang="ts">
/**
 * 主题切换组件 - 支持明亮/黑暗/跟随系统三种模式
 * 用于认证页面（登录/注册/重置密码）的右上角主题切换
 */
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import type { ThemeMode } from '@/stores/app'

const { t } = useI18n()
const appStore = useAppStore()

const themeOptions = computed(() => [
  { label: t('theme.light'), value: 'light' as ThemeMode },
  { label: t('theme.dark'), value: 'dark' as ThemeMode },
  { label: t('theme.system'), value: 'system' as ThemeMode },
])

const themeIcon = computed(() => {
  if (appStore.theme === 'dark') return 'Moon'
  if (appStore.theme === 'light') return 'Sunny'
  return 'Monitor'
})
</script>

<template>
  <el-dropdown trigger="click" @command="(cmd: ThemeMode) => appStore.setTheme(cmd)">
    <span class="theme-toggle">
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
</template>

<script lang="ts">
import { Sunny, Moon, Monitor } from '@element-plus/icons-vue'
export default {
  components: { Sunny, Moon, Monitor },
}
</script>

<style scoped>
.theme-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  cursor: pointer;
  font-size: 18px;
  color: var(--app-text-secondary);
  transition: background-color 0.2s;
}

.theme-toggle:hover {
  background-color: var(--el-fill-color-light);
  color: var(--el-color-primary);
}
</style>

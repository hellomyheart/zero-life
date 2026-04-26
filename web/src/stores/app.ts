/**
 * 应用全局状态管理（Pinia Store）
 *
 * 组件功能：
 * - 管理侧边栏折叠/展开状态
 * - 管理全局加载状态
 * - 管理主题模式（light / dark / system）
 *
 * 数据流：
 * - sidebarCollapsed: 由AppHeader的折叠按钮触发，影响AppSidebar的宽度
 * - loading: 全局加载遮罩，可用于页面切换等场景
 * - theme: 主题模式，持久化到localStorage，影响html.dark类
 */
import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

/** 主题模式类型 */
export type ThemeMode = 'light' | 'dark' | 'system'

export const useAppStore = defineStore('app', () => {
  /** 侧边栏是否折叠 */
  const sidebarCollapsed = ref(false)
  /** 全局加载状态 */
  const loading = ref(false)
  /** 移动端侧边栏是否打开（overlay模式） */
  const mobileMenuOpen = ref(false)

  /**
   * 主题模式：light / dark / system
   * 从localStorage恢复，默认system
   */
  const theme = ref<ThemeMode>((localStorage.getItem('theme') as ThemeMode) || 'system')

  /** 切换侧边栏折叠/展开状态 */
  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  /** 切换移动端侧边栏 */
  function toggleMobileMenu() {
    mobileMenuOpen.value = !mobileMenuOpen.value
  }

  /** 关闭移动端侧边栏 */
  function closeMobileMenu() {
    mobileMenuOpen.value = false
  }

  /**
   * 设置全局加载状态
   * @param v - true显示加载遮罩，false隐藏
   */
  function setLoading(v: boolean) {
    loading.value = v
  }

  /**
   * 设置主题模式
   * @param mode - light / dark / system
   */
  function setTheme(mode: ThemeMode) {
    theme.value = mode
    localStorage.setItem('theme', mode)
  }

  /**
   * 根据theme设置计算实际是否为暗色模式
   * system模式时跟随系统偏好
   */
  function applyTheme() {
    const isDark = theme.value === 'dark' ||
      (theme.value === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches)
    document.documentElement.classList.toggle('dark', isDark)
  }

  // 监听theme变化，立即应用
  watch(theme, () => applyTheme(), { immediate: true })

  // 监听系统主题变化（system模式下自动切换）
  if (typeof window !== 'undefined') {
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
      if (theme.value === 'system') {
        applyTheme()
      }
    })
  }

  return {
    sidebarCollapsed,
    loading,
    mobileMenuOpen,
    theme,
    toggleSidebar,
    toggleMobileMenu,
    closeMobileMenu,
    setLoading,
    setTheme,
    applyTheme,
  }
})

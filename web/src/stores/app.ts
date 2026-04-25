/**
 * 应用全局状态管理（Pinia Store）
 *
 * 组件功能：
 * - 管理侧边栏折叠/展开状态
 * - 管理全局加载状态
 *
 * 数据流：
 * - sidebarCollapsed: 由AppHeader的折叠按钮触发，影响AppSidebar的宽度
 * - loading: 全局加载遮罩，可用于页面切换等场景
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  /** 侧边栏是否折叠 */
  const sidebarCollapsed = ref(false)
  /** 全局加载状态 */
  const loading = ref(false)

  /** 切换侧边栏折叠/展开状态 */
  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  /**
   * 设置全局加载状态
   * @param v - true显示加载遮罩，false隐藏
   */
  function setLoading(v: boolean) {
    loading.value = v
  }

  return {
    sidebarCollapsed,
    loading,
    toggleSidebar,
    setLoading,
  }
})

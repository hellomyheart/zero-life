// 应用全局状态管理 - 管理侧边栏折叠、全局加载状态
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const sidebarCollapsed = ref(false) // 侧边栏是否折叠
  const loading = ref(false) // 全局加载状态

  // toggleSidebar 切换侧边栏折叠/展开
  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  // setLoading 设置全局加载状态
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

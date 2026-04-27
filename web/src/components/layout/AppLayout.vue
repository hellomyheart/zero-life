<script setup lang="ts">
/**
 * 主布局组件 - 包含侧边栏、顶栏和内容区
 * 响应式行为：
 * - 桌面端(>=1024px)：侧边栏固定在左侧，可折叠
 * - 平板/手机(<1024px)：侧边栏为抽屉覆盖模式，点击遮罩关闭
 */
import AppHeader from './AppHeader.vue'
import AppSidebar from './AppSidebar.vue'
import { useAppStore } from '@/stores/app'
import { ref, onMounted, onUnmounted } from 'vue'

const appStore = useAppStore()

/** 是否为桌面端（侧边栏固定模式） */
const isDesktop = ref(window.innerWidth >= 1024)

function onResize() {
  isDesktop.value = window.innerWidth >= 1024
}

onMounted(() => window.addEventListener('resize', onResize))
onUnmounted(() => window.removeEventListener('resize', onResize))

/** 关闭移动端侧边栏 */
function closeMobileMenu() {
  appStore.closeMobileMenu()
}
</script>

<template>
  <el-container class="app-layout">
    <!-- 桌面端：固定侧边栏 -->
    <AppSidebar v-if="isDesktop" />
    <!-- 移动端：抽屉覆盖侧边栏 -->
    <template v-else>
      <el-drawer
        v-model="appStore.mobileMenuOpen"
        direction="ltr"
        :show-close="false"
        :with-header="false"
        size="220px"
        class="mobile-sidebar-drawer"
        @close="closeMobileMenu"
      >
        <AppSidebar />
      </el-drawer>
    </template>

    <el-container class="main-container">
      <el-header class="app-header">
        <AppHeader />
      </el-header>
      <el-main class="app-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.app-layout {
  height: 100vh;
  width: 100vw;
}

.main-container {
  flex-direction: column;
  overflow: hidden;
}

.app-header {
  padding: 0;
  height: var(--app-header-height);
  border-bottom: 1px solid var(--app-border);
  background: var(--app-header-bg);
}

.app-main {
  background: var(--app-bg);
  overflow-y: auto;
  padding: 20px;
}

@media (max-width: 767px) {
  .app-main {
    padding: 12px;
  }
}
</style>

<style>
/* 移动端侧边栏抽屉样式（非scoped，影响el-drawer内部） */
.mobile-sidebar-drawer .el-drawer__body {
  padding: 0;
  background: var(--app-sidebar-bg);
}

.mobile-sidebar-drawer .el-drawer__body .app-sidebar {
  border-right: none;
  height: 100%;
}
</style>

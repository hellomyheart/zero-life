<script setup lang="ts">
// 主布局组件 - 包含侧边栏、顶栏和内容区
import AppHeader from './AppHeader.vue'
import AppSidebar from './AppSidebar.vue'
</script>

<template>
  <!-- 整体布局：左侧侧边栏 + 右侧主内容区 -->
  <el-container class="app-layout">
    <AppSidebar />
    <el-container class="main-container">
      <!-- 顶部导航栏 -->
      <el-header class="app-header">
        <AppHeader />
      </el-header>
      <!-- 页面内容区，通过router-view渲染当前路由页面 -->
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
  height: 56px;
  border-bottom: 1px solid var(--el-border-color-light);
  background: #fff;
}

.app-main {
  background: #f5f7fa;
  overflow-y: auto;
  padding: 20px;
}
</style>

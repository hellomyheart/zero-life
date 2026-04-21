<script setup lang="ts">
import AppHeader from './AppHeader.vue'
import AppSidebar from './AppSidebar.vue'
</script>

<template>
  <el-container class="app-layout">
    <AppSidebar />
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

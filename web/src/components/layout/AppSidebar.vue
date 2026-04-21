<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { useRoute, useRouter } from 'vue-router'
import { computed } from 'vue'
import {
  House,
  Wallet,
  List,
  Menu,
  PriceTag,
  PieChart,
  Bell,
  DataAnalysis,
  Upload,
  Setting,
  User,
  Coin,
} from '@element-plus/icons-vue'

const { t } = useI18n()
const appStore = useAppStore()
const route = useRoute()
const router = useRouter()

const activeMenu = computed(() => route.path)

const sidebarWidth = computed(() => (appStore.sidebarCollapsed ? '64px' : '220px'))

function handleSelect(index: string) {
  router.push(index)
}
</script>

<template>
  <el-aside :width="sidebarWidth" class="app-sidebar">
    <el-menu
      :default-active="activeMenu"
      :collapse="appStore.sidebarCollapsed"
      :collapse-transition="false"
      router
      @select="handleSelect"
    >
      <el-menu-item index="/">
        <el-icon><House /></el-icon>
        <template #title>{{ t('nav.dashboard') }}</template>
      </el-menu-item>

      <el-menu-item index="/accounts">
        <el-icon><Wallet /></el-icon>
        <template #title>{{ t('nav.accounts') }}</template>
      </el-menu-item>

      <el-menu-item index="/transactions">
        <el-icon><List /></el-icon>
        <template #title>{{ t('nav.transactions') }}</template>
      </el-menu-item>

      <el-menu-item index="/categories">
        <el-icon><Menu /></el-icon>
        <template #title>{{ t('nav.categories') }}</template>
      </el-menu-item>

      <el-menu-item index="/tags">
        <el-icon><PriceTag /></el-icon>
        <template #title>{{ t('nav.tags') }}</template>
      </el-menu-item>

      <el-menu-item index="/budgets">
        <el-icon><PieChart /></el-icon>
        <template #title>{{ t('nav.budgets') }}</template>
      </el-menu-item>

      <el-menu-item index="/bills">
        <el-icon><Bell /></el-icon>
        <template #title>{{ t('nav.bills') }}</template>
      </el-menu-item>

      <el-sub-menu index="/reports">
        <template #title>
          <el-icon><DataAnalysis /></el-icon>
          <span>{{ t('nav.reports') }}</span>
        </template>
        <el-menu-item index="/reports/income-expense">{{ t('report.incomeExpense') }}</el-menu-item>
        <el-menu-item index="/reports/category">{{ t('report.category') }}</el-menu-item>
        <el-menu-item index="/reports/budget">{{ t('report.budget') }}</el-menu-item>
        <el-menu-item index="/reports/net-worth">{{ t('report.netWorth') }}</el-menu-item>
        <el-menu-item index="/reports/trend">{{ t('report.trend') }}</el-menu-item>
      </el-sub-menu>

      <el-menu-item index="/imports">
        <el-icon><Upload /></el-icon>
        <template #title>{{ t('nav.imports') }}</template>
      </el-menu-item>

      <el-menu-item index="/rules">
        <el-icon><Setting /></el-icon>
        <template #title>{{ t('nav.rules') }}</template>
      </el-menu-item>

      <el-sub-menu index="/settings">
        <template #title>
          <el-icon><User /></el-icon>
          <span>{{ t('nav.settings') }}</span>
        </template>
        <el-menu-item index="/settings/profile">{{ t('setting.profile') }}</el-menu-item>
        <el-menu-item index="/settings/currencies">
          <el-icon><Coin /></el-icon>
          {{ t('setting.currencies') }}
        </el-menu-item>
      </el-sub-menu>
    </el-menu>
  </el-aside>
</template>

<style scoped>
.app-sidebar {
  background: #fff;
  border-right: 1px solid var(--el-border-color-light);
  overflow-y: auto;
  transition: width 0.3s;
}

.app-sidebar :deep(.el-menu) {
  border-right: none;
}
</style>

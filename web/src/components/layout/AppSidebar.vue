<script setup lang="ts">
// 侧边栏导航组件 - 显示应用菜单，支持折叠/展开
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
  Opportunity,
  Timer,
  Link,
  Download,
  Paperclip,
  Histogram,
  Key,
  UserFilled,
} from '@element-plus/icons-vue'

const { t } = useI18n()
const appStore = useAppStore()
const route = useRoute()
const router = useRouter()

// 当前激活的菜单项，根据路由路径自动高亮
const activeMenu = computed(() => route.path)

// 侧边栏宽度：折叠时64px，展开时220px
const sidebarWidth = computed(() => (appStore.sidebarCollapsed ? '64px' : '220px'))

// 处理菜单项点击，跳转到对应路由
function handleSelect(index: string) {
  router.push(index)
}
</script>

<template>
  <el-aside :width="sidebarWidth" class="app-sidebar">
    <!-- Element Plus菜单组件，collapse控制折叠状态 -->
    <el-menu
      :default-active="activeMenu"
      :collapse="appStore.sidebarCollapsed"
      :collapse-transition="false"
      router
      @select="handleSelect"
    >
      <!-- 仪表盘 -->
      <el-menu-item index="/">
        <el-icon><House /></el-icon>
        <template #title>{{ t('nav.dashboard') }}</template>
      </el-menu-item>

      <!-- 账户管理 -->
      <el-menu-item index="/accounts">
        <el-icon><Wallet /></el-icon>
        <template #title>{{ t('nav.accounts') }}</template>
      </el-menu-item>

      <!-- 交易记录 -->
      <el-menu-item index="/transactions">
        <el-icon><List /></el-icon>
        <template #title>{{ t('nav.transactions') }}</template>
      </el-menu-item>

      <!-- 分类管理 -->
      <el-menu-item index="/categories">
        <el-icon><Menu /></el-icon>
        <template #title>{{ t('nav.categories') }}</template>
      </el-menu-item>

      <!-- 标签管理 -->
      <el-menu-item index="/tags">
        <el-icon><PriceTag /></el-icon>
        <template #title>{{ t('nav.tags') }}</template>
      </el-menu-item>

      <!-- 预算管理 -->
      <el-menu-item index="/budgets">
        <el-icon><PieChart /></el-icon>
        <template #title>{{ t('nav.budgets') }}</template>
      </el-menu-item>

      <!-- 账单提醒 -->
      <el-menu-item index="/bills">
        <el-icon><Bell /></el-icon>
        <template #title>{{ t('nav.bills') }}</template>
      </el-menu-item>

      <!-- 存钱罐 -->
      <el-menu-item index="/piggy-banks">
        <el-icon><Opportunity /></el-icon>
        <template #title>{{ t('nav.piggyBanks') }}</template>
      </el-menu-item>

      <!-- 周期性交易 -->
      <el-menu-item index="/recurring-transactions">
        <el-icon><Timer /></el-icon>
        <template #title>{{ t('nav.recurringTransactions') }}</template>
      </el-menu-item>

      <!-- 报表子菜单 -->
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
        <el-menu-item index="/reports/tag">{{ t('report.tag') }}</el-menu-item>
      </el-sub-menu>

      <!-- 数据导入 -->
      <el-menu-item index="/imports">
        <el-icon><Upload /></el-icon>
        <template #title>{{ t('nav.imports') }}</template>
      </el-menu-item>

      <!-- 数据导出 -->
      <el-menu-item index="/exports">
        <el-icon><Download /></el-icon>
        <template #title>{{ t('nav.exports') }}</template>
      </el-menu-item>

      <!-- 自动化规则 -->
      <el-menu-item index="/rules">
        <el-icon><Setting /></el-icon>
        <template #title>{{ t('nav.rules') }}</template>
      </el-menu-item>

      <!-- Webhook -->
      <el-menu-item index="/webhooks">
        <el-icon><Link /></el-icon>
        <template #title>{{ t('nav.webhooks') }}</template>
      </el-menu-item>

      <!-- 附件管理 -->
      <el-menu-item index="/attachments">
        <el-icon><Paperclip /></el-icon>
        <template #title>{{ t('nav.attachments') }}</template>
      </el-menu-item>

      <!-- 对账 -->
      <el-menu-item index="/reconciliations">
        <el-icon><Histogram /></el-icon>
        <template #title>{{ t('nav.reconciliations') }}</template>
      </el-menu-item>

      <!-- 设置子菜单 -->
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
        <el-menu-item index="/settings/mfa">
          <el-icon><Key /></el-icon>
          {{ t('setting.mfa') }}
        </el-menu-item>
        <el-menu-item index="/settings/users">
          <el-icon><UserFilled /></el-icon>
          {{ t('setting.users') }}
        </el-menu-item>
        <el-menu-item index="/settings/preferences">
          <el-icon><Setting /></el-icon>
          {{ t('preference.title') }}
        </el-menu-item>
      </el-sub-menu>

      <!-- 对象分组 -->
      <el-menu-item index="/object-groups">
        <el-icon><Menu /></el-icon>
        <template #title>{{ t('objectGroup.title') }}</template>
      </el-menu-item>

      <!-- 交易关联 -->
      <el-menu-item index="/transaction-links">
        <el-icon><Link /></el-icon>
        <template #title>{{ t('transactionLink.title') }}</template>
      </el-menu-item>

      <!-- 定时任务 -->
      <el-menu-item index="/cron">
        <el-icon><Timer /></el-icon>
        <template #title>{{ t('nav.cron') }}</template>
      </el-menu-item>
    </el-menu>
  </el-aside>
</template>

<style scoped>
.app-sidebar {
  background: var(--app-sidebar-bg);
  border-right: 1px solid var(--app-border);
  overflow-y: auto;
  transition: width 0.3s;
}

.app-sidebar :deep(.el-menu) {
  border-right: none;
}
</style>

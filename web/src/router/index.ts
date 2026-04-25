// 路由配置 - 定义所有页面路由和导航守卫
import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// 公开路由 - 无需登录即可访问
const publicRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/auth/LoginPage.vue'), // 懒加载登录页
    meta: { requiresAuth: false }, // 不需要认证
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/pages/auth/RegisterPage.vue'), // 懒加载注册页
    meta: { requiresAuth: false },
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: () => import('@/pages/auth/ResetPasswordPage.vue'), // 懒加载重置密码页
    meta: { requiresAuth: false },
  },
]

// 受保护路由 - 需要登录才能访问
const protectedRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Dashboard',
    component: () => import('@/pages/dashboard/DashboardPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/accounts',
    name: 'AccountList',
    component: () => import('@/pages/accounts/AccountListPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/accounts/create',
    name: 'AccountCreate',
    component: () => import('@/pages/accounts/AccountFormPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/accounts/:id/edit',
    name: 'AccountEdit',
    component: () => import('@/pages/accounts/AccountFormPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/transactions',
    name: 'TransactionList',
    component: () => import('@/pages/transactions/TransactionListPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/transactions/create',
    name: 'TransactionCreate',
    component: () => import('@/pages/transactions/TransactionFormPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/transactions/:id/edit',
    name: 'TransactionEdit',
    component: () => import('@/pages/transactions/TransactionFormPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/categories',
    name: 'CategoryList',
    component: () => import('@/pages/categories/CategoryListPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/tags',
    name: 'TagList',
    component: () => import('@/pages/tags/TagListPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/budgets',
    name: 'BudgetList',
    component: () => import('@/pages/budgets/BudgetListPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/budgets/:id',
    name: 'BudgetDetail',
    component: () => import('@/pages/budgets/BudgetDetailPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/bills',
    name: 'BillList',
    component: () => import('@/pages/bills/BillListPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/reports/income-expense',
    name: 'IncomeExpenseReport',
    component: () => import('@/pages/reports/IncomeExpenseReport.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/reports/category',
    name: 'CategoryReport',
    component: () => import('@/pages/reports/CategoryReport.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/reports/budget',
    name: 'BudgetReport',
    component: () => import('@/pages/reports/BudgetReport.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/reports/net-worth',
    name: 'NetWorthReport',
    component: () => import('@/pages/reports/NetWorthReport.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/reports/trend',
    name: 'TrendReport',
    component: () => import('@/pages/reports/TrendReport.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/reports/tag',
    name: 'TagReport',
    component: () => import('@/pages/reports/TagReport.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/piggy-banks',
    name: 'PiggyBankList',
    component: () => import('@/pages/piggyBanks/PiggyBankListPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/recurring-transactions',
    name: 'RecurringTransactionList',
    component: () => import('@/pages/recurringTransactions/RecurringTransactionListPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/webhooks',
    name: 'WebhookList',
    component: () => import('@/pages/webhooks/WebhookListPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/exports',
    name: 'Export',
    component: () => import('@/pages/exports/ExportPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/attachments',
    name: 'AttachmentList',
    component: () => import('@/pages/attachments/AttachmentListPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/reconciliations',
    name: 'ReconciliationList',
    component: () => import('@/pages/reconciliations/ReconciliationListPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/imports',
    name: 'Import',
    component: () => import('@/pages/imports/ImportPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/rules',
    name: 'RuleList',
    component: () => import('@/pages/rules/RuleListPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/settings/profile',
    name: 'ProfileSettings',
    component: () => import('@/pages/settings/ProfilePage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/settings/currencies',
    name: 'CurrencySettings',
    component: () => import('@/pages/settings/CurrencySettingsPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/settings/mfa',
    name: 'MFASettings',
    component: () => import('@/pages/settings/MFASettingsPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/settings/users',
    name: 'UserManagement',
    component: () => import('@/pages/settings/UserManagementPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/settings/preferences',
    name: 'PreferenceSettings',
    component: () => import('@/pages/settings/PreferencesPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/object-groups',
    name: 'ObjectGroupList',
    component: () => import('@/pages/objectGroups/ObjectGroupListPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/transaction-links',
    name: 'TransactionLinkList',
    component: () => import('@/pages/transactionLinks/TransactionLinkListPage.vue'),
    meta: { requiresAuth: true },
  },
  { path: '/:pathMatch(.*)*', name: 'NotFound', redirect: '/dashboard' },
]

// 创建路由实例，使用HTML5 History模式
const router = createRouter({
  history: createWebHistory(),
  routes: [...publicRoutes, ...protectedRoutes], // 合并公开和受保护路由
})

// 全局前置守卫 - 控制页面访问权限
router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.meta.requiresAuth !== false // 默认需要认证

  if (requiresAuth && !authStore.isAuthenticated) {
    // 需要认证但未登录，跳转到登录页并记录原始路径
    next({ path: '/login', query: { redirect: to.fullPath } })
  } else if (!requiresAuth && authStore.isAuthenticated && to.path !== '/') {
    // 已登录用户访问公开页面，重定向到首页
    next({ path: '/' })
  } else {
    next()
  }
})

export default router

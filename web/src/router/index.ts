import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const publicRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/auth/LoginPage.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/pages/auth/RegisterPage.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: () => import('@/pages/auth/ResetPasswordPage.vue'),
    meta: { requiresAuth: false },
  },
]

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
]

const router = createRouter({
  history: createWebHistory(),
  routes: [...publicRoutes, ...protectedRoutes],
})

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.meta.requiresAuth !== false

  if (requiresAuth && !authStore.isAuthenticated) {
    next({ path: '/login', query: { redirect: to.fullPath } })
  } else if (!requiresAuth && authStore.isAuthenticated && to.path !== '/') {
    next({ path: '/' })
  } else {
    next()
  }
})

export default router

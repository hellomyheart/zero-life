// 认证状态管理 - 管理用户登录状态、Token、个人信息
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { LoginReq, RegisterReq, LoginResp, ProfileResp, UpdateProfileReq, ChangePasswordReq } from '@/types/auth'
import { login as loginApi, register as registerApi, getProfile, updateProfile as updateProfileApi, changePassword as changePasswordApi } from '@/api/auth'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('token') || '') // 访问令牌
  const refreshToken = ref<string>(localStorage.getItem('refreshToken') || '') // 刷新令牌
  const user = ref<ProfileResp | null>(null) // 当前登录用户信息

  // 是否已认证（有Token即为已登录）
  const isAuthenticated = computed(() => !!token.value)

  // 设置Token并持久化到localStorage
  function setTokens(accessToken: string, refresh: string) {
    token.value = accessToken
    refreshToken.value = refresh
    localStorage.setItem('token', accessToken)
    localStorage.setItem('refreshToken', refresh)
  }

  // 清除Token和用户信息
  function clearTokens() {
    token.value = ''
    refreshToken.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('refreshToken')
  }

  // 登录 - 调用API获取Token并加载用户信息
  async function login(data: LoginReq) {
    const res = await loginApi(data) as unknown as LoginResp
    setTokens(res.access_token, res.refresh_token)
    await loadProfile()
  }

  // 注册 - 调用API注册并自动登录
  async function register(data: RegisterReq) {
    const res = await registerApi(data) as unknown as LoginResp
    setTokens(res.access_token, res.refresh_token)
    await loadProfile()
  }

  // 加载当前用户信息
  async function loadProfile() {
    try {
      user.value = await getProfile() as unknown as ProfileResp
    } catch {
      clearTokens()
    }
  }

  // 更新用户资料
  async function updateProfile(data: UpdateProfileReq) {
    user.value = await updateProfileApi(data) as unknown as ProfileResp
  }

  // 修改密码
  async function changePassword(data: ChangePasswordReq) {
    await changePasswordApi(data)
  }

  // 退出登录 - 清除Token并跳转到登录页
  function logout() {
    clearTokens()
    router.push('/login')
  }

  return {
    token,
    refreshToken,
    user,
    isAuthenticated,
    setTokens,
    login,
    register,
    loadProfile,
    updateProfile,
    changePassword,
    logout,
  }
})

/**
 * 认证状态管理（Pinia Store）
 *
 * 组件功能：
 * - 管理用户登录状态、访问令牌和刷新令牌
 * - 提供登录、注册、退出登录等认证操作
 * - 管理当前登录用户的个人信息
 *
 * 业务流程：
 * - 登录/注册 → 获取Token → 存储到localStorage → 加载用户信息
 * - 退出登录 → 清除Token → 跳转登录页
 * - Token过期 → 请求拦截器自动刷新Token（见utils/request.ts）
 *
 * 数据流：
 * - Token持久化到localStorage，页面刷新后自动恢复登录状态
 * - 用户信息在登录成功后自动加载，失败则清除Token
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { LoginReq, RegisterReq, LoginResp, ProfileResp, UpdateProfileReq, ChangePasswordReq } from '@/types/auth'
import { login as loginApi, register as registerApi, getProfile, updateProfile as updateProfileApi, changePassword as changePasswordApi } from '@/api/auth'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  /** 访问令牌，用于API请求的身份验证 */
  const token = ref<string>(localStorage.getItem('token') || '')
  /** 刷新令牌，用于获取新的访问令牌 */
  const refreshToken = ref<string>(localStorage.getItem('refreshToken') || '')
  /** 当前登录用户信息，未登录时为null */
  const user = ref<ProfileResp | null>(null)

  /** 是否已认证（有Token即为已登录） */
  const isAuthenticated = computed(() => !!token.value)

  /**
   * 设置Token并持久化到localStorage
   * @param accessToken - 新的访问令牌
   * @param refresh - 新的刷新令牌
   */
  function setTokens(accessToken: string, refresh: string) {
    token.value = accessToken
    refreshToken.value = refresh
    localStorage.setItem('token', accessToken)
    localStorage.setItem('refreshToken', refresh)
  }

  /** 清除Token和用户信息，用于退出登录或Token失效时 */
  function clearTokens() {
    token.value = ''
    refreshToken.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('refreshToken')
  }

  /**
   * 用户登录
   * 调用登录API获取Token，存储后自动加载用户信息
   * @param data - 登录请求参数（邮箱、密码）
   */
  async function login(data: LoginReq) {
    const res = await loginApi(data) as unknown as LoginResp
    setTokens(res.access_token, res.refresh_token)
    await loadProfile()
  }

  /**
   * 用户注册
   * 调用注册API创建账户，成功后自动登录
   * @param data - 注册请求参数（邮箱、密码、昵称）
   */
  async function register(data: RegisterReq) {
    const res = await registerApi(data) as unknown as LoginResp
    setTokens(res.access_token, res.refresh_token)
    await loadProfile()
  }

  /**
   * 加载当前用户信息
   * 获取失败时不清除Token（Token可能仍然有效，只是profile接口暂时不可用）
   * 只有明确收到401且token刷新也失败时，才清除Token（由响应拦截器处理）
   */
  async function loadProfile() {
    try {
      user.value = await getProfile() as unknown as ProfileResp
    } catch {
      // 不清除Token，只标记用户信息为null
      // Token有效性由响应拦截器的401处理逻辑判断
      user.value = null
    }
  }

  /**
   * 更新用户资料
   * @param data - 更新内容（昵称、语言、时区等）
   */
  async function updateProfile(data: UpdateProfileReq) {
    user.value = await updateProfileApi(data) as unknown as ProfileResp
  }

  /**
   * 修改密码
   * @param data - 密码修改信息（旧密码 + 新密码）
   */
  async function changePassword(data: ChangePasswordReq) {
    await changePasswordApi(data)
  }

  /** 退出登录 - 清除Token和用户信息，跳转到登录页 */
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

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { LoginReq, RegisterReq, LoginResp, ProfileResp, UpdateProfileReq, ChangePasswordReq } from '@/types/auth'
import { login as loginApi, register as registerApi, refresh as refreshApi, getProfile, updateProfile as updateProfileApi, changePassword as changePasswordApi } from '@/api/auth'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const refreshToken = ref<string>(localStorage.getItem('refreshToken') || '')
  const user = ref<ProfileResp | null>(null)

  const isAuthenticated = computed(() => !!token.value)

  function setTokens(accessToken: string, refresh: string) {
    token.value = accessToken
    refreshToken.value = refresh
    localStorage.setItem('token', accessToken)
    localStorage.setItem('refreshToken', refresh)
  }

  function clearTokens() {
    token.value = ''
    refreshToken.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('refreshToken')
  }

  async function login(data: LoginReq) {
    const res = await loginApi(data) as unknown as LoginResp
    setTokens(res.access_token, res.refresh_token)
    await loadProfile()
  }

  async function register(data: RegisterReq) {
    const res = await registerApi(data) as unknown as LoginResp
    setTokens(res.access_token, res.refresh_token)
    await loadProfile()
  }

  async function loadProfile() {
    try {
      user.value = await getProfile() as unknown as ProfileResp
    } catch {
      clearTokens()
    }
  }

  async function updateProfile(data: UpdateProfileReq) {
    user.value = await updateProfileApi(data) as unknown as ProfileResp
  }

  async function changePassword(data: ChangePasswordReq) {
    await changePasswordApi(data)
  }

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

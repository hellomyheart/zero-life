// Axios HTTP请求封装 - 统一处理请求拦截、Token刷新、错误处理
import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig, InternalAxiosRequestConfig } from 'axios'
import type { ApiResponse } from '@/types/common'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

// 创建Axios实例，配置基础URL和超时时间
const service: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1', // API基础路径，从环境变量读取
  timeout: 30000, // 请求超时时间30秒
  headers: {
    'Content-Type': 'application/json',
  },
})

// Token刷新相关状态
let isRefreshing = false // 是否正在刷新Token
let pendingRequests: Array<(token: string) => void> = [] // 等待Token刷新的请求队列

// Token刷新完成后，依次执行队列中的请求
function onTokenRefreshed(token: string) {
  pendingRequests.forEach((cb) => cb(token))
  pendingRequests = []
}

// 请求拦截器 - 自动附加Token到请求头
service.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}` // 添加Bearer Token
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理业务错误码和401自动刷新Token
service.interceptors.response.use(
  (response) => {
    const res = response.data as ApiResponse<unknown>
    // 业务错误码不为0时，视为请求失败
    if (res.code && res.code !== 0) {
      return Promise.reject(new Error(res.message || 'Request failed'))
    }
    // 返回data字段，如果没有data则返回整个响应
    return (res.data !== undefined ? res.data : res) as typeof response.data
  },
  async (error) => {
    const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean }

    // 401未授权错误处理 - 尝试刷新Token
    if (error.response?.status === 401 && !originalRequest._retry) {
      const authStore = useAuthStore()

      // 没有refreshToken，直接登出
      if (!authStore.refreshToken) {
        authStore.logout()
        router.push('/login')
        return Promise.reject(error)
      }

      // 正在刷新Token时，将请求加入等待队列
      if (isRefreshing) {
        return new Promise((resolve) => {
          pendingRequests.push((token: string) => {
            if (originalRequest.headers) {
              originalRequest.headers.Authorization = `Bearer ${token}`
            }
            resolve(service(originalRequest)) // 用新Token重新发送请求
          })
        })
      }

      originalRequest._retry = true // 标记已重试，防止无限循环
      isRefreshing = true

      try {
        // 调用刷新Token接口
        const resp = await axios.post(
          `${import.meta.env.VITE_API_BASE_URL || '/api/v1'}/auth/refresh`,
          { refresh_token: authStore.refreshToken }
        )
        const data = resp.data as ApiResponse<{ access_token: string; refresh_token: string }>
        const newAccessToken = data.data.access_token
        const newRefreshToken = data.data.refresh_token

        // 更新存储中的Token
        authStore.setTokens(newAccessToken, newRefreshToken)
        // 通知等待队列中的请求
        onTokenRefreshed(newAccessToken)

        // 用新Token重新发送原始请求
        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${newAccessToken}`
        }
        return service(originalRequest)
      } catch {
        // 刷新Token失败，强制登出
        authStore.logout()
        router.push('/login')
        return Promise.reject(error)
      } finally {
        isRefreshing = false
      }
    }

    // 其他错误，提取错误信息
    const message = error.response?.data?.message || error.message || 'Network Error'
    return Promise.reject(new Error(message))
  }
)

// 封装的HTTP方法 - 简化API调用
export function get<T = unknown>(url: string, params?: Record<string, unknown>, config?: AxiosRequestConfig): Promise<T> {
  return service.get(url, { params, ...config }) as Promise<T>
}

export function post<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
  return service.post(url, data, config) as Promise<T>
}

export function put<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
  return service.put(url, data, config) as Promise<T>
}

export function del<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T> {
  return service.delete(url, config) as Promise<T>
}

export default service

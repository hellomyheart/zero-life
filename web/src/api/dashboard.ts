import { get } from '@/utils/request'
import type { DashboardResp } from '@/types/dashboard'

export function getDashboard() {
  return get<DashboardResp>('/dashboard')
}

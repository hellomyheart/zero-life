// 仪表盘相关API接口 - 首页汇总数据
import { get } from '@/utils/request'
import type { DashboardResp } from '@/types/dashboard'

export function getDashboard() {
  return get<DashboardResp>('/dashboard')
}

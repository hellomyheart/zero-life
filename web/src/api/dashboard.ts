/**
 * 仪表盘相关API接口
 * 提供首页汇总数据的获取功能，包括收支概览、预算预警、账单提醒和最近交易
 */
import { get } from '@/utils/request'
import type { DashboardResp } from '@/types/dashboard'

/**
 * 获取仪表盘汇总数据
 * 返回首页展示所需的所有汇总信息
 * @returns 仪表盘数据（总收入、总支出、净收入、总余额、预算预警、账单提醒、最近交易）
 * @endpoint GET /dashboard
 */
export function getDashboard() {
  return get<DashboardResp>('/dashboard')
}

/**
 * 汇率管理 API 接口
 * 提供汇率的增删改查和货币转换功能
 */
import request from '@/utils/request'

/**
 * 获取汇率列表
 */
export function list(params: any) {
  return request.get('/exchange-rates', { params })
}

/**
 * 获取汇率详情
 */
export function get(id: string) {
  return request.get(`/exchange-rates/${id}`)
}

/**
 * 创建汇率
 */
export function create(data: any) {
  return request.post('/exchange-rates', data)
}

/**
 * 更新汇率
 */
export function update(id: string, data: any) {
  return request.put(`/exchange-rates/${id}`, data)
}

/**
 * 删除汇率
 */
export function remove(id: string) {
  return request.delete(`/exchange-rates/${id}`)
}

/**
 * 获取最新汇率
 */
export function getLatest(from: string, to: string) {
  return request.get(`/exchange-rates/latest/${from}/${to}`)
}

/**
 * 货币转换
 */
export function convert(from: string, to: string, amount: string, date?: string) {
  const params: any = { amount }
  if (date) params.date = date
  return request.get(`/exchange-rates/convert/${from}/${to}`, { params })
}

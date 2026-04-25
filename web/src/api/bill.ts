/**
 * 账单相关API接口
 * 提供账单的增删改查功能，账单用于管理定期支付的固定支出（如房租、水电费、订阅服务等）
 * 账单会根据重复规则自动提醒下次到期时间
 */
import { get, post, put, del } from '@/utils/request'
import type { Bill, CreateBillReq, UpdateBillReq } from '@/types/bill'

/**
 * 获取所有账单列表
 * @returns 账单数组
 * @endpoint GET /bills
 */
export function list() {
  return get<Bill[]>('/bills')
}

/**
 * 获取单个账单详情
 * @param id - 账单ID
 * @returns 账单详细信息（含下次到期日和是否逾期）
 * @endpoint GET /bills/:id
 */
export function getBill(id: number) {
  return get<Bill>(`/bills/${id}`)
}

/**
 * 创建新账单
 * @param data - 创建账单请求参数（名称、金额、关联账户、关联分类、重复规则、下次到期日）
 * @returns 创建成功的账单信息
 * @endpoint POST /bills
 */
export function create(data: CreateBillReq) {
  return post<Bill>('/bills', data)
}

/**
 * 更新账单信息
 * @param id - 要更新的账单ID
 * @param data - 更新内容（名称、金额、关联账户、关联分类、重复规则、到期日，均为可选）
 * @returns 更新后的账单信息
 * @endpoint PUT /bills/:id
 */
export function update(id: number, data: UpdateBillReq) {
  return put<Bill>(`/bills/${id}`, data)
}

/**
 * 删除账单
 * @param id - 要删除的账单ID
 * @endpoint DELETE /bills/:id
 */
export function remove(id: number) {
  return del<void>(`/bills/${id}`)
}

/**
 * 存钱罐相关API接口
 * 提供存钱罐的增删改查、存取款操作和事件记录查询功能
 * 存钱罐用于设定储蓄目标，跟踪存款进度（如旅行基金、应急储备金等）
 */
import { get, post, put, del } from '@/utils/request'
import type { PiggyBank, CreatePiggyBankReq, UpdatePiggyBankReq, AddAmountReq, RemoveAmountReq } from '@/types/piggyBank'
import type { PiggyEvent } from '@/types/piggyBank'

/**
 * 获取存钱罐列表
 * 后端返回扁平数组（非分页）
 * @returns 存钱罐数组
 * @endpoint GET /piggy-banks
 */
export function list() {
  return get<PiggyBank[]>('/piggy-banks')
}

/**
 * 获取单个存钱罐详情
 * @param id - 存钱罐ID
 * @returns 存钱罐详细信息（含当前金额和完成百分比）
 * @endpoint GET /piggy-banks/:id
 */
export function getPiggyBank(id: number) { return get<PiggyBank>(`/piggy-banks/${id}`) }

/**
 * 创建新存钱罐
 * @param data - 创建请求参数（名称、关联账户ID、目标金额、可选目标日期和备注）
 * @returns 创建成功的存钱罐信息
 * @endpoint POST /piggy-banks
 */
export function create(data: CreatePiggyBankReq) { return post<PiggyBank>('/piggy-banks', data) }

/**
 * 更新存钱罐信息
 * @param id - 存钱罐ID
 * @param data - 更新内容（名称、目标金额、目标日期、备注，均为可选）
 * @returns 更新后的存钱罐信息
 * @endpoint PUT /piggy-banks/:id
 */
export function update(id: number, data: UpdatePiggyBankReq) { return put<PiggyBank>(`/piggy-banks/${id}`, data) }

/**
 * 删除存钱罐
 * @param id - 存钱罐ID
 * @endpoint DELETE /piggy-banks/:id
 */
export function remove(id: number) { return del<void>(`/piggy-banks/${id}`) }

/**
 * 向存钱罐存入金额
 * 从关联账户中转入指定金额到存钱罐
 * @param id - 存钱罐ID
 * @param data - 存入请求参数（金额、可选备注）
 * @endpoint POST /piggy-banks/:id/add
 */
export function addAmount(id: number, data: AddAmountReq) { return post<void>(`/piggy-banks/${id}/add`, data) }

/**
 * 从存钱罐取出金额
 * 从存钱罐中转出指定金额回关联账户
 * @param id - 存钱罐ID
 * @param data - 取出请求参数（金额、可选备注）
 * @endpoint POST /piggy-banks/:id/remove
 */
export function removeAmount(id: number, data: RemoveAmountReq) { return post<void>(`/piggy-banks/${id}/remove`, data) }

/**
 * 获取存钱罐的事件记录
 * 返回该存钱罐的所有存取款操作历史
 * @param id - 存钱罐ID
 * @returns 事件记录数组
 * @endpoint GET /piggy-banks/:id/events
 */
export function getEvents(id: number) { return get<PiggyEvent[]>(`/piggy-banks/${id}/events`) }

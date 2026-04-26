/**
 * Webhook管理API接口
 * 提供Webhook的增删改查和投递记录查询功能
 * Webhook用于在特定事件发生时（如创建交易、账单到期等）自动向外部URL发送HTTP通知，
 * 实现系统与其他服务的集成
 */
import { get, post, put, del } from '@/utils/request'
import type { Webhook, CreateWebhookReq, UpdateWebhookReq, WebhookDelivery } from '@/types/webhook'

/**
 * 获取Webhook列表
 * 后端返回扁平数组（非分页）
 * @returns Webhook数组
 * @endpoint GET /webhooks
 */
export function list() {
  return get<Webhook[]>('/webhooks')
}

/**
 * 获取单个Webhook详情
 * @param id - Webhook ID
 * @returns Webhook详细信息
 * @endpoint GET /webhooks/:id
 */
export function getWebhook(id: number) { return get<Webhook>(`/webhooks/${id}`) }

/**
 * 创建新Webhook
 * @param data - 创建请求参数（名称、回调URL、触发事件类型）
 * @returns 创建成功的Webhook信息
 * @endpoint POST /webhooks
 */
export function create(data: CreateWebhookReq) { return post<Webhook>('/webhooks', data) }

/**
 * 更新Webhook信息
 * @param id - Webhook ID
 * @param data - 更新内容（名称、回调URL、触发事件、是否启用，均为可选）
 * @returns 更新后的Webhook信息
 * @endpoint PUT /webhooks/:id
 */
export function update(id: number, data: UpdateWebhookReq) { return put<Webhook>(`/webhooks/${id}`, data) }

/**
 * 删除Webhook
 * @param id - Webhook ID
 * @endpoint DELETE /webhooks/:id
 */
export function remove(id: number) { return del<void>(`/webhooks/${id}`) }

/**
 * 获取Webhook的投递记录
 * 返回该Webhook所有历史通知发送记录，包含HTTP状态码和是否成功
 * @param id - Webhook ID
 * @returns 投递记录分页列表
 * @endpoint GET /webhooks/:id/deliveries
 */
export function listDeliveries(id: number) { return get<WebhookDelivery[]>(`/webhooks/${id}/deliveries`) }

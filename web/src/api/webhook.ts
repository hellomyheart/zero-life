import { get, post, put, del } from '@/utils/request'
import type { Webhook, CreateWebhookReq, UpdateWebhookReq, WebhookListReq, WebhookDelivery } from '@/types/webhook'
import type { PageResult } from '@/types/common'

export function list(params: WebhookListReq) {
  return get<PageResult<Webhook>>('/webhooks', params as Record<string, unknown>)
}
export function getWebhook(id: string) { return get<Webhook>(`/webhooks/${id}`) }
export function create(data: CreateWebhookReq) { return post<Webhook>('/webhooks', data) }
export function update(id: string, data: UpdateWebhookReq) { return put<Webhook>(`/webhooks/${id}`, data) }
export function remove(id: string) { return del<void>(`/webhooks/${id}`) }
export function listDeliveries(id: string) { return get<PageResult<WebhookDelivery>>(`/webhooks/${id}/deliveries`) }

// Webhook相关类型定义 - Webhook接口、发送事件和消息/响应
export interface Webhook {
  id: string
  name: string
  url: string
  trigger: string
  active: boolean
  secret: string
  last_delivered_at: string
  created_at: string
  updated_at: string
}

export interface WebhookDelivery {
  id: string
  webhook_id: string
  url: string
  status_code: number
  success: boolean
  error_message: string
  created_at: string
}

export interface CreateWebhookReq {
  name: string
  url: string
  trigger: string
}

export interface UpdateWebhookReq {
  name?: string
  url?: string
  trigger?: string
  active?: boolean
}

export interface WebhookListReq {
  page?: number
  page_size?: number
}

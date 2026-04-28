/**
 * Webhook相关类型定义
 * 字段名与后端 JSON tag 完全对应
 */

/**
 * Webhook触发事件类型（对应后端 oneof 验证）
 */
export enum WebhookTrigger {
  TransactionCreated = 'transaction.created',
  TransactionUpdated = 'transaction.updated',
  TransactionDeleted = 'transaction.deleted',
  RecurringTransactionExecuted = 'recurring_transaction.executed',
  BudgetCreated = 'budget.created',
  BudgetUpdated = 'budget.updated',
  BudgetDeleted = 'budget.deleted',
}

/**
 * Webhook信息接口（对应后端 WebhookResp）
 */
export interface Webhook {
  id: number
  name: string
  url: string
  trigger: WebhookTrigger
  is_active: boolean
  last_delivered_at: string | null
  created_at: string
  updated_at: string
}

/**
 * Webhook投递记录接口（对应后端 WebhookDeliveryResp）
 */
export interface WebhookDelivery {
  id: number
  webhook_id: number
  status_code: number | null
  error_message: string | null
  delivered_at: string | null
  created_at: string
}

/**
 * 创建Webhook请求接口（对应后端 CreateWebhookReq）
 */
export interface CreateWebhookReq {
  name: string
  url: string
  trigger: WebhookTrigger
}

/**
 * 更新Webhook请求接口（对应后端 UpdateWebhookReq）
 * is_active 字段名与后端一致（非 active）
 */
export interface UpdateWebhookReq {
  name?: string
  url?: string
  trigger?: WebhookTrigger
  is_active?: boolean | null
}


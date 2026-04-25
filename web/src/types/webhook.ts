/**
 * Webhook相关类型定义
 * Webhook用于在特定事件发生时自动向外部URL发送HTTP通知，实现系统与其他服务的集成
 */

/**
 * Webhook信息接口
 * 表示一个Webhook配置
 */
export interface Webhook {
  id: number                    // Webhook唯一标识
  name: string                 // Webhook名称，如"Slack通知"
  url: string                  // 回调URL，事件触发时向此地址发送POST请求
  trigger: string              // 触发事件类型，如"transaction.created"、"bill.due"
  active: boolean              // 是否启用
  secret: string               // 签名密钥，用于验证回调请求的来源真实性
  last_delivered_at: string    // 最后一次成功投递的时间
  created_at: string           // 创建时间
  updated_at: string           // 最后更新时间
}

/**
 * Webhook投递记录接口
 * 记录每次Webhook通知的发送结果
 */
export interface WebhookDelivery {
  id: number                    // 投递记录唯一标识
  webhook_id: number            // 关联的Webhook ID
  url: string                  // 请求的目标URL
  status_code: number          // HTTP响应状态码（如200、404、500等）
  success: boolean             // 是否投递成功（状态码2xx视为成功）
  error_message: string        // 错误信息（投递失败时记录原因）
  created_at: string           // 投递时间
}

/**
 * 创建Webhook请求接口
 */
export interface CreateWebhookReq {
  name: string                 // Webhook名称（必填）
  url: string                  // 回调URL（必填）
  trigger: string              // 触发事件类型（必填）
}

/**
 * 更新Webhook请求接口
 * 所有字段均为可选，只传需要修改的字段
 */
export interface UpdateWebhookReq {
  name?: string                // Webhook名称
  url?: string                 // 回调URL
  trigger?: string             // 触发事件类型
  active?: boolean             // 是否启用
}

/**
 * Webhook列表查询请求接口
 */
export interface WebhookListReq {
  page?: number                // 页码（从1开始）
  page_size?: number           // 每页数量
}

/**
 * 附件相关类型定义
 * 字段名与后端 JSON tag 完全对应
 */

/**
 * 附件信息接口（对应后端 AttachmentResp）
 * 后端不返回 title 和 notes 字段
 */
export interface Attachment {
  id: number
  attachable_type: string
  attachable_id: number
  filename: string
  mime: string
  size: number
  created_at: string
}

/**
 * 附件列表查询请求接口
 */
export interface AttachmentListReq {
  attachable_type?: string
  attachable_id?: number
  page?: number
  page_size?: number
}

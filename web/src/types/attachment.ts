// 附件相关类型定义 - 附件接口、列表类型
export interface Attachment {
  id: string
  attachable_type: string
  attachable_id: string
  filename: string
  mime: string
  size: number
  title: string
  notes: string
  created_at: string
  updated_at: string
}

export interface AttachmentListReq {
  attachable_type?: string
  attachable_id?: string
  page?: number
  page_size?: number
}

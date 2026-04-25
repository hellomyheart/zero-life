/**
 * 附件相关类型定义
 * 附件用于关联到交易、账单等实体上，保存票据、收据、合同等文件
 */

/**
 * 附件信息接口
 * 表示系统中已上传的一个附件文件
 */
export interface Attachment {
  id: number                    // 附件唯一标识
  attachable_type: string       // 关联实体类型，如"transaction"、"bill"等
  attachable_id: number         // 关联实体的ID
  filename: string             // 文件名
  mime: string                 // 文件MIME类型，如"image/png"、"application/pdf"
  size: number                 // 文件大小（字节）
  title: string                // 附件标题
  notes: string                // 附件备注说明
  created_at: string           // 上传时间
  updated_at: string           // 最后更新时间
}

/**
 * 附件列表查询请求接口
 * 支持按关联实体筛选和分页
 */
export interface AttachmentListReq {
  attachable_type?: string     // 按关联实体类型筛选（如"transaction"）
  attachable_id?: number        // 按关联实体ID筛选
  page?: number                // 页码（从1开始）
  page_size?: number           // 每页数量
}

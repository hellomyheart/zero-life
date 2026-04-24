/**
 * 附件相关API接口
 * 提供文件的上传、下载、预览和删除功能
 * 附件可以关联到交易、账单等实体上，用于保存票据、收据等文件
 */
import { get, post, del } from '@/utils/request'
import type { Attachment, AttachmentListReq } from '@/types/attachment'
import type { PageResult } from '@/types/common'

/**
 * 获取附件列表（分页）
 * 支持按关联实体类型和ID进行筛选
 * @param params - 查询参数（可选：关联实体类型、关联实体ID、页码、每页数量）
 * @returns 附件分页列表
 * @endpoint GET /attachments
 */
export function list(params: AttachmentListReq) {
  return get<PageResult<Attachment>>('/attachments', params as Record<string, unknown>)
}

/**
 * 上传附件
 * 使用 FormData 上传文件，Content-Type 为 multipart/form-data
 * @param formData - 包含文件的 FormData 对象
 * @returns 上传成功的附件信息
 * @endpoint POST /attachments/upload
 */
export function upload(formData: FormData) { return post<Attachment>('/attachments/upload', formData) }

/**
 * 下载附件
 * 以文件流方式下载指定附件
 * @param id - 附件ID
 * @endpoint GET /attachments/:id/download
 */
export function download(id: string) { return get<void>(`/attachments/${id}/download`) }

/**
 * 在线预览附件
 * 以文件流方式在浏览器中直接预览附件内容（如图片、PDF等）
 * @param id - 附件ID
 * @endpoint GET /attachments/:id/view
 */
export function view(id: string) { return get<void>(`/attachments/${id}/view`) }

/**
 * 删除附件
 * @param id - 附件ID
 * @endpoint DELETE /attachments/:id
 */
export function remove(id: string) { return del<void>(`/attachments/${id}`) }

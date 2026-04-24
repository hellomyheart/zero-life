import { get, post, del } from '@/utils/request'
import type { Attachment, AttachmentListReq } from '@/types/attachment'
import type { PageResult } from '@/types/common'

export function list(params: AttachmentListReq) {
  return get<PageResult<Attachment>>('/attachments', params as Record<string, unknown>)
}
export function upload(formData: FormData) { return post<Attachment>('/attachments/upload', formData) }
export function download(id: string) { return get<void>(`/attachments/${id}/download`) }
export function view(id: string) { return get<void>(`/attachments/${id}/view`) }
export function remove(id: string) { return del<void>(`/attachments/${id}`) }

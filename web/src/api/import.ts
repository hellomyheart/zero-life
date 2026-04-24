// 数据导入API接口 - CSV文件上传、映射和导入
import { post } from '@/utils/request'
import type { AxiosRequestConfig } from 'axios'

export function upload(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return post<unknown>('/imports/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  } as AxiosRequestConfig)
}

export function parse(data: { file_id: string; column_mapping: Record<string, string> }) {
  return post<unknown>('/imports/parse', data)
}

export function execute(data: { file_id: string; column_mapping: Record<string, string> }) {
  return post<unknown>('/imports/execute', data)
}

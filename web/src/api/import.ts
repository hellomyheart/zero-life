/**
 * 数据导入API接口
 * 提供CSV文件上传、列映射解析和执行导入的功能
 * 导入流程：1.上传文件 → 2.解析并映射列 → 3.执行导入
 */
import { post } from '@/utils/request'
import type { AxiosRequestConfig } from 'axios'

/**
 * 上传待导入的文件
 * 将CSV等格式的文件上传到服务器，返回文件ID用于后续解析
 * @param file - 要上传的文件对象
 * @returns 上传结果（包含文件ID等信息）
 * @endpoint POST /imports/upload
 */
export function upload(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return post<unknown>('/imports/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  } as AxiosRequestConfig)
}

/**
 * 解析已上传的文件
 * 根据列映射配置解析文件内容，预览数据以便确认映射是否正确
 * @param data - 解析请求参数（文件ID + 列映射配置，如 {"date": "A列", "amount": "B列"}）
 * @returns 解析结果（预览数据）
 * @endpoint POST /imports/parse
 */
export function parse(data: { file_id: string; mapping: Record<string, string> }) {
  return post<unknown>('/imports/parse', data)
}

/**
 * 执行数据导入
 * 确认列映射无误后，将文件数据正式导入到系统中
 * @param data - 导入请求参数（文件ID + 列映射配置）
 * @returns 导入结果
 * @endpoint POST /imports/execute
 */
export function execute(data: { file_id: string; mapping: Record<string, string> }) {
  return post<unknown>('/imports/execute', data)
}

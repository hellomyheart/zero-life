/**
 * 数据导入相关类型定义
 * 字段名与后端 JSON tag 完全对应
 */

/**
 * 导入上传响应接口（对应后端 ImportUploadResp）
 * 后端仅返回 file_id，不返回 columns
 */
export interface ImportUploadResp {
  file_id: string
}

/**
 * 导入预览响应接口（对应后端 ImportPreviewResp）
 */
export interface ImportPreviewResp {
  total: number
  valid: number
  invalid: number
  rows: ImportRowResp[]
}

/**
 * 导入行预览响应接口（对应后端 ImportRowResp）
 */
export interface ImportRowResp {
  index: number
  data: Record<string, string>
  is_valid: boolean
  errors?: string[]
}

/**
 * 导入结果响应接口（对应后端 ImportResultResp）
 */
export interface ImportResultResp {
  total: number
  success: number
  failed: number
  skipped: number
}

/**
 * 导入解析请求接口（对应后端 ImportParseReq）
 * 后端字段名为 mapping（非 column_mapping）
 */
export interface ImportParseReq {
  file_id: string
  mapping: Record<string, string>
}

/**
 * 导入执行请求接口（对应后端 ImportExecuteReq）
 */
export interface ImportExecuteReq {
  file_id: string
  mapping: Record<string, string>
}

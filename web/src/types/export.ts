/**
 * 数据导出相关类型定义
 * 定义数据导出时的配置参数
 */

/**
 * 导出请求接口
 * 通用导出参数，所有导出接口共用
 */
export interface ExportReq {
  start_date?: string          // 起始日期（可选，筛选该日期之后的数据）
  end_date?: string            // 结束日期（可选，筛选该日期之前的数据）
  format?: 'csv' | 'json'     // 导出格式（可选，默认csv）
}

// 数据导出相关类型定义 - 导出配置类型
export interface ExportReq {
  start_date?: string
  end_date?: string
  format?: 'csv' | 'json'
}

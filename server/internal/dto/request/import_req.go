// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// ImportParseReq 导入文件解析请求
// 用于预览上传文件的内容，确认字段映射是否正确
type ImportParseReq struct {
	FileID    string            `json:"file_id" binding:"required"` // 上传文件后返回的文件ID
	Mapping   map[string]string `json:"mapping" binding:"required"` // 字段映射，key为CSV列名，value为系统字段名
}

// ImportExecuteReq 执行导入请求
// 确认字段映射后正式执行数据导入
type ImportExecuteReq struct {
	FileID    string            `json:"file_id" binding:"required"` // 上传文件后返回的文件ID
	Mapping   map[string]string `json:"mapping" binding:"required"` // 字段映射
}

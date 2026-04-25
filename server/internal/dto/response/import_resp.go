// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

// ImportUploadResp 导入文件上传响应
type ImportUploadResp struct {
	FileID string `json:"file_id"` // 文件ID，用于后续解析和导入
}

// ImportPreviewResp 导入预览响应
// 返回文件解析后的预览数据，供用户确认字段映射
type ImportPreviewResp struct {
	Total   int             `json:"total"`   // 总行数
	Valid   int             `json:"valid"`   // 有效行数
	Invalid int             `json:"invalid"` // 无效行数
	Rows    []ImportRowResp `json:"rows"`    // 每行预览数据
}

// ImportRowResp 导入行预览响应
type ImportRowResp struct {
	Index   int              `json:"index"`             // 行号
	Data    map[string]string `json:"data"`             // 行数据
	IsValid bool             `json:"is_valid"`          // 是否有效
	Errors  []string         `json:"errors,omitempty"`  // 验证错误信息
}

// ImportResultResp 导入执行结果响应
type ImportResultResp struct {
	Total   int `json:"total"`   // 总处理行数
	Success int `json:"success"` // 成功导入行数
	Failed  int `json:"failed"`  // 失败行数
	Skipped int `json:"skipped"` // 跳过行数
}

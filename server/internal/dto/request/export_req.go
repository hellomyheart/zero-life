// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// ExportReq 数据导出请求
// 用于指定导出的日期范围和文件格式
type ExportReq struct {
	StartDate string `form:"start_date"`                            // 开始日期
	EndDate   string `form:"end_date"`                              // 结束日期
	Format    string `form:"format" binding:"required,oneof=csv json"` // 导出格式：csv或json
}

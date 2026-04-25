// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// ReportReq 报表查询请求
// 用于指定报表的时间范围和数据粒度
type ReportReq struct {
	StartDate string `form:"start_date" binding:"required"` // 开始日期
	EndDate   string `form:"end_date" binding:"required"`   // 结束日期
	Granularity string `form:"granularity"`                  // 数据粒度：day/week/month/quarter/year
}

package request

type ReportReq struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
	Granularity string `form:"granularity"` // day, week, month, quarter, year
}

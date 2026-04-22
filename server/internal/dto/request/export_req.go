package request

type ExportReq struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Format    string `form:"format" binding:"required,oneof=csv json"`
}

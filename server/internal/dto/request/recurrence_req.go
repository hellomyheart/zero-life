// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// CreateRecurrenceReq 创建周期性交易请求
// 用于设置自动重复创建的交易规则
type CreateRecurrenceReq struct {
	Title          string  `json:"title" binding:"required"`                                         // 周期性交易名称
	Type           string  `json:"type" binding:"required,oneof=deposit withdrawal transfer"`         // 交易类型
	Amount         string  `json:"amount" binding:"required"`                                        // 交易金额
	SourceID       uint64  `json:"source_id" binding:"required"`                                     // 源账户ID
	DestinationID  *uint64 `json:"destination_id"`                                                   // 目标账户ID（转账时使用）
	CategoryID     *uint64 `json:"category_id"`                                                      // 分类ID
	Description    string  `json:"description"`                                                      // 交易描述
	Notes          string  `json:"notes"`                                                            // 备注
	TagNames       string  `json:"tag_names"`                                                        // 标签名称（逗号分隔）
	RepeatFreq     string  `json:"repeat_freq" binding:"required,oneof=daily weekly monthly yearly"`  // 重复频率
	RepeatInterval int     `json:"repeat_interval" binding:"required,min=1"`                         // 重复间隔（如每2周执行一次）
	NextDate       string  `json:"next_date" binding:"required"`                                     // 下次执行日期
	EndDate        string  `json:"end_date"`                                                         // 结束日期（可选）
	MaxRepetitions *int    `json:"max_repetitions"`                                                  // 最大重复次数（可选）
}

// UpdateRecurrenceReq 更新周期性交易请求
type UpdateRecurrenceReq struct {
	Title          string  `json:"title"`
	Type           string  `json:"type" binding:"omitempty,oneof=deposit withdrawal transfer"`
	Amount         string  `json:"amount"`
	SourceID       *uint64 `json:"source_id"`
	DestinationID  *uint64 `json:"destination_id"`
	CategoryID     *uint64 `json:"category_id"`
	Description    string  `json:"description"`
	Notes          string  `json:"notes"`
	TagNames       string  `json:"tag_names"`
	RepeatFreq     string  `json:"repeat_freq" binding:"omitempty,oneof=daily weekly monthly yearly"`
	RepeatInterval *int    `json:"repeat_interval"`
	NextDate       string  `json:"next_date"`
	EndDate        string  `json:"end_date"`
	MaxRepetitions *int    `json:"max_repetitions"`
	IsActive       *bool   `json:"is_active"` // 是否启用
}

// RecurrenceListReq 周期性交易列表查询请求
type RecurrenceListReq struct {
	Page     int `form:"page,default=1"`        // 页码
	PageSize int `form:"page_size,default=20"`  // 每页数量
}

// ExecuteRecurrenceReq 执行周期性交易请求
type ExecuteRecurrenceReq struct {
	StartDate string `json:"start_date"` // 开始日期
	EndDate   string `json:"end_date"`   // 结束日期
}

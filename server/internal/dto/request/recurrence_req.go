package request

type CreateRecurrenceReq struct {
	Title          string  `json:"title" binding:"required"`
	Type           string  `json:"type" binding:"required,oneof=deposit withdrawal transfer"`
	Amount         string  `json:"amount" binding:"required"`
	SourceID       uint64  `json:"source_id" binding:"required"`
	DestinationID  *uint64 `json:"destination_id"`
	CategoryID     *uint64 `json:"category_id"`
	Description    string  `json:"description"`
	Notes          string  `json:"notes"`
	TagNames       string  `json:"tag_names"`
	RepeatFreq     string  `json:"repeat_freq" binding:"required,oneof=daily weekly monthly yearly"`
	RepeatInterval int     `json:"repeat_interval" binding:"required,min=1"`
	NextDate       string  `json:"next_date" binding:"required"`
	EndDate        string  `json:"end_date"`
	MaxRepetitions *int    `json:"max_repetitions"`
}

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
	IsActive       *bool   `json:"is_active"`
}

type RecurrenceListReq struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}

type ExecuteRecurrenceReq struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

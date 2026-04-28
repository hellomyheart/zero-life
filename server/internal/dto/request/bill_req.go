package request

type CreateBillReq struct {
	Name       string  `json:"name" binding:"required"`
	Amount     string  `json:"amount" binding:"required"`
	RepeatRule string  `json:"repeat_rule" binding:"required,oneof=daily weekly monthly yearly"`
	NextDue    string  `json:"next_due" binding:"required"`
	SourceID   *uint64 `json:"source_id"`
	CategoryID *uint64 `json:"category_id"`
	Notes      string  `json:"notes"`
}

type UpdateBillReq struct {
	Name             string  `json:"name"`
	Amount           string  `json:"amount"`
	RepeatRule       string  `json:"repeat_rule" binding:"omitempty,oneof=daily weekly monthly yearly"`
	NextDue          string  `json:"next_due"`
	SourceID         *uint64 `json:"source_id"`
	CategoryID       *uint64 `json:"category_id"`
	Notes            string  `json:"notes"`
	ClearSourceID    bool    `json:"clear_source_id"`
	ClearCategoryID  bool    `json:"clear_category_id"`
	ClearNotes       bool    `json:"clear_notes"`
}
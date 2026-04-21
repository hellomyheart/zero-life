package request

type CreateCategoryReq struct {
	Name     string  `json:"name" binding:"required"`
	ParentID *uint64 `json:"parent_id"`
	Icon     string  `json:"icon"`
	Notes    string  `json:"notes"`
}

type UpdateCategoryReq struct {
	Name      string  `json:"name"`
	Icon      string  `json:"icon"`
	Notes     string  `json:"notes"`
	SortOrder *int    `json:"sort_order"`
}

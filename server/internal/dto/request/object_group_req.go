package request

type CreateObjectGroupReq struct {
	Name           string `json:"name" binding:"required"`
	GroupableType  string `json:"groupable_type" binding:"required,oneof=transaction bill budget piggy_bank"`
	GroupableID    uint64 `json:"groupable_id" binding:"required"`
}

type UpdateObjectGroupReq struct {
	Name string `json:"name"`
}
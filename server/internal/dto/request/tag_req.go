package request

type CreateTagReq struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color"`
}

type UpdateTagReq struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

package request

type CreateLinkTypeReq struct {
	Name          string `json:"name" binding:"required"`
	Outward       string `json:"outward" binding:"required"`
	Inward        string `json:"inward" binding:"required"`
	IsDirectional bool   `json:"is_directional"`
}

type UpdateLinkTypeReq struct {
	Name          string `json:"name"`
	Outward       string `json:"outward"`
	Inward        string `json:"inward"`
	IsDirectional *bool  `json:"is_directional"`
}

type LinkTypeListReq struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}

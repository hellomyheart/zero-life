package request

type SetPreferenceReq struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type GetPreferenceReq struct {
	Key string `form:"key" binding:"required"`
}
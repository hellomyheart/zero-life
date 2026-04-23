package request

type UpdatePreferenceReq struct {
	Value string `json:"value" binding:"required"`
}

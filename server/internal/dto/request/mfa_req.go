package request

type MFAEnableReq struct{}

type MFAConfirmReq struct {
	Code string `json:"code" binding:"required"`
}

type MFADisableReq struct {
	Code string `json:"code" binding:"required"`
}

type MFAVerifyReq struct {
	MFAToken string `json:"mfa_token" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

type MFACodeReq struct {
	Code string `json:"code" binding:"required"`
}

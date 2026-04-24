package request

// MFAVerifyReq MFA验证请求
type MFAVerifyReq struct {
	Code string `json:"code" binding:"required,len=6"` // 6位MFA代码
}

package response

// MFASetupResp MFA设置响应
type MFASetupResp struct {
	Secret string `json:"secret"` // MFA密钥
	QRCode string `json:"qr_code"` // 二维码URL
}

// MFAVerifyResp MFA验证响应
type MFAVerifyResp struct {
	Verified bool   `json:"verified"` // 是否验证成功
	Message  string `json:"message"`  // 消息
}

// MFAStatusResp MFA状态响应
type MFAStatusResp struct {
	Enabled bool `json:"enabled"` // 是否已启用MFA
}

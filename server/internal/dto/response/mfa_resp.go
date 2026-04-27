package response

type MFASetupResp struct {
	Secret string `json:"secret"`
	QRCode string `json:"qr_code"`
}

type MFAVerifyResp struct {
	Verified bool   `json:"verified"`
	Message  string `json:"message"`
}

type MFAStatusResp struct {
	Enabled bool `json:"enabled"`
}

type BackupCodesResp struct {
	Codes []string `json:"codes"`
}

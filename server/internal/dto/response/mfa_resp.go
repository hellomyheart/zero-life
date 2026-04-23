package response

type MFAEnableResp struct {
	Secret       string   `json:"secret"`
	QRCodeURL    string   `json:"qr_code_url"`
	BackupCodes  []string `json:"backup_codes"`
}

type MFAVerifyResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

type MFABackupCodeResp struct {
	ID        uint64  `json:"id"`
	Code      string  `json:"code"`
	UsedAt    *string `json:"used_at,omitempty"`
	CreatedAt string  `json:"created_at"`
}

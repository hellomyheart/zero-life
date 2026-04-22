package response

type AttachmentResp struct {
	ID             uint64 `json:"id"`
	AttachableType string `json:"attachable_type"`
	AttachableID   uint64 `json:"attachable_id"`
	Filename       string `json:"filename"`
	Mime           string `json:"mime"`
	Size           int64  `json:"size"`
	CreatedAt      string `json:"created_at"`
}

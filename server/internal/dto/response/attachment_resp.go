// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

// AttachmentResp 附件响应
type AttachmentResp struct {
	ID             uint64 `json:"id"`              // 附件ID
	AttachableType string `json:"attachable_type"` // 关联对象类型
	AttachableID   uint64 `json:"attachable_id"`   // 关联对象ID
	Filename       string `json:"filename"`        // 文件名
	Mime           string `json:"mime"`            // MIME类型
	Size           int64  `json:"size"`            // 文件大小（字节）
	CreatedAt      string `json:"created_at"`      // 创建时间
}

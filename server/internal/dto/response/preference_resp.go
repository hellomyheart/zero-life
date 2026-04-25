// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

// PreferenceResp 偏好设置响应
type PreferenceResp struct {
	Key   string `json:"key"`   // 偏好键名
	Value string `json:"value"` // 偏好值
}
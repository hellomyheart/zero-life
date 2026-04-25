// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// SetPreferenceReq 设置偏好请求
type SetPreferenceReq struct {
	Key   string `json:"key" binding:"required"`   // 偏好键名
	Value string `json:"value" binding:"required"` // 偏好值
}

// GetPreferenceReq 获取偏好请求
type GetPreferenceReq struct {
	Key string `form:"key" binding:"required"` // 偏好键名
}
// Package request 定义API请求参数结构
package request

// UpdateProfileReq 更新用户资料请求
type UpdateProfileReq struct {
	Name  string `json:"name"`  // 用户名称
	Email string `json:"email"` // 邮箱地址
}

// ChangePasswordReq 修改密码请求
type ChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"` // 旧密码
	NewPassword string `json:"new_password" binding:"required"` // 新密码
}

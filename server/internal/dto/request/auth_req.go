// Package request 定义所有HTTP请求的数据传输对象（DTO）
// 包含认证相关的请求结构体，用于绑定和校验请求参数
package request

// RegisterReq 用户注册请求
// 用于新用户注册，邮箱必须唯一，密码至少8位
type RegisterReq struct {
	Email    string `json:"email" binding:"required,email"`    // 用户邮箱，必须唯一
	Password string `json:"password" binding:"required,min=8"` // 登录密码，至少8位
	Nickname string `json:"nickname" binding:"required"`       // 用户昵称
}

// LoginReq 用户登录请求
// 用于邮箱密码登录认证
type LoginReq struct {
	Email    string `json:"email" binding:"required,email"` // 用户邮箱
	Password string `json:"password" binding:"required"`    // 登录密码
}

// RefreshReq 刷新令牌请求
// 使用刷新令牌获取新的访问令牌
type RefreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"` // 刷新令牌
}

// ForgotPasswordReq 忘记密码请求
// 用于请求发送密码重置邮件
type ForgotPasswordReq struct {
	Email string `json:"email" binding:"required,email"` // 注册邮箱
}

// ResetPasswordReq 重置密码请求
// 使用邮件中的重置令牌设置新密码
type ResetPasswordReq struct {
	Token    string `json:"token" binding:"required"`         // 密码重置令牌
	Password string `json:"password" binding:"required,min=8"` // 新密码，至少8位
}

// UpdateProfileReq 更新个人信息请求
// 用于修改当前用户的昵称、语言和时区
type UpdateProfileReq struct {
	Nickname string `json:"nickname"`                                // 用户昵称
	Language string `json:"language" binding:"omitempty,oneof=zh-CN en-US"` // 界面语言，仅支持zh-CN和en-US
	Timezone string `json:"timezone"`                                // 用户时区
}

// ChangePasswordReq 修改密码请求
// 用于已登录用户修改密码，需验证旧密码
type ChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`  // 当前密码
	NewPassword string `json:"new_password" binding:"required,min=8"` // 新密码，至少8位
}

// MFALoginVerifyReq MFA登录二次验证请求
// 用户启用MFA后，登录时需提交此请求完成二次验证
type MFALoginVerifyReq struct {
	MFAToken string `json:"mfa_token" binding:"required"` // 登录时返回的临时MFA验证令牌
	Code     string `json:"code" binding:"required,len=6"` // TOTP 6位验证码
}

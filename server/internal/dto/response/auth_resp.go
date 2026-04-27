// Package response 定义所有HTTP响应的数据传输对象（DTO）
// 包含认证相关的响应结构体，用于封装返回给客户端的数据
package response

import "time"

// LoginResp 登录响应
// 返回访问令牌、刷新令牌和过期时间
// 当用户启用了MFA时，mfa_required为true，access_token为临时MFA验证令牌（5分钟有效）
// 前端需引导用户输入TOTP码，再调用 /auth/mfa-verify 完成二次验证换取真正的TokenPair
type LoginResp struct {
	MFARequired  bool      `json:"mfa_required"`  // 是否需要MFA二次验证
	AccessToken  string    `json:"access_token"`  // 访问令牌（MFA时为临时验证令牌）
	RefreshToken string    `json:"refresh_token"` // 刷新令牌（MFA时为空）
	ExpiresAt    time.Time `json:"expires_at"`    // 令牌过期时间
}

// ProfileResp 用户个人信息响应
type ProfileResp struct {
	ID        uint64    `json:"id"`         // 用户ID
	Email     string    `json:"email"`      // 用户邮箱
	Nickname  string    `json:"nickname"`   // 用户昵称
	Language  string    `json:"language"`   // 界面语言
	Timezone  string    `json:"timezone"`   // 时区
	CreatedAt time.Time `json:"created_at"` // 注册时间
}

// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// AdminUserResp 管理员用户信息响应
type AdminUserResp struct {
	ID         uint64    `json:"id"`          // 用户ID
	Email      string    `json:"email"`       // 用户邮箱
	Nickname   string    `json:"nickname"`    // 用户昵称
	Role       string    `json:"role"`        // 用户角色
	Language   string    `json:"language"`    // 界面语言
	Timezone   string    `json:"timezone"`    // 时区
	MFAEnabled bool     `json:"mfa_enabled"` // 是否启用MFA
	CreatedAt  time.Time `json:"created_at"`  // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`  // 更新时间
}

// AdminConfigurationResp 管理员配置项响应
type AdminConfigurationResp struct {
	ID        uint64    `json:"id"`         // 配置ID
	Name      string    `json:"name"`       // 配置名称
	Value     string    `json:"value"`      // 配置值
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// AdminTestEmailResp 管理员邮件测试响应
type AdminTestEmailResp struct {
	Success bool   `json:"success"` // 是否发送成功
	Message string `json:"message"` // 结果消息
}

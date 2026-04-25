// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// AdminListUsersReq 管理员查询用户列表请求
type AdminListUsersReq struct {
	Page     int    `form:"page"`      // 页码
	PageSize int    `form:"page_size"` // 每页数量
	Search   string `form:"search"`    // 搜索关键词
}

// AdminUpdateUserReq 管理员更新用户请求
type AdminUpdateUserReq struct {
	Nickname string `json:"nickname"`                           // 用户昵称
	Role     string `json:"role" binding:"omitempty,oneof=user admin"` // 用户角色
	Language string `json:"language"`                          // 界面语言
	Timezone string `json:"timezone"`                          // 时区
}

// AdminInviteUserReq 管理员邀请用户请求
type AdminInviteUserReq struct {
	Email    string `json:"email" binding:"required,email"`     // 邀请邮箱
	Nickname string `json:"nickname" binding:"required"`        // 用户昵称
	Role     string `json:"role" binding:"omitempty,oneof=user admin"` // 用户角色
}

// AdminUpdateConfigurationReq 管理员更新配置请求
type AdminUpdateConfigurationReq struct {
	Value string `json:"value" binding:"required"` // 配置值
}

// AdminTestEmailReq 管理员测试邮件请求
type AdminTestEmailReq struct {
	Email string `json:"email" binding:"required,email"` // 测试邮箱地址
}

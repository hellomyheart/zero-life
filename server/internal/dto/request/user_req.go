package request

// UserListReq 用户列表请求
type UserListReq struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}

// UpdateUserReq 更新用户请求
type UpdateUserReq struct {
	Nickname string `json:"nickname"`
	Language string `json:"language"`
	Timezone string `json:"timezone"`
}

// ChangeRoleReq 修改角色请求
type ChangeRoleReq struct {
	Role string `json:"role" binding:"required,oneof=user admin"`
}

// AdminResetPasswordReq 管理员重置密码请求
type AdminResetPasswordReq struct {
	Password string `json:"password" binding:"required,min=8"`
}

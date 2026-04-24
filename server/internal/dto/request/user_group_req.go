// Package request 定义API请求参数结构
package request

// CreateUserGroupReq 创建用户组请求
type CreateUserGroupReq struct {
	Title string `json:"title" binding:"required"` // 用户组标题
}

// UpdateUserGroupReq 更新用户组请求
type UpdateUserGroupReq struct {
	Title string `json:"title"` // 用户组标题
}

// UserGroupListReq 用户组列表查询请求
type UserGroupListReq struct {
	Page     int `form:"page"`      // 页码
	PageSize int `form:"page_size"` // 每页记录数
}

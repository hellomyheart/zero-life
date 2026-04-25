// Package request 定义所有HTTP请求的数据传输对象（DTO）
package request

// CreateAccountReq 创建账户请求
// 用于创建新的资产/支出/收入/负债账户
type CreateAccountReq struct {
	Name           string `json:"name" binding:"required"`                                          // 账户名称，同类型下不能重复
	Type           string `json:"type" binding:"required,oneof=asset expense revenue liability"`     // 账户类型：asset(资产)/expense(支出)/revenue(收入)/liability(负债)
	CurrencyID     uint64 `json:"currency_id" binding:"required"`                                   // 关联货币ID
	InitialBalance string `json:"initial_balance" binding:"required"`                               // 初始余额
	Notes          string `json:"notes"`                                                            // 备注
	IsVirtual      bool   `json:"is_virtual"`                                                       // 是否为虚拟账户（虚拟账户不影响总资产）
}

// UpdateAccountReq 更新账户请求
// 仅允许修改名称、备注和虚拟属性，账户类型和货币不可更改
type UpdateAccountReq struct {
	Name      string `json:"name"`                  // 账户名称
	Notes     string `json:"notes"`                 // 备注
	IsVirtual *bool  `json:"is_virtual"`            // 是否为虚拟账户，使用指针以区分未传值和false
}

// AccountListReq 账户列表查询请求
// 用于分页查询和过滤账户列表
type AccountListReq struct {
	Page     int    `form:"page,default=1"`        // 页码，默认第1页
	PageSize int    `form:"page_size,default=20"`  // 每页数量，默认20
	Type     string `form:"type"`                  // 按账户类型过滤
	Search   string `form:"search"`                // 按名称搜索
	Sort     string `form:"sort,default=name"`     // 排序字段，默认按名称排序
}

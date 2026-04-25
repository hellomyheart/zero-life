// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// RuleGroupResp 规则组响应
type RuleGroupResp struct {
	ID        uint64    `json:"id"`         // 规则组ID
	UserID    uint64    `json:"user_id"`    // 用户ID
	Name      string    `json:"name"`       // 规则组名称
	Order     int       `json:"order"`      // 排序序号
	IsActive  bool      `json:"is_active"`  // 是否启用
	RuleCount int64     `json:"rule_count"` // 组内规则数量
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// RuleGroupExecuteResultResp 规则组执行结果响应
type RuleGroupExecuteResultResp struct {
	GroupID      int      `json:"group_id"`      // 规则组ID
	MatchedCount int      `json:"matched_count"` // 匹配的交易数
	SuccessCount int      `json:"success_count"` // 成功执行的动作数
	FailCount    int      `json:"fail_count"`    // 失败的动作数
	Errors       []string `json:"errors,omitempty"` // 错误信息列表
}

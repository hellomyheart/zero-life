// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

import "time"

// RecurrenceResp 周期性交易响应
type RecurrenceResp struct {
	ID             uint64     `json:"id"`                          // 周期性交易ID
	UserID         uint64     `json:"user_id"`                     // 用户ID
	Title          string     `json:"title"`                       // 标题
	Type           string     `json:"type"`                        // 交易类型
	Amount         string     `json:"amount"`                      // 交易金额
	SourceID       uint64     `json:"source_id"`                   // 源账户ID
	SourceName     string     `json:"source_name"`                 // 源账户名称
	DestinationID  *uint64    `json:"destination_id"`              // 目标账户ID
	DestinationName string    `json:"destination_name,omitempty"`  // 目标账户名称
	CategoryID     *uint64    `json:"category_id"`                 // 分类ID
	Description    string     `json:"description"`                 // 交易描述
	Notes          string     `json:"notes"`                       // 备注
	TagNames       string     `json:"tag_names"`                   // 标签名称
	RepeatFreq     string     `json:"repeat_freq"`                 // 重复频率
	RepeatInterval int        `json:"repeat_interval"`             // 重复间隔
	NextDate       time.Time  `json:"next_date"`                   // 下次执行日期
	EndDate        *time.Time `json:"end_date"`                    // 结束日期
	Repetitions    int        `json:"repetitions"`                 // 已执行次数
	MaxRepetitions *int       `json:"max_repetitions"`             // 最大重复次数
	IsActive       bool       `json:"is_active"`                   // 是否启用
	CreatedAt      time.Time  `json:"created_at"`                  // 创建时间
	UpdatedAt      time.Time  `json:"updated_at"`                  // 更新时间
}

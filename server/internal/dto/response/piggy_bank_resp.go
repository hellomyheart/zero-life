// Package response 定义所有HTTP响应的数据传输对象（DTO）
package response

// PiggyBankResp 存钱罐响应
type PiggyBankResp struct {
	ID               uint64       `json:"id"`                          // 存钱罐ID
	Name             string       `json:"name"`                        // 存钱罐名称
	TargetAmount     string       `json:"target_amount"`               // 目标金额
	CurrentAmount    string       `json:"current_amount"`              // 当前已存金额
	AccountID        uint64       `json:"account_id"`                  // 关联账户ID
	Account          *AccountResp `json:"account,omitempty"`           // 关联账户信息
	TargetDate       *string      `json:"target_date,omitempty"`       // 目标达成日期
	Notes            string       `json:"notes"`                       // 备注
	Percentage       float64      `json:"percentage"`                  // 完成百分比
	AvailableDeposit string       `json:"available_deposit"`           // 可存入金额（考虑账户余额和目标金额限制）
	CreatedAt        string       `json:"created_at"`                  // 创建时间
	UpdatedAt        string       `json:"updated_at"`                  // 更新时间
}

// PiggyEventResp 存钱罐事件响应
// 记录每次存入或取出的操作
type PiggyEventResp struct {
	ID            uint64  `json:"id"`                          // 事件ID
	PiggyBankID   uint64  `json:"piggy_bank_id"`               // 存钱罐ID
	Amount        string  `json:"amount"`                      // 存取金额
	TransactionID *uint64 `json:"transaction_id,omitempty"`    // 关联交易ID
	Note          string  `json:"note"`                        // 操作备注
	CreatedAt     string  `json:"created_at"`                  // 创建时间
}

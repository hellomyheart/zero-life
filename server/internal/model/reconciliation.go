// Package model 定义了系统的数据模型，对应数据库表结构
// 包含所有业务实体：账户、交易、分类、标签、预算、账单、对账等
package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Reconciliation 对账模型，对应 reconciliations 表
// 
// 功能说明：
// - 记录账户对账信息，核对账面余额与银行实际余额
// - 发现并记录差异，确保账务准确
// - 支持分期对账（按日期范围）
//
// 使用场景：
// - 银行对账：核对银行账户与银行对账单
// - 信用卡对账：核对信用卡账单与实际消费
// - 现金盘点：核对现金账户与实际库存现金
//
// 对账流程：
// 1. 选择对账账户（如：招商银行卡）
// 2. 设定对账期间（如：2026-04-01 至 2026-04-30）
// 3. 输入期初余额（银行对账单上的期初余额）
// 4. 输入期末余额（银行对账单上的期末余额）
// 5. 逐笔核对交易，标记已匹配的交易
// 6. 系统自动计算差额
// 7. 提交对账结果
//
// 关键字段说明：
// - StartBalance: 期初余额（银行对账单）
// - EndBalance: 期末余额（银行对账单）
// - SubmittedBalance: 用户提交的余额（系统账面余额）
// - Difference: 差额 = SubmittedBalance - EndBalance
//
// 差额分析：
// - Difference = 0: 完全匹配
// - Difference > 0: 账面多于实际（可能有未达账项）
// - Difference < 0: 账面少于实际（可能有遗漏交易）
//
// 示例：
//   账户：招商银行卡
//   期间：2026-04-01 ~ 2026-04-30
//   期初余额：10,000 元
//   期末余额：15,000 元
//   账面余额：14,950 元
//   差额：-50 元（需要查找原因）
type Reconciliation struct {
	// ID 对账记录唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	
	// UserID 所属用户 ID
	// 每个用户有自己的对账记录，用户间数据隔离
	// gorm:"index" 加速按用户查询
	UserID uint64 `gorm:"not null;index" json:"user_id"`
	
	// AccountID 对账账户 ID
	// 关联到 accounts 表
	// 如：对"招商银行卡"进行对账
	// gorm:"index" 加速按账户查询
	AccountID uint64 `gorm:"not null;index" json:"account_id"`
	
	// StartDate 对账开始日期
	// 对账期间的起始日期
	// 如：2026-04-01
	StartDate time.Time `gorm:"not null" json:"start_date"`
	
	// EndDate 对账结束日期
	// 对账期间的结束日期
	// 如：2026-04-30
	EndDate time.Time `gorm:"not null" json:"end_date"`
	
	// StartBalance 期初余额
	// 对账期间开始时的账户余额（银行对账单上的余额）
	// 使用 decimal 类型确保金额精度
	// gorm:"type:decimal(19,4)" 支持最大 15 位整数 +4 位小数
	StartBalance decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"start_balance"`
	
	// EndBalance 期末余额
	// 对账期间结束时的账户余额（银行对账单上的余额）
	// 这是"实际余额"，用于与系统"账面余额"对比
	EndBalance decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"end_balance"`
	
	// SubmittedBalance 用户提交的余额（系统账面余额）
	// 系统中该账户在对账期间结束时的余额
	// 这是"账面余额"，由系统自动计算
	// 
	// 计算公式：
	//   SubmittedBalance = StartBalance + 期间收入 - 期间支出
	SubmittedBalance decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"submitted_balance"`
	
	// Difference 差额
	// 计算公式：Difference = SubmittedBalance - EndBalance
	// 
	// 结果分析：
	// - Difference = 0: 完全匹配，对账成功
	// - Difference > 0: 账面多于实际，可能有未达账项或重复记账
	// - Difference < 0: 账面少于实际，可能有遗漏交易或金额错误
	//
	// 示例：
	//   SubmittedBalance = 14,950
	//   EndBalance = 15,000
	//   Difference = -50（需要查找 50 元的差异原因）
	Difference decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"difference"`
	
	// CreatedAt 对账记录创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
}

// TableName 指定 Reconciliation 模型对应的数据库表名为 reconciliations
func (Reconciliation) TableName() string { return "reconciliations" }

// TransactionReconciliation 交易对账模型，对应 transaction_reconciliations 表
// 
// 功能说明：
// - 记录单次对账会话的详细信息
// - 跟踪对账状态（开放/关闭）
// - 存储对账过程中的余额信息
//
// 与 Reconciliation 的区别：
// - Reconciliation: 简化的对账记录
// - TransactionReconciliation: 详细的对账会话，包含状态跟踪
//
// 对账状态：
// - open: 对账进行中，可以继续匹配交易
// - closed: 对账已完成，不允许修改
type TransactionReconciliation struct {
	// ID 对账会话唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	
	// AccountID 对账账户 ID
	AccountID uint64 `gorm:"not null;index" json:"account_id"`
	
	// StartDate 对账开始日期
	StartDate time.Time `gorm:"not null" json:"start_date"`
	
	// EndDate 对账结束日期
	EndDate time.Time `gorm:"not null" json:"end_date"`
	
	// StartingBalance 期初余额（字符串格式）
	// 使用字符串避免精度丢失
	StartingBalance string `gorm:"type:decimal(19,4)" json:"starting_balance"`
	
	// EndingBalance 期末余额（字符串格式）
	// 银行对账单上的期末余额
	EndingBalance string `gorm:"type:decimal(19,4)" json:"ending_balance"`
	
	// BookBalance 账面余额（字符串格式）
	// 系统中计算的余额
	BookBalance string `gorm:"type:decimal(19,4)" json:"book_balance"`
	
	// Difference 差额（字符串格式）
	// BookBalance - EndingBalance
	Difference string `gorm:"type:decimal(19,4)" json:"difference"`
	
	// Status 对账状态
	// 可选值：
	//   - "open": 对账进行中
	//   - "closed": 对账已完成
	// 默认值：open
	Status string `gorm:"default:open" json:"status"`
	
	// CreatedAt 对账会话创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	
	// UpdatedAt 对账会话最后更新时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

// TableName 指定 TransactionReconciliation 模型对应的数据库表名为 transaction_reconciliations
func (TransactionReconciliation) TableName() string { return "transaction_reconciliations" }

// ReconciliationEntry 对账条目模型，对应 reconciliation_entries 表
// 
// 功能说明：
// - 记录对账过程中每笔交易的匹配情况
// - 逐笔核对交易，标记是否匹配
// - 记录金额差异（如有）
//
// 使用场景：
// - 标记某笔交易已在银行对账单中找到
// - 记录系统交易与银行交易的金额差异
// - 统计已匹配和未匹配的交易数量
//
// 匹配流程：
// 1. 列出对账期间内的所有交易
// 2. 逐笔与银行对账单核对
// 3. 找到对应交易则标记 Matched=true
// 4. 如有金额差异，记录 AmountDifference
// 5. 统计未匹配交易，查找原因
type ReconciliationEntry struct {
	// ID 对账条目唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	
	// ReconciliationID 所属对账记录 ID
	// 关联到 transaction_reconciliations 表
	// gorm:"not null;index" 加速按对账记录查询
	ReconciliationID uint64 `gorm:"not null;index" json:"reconciliation_id"`
	
	// TransactionID 交易 ID
	// 关联到 transactions 表
	// 这是需要核对的具体交易
	// gorm:"not null;index" 加速按交易查询
	TransactionID uint64 `gorm:"not null;index" json:"transaction_id"`
	
	// Matched 是否匹配
	// true: 已在银行对账单中找到对应交易
	// false: 未找到对应交易（可能是未达账项）
	// 默认值：false
	Matched bool `gorm:"default:false" json:"matched"`
	
	// AmountDifference 金额差异
	// 系统交易金额与银行交易金额的差额
	// 如：系统记录 100 元，银行记录 99.5 元，差异 0.5 元
	// gorm:"type:decimal(19,4)" 支持高精度金额
	AmountDifference string `gorm:"type:decimal(19,4)" json:"amount_difference"`
	
	// CreatedAt 对账条目创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
}

// TableName 指定 ReconciliationEntry 模型对应的数据库表名为 reconciliation_entries
func (ReconciliationEntry) TableName() string { return "reconciliation_entries" }

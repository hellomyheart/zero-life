// Package repository 提供数据访问层，封装所有与数据库的交互逻辑。
// 每个 Repository 结构体对应一个数据库表，提供增删改查等基本操作。
// 所有 Repository 都依赖 GORM 作为 ORM 框架，通过 *gorm.DB 执行数据库操作。
package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// AccountRepository 账户仓库，负责账户的数据访问。
// 账户（Account）是用户管理资金的容器，如银行卡、现金、信用卡等。
// 每个账户关联一种货币（Currency），并记录当前余额（current_balance）。
// 所有方法都通过 user_id 参数进行数据隔离，确保用户只能访问自己的账户。
type AccountRepository struct {
	db *gorm.DB
}

// NewAccountRepository 创建账户仓库实例。
// 参数 db: GORM 数据库连接实例。
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// Create 创建一条新的账户记录。
// 执行 SQL: INSERT INTO accounts (...)
// 参数 account: 要创建的账户对象，GORM 会自动填充 ID、CreatedAt 等字段。
// 返回: 创建失败时返回错误。
func (r *AccountRepository) Create(account *model.Account) error {
	return r.db.Create(account).Error
}

// GetByID 根据 ID 和用户 ID 获取单条账户，并预加载关联的货币信息。
// 同时验证该账户属于指定用户（权限校验），确保数据隔离。
// 执行 SQL:
//   主查询: SELECT * FROM accounts WHERE id = ? AND user_id = ? LIMIT 1
//   预加载: SELECT * FROM currencies WHERE id IN (...)  (Currency)
// 参数 id: 账户 ID。
// 参数 userID: 当前登录用户 ID，用于权限校验，确保用户只能查看自己的账户。
// 返回: 包含货币信息的账户对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *AccountRepository) GetByID(id, userID uint64) (*model.Account, error) {
	var account model.Account
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).Preload("Currency").First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

// List 分页获取指定用户的账户列表，支持按类型过滤、名称搜索和多种排序方式。
// 所有查询都基于 user_id 进行数据隔离，确保用户只能看到自己的账户。
// 执行 SQL: SELECT * FROM accounts WHERE user_id = ? [AND type = ?] [AND name LIKE ?] ORDER BY ... LIMIT ? OFFSET ?
// 参数 userID: 用户 ID，用于数据隔离。
// 参数 accountType: 可选的账户类型过滤，如 "asset"（资产）、"liability"（负债）等，为空时不过滤。
// 参数 search: 可选的名称搜索关键词，使用 LIKE 模糊匹配，为空时不过滤。
// 参数 sort: 排序方式，支持 "name"（名称升序）、"-name"（名称降序）、"balance"（余额升序）、"-balance"（余额降序），默认按名称升序。
// 参数 offset: 分页偏移量（跳过的记录数）。
// 参数 limit: 每页记录数。
// 返回: 账户列表（含预加载的货币信息）。
func (r *AccountRepository) List(userID uint64, accountType, search, sort string, offset, limit int) ([]model.Account, error) {
	var accounts []model.Account
	// 基于 user_id 过滤，确保数据隔离：用户只能查看自己的账户
	query := r.db.Where("user_id = ?", userID)

	// 按账户类型过滤（如资产账户、负债账户等）
	if accountType != "" {
		query = query.Where("type = ?", accountType)
	}
	// 按名称模糊搜索
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Apply sorting
	// 排序规则：支持按名称或余额排序，前缀 "-" 表示降序
	switch sort {
	case "name":
		query = query.Order("name ASC")
	case "-name":
		query = query.Order("name DESC")
	case "balance":
		query = query.Order("current_balance ASC")
	case "-balance":
		query = query.Order("current_balance DESC")
	default:
		query = query.Order("name ASC")
	}

	// 预加载关联的货币信息，避免 N+1 查询问题
	if err := query.Preload("Currency").Offset(offset).Limit(limit).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

// Count 统计指定用户符合过滤条件的账户总数，用于分页计算。
// 过滤条件与 List 方法一致，确保分页总数与列表数据一致。
// 执行 SQL: SELECT COUNT(*) FROM accounts WHERE user_id = ? [AND type = ?] [AND name LIKE ?]
// 参数 userID: 用户 ID，用于数据隔离。
// 参数 accountType: 可选的账户类型过滤，为空时不过滤。
// 参数 search: 可选的名称搜索关键词，为空时不过滤。
// 返回: 符合条件的账户总数。
func (r *AccountRepository) Count(userID uint64, accountType, search string) (int64, error) {
	var count int64
	query := r.db.Model(&model.Account{}).Where("user_id = ?", userID)

	if accountType != "" {
		query = query.Where("type = ?", accountType)
	}
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Update 更新账户信息。GORM 的 Save 方法会更新所有字段。
// 执行 SQL: UPDATE accounts SET ... WHERE id = ?
// 参数 account: 要更新的账户对象（必须包含 ID 字段）。
// 返回: 更新失败时返回错误。
func (r *AccountRepository) Update(account *model.Account) error {
	return r.db.Save(account).Error
}

// Delete 根据 ID 和用户 ID 删除账户，同时验证用户权限。
// 执行 SQL: DELETE FROM accounts WHERE id = ? AND user_id = ?
// 参数 id: 账户 ID。
// 参数 userID: 当前登录用户 ID，确保只能删除自己的账户。
// 返回: 删除失败时返回错误。
func (r *AccountRepository) Delete(id, userID uint64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Account{}).Error
}

// HasTransactions 检查指定账户是否有关联的交易记录。
// 匹配源账户（source_id）或目标账户（destination_id），用于删除账户前的校验。
// 如果账户存在关联交易，应阻止删除，以保持数据完整性。
// 执行 SQL: SELECT COUNT(*) FROM transactions WHERE user_id = ? AND (source_id = ? OR destination_id = ?)
// 参数 accountID: 账户 ID。
// 参数 userID: 当前登录用户 ID，用于数据隔离。
// 返回: 如果有关联交易返回 true，否则返回 false。
func (r *AccountRepository) HasTransactions(accountID, userID uint64) (bool, error) {
	var count int64
	// 查询该账户作为源账户或目标账户的交易数量
	if err := r.db.Model(&model.Transaction{}).
		Where("user_id = ? AND (source_id = ? OR destination_id = ?)", userID, accountID, accountID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountByCurrency 统计使用指定货币的账户数量
// 用于禁用货币前的校验：如果有账户正在使用该货币，则不允许禁用
// 执行 SQL: SELECT COUNT(*) FROM accounts WHERE currency_id = ?
// 参数 currencyID: 货币ID
// 返回: 使用该货币的账户数量，查询失败时返回错误
func (r *AccountRepository) CountByCurrency(currencyID uint64) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Account{}).Where("currency_id = ?", currencyID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

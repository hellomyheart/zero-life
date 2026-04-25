package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// CurrencyRepository 货币仓库，负责货币和汇率的数据访问。
// 货币（Currency）是全局资源，不属于某个用户，用于定义系统支持的货币类型。
// 汇率（ExchangeRate）记录不同货币之间的兑换比率。
// 每个货币可以设置为默认货币（is_default），系统会使用默认货币作为基准。
type CurrencyRepository struct {
	db *gorm.DB
}

// NewCurrencyRepository 创建货币仓库实例。
// 参数 db: GORM 数据库连接实例。
func NewCurrencyRepository(db *gorm.DB) *CurrencyRepository {
	return &CurrencyRepository{db: db}
}

// Transaction 在事务中执行数据库操作，用于需要原子性的场景（如设置默认货币）。
// 参数 fn: 事务内执行的函数，接收 *gorm.DB 作为参数。
// 返回: 事务执行失败时返回错误。
func (r *CurrencyRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

// List 获取所有货币列表，按 ID 升序排列。
// 执行 SQL: SELECT * FROM currencies ORDER BY id ASC
// 返回: 货币列表。
func (r *CurrencyRepository) List() ([]model.Currency, error) {
	var currencies []model.Currency
	if err := r.db.Order("id ASC").Find(&currencies).Error; err != nil {
		return nil, err
	}
	return currencies, nil
}

// GetByID 根据 ID 获取单条货币记录。
// 执行 SQL: SELECT * FROM currencies WHERE id = ? LIMIT 1
// 参数 id: 货币 ID。
// 返回: 找到的货币对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *CurrencyRepository) GetByID(id uint64) (*model.Currency, error) {
	var currency model.Currency
	if err := r.db.First(&currency, id).Error; err != nil {
		return nil, err
	}
	return &currency, nil
}

// Update 更新货币信息。GORM 的 Save 方法会更新所有字段。
// 执行 SQL: UPDATE currencies SET ... WHERE id = ?
// 参数 currency: 要更新的货币对象（必须包含 ID 字段）。
// 返回: 更新失败时返回错误。
func (r *CurrencyRepository) Update(currency *model.Currency) error {
	return r.db.Save(currency).Error
}

// GetByCode 根据货币代码查找货币，如 "CNY"、"USD"、"EUR" 等。
// 执行 SQL: SELECT * FROM currencies WHERE code = ? LIMIT 1
// 参数 code: 货币代码（ISO 4217 标准）。
// 返回: 找到的货币对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *CurrencyRepository) GetByCode(code string) (*model.Currency, error) {
	var currency model.Currency
	if err := r.db.Where("code = ?", code).First(&currency).Error; err != nil {
		return nil, err
	}
	return &currency, nil
}

// GetDefault 获取系统默认货币。默认货币用于金额汇总和显示的基准。
// 执行 SQL: SELECT * FROM currencies WHERE is_default = true LIMIT 1
// 返回: 默认货币对象；未设置默认货币时返回 gorm.ErrRecordNotFound 错误。
func (r *CurrencyRepository) GetDefault() (*model.Currency, error) {
	var currency model.Currency
	if err := r.db.Where("is_default = ?", true).First(&currency).Error; err != nil {
		return nil, err
	}
	return &currency, nil
}

// ListExchangeRates 获取所有汇率记录，按 ID 升序排列。
// 执行 SQL: SELECT * FROM exchange_rates ORDER BY id ASC
// 返回: 汇率列表。
func (r *CurrencyRepository) ListExchangeRates() ([]model.ExchangeRate, error) {
	var rates []model.ExchangeRate
	if err := r.db.Order("id ASC").Find(&rates).Error; err != nil {
		return nil, err
	}
	return rates, nil
}

// UpsertExchangeRate 创建或更新汇率记录。GORM 的 Save 方法：
// 如果对象包含主键 ID 且记录存在则更新，否则创建新记录。
// 执行 SQL: INSERT INTO exchange_rates (...) / UPDATE exchange_rates SET ... WHERE id = ?
// 参数 rate: 汇率对象，如果包含 ID 则尝试更新，否则创建。
// 返回: 操作失败时返回错误。
func (r *CurrencyRepository) UpsertExchangeRate(rate *model.ExchangeRate) error {
	return r.db.Save(rate).Error
}

// Create 创建一条新的货币记录。
// 执行 SQL: INSERT INTO currencies (...)
// 参数 currency: 要创建的货币对象。
// 返回: 创建失败时返回错误。
func (r *CurrencyRepository) Create(currency *model.Currency) error {
	return r.db.Create(currency).Error
}

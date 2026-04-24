// Package repository 数据访问层，封装数据库操作
package repository

import (
	"time"

	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// ExchangeRateRepository 汇率数据访问对象
// 提供汇率数据的增删改查操作
type ExchangeRateRepository struct {
	db *gorm.DB
}

// NewExchangeRateRepository 创建汇率数据访问对象实例
// 参数：
//   - db: 数据库连接
// 返回：
//   - *ExchangeRateRepository: 汇率数据访问对象实例
func NewExchangeRateRepository(db *gorm.DB) *ExchangeRateRepository {
	return &ExchangeRateRepository{db: db}
}

// Create 创建汇率记录
// 参数：
//   - rate: 汇率模型
// 返回：
//   - error: 错误信息
func (r *ExchangeRateRepository) Create(rate *model.ExchangeRate) error {
	return r.db.Create(rate).Error
}

// GetByID 根据ID获取汇率记录
// 参数：
//   - id: 汇率ID
//   - userID: 用户ID
// 返回：
//   - *model.ExchangeRate: 汇率记录
//   - error: 错误信息
func (r *ExchangeRateRepository) GetByID(id, userID uint64) (*model.ExchangeRate, error) {
	var rate model.ExchangeRate
	err := r.db.Where("id = ? AND user_id = ?", id, userID).
		Preload("FromCurrency").
		Preload("ToCurrency").
		First(&rate).Error
	return &rate, err
}

// List 获取汇率列表
// 参数：
//   - userID: 用户ID
//   - fromCurrencyID: 源货币ID（可选，0表示不过滤）
//   - toCurrencyID: 目标货币ID（可选，0表示不过滤）
//   - startDate: 开始日期（可选）
//   - endDate: 结束日期（可选）
//   - offset: 偏移量
//   - limit: 限制数量
// 返回：
//   - []model.ExchangeRate: 汇率列表
//   - error: 错误信息
func (r *ExchangeRateRepository) List(userID uint64, fromCurrencyID, toCurrencyID uint64, startDate, endDate *time.Time, offset, limit int) ([]model.ExchangeRate, error) {
	query := r.db.Where("user_id = ?", userID)

	if fromCurrencyID > 0 {
		query = query.Where("from_currency_id = ?", fromCurrencyID)
	}
	if toCurrencyID > 0 {
		query = query.Where("to_currency_id = ?", toCurrencyID)
	}
	if startDate != nil {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("date <= ?", endDate)
	}

	var rates []model.ExchangeRate
	err := query.Preload("FromCurrency").
		Preload("ToCurrency").
		Order("date DESC").
		Offset(offset).
		Limit(limit).
		Find(&rates).Error
	return rates, err
}

// Count 统计汇率记录数量
// 参数：
//   - userID: 用户ID
//   - fromCurrencyID: 源货币ID（可选，0表示不过滤）
//   - toCurrencyID: 目标货币ID（可选，0表示不过滤）
// 返回：
//   - int64: 记录数量
//   - error: 错误信息
func (r *ExchangeRateRepository) Count(userID uint64, fromCurrencyID, toCurrencyID uint64) (int64, error) {
	query := r.db.Model(&model.ExchangeRate{}).Where("user_id = ?", userID)

	if fromCurrencyID > 0 {
		query = query.Where("from_currency_id = ?", fromCurrencyID)
	}
	if toCurrencyID > 0 {
		query = query.Where("to_currency_id = ?", toCurrencyID)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}

// Update 更新汇率记录
// 参数：
//   - rate: 汇率模型
// 返回：
//   - error: 错误信息
func (r *ExchangeRateRepository) Update(rate *model.ExchangeRate) error {
	return r.db.Save(rate).Error
}

// Delete 删除汇率记录
// 参数：
//   - id: 汇率ID
//   - userID: 用户ID
// 返回：
//   - error: 错误信息
func (r *ExchangeRateRepository) Delete(id, userID uint64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.ExchangeRate{}).Error
}

// GetSpecificRateOnDate 获取特定日期的汇率
// 参数：
//   - userID: 用户ID
//   - fromCurrencyID: 源货币ID
//   - toCurrencyID: 目标货币ID
//   - date: 日期
// 返回：
//   - *model.ExchangeRate: 汇率记录
//   - error: 错误信息
func (r *ExchangeRateRepository) GetSpecificRateOnDate(userID, fromCurrencyID, toCurrencyID uint64, date time.Time) (*model.ExchangeRate, error) {
	var rate model.ExchangeRate
	err := r.db.Where("user_id = ? AND from_currency_id = ? AND to_currency_id = ? AND date = ?",
		userID, fromCurrencyID, toCurrencyID, date).
		Preload("FromCurrency").
		Preload("ToCurrency").
		First(&rate).Error
	return &rate, err
}

// GetLatestRate 获取最新汇率
// 参数：
//   - userID: 用户ID
//   - fromCurrencyID: 源货币ID
//   - toCurrencyID: 目标货币ID
// 返回：
//   - *model.ExchangeRate: 汇率记录
//   - error: 错误信息
func (r *ExchangeRateRepository) GetLatestRate(userID, fromCurrencyID, toCurrencyID uint64) (*model.ExchangeRate, error) {
	var rate model.ExchangeRate
	err := r.db.Where("user_id = ? AND from_currency_id = ? AND to_currency_id = ?",
		userID, fromCurrencyID, toCurrencyID).
		Preload("FromCurrency").
		Preload("ToCurrency").
		Order("date DESC").
		First(&rate).Error
	return &rate, err
}

// DeleteRates 删除指定货币对的所有汇率记录
// 参数：
//   - userID: 用户ID
//   - fromCurrencyID: 源货币ID
//   - toCurrencyID: 目标货币ID
// 返回：
//   - error: 错误信息
func (r *ExchangeRateRepository) DeleteRates(userID, fromCurrencyID, toCurrencyID uint64) error {
	return r.db.Where("user_id = ? AND from_currency_id = ? AND to_currency_id = ?",
		userID, fromCurrencyID, toCurrencyID).
		Delete(&model.ExchangeRate{}).Error
}

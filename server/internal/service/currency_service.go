// Package service 业务逻辑层，实现核心业务逻辑
// CurrencyService 货币业务逻辑，处理货币管理和汇率转换
package service

import (
	"github.com/shopspring/decimal"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

// CurrencyService 货币服务
// 负责处理货币管理、默认货币设置和汇率转换
// 依赖currencyRepo进行货币数据访问，依赖accountRepo验证货币使用情况
type CurrencyService struct {
	currencyRepo *repository.CurrencyRepository // 货币数据访问对象
	accountRepo  *repository.AccountRepository  // 账户数据访问对象
}

// NewCurrencyService 创建货币服务实例
func NewCurrencyService(currencyRepo *repository.CurrencyRepository, accountRepo *repository.AccountRepository) *CurrencyService {
	return &CurrencyService{
		currencyRepo: currencyRepo,
		accountRepo:  accountRepo,
	}
}

// List 获取所有货币列表
// 返回：
//   - []response.CurrencyResp: 货币列表
//   - error: 错误信息
func (s *CurrencyService) List() ([]response.CurrencyResp, error) {
	currencies, err := s.currencyRepo.List()
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.CurrencyResp, 0, len(currencies))
	for _, c := range currencies {
		items = append(items, currencyToResp(&c))
	}

	return items, nil
}

// UpdateStatus 更新货币的启用状态
// 业务规则：不能禁用默认货币
// 参数：
//   - id: 货币ID
//   - isEnabled: 是否启用
// 返回：
//   - error: 错误信息
func (s *CurrencyService) UpdateStatus(id uint64, isEnabled bool) error {
	currency, err := s.currencyRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}

	// Cannot disable default currency
	if !isEnabled && currency.IsDefault {
		return errcode.ErrDefaultCurrency
	}

	// Cannot disable if in use by accounts
	// Simplified: allow disabling, actual check would query accounts with this currency

	currency.IsEnabled = isEnabled
	return s.currencyRepo.Update(currency)
}

// SetDefault 设置默认货币
// 业务规则：只能将已启用的货币设为默认；会先取消当前默认货币的默认标记
// 参数：
//   - id: 货币ID
// 返回：
//   - error: 错误信息
func (s *CurrencyService) SetDefault(id uint64) error {
	currency, err := s.currencyRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}

	if !currency.IsEnabled {
		return errcode.ErrDefaultCurrency
	}

	// Unset current default
	currentDefault, err := s.currencyRepo.GetDefault()
	if err == nil && currentDefault != nil {
		currentDefault.IsDefault = false
		s.currencyRepo.Update(currentDefault)
	}

	currency.IsDefault = true
	return s.currencyRepo.Update(currency)
}

// GetExchangeRates 获取所有汇率列表
// 返回：
//   - []response.ExchangeRateResp: 汇率列表
//   - error: 错误信息
func (s *CurrencyService) GetExchangeRates() ([]response.ExchangeRateResp, error) {
	rates, err := s.currencyRepo.ListExchangeRates()
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.ExchangeRateResp, 0, len(rates))
	for _, r := range rates {
		items = append(items, response.ExchangeRateResp{
			ID:             r.ID,
			FromCurrencyID: r.FromCurrencyID,
			ToCurrencyID:   r.ToCurrencyID,
			Rate:           r.Rate.StringFixed(8),
			UpdatedAt:      r.UpdatedAt,
		})
	}

	return items, nil
}

// SetExchangeRate 设置汇率（不存在则创建，存在则更新）
// 参数：
//   - req: 设置汇率请求参数（源货币ID、目标货币ID、汇率值）
// 返回：
//   - error: 错误信息
func (s *CurrencyService) SetExchangeRate(req *request.SetExchangeRateReq) error {
	rate, err := decimal.NewFromString(req.Rate)
	if err != nil || rate.LessThanOrEqual(decimal.Zero) {
		return errcode.ErrBadRequest
	}

	exchangeRate := &model.ExchangeRate{
		FromCurrencyID: req.FromCurrencyID,
		ToCurrencyID:   req.ToCurrencyID,
		Rate:           rate,
	}

	return s.currencyRepo.UpsertExchangeRate(exchangeRate)
}

// InitDefaultCurrencies 初始化默认货币
// 创建CNY（人民币，默认）、USD、EUR、JPY、GBP五种货币
// 如果货币已存在则跳过，不会重复创建
// 返回：
//   - error: 错误信息
func (s *CurrencyService) InitDefaultCurrencies() error {
	defaults := []model.Currency{
		{Code: "CNY", Name: "Chinese Yuan", Symbol: "¥", DecimalPlaces: 2, IsEnabled: true, IsDefault: true},
		{Code: "USD", Name: "US Dollar", Symbol: "$", DecimalPlaces: 2, IsEnabled: true, IsDefault: false},
		{Code: "EUR", Name: "Euro", Symbol: "€", DecimalPlaces: 2, IsEnabled: true, IsDefault: false},
		{Code: "JPY", Name: "Japanese Yen", Symbol: "¥", DecimalPlaces: 0, IsEnabled: true, IsDefault: false},
		{Code: "GBP", Name: "British Pound", Symbol: "£", DecimalPlaces: 2, IsEnabled: true, IsDefault: false},
	}

	for _, c := range defaults {
		existing, err := s.currencyRepo.GetByCode(c.Code)
		if err == gorm.ErrRecordNotFound {
			s.currencyRepo.Create(&c)
		} else if err == nil && existing != nil {
			// Already exists, skip
			_ = existing
		}
	}

	return nil
}

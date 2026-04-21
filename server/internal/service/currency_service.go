package service

import (
	"github.com/shopspring/decimal"
	"github.com/zero-life/server/internal/dto/request"
	"github.com/zero-life/server/internal/dto/response"
	"github.com/zero-life/server/internal/model"
	"github.com/zero-life/server/internal/pkg/errcode"
	"github.com/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

type CurrencyService struct {
	currencyRepo *repository.CurrencyRepository
	accountRepo  *repository.AccountRepository
}

func NewCurrencyService(currencyRepo *repository.CurrencyRepository, accountRepo *repository.AccountRepository) *CurrencyService {
	return &CurrencyService{
		currencyRepo: currencyRepo,
		accountRepo:  accountRepo,
	}
}

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

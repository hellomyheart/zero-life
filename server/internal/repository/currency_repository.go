package repository

import (
	"github.com/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type CurrencyRepository struct {
	db *gorm.DB
}

func NewCurrencyRepository(db *gorm.DB) *CurrencyRepository {
	return &CurrencyRepository{db: db}
}

func (r *CurrencyRepository) List() ([]model.Currency, error) {
	var currencies []model.Currency
	if err := r.db.Order("id ASC").Find(&currencies).Error; err != nil {
		return nil, err
	}
	return currencies, nil
}

func (r *CurrencyRepository) GetByID(id uint64) (*model.Currency, error) {
	var currency model.Currency
	if err := r.db.First(&currency, id).Error; err != nil {
		return nil, err
	}
	return &currency, nil
}

func (r *CurrencyRepository) Update(currency *model.Currency) error {
	return r.db.Save(currency).Error
}

func (r *CurrencyRepository) GetByCode(code string) (*model.Currency, error) {
	var currency model.Currency
	if err := r.db.Where("code = ?", code).First(&currency).Error; err != nil {
		return nil, err
	}
	return &currency, nil
}

func (r *CurrencyRepository) GetDefault() (*model.Currency, error) {
	var currency model.Currency
	if err := r.db.Where("is_default = ?", true).First(&currency).Error; err != nil {
		return nil, err
	}
	return &currency, nil
}

func (r *CurrencyRepository) ListExchangeRates() ([]model.ExchangeRate, error) {
	var rates []model.ExchangeRate
	if err := r.db.Order("id ASC").Find(&rates).Error; err != nil {
		return nil, err
	}
	return rates, nil
}

func (r *CurrencyRepository) UpsertExchangeRate(rate *model.ExchangeRate) error {
	return r.db.Save(rate).Error
}

func (r *CurrencyRepository) Create(currency *model.Currency) error {
	return r.db.Create(currency).Error
}

package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type CurrencyRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewCurrencyRepository(readDB, writeDB *gorm.DB) *CurrencyRepository {
	return &CurrencyRepository{readDB: readDB, writeDB: writeDB}
}

func (r *CurrencyRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.writeDB.Transaction(fn)
}

func (r *CurrencyRepository) List() ([]model.Currency, error) {
	var currencies []model.Currency
	if err := r.readDB.Order("id ASC").Find(&currencies).Error; err != nil {
		return nil, err
	}
	return currencies, nil
}

func (r *CurrencyRepository) GetByID(id uint64) (*model.Currency, error) {
	var currency model.Currency
	if err := r.readDB.First(&currency, id).Error; err != nil {
		return nil, err
	}
	return &currency, nil
}

func (r *CurrencyRepository) Update(currency *model.Currency) error {
	return r.writeDB.Save(currency).Error
}

func (r *CurrencyRepository) GetByCode(code string) (*model.Currency, error) {
	var currency model.Currency
	if err := r.readDB.Where("code = ?", code).First(&currency).Error; err != nil {
		return nil, err
	}
	return &currency, nil
}

func (r *CurrencyRepository) GetDefault() (*model.Currency, error) {
	var currency model.Currency
	if err := r.readDB.Where("is_default = ?", true).First(&currency).Error; err != nil {
		return nil, err
	}
	return &currency, nil
}

func (r *CurrencyRepository) ListExchangeRates() ([]model.ExchangeRate, error) {
	var rates []model.ExchangeRate
	if err := r.readDB.Order("id ASC").Find(&rates).Error; err != nil {
		return nil, err
	}
	return rates, nil
}

func (r *CurrencyRepository) UpsertExchangeRate(rate *model.ExchangeRate) error {
	return r.writeDB.Save(rate).Error
}

func (r *CurrencyRepository) Create(currency *model.Currency) error {
	return r.writeDB.Create(currency).Error
}

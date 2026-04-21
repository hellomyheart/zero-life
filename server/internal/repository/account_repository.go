package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(account *model.Account) error {
	return r.db.Create(account).Error
}

func (r *AccountRepository) GetByID(id, userID uint64) (*model.Account, error) {
	var account model.Account
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).Preload("Currency").First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) List(userID uint64, accountType, search, sort string, offset, limit int) ([]model.Account, error) {
	var accounts []model.Account
	query := r.db.Where("user_id = ?", userID)

	if accountType != "" {
		query = query.Where("type = ?", accountType)
	}
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Apply sorting
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

	if err := query.Preload("Currency").Offset(offset).Limit(limit).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

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

func (r *AccountRepository) Update(account *model.Account) error {
	return r.db.Save(account).Error
}

func (r *AccountRepository) Delete(id, userID uint64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Account{}).Error
}

func (r *AccountRepository) HasTransactions(accountID, userID uint64) (bool, error) {
	var count int64
	if err := r.db.Model(&model.Transaction{}).
		Where("user_id = ? AND (source_id = ? OR destination_id = ?)", userID, accountID, accountID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

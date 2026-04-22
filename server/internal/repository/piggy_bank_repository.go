package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type PiggyBankRepository struct {
	db *gorm.DB
}

func NewPiggyBankRepository(db *gorm.DB) *PiggyBankRepository {
	return &PiggyBankRepository{db: db}
}

func (r *PiggyBankRepository) Create(piggyBank *model.PiggyBank) error {
	return r.db.Create(piggyBank).Error
}

func (r *PiggyBankRepository) GetByID(id, userID uint64) (*model.PiggyBank, error) {
	var piggyBank model.PiggyBank
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).
		Preload("Account").
		Preload("Events", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		First(&piggyBank).Error; err != nil {
		return nil, err
	}
	return &piggyBank, nil
}

func (r *PiggyBankRepository) List(userID uint64) ([]model.PiggyBank, error) {
	var piggyBanks []model.PiggyBank
	if err := r.db.Where("user_id = ?", userID).
		Preload("Account").
		Order("name ASC").Find(&piggyBanks).Error; err != nil {
		return nil, err
	}
	return piggyBanks, nil
}

func (r *PiggyBankRepository) Update(piggyBank *model.PiggyBank) error {
	return r.db.Save(piggyBank).Error
}

func (r *PiggyBankRepository) Delete(id, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("piggy_bank_id = ?", id).Delete(&model.PiggyEvent{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&model.PiggyBank{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *PiggyBankRepository) CreateEvent(event *model.PiggyEvent) error {
	return r.db.Create(event).Error
}

func (r *PiggyBankRepository) ListEvents(piggyBankID uint64) ([]model.PiggyEvent, error) {
	var events []model.PiggyEvent
	if err := r.db.Where("piggy_bank_id = ?", piggyBankID).
		Order("created_at DESC").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

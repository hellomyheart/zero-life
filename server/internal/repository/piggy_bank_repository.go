package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type PiggyBankRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewPiggyBankRepository(readDB, writeDB *gorm.DB) *PiggyBankRepository {
	return &PiggyBankRepository{readDB: readDB, writeDB: writeDB}
}

func (r *PiggyBankRepository) Create(piggyBank *model.PiggyBank) error {
	return r.writeDB.Create(piggyBank).Error
}

func (r *PiggyBankRepository) GetByID(id, userID uint64) (*model.PiggyBank, error) {
	var piggyBank model.PiggyBank
	if err := r.readDB.Where("id = ? AND user_id = ?", id, userID).
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
	if err := r.readDB.Where("user_id = ?", userID).
		Preload("Account").
		Order("name ASC").Find(&piggyBanks).Error; err != nil {
		return nil, err
	}
	return piggyBanks, nil
}

func (r *PiggyBankRepository) Update(piggyBank *model.PiggyBank) error {
	return r.writeDB.Save(piggyBank).Error
}

func (r *PiggyBankRepository) Delete(id, userID uint64) error {
	return r.writeDB.Transaction(func(tx *gorm.DB) error {
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
	return r.writeDB.Create(event).Error
}

func (r *PiggyBankRepository) ListEvents(piggyBankID uint64) ([]model.PiggyEvent, error) {
	var events []model.PiggyEvent
	if err := r.readDB.Where("piggy_bank_id = ?", piggyBankID).
		Order("created_at DESC").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *PiggyBankRepository) Reorder(userID uint64, orders map[uint64]int) error {
	return r.writeDB.Transaction(func(tx *gorm.DB) error {
		for id, order := range orders {
			if err := tx.Model(&model.PiggyBank{}).Where("id = ? AND user_id = ?", id, userID).Update("order", order).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PiggyBankRepository) DeleteEvents(piggyBankID uint64) error {
	return r.writeDB.Where("piggy_bank_id = ?", piggyBankID).Delete(&model.PiggyEvent{}).Error
}

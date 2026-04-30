package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/shopspring/decimal"
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
		First(&piggyBank).Error; err != nil {
		return nil, err
	}
	return &piggyBank, nil
}

func (r *PiggyBankRepository) List(userID uint64) ([]model.PiggyBank, error) {
	var piggyBanks []model.PiggyBank
	if err := r.readDB.Where("user_id = ?", userID).
		Preload("Account").
		Order("`order` ASC, name ASC").Find(&piggyBanks).Error; err != nil {
		return nil, err
	}
	return piggyBanks, nil
}

func (r *PiggyBankRepository) Update(piggyBank *model.PiggyBank) error {
	return r.writeDB.Save(piggyBank).Error
}

func (r *PiggyBankRepository) UpdateWithDB(db *gorm.DB, piggyBank *model.PiggyBank) error {
	return db.Save(piggyBank).Error
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

func (r *PiggyBankRepository) CreateEventWithDB(db *gorm.DB, event *model.PiggyEvent) error {
	return db.Create(event).Error
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

func (r *PiggyBankRepository) DeleteEventsWithDB(db *gorm.DB, piggyBankID uint64) error {
	return db.Where("piggy_bank_id = ?", piggyBankID).Delete(&model.PiggyEvent{}).Error
}

func (r *PiggyBankRepository) AddAmountWithDB(db *gorm.DB, id, userID uint64, amount, targetAmount decimal.Decimal) (bool, error) {
	result := db.Model(&model.PiggyBank{}).
		Where("id = ? AND user_id = ? AND current_amount + ? <= target_amount", id, userID, amount).
		Update("current_amount", gorm.Expr("current_amount + ?", amount))
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *PiggyBankRepository) RemoveAmountWithDB(db *gorm.DB, id, userID uint64, amount decimal.Decimal) (bool, error) {
	result := db.Model(&model.PiggyBank{}).
		Where("id = ? AND user_id = ? AND current_amount >= ?", id, userID, amount).
		Update("current_amount", gorm.Expr("current_amount - ?", amount))
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *PiggyBankRepository) ResetAmountWithDB(db *gorm.DB, id, userID uint64) error {
	return db.Model(&model.PiggyBank{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("current_amount", decimal.Zero).Error
}

func (r *PiggyBankRepository) SumCurrentAmountByAccount(accountID, userID uint64) (decimal.Decimal, error) {
	var sum *decimal.Decimal
	if err := r.readDB.Model(&model.PiggyBank{}).
		Select("COALESCE(SUM(current_amount), 0)").
		Where("account_id = ? AND user_id = ?", accountID, userID).
		Scan(&sum).Error; err != nil {
		return decimal.Zero, err
	}
	if sum == nil {
		return decimal.Zero, nil
	}
	return *sum, nil
}

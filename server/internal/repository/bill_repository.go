package repository

import (
	"time"

	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type BillRepository struct {
	db *gorm.DB
}

func NewBillRepository(db *gorm.DB) *BillRepository {
	return &BillRepository{db: db}
}

func (r *BillRepository) Create(bill *model.Bill) error {
	return r.db.Create(bill).Error
}

func (r *BillRepository) GetByID(id, userID uint64) (*model.Bill, error) {
	var bill model.Bill
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&bill).Error; err != nil {
		return nil, err
	}
	return &bill, nil
}

func (r *BillRepository) List(userID uint64) ([]model.Bill, error) {
	var bills []model.Bill
	if err := r.db.Where("user_id = ?", userID).Order("next_due ASC").Find(&bills).Error; err != nil {
		return nil, err
	}
	return bills, nil
}

func (r *BillRepository) Update(bill *model.Bill) error {
	return r.db.Save(bill).Error
}

func (r *BillRepository) Delete(id, userID uint64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Bill{}).Error
}

func (r *BillRepository) GetUpcoming(userID uint64, days int) ([]model.Bill, error) {
	var bills []model.Bill
	now := time.Now()
	endDate := now.AddDate(0, 0, days)
	if err := r.db.Where("user_id = ? AND next_due BETWEEN ? AND ?", userID, now, endDate).
		Order("next_due ASC").Find(&bills).Error; err != nil {
		return nil, err
	}
	return bills, nil
}

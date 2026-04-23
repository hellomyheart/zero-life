package repository

import (
	"time"

	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type RecurrenceRepository struct {
	db *gorm.DB
}

func NewRecurrenceRepository(db *gorm.DB) *RecurrenceRepository {
	return &RecurrenceRepository{db: db}
}

func (r *RecurrenceRepository) Create(rec *model.Recurrence) error {
	return r.db.Create(rec).Error
}

func (r *RecurrenceRepository) GetByID(id, userID uint64) (*model.Recurrence, error) {
	var rec model.Recurrence
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&rec).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *RecurrenceRepository) List(userID uint64, offset, limit int) ([]model.Recurrence, error) {
	var recs []model.Recurrence
	if err := r.db.Where("user_id = ?", userID).Order("next_date ASC").Offset(offset).Limit(limit).Find(&recs).Error; err != nil {
		return nil, err
	}
	return recs, nil
}

func (r *RecurrenceRepository) Count(userID uint64) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Recurrence{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RecurrenceRepository) Update(rec *model.Recurrence) error {
	return r.db.Save(rec).Error
}

func (r *RecurrenceRepository) Delete(id, userID uint64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Recurrence{}).Error
}

func (r *RecurrenceRepository) GetDueRecurrences(today time.Time) ([]model.Recurrence, error) {
	var recs []model.Recurrence
	if err := r.db.Where("is_active = ? AND next_date <= ?", true, today).
		Order("next_date ASC").Find(&recs).Error; err != nil {
		return nil, err
	}
	return recs, nil
}

package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type PreferenceRepository struct {
	db *gorm.DB
}

func NewPreferenceRepository(db *gorm.DB) *PreferenceRepository {
	return &PreferenceRepository{db: db}
}

func (r *PreferenceRepository) Get(userID uint64, name string) (*model.Preference, error) {
	var pref model.Preference
	if err := r.db.Where("user_id = ? AND name = ?", userID, name).First(&pref).Error; err != nil {
		return nil, err
	}
	return &pref, nil
}

func (r *PreferenceRepository) Set(userID uint64, name, value string) error {
	var pref model.Preference
	err := r.db.Where("user_id = ? AND name = ?", userID, name).First(&pref).Error
	if err == gorm.ErrRecordNotFound {
		pref = model.Preference{
			UserID: userID,
			Name:   name,
			Value:  value,
		}
		return r.db.Create(&pref).Error
	}
	if err != nil {
		return err
	}
	pref.Value = value
	return r.db.Save(&pref).Error
}

func (r *PreferenceRepository) List(userID uint64) ([]model.Preference, error) {
	var prefs []model.Preference
	if err := r.db.Where("user_id = ?", userID).Order("name ASC").Find(&prefs).Error; err != nil {
		return nil, err
	}
	return prefs, nil
}

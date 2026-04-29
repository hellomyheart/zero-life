package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type PreferenceRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewPreferenceRepository(readDB, writeDB *gorm.DB) *PreferenceRepository {
	return &PreferenceRepository{readDB: readDB, writeDB: writeDB}
}

func (r *PreferenceRepository) Get(userID uint64, key string) (*model.Preference, error) {
	var pref model.Preference
	if err := r.readDB.Where("user_id = ? AND key = ?", userID, key).First(&pref).Error; err != nil {
		return nil, err
	}
	return &pref, nil
}

func (r *PreferenceRepository) Set(userID uint64, key, value string) error {
	var pref model.Preference
	err := r.writeDB.Where("user_id = ? AND key = ?", userID, key).First(&pref).Error
	if err == gorm.ErrRecordNotFound {
		pref = model.Preference{UserID: userID, Key: key, Value: value}
		return r.writeDB.Create(&pref).Error
	}
	if err != nil {
		return err
	}
	pref.Value = value
	return r.writeDB.Save(&pref).Error
}

func (r *PreferenceRepository) List(userID uint64) ([]model.Preference, error) {
	var prefs []model.Preference
	if err := r.readDB.Where("user_id = ?", userID).Find(&prefs).Error; err != nil {
		return nil, err
	}
	return prefs, nil
}

func (r *PreferenceRepository) Delete(userID uint64, key string) error {
	return r.writeDB.Where("user_id = ? AND key = ?", userID, key).Delete(&model.Preference{}).Error
}

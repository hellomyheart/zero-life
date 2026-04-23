package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type ConfigurationRepository struct {
	db *gorm.DB
}

func NewConfigurationRepository(db *gorm.DB) *ConfigurationRepository {
	return &ConfigurationRepository{db: db}
}

func (r *ConfigurationRepository) Get(name string) (*model.Configuration, error) {
	var cfg model.Configuration
	if err := r.db.Where("name = ?", name).First(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (r *ConfigurationRepository) Set(cfg *model.Configuration) error {
	var existing model.Configuration
	err := r.db.Where("name = ?", cfg.Name).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.Create(cfg).Error
	}
	if err != nil {
		return err
	}
	existing.Value = cfg.Value
	return r.db.Save(&existing).Error
}

func (r *ConfigurationRepository) List() ([]model.Configuration, error) {
	var configs []model.Configuration
	if err := r.db.Order("name ASC").Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

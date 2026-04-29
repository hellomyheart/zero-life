package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type ConfigurationRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewConfigurationRepository(readDB, writeDB *gorm.DB) *ConfigurationRepository {
	return &ConfigurationRepository{readDB: readDB, writeDB: writeDB}
}

func (r *ConfigurationRepository) Get(name string) (*model.Configuration, error) {
	var cfg model.Configuration
	if err := r.readDB.Where("name = ?", name).First(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (r *ConfigurationRepository) Set(cfg *model.Configuration) error {
	var existing model.Configuration
	err := r.writeDB.Where("name = ?", cfg.Name).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.writeDB.Create(cfg).Error
	}
	if err != nil {
		return err
	}
	existing.Value = cfg.Value
	return r.writeDB.Save(&existing).Error
}

func (r *ConfigurationRepository) List() ([]model.Configuration, error) {
	var configs []model.Configuration
	if err := r.readDB.Order("name ASC").Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

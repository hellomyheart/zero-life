package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type LinkTypeRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewLinkTypeRepository(readDB, writeDB *gorm.DB) *LinkTypeRepository {
	return &LinkTypeRepository{readDB: readDB, writeDB: writeDB}
}

func (r *LinkTypeRepository) Create(linkType *model.LinkType) error {
	return r.writeDB.Create(linkType).Error
}

func (r *LinkTypeRepository) GetByID(id uint64) (*model.LinkType, error) {
	var linkType model.LinkType
	if err := r.readDB.Where("id = ?", id).First(&linkType).Error; err != nil {
		return nil, err
	}
	return &linkType, nil
}

func (r *LinkTypeRepository) List(offset, limit int) ([]model.LinkType, error) {
	var linkTypes []model.LinkType
	if err := r.readDB.Order("id ASC").Offset(offset).Limit(limit).Find(&linkTypes).Error; err != nil {
		return nil, err
	}
	return linkTypes, nil
}

func (r *LinkTypeRepository) Count() (int64, error) {
	var count int64
	if err := r.readDB.Model(&model.LinkType{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *LinkTypeRepository) Update(linkType *model.LinkType) error {
	return r.writeDB.Save(linkType).Error
}

func (r *LinkTypeRepository) Delete(id uint64) error {
	return r.writeDB.Where("id = ?", id).Delete(&model.LinkType{}).Error
}

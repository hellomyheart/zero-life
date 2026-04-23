package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type LinkTypeRepository struct {
	db *gorm.DB
}

func NewLinkTypeRepository(db *gorm.DB) *LinkTypeRepository {
	return &LinkTypeRepository{db: db}
}

func (r *LinkTypeRepository) Create(linkType *model.LinkType) error {
	return r.db.Create(linkType).Error
}

func (r *LinkTypeRepository) GetByID(id uint64) (*model.LinkType, error) {
	var linkType model.LinkType
	if err := r.db.Where("id = ?", id).First(&linkType).Error; err != nil {
		return nil, err
	}
	return &linkType, nil
}

func (r *LinkTypeRepository) List(offset, limit int) ([]model.LinkType, error) {
	var linkTypes []model.LinkType
	if err := r.db.Order("id ASC").Offset(offset).Limit(limit).Find(&linkTypes).Error; err != nil {
		return nil, err
	}
	return linkTypes, nil
}

func (r *LinkTypeRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&model.LinkType{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *LinkTypeRepository) Update(linkType *model.LinkType) error {
	return r.db.Save(linkType).Error
}

func (r *LinkTypeRepository) Delete(id uint64) error {
	return r.db.Where("id = ?", id).Delete(&model.LinkType{}).Error
}

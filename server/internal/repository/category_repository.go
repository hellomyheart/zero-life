package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) GetByID(id, userID uint64) (*model.Category, error) {
	var category model.Category
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) List(userID uint64) ([]model.Category, error) {
	var categories []model.Category
	if err := r.db.Where("user_id = ?", userID).Order("sort_order ASC, id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

func (r *CategoryRepository) Delete(id, userID uint64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Category{}).Error
}

func (r *CategoryRepository) GetSubCategories(parentID, userID uint64) ([]model.Category, error) {
	var categories []model.Category
	if err := r.db.Where("parent_id = ? AND user_id = ?", parentID, userID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

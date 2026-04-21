package repository

import (
	"github.com/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

func (r *TagRepository) GetByID(id, userID uint64) (*model.Tag, error) {
	var tag model.Tag
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepository) List(userID uint64) ([]model.Tag, error) {
	var tags []model.Tag
	if err := r.db.Where("user_id = ?", userID).Order("id ASC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepository) Update(tag *model.Tag) error {
	return r.db.Save(tag).Error
}

func (r *TagRepository) Delete(id, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete transaction associations
		if err := tx.Where("tag_id = ?", id).Delete(&model.TransactionTag{}).Error; err != nil {
			return err
		}
		// Delete the tag
		if err := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Tag{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *TagRepository) CountTransactions(tagID uint64) (int64, error) {
	var count int64
	if err := r.db.Model(&model.TransactionTag{}).Where("tag_id = ?", tagID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

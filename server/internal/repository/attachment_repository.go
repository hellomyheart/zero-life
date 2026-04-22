package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type AttachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) *AttachmentRepository {
	return &AttachmentRepository{db: db}
}

func (r *AttachmentRepository) Create(attachment *model.Attachment) error {
	return r.db.Create(attachment).Error
}

func (r *AttachmentRepository) GetByID(id, userID uint64) (*model.Attachment, error) {
	var attachment model.Attachment
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&attachment).Error; err != nil {
		return nil, err
	}
	return &attachment, nil
}

func (r *AttachmentRepository) List(userID uint64, attachableType string, attachableID uint64) ([]model.Attachment, error) {
	var attachments []model.Attachment
	query := r.db.Where("user_id = ?", userID)
	if attachableType != "" {
		query = query.Where("attachable_type = ?", attachableType)
	}
	if attachableID != 0 {
		query = query.Where("attachable_id = ?", attachableID)
	}
	if err := query.Order("created_at DESC").Find(&attachments).Error; err != nil {
		return nil, err
	}
	return attachments, nil
}

func (r *AttachmentRepository) Delete(id, userID uint64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Attachment{}).Error
}

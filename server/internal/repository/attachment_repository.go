package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type AttachmentRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewAttachmentRepository(readDB, writeDB *gorm.DB) *AttachmentRepository {
	return &AttachmentRepository{readDB: readDB, writeDB: writeDB}
}

func (r *AttachmentRepository) Create(attachment *model.Attachment) error {
	return r.writeDB.Create(attachment).Error
}

func (r *AttachmentRepository) GetByID(id, userID uint64) (*model.Attachment, error) {
	var attachment model.Attachment
	if err := r.readDB.Where("id = ? AND user_id = ?", id, userID).First(&attachment).Error; err != nil {
		return nil, err
	}
	return &attachment, nil
}

func (r *AttachmentRepository) List(userID uint64, attachableType string, attachableID uint64) ([]model.Attachment, error) {
	var attachments []model.Attachment
	query := r.readDB.Where("user_id = ?", userID)
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
	return r.writeDB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Attachment{}).Error
}

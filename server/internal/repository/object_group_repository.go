package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type ObjectGroupRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewObjectGroupRepository(readDB, writeDB *gorm.DB) *ObjectGroupRepository {
	return &ObjectGroupRepository{readDB: readDB, writeDB: writeDB}
}

func (r *ObjectGroupRepository) Create(og *model.ObjectGroup) error {
	return r.writeDB.Create(og).Error
}

func (r *ObjectGroupRepository) GetByID(id, userID uint64) (*model.ObjectGroup, error) {
	var og model.ObjectGroup
	if err := r.readDB.Where("id = ? AND user_id = ?", id, userID).First(&og).Error; err != nil {
		return nil, err
	}
	return &og, nil
}

func (r *ObjectGroupRepository) List(userID uint64, groupableType string) ([]model.ObjectGroup, error) {
	var groups []model.ObjectGroup
	query := r.readDB.Where("user_id = ?", userID)
	if groupableType != "" {
		query = query.Where("groupable_type = ?", groupableType)
	}
	if err := query.Order("name ASC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *ObjectGroupRepository) Update(og *model.ObjectGroup) error {
	return r.writeDB.Save(og).Error
}

func (r *ObjectGroupRepository) Delete(id, userID uint64) error {
	return r.writeDB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.ObjectGroup{}).Error
}

func (r *ObjectGroupRepository) GetByGroupable(userID uint64, groupableType string, groupableID uint64) (*model.ObjectGroup, error) {
	var og model.ObjectGroup
	if err := r.readDB.Where("user_id = ? AND groupable_type = ? AND groupable_id = ?", userID, groupableType, groupableID).
		First(&og).Error; err != nil {
		return nil, err
	}
	return &og, nil
}

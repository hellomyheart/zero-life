package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type RuleGroupRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewRuleGroupRepository(readDB, writeDB *gorm.DB) *RuleGroupRepository {
	return &RuleGroupRepository{readDB: readDB, writeDB: writeDB}
}

func (r *RuleGroupRepository) Create(group *model.RuleGroup) error {
	return r.writeDB.Create(group).Error
}

func (r *RuleGroupRepository) GetByID(id, userID uint64) (*model.RuleGroup, error) {
	var group model.RuleGroup
	if err := r.readDB.Where("id = ? AND user_id = ?", id, userID).First(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *RuleGroupRepository) List(userID uint64) ([]model.RuleGroup, error) {
	var groups []model.RuleGroup
	if err := r.readDB.Where("user_id = ?", userID).Order("`order` ASC, id ASC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *RuleGroupRepository) Update(group *model.RuleGroup) error {
	return r.writeDB.Save(group).Error
}

func (r *RuleGroupRepository) Delete(id, userID uint64) error {
	return r.writeDB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.RuleGroup{}).Error
}

func (r *RuleGroupRepository) HasRules(groupID, userID uint64) (bool, error) {
	var count int64
	if err := r.readDB.Model(&model.Rule{}).Where("group_id = ? AND user_id = ?", groupID, userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RuleGroupRepository) GetRuleCount(groupID, userID uint64) (int64, error) {
	var count int64
	if err := r.readDB.Model(&model.Rule{}).Where("group_id = ? AND user_id = ?", groupID, userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

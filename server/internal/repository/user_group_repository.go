// Package repository 数据访问层，封装数据库操作
package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// UserGroupRepository 用户组数据访问对象
type UserGroupRepository struct {
	db *gorm.DB
}

// NewUserGroupRepository 创建用户组数据访问对象实例
func NewUserGroupRepository(db *gorm.DB) *UserGroupRepository {
	return &UserGroupRepository{db: db}
}

// Create 创建用户组
func (r *UserGroupRepository) Create(group *model.UserGroup) error {
	return r.db.Create(group).Error
}

// GetByID 根据ID获取用户组
func (r *UserGroupRepository) GetByID(id uint64) (*model.UserGroup, error) {
	var group model.UserGroup
	err := r.db.First(&group, id).Error
	return &group, err
}

// List 获取用户组列表
func (r *UserGroupRepository) List(userID uint64, offset, limit int) ([]model.UserGroup, error) {
	var groups []model.UserGroup
	err := r.db.Joins("JOIN user_group_members ON user_group_members.user_group_id = user_groups.id").
		Where("user_group_members.user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Find(&groups).Error
	return groups, err
}

// Update 更新用户组
func (r *UserGroupRepository) Update(group *model.UserGroup) error {
	return r.db.Save(group).Error
}

// Delete 删除用户组
func (r *UserGroupRepository) Delete(id uint64) error {
	return r.db.Delete(&model.UserGroup{}, id).Error
}

// AddMember 添加成员
func (r *UserGroupRepository) AddMember(member *model.UserGroupMember) error {
	return r.db.Create(member).Error
}

// RemoveMember 移除成员
func (r *UserGroupRepository) RemoveMember(groupID, userID uint64) error {
	return r.db.Where("user_group_id = ? AND user_id = ?", groupID, userID).
		Delete(&model.UserGroupMember{}).Error
}

// GetMembers 获取用户组成员列表
func (r *UserGroupRepository) GetMembers(groupID uint64) ([]model.UserGroupMember, error) {
	var members []model.UserGroupMember
	err := r.db.Where("user_group_id = ?", groupID).
		Preload("User").
		Find(&members).Error
	return members, err
}

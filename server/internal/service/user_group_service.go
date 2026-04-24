// Package service 业务逻辑层，实现核心业务逻辑
package service

import (
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/pkg/pagination"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

// UserGroupService 用户组业务服务
type UserGroupService struct {
	groupRepo *repository.UserGroupRepository
}

// NewUserGroupService 创建用户组服务实例
func NewUserGroupService(groupRepo *repository.UserGroupRepository) *UserGroupService {
	return &UserGroupService{groupRepo: groupRepo}
}

// Create 创建用户组
func (s *UserGroupService) Create(userID uint64, req *request.CreateUserGroupReq) (*response.UserGroupResp, error) {
	group := &model.UserGroup{
		Title: req.Title,
	}

	if err := s.groupRepo.Create(group); err != nil {
		return nil, errcode.ErrInternal
	}

	// 创建者自动成为所有者
	member := &model.UserGroupMember{
		UserGroupID: group.ID,
		UserID:      userID,
		IsOwner:     true,
	}
	if err := s.groupRepo.AddMember(member); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(group), nil
}

// Get 获取用户组详情
func (s *UserGroupService) Get(id uint64) (*response.UserGroupResp, error) {
	group, err := s.groupRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(group), nil
}

// List 获取用户组列表
func (s *UserGroupService) List(userID uint64, req *request.UserGroupListReq) (*pagination.Result, error) {
	params := pagination.Params{Page: req.Page, PageSize: req.PageSize}
	params.Normalize()

	groups, err := s.groupRepo.List(userID, params.Offset(), params.PageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.UserGroupResp, 0, len(groups))
	for _, g := range groups {
		items = append(items, *s.toResp(&g))
	}

	return pagination.NewResult(items, int64(len(items)), params), nil
}

// Update 更新用户组
func (s *UserGroupService) Update(id uint64, req *request.UpdateUserGroupReq) (*response.UserGroupResp, error) {
	group, err := s.groupRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Title != "" {
		group.Title = req.Title
	}

	if err := s.groupRepo.Update(group); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(group), nil
}

// Delete 删除用户组
func (s *UserGroupService) Delete(id uint64) error {
	return s.groupRepo.Delete(id)
}

// toResp 将用户组模型转换为响应对象
func (s *UserGroupService) toResp(group *model.UserGroup) *response.UserGroupResp {
	return &response.UserGroupResp{
		ID:        group.ID,
		Title:     group.Title,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
	}
}

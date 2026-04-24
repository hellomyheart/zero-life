// Package service 业务逻辑层，实现核心业务逻辑
// ObjectGroupService 对象分组业务逻辑，管理实体分组排序
package service

import (
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

type ObjectGroupService struct {
	ogRepo *repository.ObjectGroupRepository
}

func NewObjectGroupService(ogRepo *repository.ObjectGroupRepository) *ObjectGroupService {
	return &ObjectGroupService{ogRepo: ogRepo}
}

func (s *ObjectGroupService) Create(userID uint64, req *request.CreateObjectGroupReq) (*response.ObjectGroupResp, error) {
	og := &model.ObjectGroup{
		UserID:        userID,
		Name:          req.Name,
		GroupableType: req.GroupableType,
		GroupableID:   req.GroupableID,
	}

	if err := s.ogRepo.Create(og); err != nil {
		return nil, errcode.ErrInternal
	}

	created, err := s.ogRepo.GetByID(og.ID, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(created), nil
}

func (s *ObjectGroupService) Get(userID, id uint64) (*response.ObjectGroupResp, error) {
	og, err := s.ogRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(og), nil
}

func (s *ObjectGroupService) List(userID uint64, groupableType string) ([]response.ObjectGroupResp, error) {
	groups, err := s.ogRepo.List(userID, groupableType)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.ObjectGroupResp, 0, len(groups))
	for _, g := range groups {
		items = append(items, *s.toResp(&g))
	}
	return items, nil
}

func (s *ObjectGroupService) Update(userID, id uint64, req *request.UpdateObjectGroupReq) (*response.ObjectGroupResp, error) {
	og, err := s.ogRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Name != "" {
		og.Name = req.Name
	}

	if err := s.ogRepo.Update(og); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(og), nil
}

func (s *ObjectGroupService) Delete(userID, id uint64) error {
	_, err := s.ogRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}
	return s.ogRepo.Delete(id, userID)
}

func (s *ObjectGroupService) toResp(og *model.ObjectGroup) *response.ObjectGroupResp {
	return &response.ObjectGroupResp{
		ID:            og.ID,
		Name:          og.Name,
		GroupableType: og.GroupableType,
		GroupableID:   og.GroupableID,
		CreatedAt:     og.CreatedAt,
		UpdatedAt:     og.UpdatedAt,
	}
}
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

// ObjectGroupService 对象分组服务
// 管理实体的分组排序，支持将账户、分类等实体按自定义顺序排列
// 依赖ogRepo进行对象分组数据访问
type ObjectGroupService struct {
	ogRepo *repository.ObjectGroupRepository // 对象分组数据访问对象
}

// NewObjectGroupService 创建对象分组服务实例
func NewObjectGroupService(ogRepo *repository.ObjectGroupRepository) *ObjectGroupService {
	return &ObjectGroupService{ogRepo: ogRepo}
}

// Create 创建对象分组
// 将实体（如账户、分类）添加到分组中，用于自定义排序
// 参数：
//   - userID: 用户ID
//   - req: 创建请求参数（名称、实体类型、实体ID）
// 返回：
//   - *response.ObjectGroupResp: 创建成功的分组信息
//   - error: 错误信息
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

// Get 获取单个对象分组详情
// 参数：
//   - userID: 用户ID
//   - id: 分组ID
// 返回：
//   - *response.ObjectGroupResp: 分组信息
//   - error: 错误信息
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

// List 获取指定实体类型的分组列表
// 参数：
//   - userID: 用户ID
//   - groupableType: 实体类型（如"account"、"category"）
// 返回：
//   - []response.ObjectGroupResp: 分组列表
//   - error: 错误信息
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

// Update 更新对象分组（仅支持更新名称）
// 参数：
//   - userID: 用户ID
//   - id: 分组ID
//   - req: 更新请求参数
// 返回：
//   - *response.ObjectGroupResp: 更新后的分组信息
//   - error: 错误信息
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

// Delete 删除对象分组
// 参数：
//   - userID: 用户ID
//   - id: 分组ID
// 返回：
//   - error: 错误信息
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

// toResp 将对象分组模型转换为响应对象
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
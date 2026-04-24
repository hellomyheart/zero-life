// Package service 业务逻辑层，实现核心业务逻辑
// LinkTypeService 关联类型业务逻辑，管理交易关联类型和增删改查
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

// LinkTypeService 关联类型服务
// 管理交易关联类型（如"相关"、"报销"等），关联类型定义了交易之间关系的语义
// 关联类型包含正向（outward）和反向（inward）描述，以及是否为有方向关系
// 依赖linkTypeRepo进行关联类型数据访问
type LinkTypeService struct {
	linkTypeRepo *repository.LinkTypeRepository // 关联类型数据访问对象
}

// NewLinkTypeService 创建关联类型服务实例
func NewLinkTypeService(linkTypeRepo *repository.LinkTypeRepository) *LinkTypeService {
	return &LinkTypeService{linkTypeRepo: linkTypeRepo}
}

// Create 创建关联类型
// 参数：
//   - req: 创建请求参数（名称、正向描述、反向描述、是否有方向）
// 返回：
//   - *response.LinkTypeResp: 创建成功的关联类型信息
//   - error: 错误信息
func (s *LinkTypeService) Create(req *request.CreateLinkTypeReq) (*response.LinkTypeResp, error) {
	linkType := &model.LinkType{
		Name:          req.Name,
		Outward:       req.Outward,
		Inward:        req.Inward,
		IsDirectional: req.IsDirectional,
	}

	if err := s.linkTypeRepo.Create(linkType); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(linkType), nil
}

// Get 获取单个关联类型详情
// 参数：
//   - id: 关联类型ID
// 返回：
//   - *response.LinkTypeResp: 关联类型信息
//   - error: 错误信息
func (s *LinkTypeService) Get(id uint64) (*response.LinkTypeResp, error) {
	linkType, err := s.linkTypeRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(linkType), nil
}

// List 获取关联类型分页列表
// 参数：
//   - req: 列表查询参数（含分页）
// 返回：
//   - *pagination.Result: 分页结果
//   - error: 错误信息
func (s *LinkTypeService) List(req *request.LinkTypeListReq) (*pagination.Result, error) {
	params := pagination.Params{Page: req.Page, PageSize: req.PageSize}
	params.Normalize()

	linkTypes, err := s.linkTypeRepo.List(params.Offset(), params.PageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	total, err := s.linkTypeRepo.Count()
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.LinkTypeResp, 0, len(linkTypes))
	for _, lt := range linkTypes {
		items = append(items, *s.toResp(&lt))
	}

	return pagination.NewResult(items, total, params), nil
}

// Update 更新关联类型信息
// 支持部分更新：名称、正向描述、反向描述、是否有方向
// 参数：
//   - id: 关联类型ID
//   - req: 更新请求参数
// 返回：
//   - *response.LinkTypeResp: 更新后的关联类型信息
//   - error: 错误信息
func (s *LinkTypeService) Update(id uint64, req *request.UpdateLinkTypeReq) (*response.LinkTypeResp, error) {
	linkType, err := s.linkTypeRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Name != "" {
		linkType.Name = req.Name
	}
	if req.Outward != "" {
		linkType.Outward = req.Outward
	}
	if req.Inward != "" {
		linkType.Inward = req.Inward
	}
	if req.IsDirectional != nil {
		linkType.IsDirectional = *req.IsDirectional
	}

	if err := s.linkTypeRepo.Update(linkType); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(linkType), nil
}

// Delete 删除关联类型
// 参数：
//   - id: 关联类型ID
// 返回：
//   - error: 错误信息
func (s *LinkTypeService) Delete(id uint64) error {
	_, err := s.linkTypeRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}
	return s.linkTypeRepo.Delete(id)
}

// toResp 将关联类型模型转换为响应对象
func (s *LinkTypeService) toResp(lt *model.LinkType) *response.LinkTypeResp {
	return &response.LinkTypeResp{
		ID:            lt.ID,
		Name:          lt.Name,
		Outward:       lt.Outward,
		Inward:        lt.Inward,
		IsDirectional: lt.IsDirectional,
		CreatedAt:     lt.CreatedAt,
		UpdatedAt:     lt.UpdatedAt,
	}
}

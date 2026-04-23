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

type LinkTypeService struct {
	linkTypeRepo *repository.LinkTypeRepository
}

func NewLinkTypeService(linkTypeRepo *repository.LinkTypeRepository) *LinkTypeService {
	return &LinkTypeService{linkTypeRepo: linkTypeRepo}
}

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

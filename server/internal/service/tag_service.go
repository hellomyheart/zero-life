package service

import (
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

type TagService struct {
	tagRepo *repository.TagRepository
}

func NewTagService(tagRepo *repository.TagRepository) *TagService {
	return &TagService{tagRepo: tagRepo}
}

func (s *TagService) Create(userID uint64, req *request.CreateTagReq) (*response.TagResp, error) {
	// Check name uniqueness
	tags, err := s.tagRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	for _, t := range tags {
		if t.Name == req.Name {
			return nil, errcode.ErrTagNameExists
		}
	}

	tag := &model.Tag{
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
	}
	if tag.Color == "" {
		tag.Color = "#409EFF"
	}

	if err := s.tagRepo.Create(tag); err != nil {
		return nil, errcode.ErrInternal
	}

	count, _ := s.tagRepo.CountTransactions(tag.ID)

	return &response.TagResp{
		ID:               tag.ID,
		Name:             tag.Name,
		Color:            tag.Color,
		TransactionCount: count,
		CreatedAt:        tag.CreatedAt,
		UpdatedAt:        tag.UpdatedAt,
	}, nil
}

func (s *TagService) List(userID uint64) ([]response.TagResp, error) {
	tags, err := s.tagRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.TagResp, 0, len(tags))
	for _, t := range tags {
		count, _ := s.tagRepo.CountTransactions(t.ID)
		items = append(items, response.TagResp{
			ID:               t.ID,
			Name:             t.Name,
			Color:            t.Color,
			TransactionCount: count,
			CreatedAt:        t.CreatedAt,
			UpdatedAt:        t.UpdatedAt,
		})
	}

	return items, nil
}

func (s *TagService) Update(userID, id uint64, req *request.UpdateTagReq) (*response.TagResp, error) {
	tag, err := s.tagRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Name != "" {
		tag.Name = req.Name
	}
	if req.Color != "" {
		tag.Color = req.Color
	}

	if err := s.tagRepo.Update(tag); err != nil {
		return nil, errcode.ErrInternal
	}

	count, _ := s.tagRepo.CountTransactions(tag.ID)

	return &response.TagResp{
		ID:               tag.ID,
		Name:             tag.Name,
		Color:            tag.Color,
		TransactionCount: count,
		CreatedAt:        tag.CreatedAt,
		UpdatedAt:        tag.UpdatedAt,
	}, nil
}

func (s *TagService) Delete(userID, id uint64) error {
	_, err := s.tagRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}

	return s.tagRepo.Delete(id, userID)
}

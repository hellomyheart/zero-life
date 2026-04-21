package service

import (
	"github.com/zero-life/server/internal/dto/request"
	"github.com/zero-life/server/internal/dto/response"
	"github.com/zero-life/server/internal/model"
	"github.com/zero-life/server/internal/pkg/errcode"
	"github.com/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
	db           *gorm.DB
}

func NewCategoryService(categoryRepo *repository.CategoryRepository, db *gorm.DB) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		db:           db,
	}
}

func (s *CategoryService) Create(userID uint64, req *request.CreateCategoryReq) (*response.CategoryResp, error) {
	// Check name uniqueness
	categories, err := s.categoryRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	for _, c := range categories {
		if c.Name == req.Name {
			return nil, errcode.ErrCategoryNameExists
		}
	}

	// Check two-level limit
	if req.ParentID != nil {
		parent, err := s.categoryRepo.GetByID(*req.ParentID, userID)
		if err != nil {
			return nil, errcode.ErrNotFound
		}
		if parent.ParentID != nil {
			return nil, errcode.ErrCategoryTooDeep
		}
	}

	category := &model.Category{
		UserID:   userID,
		Name:     req.Name,
		ParentID: req.ParentID,
		Icon:     req.Icon,
		Notes:    req.Notes,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(category), nil
}

func (s *CategoryService) List(userID uint64) ([]response.CategoryResp, error) {
	categories, err := s.categoryRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.buildTree(categories), nil
}

func (s *CategoryService) Update(userID, id uint64, req *request.UpdateCategoryReq) (*response.CategoryResp, error) {
	category, err := s.categoryRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}
	if req.Notes != "" {
		category.Notes = req.Notes
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(category), nil
}

func (s *CategoryService) Delete(userID, id uint64) error {
	_, err := s.categoryRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}

	return s.db.Transaction(func(dbTx *gorm.DB) error {
		// Get sub-categories
		subCategories, err := s.categoryRepo.GetSubCategories(id, userID)
		if err != nil {
			return err
		}

		// Delete sub-categories
		for _, sub := range subCategories {
			if err := s.categoryRepo.Delete(sub.ID, userID); err != nil {
				return err
			}
		}

		// Nullify category_id on associated transactions
		if err := dbTx.Model(&model.Transaction{}).Where("category_id = ?", id).Update("category_id", nil).Error; err != nil {
			return err
		}

		// Delete the category
		return s.categoryRepo.Delete(id, userID)
	})
}

func (s *CategoryService) toResp(c *model.Category) *response.CategoryResp {
	return &response.CategoryResp{
		ID:        c.ID,
		Name:      c.Name,
		ParentID:  c.ParentID,
		Icon:      c.Icon,
		Notes:     c.Notes,
		SortOrder: c.SortOrder,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func (s *CategoryService) buildTree(categories []model.Category) []response.CategoryResp {
	// Build a map for quick lookup
	nodeMap := make(map[uint64]*response.CategoryResp)
	for _, c := range categories {
		nodeMap[c.ID] = s.toResp(&c)
	}

	var roots []response.CategoryResp
	for _, c := range categories {
		node := nodeMap[c.ID]
		if c.ParentID == nil {
			roots = append(roots, *node)
		} else {
			if parent, ok := nodeMap[*c.ParentID]; ok {
				parent.Children = append(parent.Children, *node)
			}
		}
	}

	return roots
}

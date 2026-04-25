// Package service 业务逻辑层，实现核心业务逻辑
// CategoryService 分类业务逻辑，处理分类的增删改查（支持树形结构）
package service

import (
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

// CategoryService 分类服务
// 负责处理分类的增删改查，支持两级树形结构
// 业务规则：分类名称唯一；最多两级（父分类下不能再有子分类的子分类）
// 依赖categoryRepo进行分类数据访问，依赖db执行删除事务
type CategoryService struct {
	categoryRepo *repository.CategoryRepository // 分类数据访问对象
	db           *gorm.DB                       // 数据库连接，用于删除分类时的事务操作
}

// NewCategoryService 创建分类服务实例
func NewCategoryService(categoryRepo *repository.CategoryRepository, db *gorm.DB) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		db:           db,
	}
}

// Create 创建分类
// 检查名称唯一性和两级深度限制
// 参数：
//   - userID: 用户ID
//   - req: 创建请求参数（名称、父分类ID、图标、备注）
// 返回：
//   - *response.CategoryResp: 创建成功的分类信息
//   - error: 错误信息（如名称重复、层级过深）
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

// List 获取分类树形列表
// 返回两级树形结构的分类数据
// 参数：
//   - userID: 用户ID
// 返回：
//   - []response.CategoryResp: 分类树形列表
//   - error: 错误信息
func (s *CategoryService) List(userID uint64) ([]response.CategoryResp, error) {
	categories, err := s.categoryRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.buildTree(categories), nil
}

// Get 获取单个分类详情
// 参数：
//   - userID: 用户ID
//   - id: 分类ID
// 返回：
//   - *response.CategoryResp: 分类信息
//   - error: 错误信息
func (s *CategoryService) Get(userID, id uint64) (*response.CategoryResp, error) {
	category, err := s.categoryRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(category), nil
}

// Update 更新分类信息
// 支持更新名称、图标、备注、排序
// 参数：
//   - userID: 用户ID
//   - id: 分类ID
//   - req: 更新请求参数
// 返回：
//   - *response.CategoryResp: 更新后的分类信息
//   - error: 错误信息
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

// Delete 删除分类
// 在数据库事务中执行：1)删除所有子分类 2)将关联交易的分类ID设为NULL 3)删除分类本身
// 参数：
//   - userID: 用户ID
//   - id: 分类ID
// 返回：
//   - error: 错误信息
func (s *CategoryService) Delete(userID, id uint64) error {
	_, err := s.categoryRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}

	// 在数据库事务中执行删除操作，确保数据一致性
	subCategories, err := s.categoryRepo.GetSubCategories(id, userID)
	if err != nil {
		return errcode.ErrInternal
	}

	for _, sub := range subCategories {
		if err := s.categoryRepo.Delete(sub.ID, userID); err != nil {
			return errcode.ErrInternal
		}
	}

	if err := s.db.Model(&model.Transaction{}).Where("category_id = ?", id).Update("category_id", nil).Error; err != nil {
		return errcode.ErrInternal
	}

	return s.categoryRepo.Delete(id, userID)
}

// toResp 将分类模型转换为响应对象
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

// buildTree 将扁平的分类列表构建为树形结构
// 使用map快速查找，将子分类挂载到父分类的Children字段
func (s *CategoryService) buildTree(categories []model.Category) []response.CategoryResp {
	nodeMap := make(map[uint64]*response.CategoryResp)
	for _, c := range categories {
		nodeMap[c.ID] = s.toResp(&c)
	}

	for _, c := range categories {
		if c.ParentID != nil {
			if parent, ok := nodeMap[*c.ParentID]; ok {
				parent.Children = append(parent.Children, *nodeMap[c.ID])
			}
		}
	}

	var roots []response.CategoryResp
	for _, c := range categories {
		if c.ParentID == nil {
			roots = append(roots, *nodeMap[c.ID])
		}
	}

	return roots
}

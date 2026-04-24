package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// CategoryRepository 分类仓库，负责交易分类的数据访问。
// 分类（Category）用于对交易进行归类，如"餐饮"、"交通"、"工资"等。
// 分类支持两级层次结构：父分类和子分类（通过 parent_id 关联）。
// 例如："日常开销"（父）→ "餐饮"（子）、"交通"（子）。
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓库实例。
// 参数 db: GORM 数据库连接实例。
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create 创建一条新的分类记录。
// 执行 SQL: INSERT INTO categories (...)
// 参数 category: 要创建的分类对象。
// 返回: 创建失败时返回错误。
func (r *CategoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

// GetByID 根据 ID 和用户 ID 获取单条分类。
// 执行 SQL: SELECT * FROM categories WHERE id = ? AND user_id = ? LIMIT 1
// 参数 id: 分类 ID。
// 参数 userID: 当前登录用户 ID，用于权限校验。
// 返回: 找到的分类对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *CategoryRepository) GetByID(id, userID uint64) (*model.Category, error) {
	var category model.Category
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// List 获取指定用户的所有分类，按排序字段和 ID 升序排列。
// sort_order 字段允许用户自定义分类的显示顺序。
// 执行 SQL: SELECT * FROM categories WHERE user_id = ? ORDER BY sort_order ASC, id ASC
// 参数 userID: 用户 ID。
// 返回: 分类列表。
func (r *CategoryRepository) List(userID uint64) ([]model.Category, error) {
	var categories []model.Category
	if err := r.db.Where("user_id = ?", userID).Order("sort_order ASC, id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// Update 更新分类。GORM 的 Save 方法会更新所有字段。
// 执行 SQL: UPDATE categories SET ... WHERE id = ?
// 参数 category: 要更新的分类对象（必须包含 ID 字段）。
// 返回: 更新失败时返回错误。
func (r *CategoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

// Delete 根据 ID 和用户 ID 删除分类。
// 注意：此方法不处理子分类和关联交易，调用方需自行确保数据完整性。
// 执行 SQL: DELETE FROM categories WHERE id = ? AND user_id = ?
// 参数 id: 分类 ID。
// 参数 userID: 当前登录用户 ID。
// 返回: 删除失败时返回错误。
func (r *CategoryRepository) Delete(id, userID uint64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Category{}).Error
}

// GetSubCategories 获取指定父分类下的所有子分类。
// 执行 SQL: SELECT * FROM categories WHERE parent_id = ? AND user_id = ?
// 参数 parentID: 父分类 ID。
// 参数 userID: 当前登录用户 ID。
// 返回: 子分类列表。
func (r *CategoryRepository) GetSubCategories(parentID, userID uint64) ([]model.Category, error) {
	var categories []model.Category
	if err := r.db.Where("parent_id = ? AND user_id = ?", parentID, userID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

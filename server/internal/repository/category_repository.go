package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// CategoryRepository 分类仓库，负责交易分类的数据访问。
// 分类（Category）用于对交易进行归类，如"餐饮"、"交通"、"工资"等。
// 分类支持最多5级层次结构：通过 parent_id 关联形成树形结构。
// 例如："日常开销"（1级）→ "餐饮"（2级）→ "外卖"（3级）→ ...
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

// GetDescendantIDs 获取指定分类的所有子孙分类ID（包含自身）
// 递归查询：先查直接子分类，再对每个子分类递归查询，最终返回所有层级的后代ID
// 参数 ids: 起始分类ID列表
// 参数 userID: 用户ID
// 返回: 包含自身及所有子孙分类的ID列表
func (r *CategoryRepository) GetDescendantIDs(ids []uint64, userID uint64) ([]uint64, error) {
	if len(ids) == 0 {
		return ids, nil
	}
	result := make([]uint64, 0, len(ids)*2)
	result = append(result, ids...)

	// 逐层向下查找子分类
	currentLevel := ids
	for len(currentLevel) > 0 {
		var children []model.Category
		if err := r.db.Where("parent_id IN ? AND user_id = ?", currentLevel, userID).Find(&children).Error; err != nil {
			return nil, err
		}
		if len(children) == 0 {
			break
		}
		currentLevel = make([]uint64, 0, len(children))
		for _, c := range children {
			currentLevel = append(currentLevel, c.ID)
			result = append(result, c.ID)
		}
	}

	return result, nil
}

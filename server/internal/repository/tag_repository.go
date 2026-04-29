package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// TagRepository 标签仓库，负责标签的数据访问。
// 标签（Tag）用于对交易进行分类标记，与交易是多对多关系（通过 transaction_tags 关联表）。
// 例如：#餐饮、#交通、#日用品 等标签可以灵活地标记交易。
type TagRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

// NewTagRepository 创建标签仓库实例。
// 参数 readDB: 读库（只读连接池，并发安全）
// 参数 writeDB: 写库（单连接，串行保证安全）
func NewTagRepository(readDB, writeDB *gorm.DB) *TagRepository {
	return &TagRepository{readDB: readDB, writeDB: writeDB}
}

// Create 创建一条新的标签记录。
// 执行 SQL: INSERT INTO tags (...)
// 参数 tag: 要创建的标签对象。
// 返回: 创建失败时返回错误。
func (r *TagRepository) Create(tag *model.Tag) error {
	return r.writeDB.Create(tag).Error
}

// GetByID 根据 ID 和用户 ID 获取单条标签。
// 执行 SQL: SELECT * FROM tags WHERE id = ? AND user_id = ? LIMIT 1
// 参数 id: 标签 ID。
// 参数 userID: 当前登录用户 ID，用于权限校验。
// 返回: 找到的标签对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *TagRepository) GetByID(id, userID uint64) (*model.Tag, error) {
	var tag model.Tag
	if err := r.readDB.Where("id = ? AND user_id = ?", id, userID).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// List 获取指定用户的所有标签，按 ID 升序排列。
// 执行 SQL: SELECT * FROM tags WHERE user_id = ? ORDER BY id ASC
// 参数 userID: 用户 ID。
// 返回: 标签列表。
func (r *TagRepository) List(userID uint64) ([]model.Tag, error) {
	var tags []model.Tag
	if err := r.readDB.Where("user_id = ?", userID).Order("id ASC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// Update 更新标签。GORM 的 Save 方法会更新所有字段。
// 执行 SQL: UPDATE tags SET ... WHERE id = ?
// 参数 tag: 要更新的标签对象（必须包含 ID 字段）。
// 返回: 更新失败时返回错误。
func (r *TagRepository) Update(tag *model.Tag) error {
	return r.writeDB.Save(tag).Error
}

// Delete 删除标签及其关联的所有交易-标签关系。使用数据库事务确保原子性。
// 删除顺序：1.删除 transaction_tags 中的关联记录 → 2.删除标签本身
// 执行 SQL（事务内）:
//   1. DELETE FROM transaction_tags WHERE tag_id = ?
//   2. DELETE FROM tags WHERE id = ? AND user_id = ?
// 参数 id: 标签 ID。
// 参数 userID: 当前登录用户 ID，确保只能删除自己的标签。
// 返回: 删除失败时返回错误，事务回滚。
func (r *TagRepository) Delete(id, userID uint64) error {
	return r.writeDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("tag_id = ?", id).Delete(&model.TransactionTag{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Tag{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// CountTransactions 统计使用指定标签的交易数量，用于删除标签前的提示。
// 执行 SQL: SELECT COUNT(*) FROM transaction_tags WHERE tag_id = ?
// 参数 tagID: 标签 ID。
// 返回: 使用该标签的交易数量。
func (r *TagRepository) CountTransactions(tagID uint64) (int64, error) {
	var count int64
	if err := r.readDB.Model(&model.TransactionTag{}).Where("tag_id = ?", tagID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetSubTags 获取指定父标签下的所有子标签。
// 执行 SQL: SELECT * FROM tags WHERE parent_id = ? AND user_id = ?
// 参数 parentID: 父标签 ID。
// 参数 userID: 当前登录用户 ID。
// 返回: 子标签列表。
func (r *TagRepository) GetSubTags(parentID, userID uint64) ([]model.Tag, error) {
	var tags []model.Tag
	if err := r.readDB.Where("parent_id = ? AND user_id = ?", parentID, userID).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// GetDescendantIDs 获取指定标签的所有子孙标签ID（包含自身）
// 递归查询：先查直接子标签，再对每个子标签递归查询，最终返回所有层级的后代ID
// 参数 ids: 起始标签ID列表
// 参数 userID: 用户ID
// 返回: 包含自身及所有子孙标签的ID列表
func (r *TagRepository) GetDescendantIDs(ids []uint64, userID uint64) ([]uint64, error) {
	if len(ids) == 0 {
		return ids, nil
	}
	result := make([]uint64, 0, len(ids)*2)
	result = append(result, ids...)

	currentLevel := ids
	for len(currentLevel) > 0 {
		var children []model.Tag
		if err := r.readDB.Where("parent_id IN ? AND user_id = ?", currentLevel, userID).Find(&children).Error; err != nil {
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

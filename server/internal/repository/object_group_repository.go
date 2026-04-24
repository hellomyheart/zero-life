package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// ObjectGroupRepository 对象组仓库，负责对象分组的数据访问。
// 对象组（ObjectGroup）用于将不同类型的对象（如账户、分类等）进行分组管理。
// 通过 groupable_type 和 groupable_id 实现多态关联，即一个对象组可以关联不同类型的对象。
// 例如：将多个资产账户归入"投资"组，将多个支出分类归入"日常开销"组。
type ObjectGroupRepository struct {
	db *gorm.DB
}

// NewObjectGroupRepository 创建对象组仓库实例。
// 参数 db: GORM 数据库连接实例。
func NewObjectGroupRepository(db *gorm.DB) *ObjectGroupRepository {
	return &ObjectGroupRepository{db: db}
}

// Create 创建一条新的对象组记录。
// 执行 SQL: INSERT INTO object_groups (...)
// 参数 og: 要创建的对象组对象。
// 返回: 创建失败时返回错误。
func (r *ObjectGroupRepository) Create(og *model.ObjectGroup) error {
	return r.db.Create(og).Error
}

// GetByID 根据 ID 和用户 ID 获取单条对象组。
// 执行 SQL: SELECT * FROM object_groups WHERE id = ? AND user_id = ? LIMIT 1
// 参数 id: 对象组 ID。
// 参数 userID: 当前登录用户 ID，用于权限校验。
// 返回: 找到的对象组对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *ObjectGroupRepository) GetByID(id, userID uint64) (*model.ObjectGroup, error) {
	var og model.ObjectGroup
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&og).Error; err != nil {
		return nil, err
	}
	return &og, nil
}

// List 获取指定用户的对象组列表，可选按对象类型过滤，按名称升序排列。
// 执行 SQL: SELECT * FROM object_groups WHERE user_id = ? [AND groupable_type = ?] ORDER BY name ASC
// 参数 userID: 用户 ID。
// 参数 groupableType: 可选的对象类型过滤，如 "Account"、"Category" 等，为空时不过滤。
// 返回: 对象组列表。
func (r *ObjectGroupRepository) List(userID uint64, groupableType string) ([]model.ObjectGroup, error) {
	var groups []model.ObjectGroup
	query := r.db.Where("user_id = ?", userID)
	if groupableType != "" {
		query = query.Where("groupable_type = ?", groupableType)
	}
	if err := query.Order("name ASC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// Update 更新对象组。GORM 的 Save 方法会更新所有字段。
// 执行 SQL: UPDATE object_groups SET ... WHERE id = ?
// 参数 og: 要更新的对象组对象（必须包含 ID 字段）。
// 返回: 更新失败时返回错误。
func (r *ObjectGroupRepository) Update(og *model.ObjectGroup) error {
	return r.db.Save(og).Error
}

// Delete 根据 ID 和用户 ID 删除对象组。
// 执行 SQL: DELETE FROM object_groups WHERE id = ? AND user_id = ?
// 参数 id: 对象组 ID。
// 参数 userID: 当前登录用户 ID。
// 返回: 删除失败时返回错误。
func (r *ObjectGroupRepository) Delete(id, userID uint64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.ObjectGroup{}).Error
}

// GetByGroupable 根据用户 ID、对象类型和对象 ID 查找关联的对象组。
// 用于查找某个具体对象（如某个账户）所属的分组。
// 执行 SQL: SELECT * FROM object_groups WHERE user_id = ? AND groupable_type = ? AND groupable_id = ? LIMIT 1
// 参数 userID: 用户 ID。
// 参数 groupableType: 对象类型，如 "Account"。
// 参数 groupableID: 对象 ID。
// 返回: 找到的对象组对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *ObjectGroupRepository) GetByGroupable(userID uint64, groupableType string, groupableID uint64) (*model.ObjectGroup, error) {
	var og model.ObjectGroup
	if err := r.db.Where("user_id = ? AND groupable_type = ? AND groupable_id = ?", userID, groupableType, groupableID).
		First(&og).Error; err != nil {
		return nil, err
	}
	return &og, nil
}

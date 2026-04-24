package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// LinkTypeRepository 链接类型仓库，负责交易链接类型的数据访问。
// 链接类型（LinkType）定义了交易之间关联的类型，如"转账"、"退款"等。
// 注意：链接类型是全局资源，不属于某个用户，因此方法中没有 userID 参数。
type LinkTypeRepository struct {
	db *gorm.DB
}

// NewLinkTypeRepository 创建链接类型仓库实例。
// 参数 db: GORM 数据库连接实例。
func NewLinkTypeRepository(db *gorm.DB) *LinkTypeRepository {
	return &LinkTypeRepository{db: db}
}

// Create 创建一条新的链接类型记录。
// 执行 SQL: INSERT INTO link_types (...)
// 参数 linkType: 要创建的链接类型对象。
// 返回: 创建失败时返回错误。
func (r *LinkTypeRepository) Create(linkType *model.LinkType) error {
	return r.db.Create(linkType).Error
}

// GetByID 根据 ID 获取单条链接类型。注意：无用户权限校验，因为链接类型是全局的。
// 执行 SQL: SELECT * FROM link_types WHERE id = ? LIMIT 1
// 参数 id: 链接类型 ID。
// 返回: 找到的链接类型对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *LinkTypeRepository) GetByID(id uint64) (*model.LinkType, error) {
	var linkType model.LinkType
	if err := r.db.Where("id = ?", id).First(&linkType).Error; err != nil {
		return nil, err
	}
	return &linkType, nil
}

// List 分页获取所有链接类型，按 ID 升序排列。
// 执行 SQL: SELECT * FROM link_types ORDER BY id ASC LIMIT ? OFFSET ?
// 参数 offset: 分页偏移量。
// 参数 limit: 每页记录数。
// 返回: 链接类型列表。
func (r *LinkTypeRepository) List(offset, limit int) ([]model.LinkType, error) {
	var linkTypes []model.LinkType
	if err := r.db.Order("id ASC").Offset(offset).Limit(limit).Find(&linkTypes).Error; err != nil {
		return nil, err
	}
	return linkTypes, nil
}

// Count 统计链接类型总数，用于分页计算。
// 执行 SQL: SELECT COUNT(*) FROM link_types
// 返回: 链接类型总数。
func (r *LinkTypeRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&model.LinkType{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Update 更新链接类型。GORM 的 Save 方法会更新所有字段。
// 执行 SQL: UPDATE link_types SET ... WHERE id = ?
// 参数 linkType: 要更新的链接类型对象（必须包含 ID 字段）。
// 返回: 更新失败时返回错误。
func (r *LinkTypeRepository) Update(linkType *model.LinkType) error {
	return r.db.Save(linkType).Error
}

// Delete 根据 ID 删除链接类型。
// 执行 SQL: DELETE FROM link_types WHERE id = ?
// 参数 id: 链接类型 ID。
// 返回: 删除失败时返回错误。
func (r *LinkTypeRepository) Delete(id uint64) error {
	return r.db.Where("id = ?", id).Delete(&model.LinkType{}).Error
}

package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// RuleGroupRepository 规则组仓库，负责规则分组的数据访问。
// 规则组（RuleGroup）用于将多条规则按逻辑分组管理，例如"自动分类规则"、"自动标签规则"等。
// 规则组有排序字段（order），决定组的执行顺序。
type RuleGroupRepository struct {
	db *gorm.DB
}

// NewRuleGroupRepository 创建规则组仓库实例。
// 参数 db: GORM 数据库连接实例。
func NewRuleGroupRepository(db *gorm.DB) *RuleGroupRepository {
	return &RuleGroupRepository{db: db}
}

// Create 创建一条新的规则组记录。
// 执行 SQL: INSERT INTO rule_groups (...)
// 参数 group: 要创建的规则组对象。
// 返回: 创建失败时返回错误。
func (r *RuleGroupRepository) Create(group *model.RuleGroup) error {
	return r.db.Create(group).Error
}

// GetByID 根据 ID 和用户 ID 获取单条规则组。
// 执行 SQL: SELECT * FROM rule_groups WHERE id = ? AND user_id = ? LIMIT 1
// 参数 id: 规则组 ID。
// 参数 userID: 当前登录用户 ID，用于权限校验。
// 返回: 找到的规则组对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *RuleGroupRepository) GetByID(id, userID uint64) (*model.RuleGroup, error) {
	var group model.RuleGroup
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// List 获取指定用户的所有规则组，按排序字段和 ID 升序排列。
// 注意：`order` 是 MySQL 保留字，需要用反引号转义。
// 执行 SQL: SELECT * FROM rule_groups WHERE user_id = ? ORDER BY `order` ASC, id ASC
// 参数 userID: 用户 ID。
// 返回: 规则组列表。
func (r *RuleGroupRepository) List(userID uint64) ([]model.RuleGroup, error) {
	var groups []model.RuleGroup
	if err := r.db.Where("user_id = ?", userID).Order("`order` ASC, id ASC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// Update 更新规则组。GORM 的 Save 方法会更新所有字段。
// 执行 SQL: UPDATE rule_groups SET ... WHERE id = ?
// 参数 group: 要更新的规则组对象（必须包含 ID 字段）。
// 返回: 更新失败时返回错误。
func (r *RuleGroupRepository) Update(group *model.RuleGroup) error {
	return r.db.Save(group).Error
}

// Delete 根据 ID 和用户 ID 删除规则组。
// 注意：此方法不删除组内的规则，调用方需自行处理或使用 HasRules 检查。
// 执行 SQL: DELETE FROM rule_groups WHERE id = ? AND user_id = ?
// 参数 id: 规则组 ID。
// 参数 userID: 当前登录用户 ID。
// 返回: 删除失败时返回错误。
func (r *RuleGroupRepository) Delete(id, userID uint64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.RuleGroup{}).Error
}

// HasRules 检查规则组下是否还有规则，用于删除前的校验。
// 执行 SQL: SELECT COUNT(*) FROM rules WHERE group_id = ? AND user_id = ?
// 参数 groupID: 规则组 ID。
// 参数 userID: 用户 ID。
// 返回: 如果有规则返回 true，否则返回 false。
func (r *RuleGroupRepository) HasRules(groupID, userID uint64) (bool, error) {
	var count int64
	if err := r.db.Model(&model.Rule{}).Where("group_id = ? AND user_id = ?", groupID, userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetRuleCount 获取规则组下的规则数量，用于展示。
// 执行 SQL: SELECT COUNT(*) FROM rules WHERE group_id = ? AND user_id = ?
// 参数 groupID: 规则组 ID。
// 参数 userID: 用户 ID。
// 返回: 规则数量。
func (r *RuleGroupRepository) GetRuleCount(groupID, userID uint64) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Rule{}).Where("group_id = ? AND user_id = ?", groupID, userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

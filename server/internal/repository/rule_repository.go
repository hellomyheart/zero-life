package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// RuleRepository 规则仓库，负责规则及其条件和操作的数据访问。
// 规则（Rule）用于自动处理交易，当交易满足条件（RuleCondition）时执行操作（RuleAction）。
// 例如：当交易描述包含"星巴克"时，自动将其分类到"餐饮"类别。
// 规则属于规则组（RuleGroup），按优先级（priority）排序执行。
type RuleRepository struct {
	db *gorm.DB
}

// NewRuleRepository 创建规则仓库实例。
// 参数 db: GORM 数据库连接实例。
func NewRuleRepository(db *gorm.DB) *RuleRepository {
	return &RuleRepository{db: db}
}

// Create 创建规则及其关联的条件和操作。使用数据库事务确保原子性。
// 先创建规则记录获取 ID，再将规则 ID 赋给条件和操作，最后批量创建。
// 执行 SQL（事务内）:
//   1. INSERT INTO rules (...)
//   2. INSERT INTO rule_conditions (rule_id, ...) VALUES (?, ...), ...
//   3. INSERT INTO rule_actions (rule_id, ...) VALUES (?, ...), ...
// 参数 rule: 要创建的规则对象，创建后 GORM 会自动填充 ID。
// 参数 conditions: 规则的条件列表，创建时会自动设置 RuleID。
// 参数 actions: 规则的操作列表，创建时会自动设置 RuleID。
// 返回: 创建失败时返回错误，事务回滚。
func (r *RuleRepository) Create(rule *model.Rule, conditions []model.RuleCondition, actions []model.RuleAction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(rule).Error; err != nil {
			return err
		}
		for i := range conditions {
			conditions[i].RuleID = rule.ID
		}
		if len(conditions) > 0 {
			if err := tx.Create(&conditions).Error; err != nil {
				return err
			}
		}
		for i := range actions {
			actions[i].RuleID = rule.ID
		}
		if len(actions) > 0 {
			if err := tx.Create(&actions).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetByID 根据 ID 和用户 ID 获取单条规则，并预加载其条件和操作。
// 执行 SQL:
//   主查询: SELECT * FROM rules WHERE id = ? AND user_id = ? LIMIT 1
//   预加载: SELECT * FROM rule_conditions WHERE rule_id = ?
//           SELECT * FROM rule_actions WHERE rule_id = ?
// 参数 id: 规则 ID。
// 参数 userID: 当前登录用户 ID，用于权限校验。
// 返回: 包含完整条件和操作的规则对象。
func (r *RuleRepository) GetByID(id, userID uint64) (*model.Rule, error) {
	var rule model.Rule
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).
		Preload("Conditions").
		Preload("Actions").
		First(&rule).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

// List 获取指定用户的所有规则，按优先级和 ID 升序排列，并预加载条件和操作。
// 执行 SQL: SELECT * FROM rules WHERE user_id = ? ORDER BY priority ASC, id ASC
// 返回: 包含完整条件和操作的规则列表。
func (r *RuleRepository) List(userID uint64) ([]model.Rule, error) {
	var rules []model.Rule
	if err := r.db.Where("user_id = ?", userID).
		Preload("Conditions").
		Preload("Actions").
		Order("priority ASC, id ASC").Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

// Update 更新规则及其条件和操作。使用数据库事务确保原子性。
// 策略：先更新规则本身，再全量替换条件和操作（删除旧的，插入新的）。
// 条件和操作的 ID 会被重置为 0，GORM 会创建新记录而非更新旧记录。
// 执行 SQL（事务内）:
//   1. UPDATE rules SET ... WHERE id = ?
//   2. DELETE FROM rule_conditions WHERE rule_id = ?
//   3. INSERT INTO rule_conditions (rule_id, ...) VALUES (?, ...), ...
//   4. DELETE FROM rule_actions WHERE rule_id = ?
//   5. INSERT INTO rule_actions (rule_id, ...) VALUES (?, ...), ...
// 参数 rule: 要更新的规则对象。
// 参数 conditions: 新的条件列表（替换原有条件）。
// 参数 actions: 新的操作列表（替换原有操作）。
// 返回: 更新失败时返回错误，事务回滚。
func (r *RuleRepository) Update(rule *model.Rule, conditions []model.RuleCondition, actions []model.RuleAction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(rule).Error; err != nil {
			return err
		}
		if err := tx.Where("rule_id = ?", rule.ID).Delete(&model.RuleCondition{}).Error; err != nil {
			return err
		}
		for i := range conditions {
			conditions[i].RuleID = rule.ID
			conditions[i].ID = 0
		}
		if len(conditions) > 0 {
			if err := tx.Create(&conditions).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("rule_id = ?", rule.ID).Delete(&model.RuleAction{}).Error; err != nil {
			return err
		}
		for i := range actions {
			actions[i].RuleID = rule.ID
			actions[i].ID = 0
		}
		if len(actions) > 0 {
			if err := tx.Create(&actions).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Delete 删除规则及其关联的条件和操作。使用数据库事务确保原子性。
// 删除顺序：1.条件 → 2.操作 → 3.规则本身
// 执行 SQL（事务内）:
//   1. DELETE FROM rule_conditions WHERE rule_id = ?
//   2. DELETE FROM rule_actions WHERE rule_id = ?
//   3. DELETE FROM rules WHERE id = ? AND user_id = ?
// 参数 id: 规则 ID。
// 参数 userID: 当前登录用户 ID。
// 返回: 删除失败时返回错误，事务回滚。
func (r *RuleRepository) Delete(id, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("rule_id = ?", id).Delete(&model.RuleCondition{}).Error; err != nil {
			return err
		}
		if err := tx.Where("rule_id = ?", id).Delete(&model.RuleAction{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Rule{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetEnabledRules 获取指定用户所有已启用的规则，用于规则引擎执行。
// 按优先级升序排列，优先级数字越小越先执行。
// 执行 SQL: SELECT * FROM rules WHERE user_id = ? AND is_enabled = true ORDER BY priority ASC, id ASC
// 参数 userID: 用户 ID。
// 返回: 已启用的规则列表（含预加载的条件和操作）。
func (r *RuleRepository) GetEnabledRules(userID uint64) ([]model.Rule, error) {
	var rules []model.Rule
	if err := r.db.Where("user_id = ? AND is_enabled = ?", userID, true).
		Preload("Conditions").
		Preload("Actions").
		Order("priority ASC, id ASC").Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

// GetEnabledRulesByGroupID 获取指定规则组中已启用的规则，用于分组执行规则。
// 执行 SQL: SELECT * FROM rules WHERE group_id = ? AND user_id = ? AND is_enabled = true ORDER BY priority ASC, id ASC
// 参数 groupID: 规则组 ID。
// 参数 userID: 用户 ID。
// 返回: 该规则组中已启用的规则列表（含预加载的条件和操作）。
func (r *RuleRepository) GetEnabledRulesByGroupID(groupID, userID uint64) ([]model.Rule, error) {
	var rules []model.Rule
	if err := r.db.Where("group_id = ? AND user_id = ? AND is_enabled = ?", groupID, userID, true).
		Preload("Conditions").
		Preload("Actions").
		Order("priority ASC, id ASC").Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

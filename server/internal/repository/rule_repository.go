package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type RuleRepository struct {
	db *gorm.DB
}

func NewRuleRepository(db *gorm.DB) *RuleRepository {
	return &RuleRepository{db: db}
}

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

func (r *RuleRepository) Update(rule *model.Rule, conditions []model.RuleCondition, actions []model.RuleAction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(rule).Error; err != nil {
			return err
		}
		// Replace conditions
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
		// Replace actions
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

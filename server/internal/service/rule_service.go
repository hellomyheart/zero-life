package service

import (
	"strings"

	"github.com/zero-life/server/internal/dto/request"
	"github.com/zero-life/server/internal/dto/response"
	"github.com/zero-life/server/internal/model"
	"github.com/zero-life/server/internal/pkg/errcode"
	"github.com/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

type RuleService struct {
	ruleRepo *repository.RuleRepository
}

func NewRuleService(ruleRepo *repository.RuleRepository) *RuleService {
	return &RuleService{ruleRepo: ruleRepo}
}

func (s *RuleService) Create(userID uint64, req *request.CreateRuleReq) (*response.RuleResp, error) {
	rule := &model.Rule{
		UserID:    userID,
		Name:      req.Name,
		Priority:  req.Priority,
		IsEnabled: req.IsEnabled,
		LogicType: model.LogicType(req.LogicType),
	}

	conditions := make([]model.RuleCondition, 0, len(req.Conditions))
	for _, c := range req.Conditions {
		conditions = append(conditions, model.RuleCondition{
			Field:    model.ConditionField(c.Field),
			Operator: model.ConditionOperator(c.Operator),
			Value:    c.Value,
		})
	}

	actions := make([]model.RuleAction, 0, len(req.Actions))
	for _, a := range req.Actions {
		actions = append(actions, model.RuleAction{
			Type:  model.ActionType(a.Type),
			Value: a.Value,
		})
	}

	if err := s.ruleRepo.Create(rule, conditions, actions); err != nil {
		return nil, errcode.ErrInternal
	}

	created, err := s.ruleRepo.GetByID(rule.ID, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(created), nil
}

func (s *RuleService) Get(userID, id uint64) (*response.RuleResp, error) {
	rule, err := s.ruleRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(rule), nil
}

func (s *RuleService) List(userID uint64) ([]response.RuleResp, error) {
	rules, err := s.ruleRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.RuleResp, 0, len(rules))
	for _, r := range rules {
		items = append(items, *s.toResp(&r))
	}

	return items, nil
}

func (s *RuleService) Update(userID, id uint64, req *request.UpdateRuleReq) (*response.RuleResp, error) {
	rule, err := s.ruleRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Name != "" {
		rule.Name = req.Name
	}
	if req.Priority != nil {
		rule.Priority = *req.Priority
	}
	if req.IsEnabled != nil {
		rule.IsEnabled = *req.IsEnabled
	}
	if req.LogicType != "" {
		rule.LogicType = model.LogicType(req.LogicType)
	}

	conditions := make([]model.RuleCondition, 0, len(req.Conditions))
	for _, c := range req.Conditions {
		conditions = append(conditions, model.RuleCondition{
			Field:    model.ConditionField(c.Field),
			Operator: model.ConditionOperator(c.Operator),
			Value:    c.Value,
		})
	}

	actions := make([]model.RuleAction, 0, len(req.Actions))
	for _, a := range req.Actions {
		actions = append(actions, model.RuleAction{
			Type:  model.ActionType(a.Type),
			Value: a.Value,
		})
	}

	if err := s.ruleRepo.Update(rule, conditions, actions); err != nil {
		return nil, errcode.ErrInternal
	}

	updated, err := s.ruleRepo.GetByID(id, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(updated), nil
}

func (s *RuleService) Delete(userID, id uint64) error {
	_, err := s.ruleRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}

	return s.ruleRepo.Delete(id, userID)
}

func (s *RuleService) ToggleStatus(userID, id uint64) (*response.RuleResp, error) {
	rule, err := s.ruleRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	rule.IsEnabled = !rule.IsEnabled

	conditions := make([]model.RuleCondition, 0, len(rule.Conditions))
	for _, c := range rule.Conditions {
		conditions = append(conditions, model.RuleCondition{
			Field:    c.Field,
			Operator: c.Operator,
			Value:    c.Value,
		})
	}

	actions := make([]model.RuleAction, 0, len(rule.Actions))
	for _, a := range rule.Actions {
		actions = append(actions, model.RuleAction{
			Type:  a.Type,
			Value: a.Value,
		})
	}

	if err := s.ruleRepo.Update(rule, conditions, actions); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(rule), nil
}

func (s *RuleService) Execute(userID, id uint64, req *request.ExecuteRuleReq) (*response.RuleResp, error) {
	rule, err := s.ruleRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	// Execute the rule against matching transactions
	// This is a simplified implementation
	_ = rule
	_ = req

	return s.toResp(rule), nil
}

// matchCondition checks if a transaction matches a single condition
func (s *RuleService) matchCondition(condition model.RuleCondition, description string, amount string, sourceAccount string) bool {
	var fieldValue string
	switch condition.Field {
	case model.ConditionFieldDescription:
		fieldValue = description
	case model.ConditionFieldAmount:
		fieldValue = amount
	case model.ConditionFieldSourceAccount:
		fieldValue = sourceAccount
	}

	switch condition.Operator {
	case model.OperatorContains:
		return strings.Contains(strings.ToLower(fieldValue), strings.ToLower(condition.Value))
	case model.OperatorEquals:
		return strings.EqualFold(fieldValue, condition.Value)
	case model.OperatorStartsWith:
		return strings.HasPrefix(strings.ToLower(fieldValue), strings.ToLower(condition.Value))
	case model.OperatorEndsWith:
		return strings.HasSuffix(strings.ToLower(fieldValue), strings.ToLower(condition.Value))
	}

	return false
}

func (s *RuleService) toResp(r *model.Rule) *response.RuleResp {
	conditions := make([]response.RuleConditionResp, 0, len(r.Conditions))
	for _, c := range r.Conditions {
		conditions = append(conditions, response.RuleConditionResp{
			ID:       c.ID,
			Field:    string(c.Field),
			Operator: string(c.Operator),
			Value:    c.Value,
		})
	}

	actions := make([]response.RuleActionResp, 0, len(r.Actions))
	for _, a := range r.Actions {
		actions = append(actions, response.RuleActionResp{
			ID:    a.ID,
			Type:  string(a.Type),
			Value: a.Value,
		})
	}

	return &response.RuleResp{
		ID:         r.ID,
		Name:       r.Name,
		Priority:   r.Priority,
		IsEnabled:  r.IsEnabled,
		LogicType:  string(r.LogicType),
		Conditions: conditions,
		Actions:    actions,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

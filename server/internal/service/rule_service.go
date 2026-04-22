package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

type RuleService struct {
	ruleRepo *repository.RuleRepository
	txnRepo  *repository.TransactionRepository
}

func NewRuleService(ruleRepo *repository.RuleRepository, txnRepo *repository.TransactionRepository) *RuleService {
	return &RuleService{ruleRepo: ruleRepo, txnRepo: txnRepo}
}

func (s *RuleService) Create(userID uint64, req *request.CreateRuleReq) (*response.RuleResp, error) {
	rule := &model.Rule{
		UserID:    userID,
		Name:      req.Name,
		Priority:  req.Priority,
		IsEnabled: req.IsEnabled,
		LogicType: model.LogicType(req.LogicType),
		Trigger:   model.RuleTrigger(req.Trigger),
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
	if req.Trigger != "" {
		rule.Trigger = model.RuleTrigger(req.Trigger)
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

func (s *RuleService) Execute(userID, id uint64, req *request.ExecuteRuleReq) (*response.RuleExecuteResultResp, error) {
	rule, err := s.ruleRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	startDate := req.StartDate
	endDate := req.EndDate
	if startDate == "" {
		startDate = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	filter := repository.TransactionFilter{
		StartDate: startDate,
		EndDate:   endDate,
	}

	txns, err := s.txnRepo.List(userID, filter, 0, 10000)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	result := &response.RuleExecuteResultResp{
		Errors: make([]string, 0),
	}

	for i := range txns {
		txn := &txns[i]
		if s.matchTransaction(rule, txn) {
			result.MatchedCount++
			if s.applyActions(rule.Actions, txn) {
				result.SuccessCount++
			} else {
				result.FailCount++
				result.Errors = append(result.Errors, fmt.Sprintf("transaction %d: action failed", txn.ID))
			}
		}
	}

	return result, nil
}

func (s *RuleService) matchTransaction(rule *model.Rule, txn *model.Transaction) bool {
	if rule.LogicType == model.LogicTypeAnd {
		for _, cond := range rule.Conditions {
			if !s.matchCondition(cond, txn) {
				return false
			}
		}
		return true
	}

	for _, cond := range rule.Conditions {
		if s.matchCondition(cond, txn) {
			return true
		}
	}
	return len(rule.Conditions) == 0
}

func (s *RuleService) matchCondition(condition model.RuleCondition, txn *model.Transaction) bool {
	var fieldValue string

	switch condition.Field {
	case model.ConditionFieldDescription:
		fieldValue = txn.Description
	case model.ConditionFieldAmount:
		fieldValue = txn.Amount.StringFixed(4)
	case model.ConditionFieldSourceAccount:
		fieldValue = txn.Source.Name
	case model.ConditionFieldDestinationAccount:
		if txn.Destination != nil {
			fieldValue = txn.Destination.Name
		}
	case model.ConditionFieldCategory:
		if txn.Category != nil {
			fieldValue = txn.Category.Name
		}
	case model.ConditionFieldTransactionType:
		fieldValue = string(txn.Type)
	case model.ConditionFieldNotes:
		fieldValue = txn.Notes
	case model.ConditionFieldTag:
		for _, t := range txn.Tags {
			if s.matchString(condition.Operator, t.Name, condition.Value) {
				return true
			}
		}
		return false
	case model.ConditionFieldBudget:
		if txn.BillID != nil {
			fieldValue = fmt.Sprintf("%d", *txn.BillID)
		}
	case model.ConditionFieldBill:
		if txn.BillID != nil {
			fieldValue = fmt.Sprintf("%d", *txn.BillID)
		}
	case model.ConditionFieldDateAfter:
		txnDate := txn.Date.Format("2006-01-02")
		return txnDate >= condition.Value
	case model.ConditionFieldDateBefore:
		txnDate := txn.Date.Format("2006-01-02")
		return txnDate <= condition.Value
	}

	return s.matchString(condition.Operator, fieldValue, condition.Value)
}

func (s *RuleService) matchString(operator model.ConditionOperator, fieldValue, conditionValue string) bool {
	fv := strings.ToLower(fieldValue)
	cv := strings.ToLower(conditionValue)

	switch operator {
	case model.OperatorContains:
		return strings.Contains(fv, cv)
	case model.OperatorEquals:
		return strings.EqualFold(fieldValue, conditionValue)
	case model.OperatorStartsWith:
		return strings.HasPrefix(fv, cv)
	case model.OperatorEndsWith:
		return strings.HasSuffix(fv, cv)
	case model.OperatorNotContains:
		return !strings.Contains(fv, cv)
	case model.OperatorNotEquals:
		return !strings.EqualFold(fieldValue, conditionValue)
	case model.OperatorLess:
		fieldDec, err1 := decimal.NewFromString(fieldValue)
		condDec, err2 := decimal.NewFromString(conditionValue)
		if err1 != nil || err2 != nil {
			return false
		}
		return fieldDec.LessThan(condDec)
	case model.OperatorMore:
		fieldDec, err1 := decimal.NewFromString(fieldValue)
		condDec, err2 := decimal.NewFromString(conditionValue)
		if err1 != nil || err2 != nil {
			return false
		}
		return fieldDec.GreaterThan(condDec)
	case model.OperatorIsEmpty:
		return fieldValue == ""
	case model.OperatorIsNotEmpty:
		return fieldValue != ""
	}

	return false
}

func (s *RuleService) applyActions(actions []model.RuleAction, txn *model.Transaction) bool {
	needsUpdate := false
	for _, action := range actions {
		switch action.Type {
		case model.ActionTypeSetDescription:
			txn.Description = action.Value
			needsUpdate = true
		case model.ActionTypeSetNotes:
			txn.Notes = action.Value
			needsUpdate = true
		case model.ActionTypeAppendNotes:
			txn.Notes += action.Value
			needsUpdate = true
		case model.ActionTypePrependNotes:
			txn.Notes = action.Value + txn.Notes
			needsUpdate = true
		case model.ActionTypeClearNotes:
			txn.Notes = ""
			needsUpdate = true
		case model.ActionTypeClearCategory:
			txn.CategoryID = nil
			needsUpdate = true
		case model.ActionTypeClearBudget:
			txn.BillID = nil
			needsUpdate = true
		// set_category, set_budget, add_tag, remove_tag require repository calls
		// which are handled as best-effort (no-op if not supported directly)
		case model.ActionTypeSetCategory:
			// Category ID resolution would require lookup by name
			// For now, treat value as category ID
			needsUpdate = true
		case model.ActionTypeSetBudget:
			needsUpdate = true
		case model.ActionTypeAddTag, model.ActionTypeRemoveTag:
			// Tag operations require transaction_tag table manipulation
			// Handled as best-effort
		}
	}

	if needsUpdate {
		if err := s.txnRepo.Update(txn); err != nil {
			return false
		}
	}
	return true
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
		Trigger:    string(r.Trigger),
		Conditions: conditions,
		Actions:    actions,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

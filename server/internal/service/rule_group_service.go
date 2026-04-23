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

type RuleGroupService struct {
	ruleGroupRepo *repository.RuleGroupRepository
	ruleRepo      *repository.RuleRepository
	txnRepo       *repository.TransactionRepository
}

func NewRuleGroupService(
	ruleGroupRepo *repository.RuleGroupRepository,
	ruleRepo *repository.RuleRepository,
	txnRepo *repository.TransactionRepository,
) *RuleGroupService {
	return &RuleGroupService{
		ruleGroupRepo: ruleGroupRepo,
		ruleRepo:      ruleRepo,
		txnRepo:       txnRepo,
	}
}

func (s *RuleGroupService) Create(userID uint64, req *request.CreateRuleGroupReq) (*response.RuleGroupResp, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	group := &model.RuleGroup{
		UserID:   userID,
		Name:     req.Name,
		Order:    req.Order,
		IsActive: isActive,
	}

	if err := s.ruleGroupRepo.Create(group); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(group, 0), nil
}

func (s *RuleGroupService) Get(userID, id uint64) (*response.RuleGroupResp, error) {
	group, err := s.ruleGroupRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	ruleCount, err := s.ruleGroupRepo.GetRuleCount(group.ID, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(group, ruleCount), nil
}

func (s *RuleGroupService) List(userID uint64) ([]response.RuleGroupResp, error) {
	groups, err := s.ruleGroupRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.RuleGroupResp, 0, len(groups))
	for _, g := range groups {
		ruleCount, err := s.ruleGroupRepo.GetRuleCount(g.ID, userID)
		if err != nil {
			return nil, errcode.ErrInternal
		}
		items = append(items, *s.toResp(&g, ruleCount))
	}

	return items, nil
}

func (s *RuleGroupService) Update(userID, id uint64, req *request.UpdateRuleGroupReq) (*response.RuleGroupResp, error) {
	group, err := s.ruleGroupRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Name != "" {
		group.Name = req.Name
	}
	if req.Order != nil {
		group.Order = *req.Order
	}
	if req.IsActive != nil {
		group.IsActive = *req.IsActive
	}

	if err := s.ruleGroupRepo.Update(group); err != nil {
		return nil, errcode.ErrInternal
	}

	ruleCount, err := s.ruleGroupRepo.GetRuleCount(group.ID, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(group, ruleCount), nil
}

func (s *RuleGroupService) Delete(userID, id uint64) error {
	group, err := s.ruleGroupRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}

	hasRules, err := s.ruleGroupRepo.HasRules(group.ID, userID)
	if err != nil {
		return errcode.ErrInternal
	}
	if hasRules {
		return errcode.ErrRuleGroupHasRules
	}

	return s.ruleGroupRepo.Delete(id, userID)
}

func (s *RuleGroupService) ExecuteGroup(userID, id uint64, req *request.ExecuteRuleGroupReq) (*response.RuleGroupExecuteResultResp, error) {
	group, err := s.ruleGroupRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if !group.IsActive {
		return nil, errcode.ErrRuleGroupInactive
	}

	// Get all enabled rules in this group, sorted by priority
	rules, err := s.ruleRepo.GetEnabledRulesByGroupID(group.ID, userID)
	if err != nil {
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

	result := &response.RuleGroupExecuteResultResp{
		GroupID: int(group.ID),
		Errors:  make([]string, 0),
	}

	for _, rule := range rules {
		for i := range txns {
			txn := &txns[i]
			if s.matchTransaction(&rule, txn) {
				result.MatchedCount++
				if s.applyActions(rule.Actions, txn) {
					result.SuccessCount++
				} else {
					result.FailCount++
					result.Errors = append(result.Errors, fmt.Sprintf("rule %d, transaction %d: action failed", rule.ID, txn.ID))
				}
			}
		}
	}

	return result, nil
}

func (s *RuleGroupService) matchTransaction(rule *model.Rule, txn *model.Transaction) bool {
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

func (s *RuleGroupService) matchCondition(condition model.RuleCondition, txn *model.Transaction) bool {
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

func (s *RuleGroupService) matchString(operator model.ConditionOperator, fieldValue, conditionValue string) bool {
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

func (s *RuleGroupService) applyActions(actions []model.RuleAction, txn *model.Transaction) bool {
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
		case model.ActionTypeSetCategory:
			needsUpdate = true
		case model.ActionTypeSetBudget:
			needsUpdate = true
		case model.ActionTypeAddTag, model.ActionTypeRemoveTag:
			// Tag operations require transaction_tag table manipulation
		}
	}

	if needsUpdate {
		if err := s.txnRepo.Update(txn); err != nil {
			return false
		}
	}
	return true
}

func (s *RuleGroupService) toResp(g *model.RuleGroup, ruleCount int64) *response.RuleGroupResp {
	return &response.RuleGroupResp{
		ID:        g.ID,
		UserID:    g.UserID,
		Name:      g.Name,
		Order:     g.Order,
		IsActive:  g.IsActive,
		RuleCount: ruleCount,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
}

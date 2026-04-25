// Package service 业务逻辑层，实现核心业务逻辑
// RuleService 规则业务逻辑，处理规则匹配和动作执行
package service

import (
	"fmt"
	"log"
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

// RuleService 规则服务
// 负责处理规则的增删改查、条件匹配、动作执行和规则触发
// 规则由条件（conditions）和动作（actions）组成，支持AND/OR逻辑组合
// 依赖ruleRepo进行规则数据访问，依赖txnRepo查询和更新交易
// 依赖categoryRepo按名称查找分类，依赖tagRepo按名称查找标签，依赖budgetRepo按名称查找预算
type RuleService struct {
	ruleRepo     *repository.RuleRepository         // 规则数据访问对象
	txnRepo      *repository.TransactionRepository  // 交易数据访问对象
	categoryRepo *repository.CategoryRepository     // 分类数据访问对象
	tagRepo      *repository.TagRepository          // 标签数据访问对象
	budgetRepo   *repository.BudgetRepository       // 预算数据访问对象
}

// NewRuleService 创建规则服务实例
func NewRuleService(
	ruleRepo *repository.RuleRepository,
	txnRepo *repository.TransactionRepository,
	categoryRepo *repository.CategoryRepository,
	budgetRepo *repository.BudgetRepository,
	tagRepo *repository.TagRepository,
) *RuleService {
	return &RuleService{
		ruleRepo:     ruleRepo,
		txnRepo:      txnRepo,
		categoryRepo: categoryRepo,
		budgetRepo:   budgetRepo,
		tagRepo:      tagRepo,
	}
}

// Create 创建规则
// 同时创建规则的条件列表和动作列表
// 参数：
//   - userID: 用户ID
//   - req: 创建请求参数（名称、优先级、启用状态、逻辑类型、触发类型、条件列表、动作列表）
// 返回：
//   - *response.RuleResp: 创建成功的规则信息
//   - error: 错误信息
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

// Get 获取单个规则详情（含条件和动作）
// 参数：
//   - userID: 用户ID
//   - id: 规则ID
// 返回：
//   - *response.RuleResp: 规则信息
//   - error: 错误信息
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

// List 获取用户所有规则列表
// 参数：
//   - userID: 用户ID
// 返回：
//   - []response.RuleResp: 规则列表
//   - error: 错误信息
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

// Update 更新规则（同时替换条件和动作列表）
// 参数：
//   - userID: 用户ID
//   - id: 规则ID
//   - req: 更新请求参数
// 返回：
//   - *response.RuleResp: 更新后的规则信息
//   - error: 错误信息
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

// Delete 删除规则
// 参数：
//   - userID: 用户ID
//   - id: 规则ID
// 返回：
//   - error: 错误信息
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

// ToggleStatus 切换规则的启用/禁用状态
// 参数：
//   - userID: 用户ID
//   - id: 规则ID
// 返回：
//   - *response.RuleResp: 更新后的规则信息
//   - error: 错误信息
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

// Execute 手动执行规则
// 在指定时间范围内查找匹配条件的交易，对匹配的交易应用规则动作
// 参数：
//   - userID: 用户ID
//   - id: 规则ID
//   - req: 执行请求参数（含日期范围，默认最近一个月）
// 返回：
//   - *response.RuleExecuteResultResp: 执行结果（匹配数、成功数、失败数）
//   - error: 错误信息
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
			if s.applyActions(userID, rule.Actions, txn) {
				result.SuccessCount++
			} else {
				result.FailCount++
				result.Errors = append(result.Errors, fmt.Sprintf("transaction %d: action failed", txn.ID))
			}
		}
	}

	return result, nil
}

// matchTransaction 判断交易是否匹配规则的所有条件
// AND逻辑：所有条件都满足才匹配；OR逻辑：任一条件满足即匹配
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

// matchCondition 判断交易是否匹配单个条件
// 根据条件字段（描述、金额、账户、分类、标签等）提取交易字段值，再用运算符比较
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

// matchString 根据运算符比较字段值和条件值
// 支持：包含、等于、开头是、结尾是、不包含、不等于、小于、大于、为空、不为空
// 字符串比较不区分大小写，数值比较使用decimal精确计算
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

// applyActions 对交易应用规则动作列表
// 支持的动作：设置描述、设置/追加/前置/清空备注、清空/设置分类、清空预算、添加/移除标签
// 如果有任何变更，使用UpdateWithTags同时更新交易和标签关联
func (s *RuleService) applyActions(userID uint64, actions []model.RuleAction, txn *model.Transaction) bool {
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
			// Find category by name from user's category tree
			categories, err := s.categoryRepo.List(userID)
			if err != nil {
				log.Printf("[RuleService] set_category: failed to list categories for user %d: %v", userID, err)
				continue
			}
			found := false
			for _, cat := range categories {
				if strings.EqualFold(cat.Name, action.Value) {
					txn.CategoryID = &cat.ID
					needsUpdate = true
					found = true
					break
				}
			}
			if !found {
				log.Printf("[RuleService] set_category: category '%s' not found for user %d", action.Value, userID)
			}
		case model.ActionTypeSetBudget:
			// Transaction model does not have BudgetID field, skip for now
			log.Printf("[RuleService] set_budget: action skipped (Transaction model has no BudgetID field), budget name='%s'", action.Value)
		case model.ActionTypeAddTag:
			// Find tag by name, add to transaction if not already present
			tags, err := s.tagRepo.List(userID)
			if err != nil {
				log.Printf("[RuleService] add_tag: failed to list tags for user %d: %v", userID, err)
				continue
			}
			for _, tag := range tags {
				if strings.EqualFold(tag.Name, action.Value) {
					// Check if tag already on transaction
					alreadyHas := false
					for _, t := range txn.Tags {
						if t.ID == tag.ID {
							alreadyHas = true
							break
						}
					}
					if !alreadyHas {
						txn.Tags = append(txn.Tags, tag)
						needsUpdate = true
					}
					break
				}
			}
		case model.ActionTypeRemoveTag:
			// Remove tag by name from transaction
			newTags := make([]model.Tag, 0, len(txn.Tags))
			removed := false
			for _, t := range txn.Tags {
				if strings.EqualFold(t.Name, action.Value) {
					removed = true
					continue
				}
				newTags = append(newTags, t)
			}
			if removed {
				txn.Tags = newTags
				needsUpdate = true
			}
		}
	}

	if needsUpdate {
		// Collect tag IDs for UpdateWithTags
		tagIDs := make([]uint64, 0, len(txn.Tags))
		for _, t := range txn.Tags {
			tagIDs = append(tagIDs, t.ID)
		}
		if err := s.txnRepo.UpdateWithTags(txn, tagIDs); err != nil {
			log.Printf("[RuleService] applyActions: failed to update transaction %d: %v", txn.ID, err)
			return false
		}
	}
	return true
}

// RuleTrigger 规则触发器接口
// 定义了在交易变更后触发规则的方法，用于解耦TransactionService和RuleService避免循环依赖
type RuleTrigger interface {
	TriggerRules(userID uint64, txn *model.Transaction, triggerType string) error
}

// TriggerRules 触发规则引擎
// 获取用户所有启用的规则，按触发类型过滤，逐条匹配交易并执行动作
// 单条规则失败不会阻止后续规则执行
// 参数：
//   - userID: 用户ID
//   - txn: 触发规则的目标交易
//   - triggerType: 触发类型（on_create/on_update）
// 返回：
//   - error: 错误信息
func (s *RuleService) TriggerRules(userID uint64, txn *model.Transaction, triggerType string) error {
	// 获取用户所有启用的规则，按优先级排序
	rules, err := s.ruleRepo.GetEnabledRules(userID)
	if err != nil {
		log.Printf("[RuleService] TriggerRules: failed to get enabled rules for user %d: %v", userID, err)
		return err
	}

	// Filter rules by trigger type and sort by group order then priority
	// GetEnabledRules already sorts by priority ASC, id ASC
	// We further filter by trigger type
	// 按触发类型过滤规则（on_create/on_update）
	matchingRules := make([]model.Rule, 0)
	for _, rule := range rules {
		if rule.Trigger == model.RuleTrigger(triggerType) {
			matchingRules = append(matchingRules, rule)
		}
	}

	// Sort by group order then priority
	// Since we don't have group info preloaded in the rule list,
	// we sort by priority (already sorted from repo) and process in order.
	// For full group-order sorting, we would need to join with rule_groups table.
	// 按优先级顺序执行匹配的规则，单条规则失败不阻止后续规则执行
	for _, rule := range matchingRules {
		if s.matchTransaction(&rule, txn) {
			if !s.applyActions(userID, rule.Actions, txn) {
				log.Printf("[RuleService] TriggerRules: rule %d failed for transaction %d, continuing", rule.ID, txn.ID)
			}
		}
	}

	return nil
}

// toResp 将规则模型转换为响应对象（含条件和动作列表）
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

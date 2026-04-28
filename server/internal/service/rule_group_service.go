// Package service 业务逻辑层，实现核心业务逻辑
// RuleGroupService 规则组业务逻辑，处理规则组的增删改查和组执行
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

// RuleGroupService 规则组服务
// 依赖ruleGroupRepo/ruleRepo进行规则数据访问，依赖txnRepo查询和更新交易
// 依赖categoryRepo按名称查找分类，依赖tagRepo按名称查找标签，依赖budgetRepo按名称查找预算
type RuleGroupService struct {
	ruleGroupRepo *repository.RuleGroupRepository
	ruleRepo      *repository.RuleRepository
	txnRepo       *repository.TransactionRepository
	categoryRepo  *repository.CategoryRepository
	tagRepo       *repository.TagRepository
	budgetRepo    *repository.BudgetRepository
}

// NewRuleGroupService 创建规则组服务实例
func NewRuleGroupService(
	ruleGroupRepo *repository.RuleGroupRepository,
	ruleRepo *repository.RuleRepository,
	txnRepo *repository.TransactionRepository,
	categoryRepo *repository.CategoryRepository,
	tagRepo *repository.TagRepository,
	budgetRepo *repository.BudgetRepository,
) *RuleGroupService {
	return &RuleGroupService{
		ruleGroupRepo: ruleGroupRepo,
		ruleRepo:      ruleRepo,
		txnRepo:       txnRepo,
		categoryRepo:  categoryRepo,
		tagRepo:       tagRepo,
		budgetRepo:    budgetRepo,
	}
}

// Create 创建规则组
// 参数：
//   - userID: 用户ID
//   - req: 创建请求参数（名称、排序、启用状态）
// 返回：
//   - *response.RuleGroupResp: 创建成功的规则组信息
//   - error: 错误信息
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

// Get 获取单个规则组详情（含规则数量）
// 参数：
//   - userID: 用户ID
//   - id: 规则组ID
// 返回：
//   - *response.RuleGroupResp: 规则组信息
//   - error: 错误信息
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

// List 获取用户所有规则组列表（每个规则组含规则数量）
// 参数：
//   - userID: 用户ID
// 返回：
//   - []response.RuleGroupResp: 规则组列表
//   - error: 错误信息
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

// Update 更新规则组信息
// 支持更新名称、排序、启用状态
// 参数：
//   - userID: 用户ID
//   - id: 规则组ID
//   - req: 更新请求参数
// 返回：
//   - *response.RuleGroupResp: 更新后的规则组信息
//   - error: 错误信息
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

// Delete 删除规则组
// 业务规则：如果规则组内还有规则，则不允许删除（需先移除或删除组内规则）
// 参数：
//   - userID: 用户ID
//   - id: 规则组ID
// 返回：
//   - error: 错误信息（如规则组不存在、组内有规则）
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

// ExecuteGroup 执行规则组：对匹配的交易应用规则操作
// 业务流程：
// 1. 获取规则组信息，检查是否启用
// 2. 获取规则组内所有启用的规则（按优先级排序）
// 3. 获取指定时间范围内的交易（默认最近一个月）
// 4. 对每笔交易逐条匹配规则条件，匹配成功则执行规则动作
// 参数：
//   - userID: 用户ID
//   - id: 规则组ID
//   - req: 执行请求参数（含日期范围）
// 返回：
//   - *response.RuleGroupExecuteResultResp: 执行结果（匹配数、成功数、失败数）
//   - error: 错误信息
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

	// 获取规则组内所有启用的规则，按优先级排序
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

	txns, err := s.txnRepo.ListAll(userID, filter)
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
				if s.applyActions(rule.Actions, txn, userID) {
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

// matchTransaction 判断交易是否匹配规则的所有条件
// AND逻辑：所有条件都满足才匹配；OR逻辑：任一条件满足即匹配
// 参数：
//   - rule: 规则（含条件和逻辑类型）
//   - txn: 待匹配的交易
// 返回：
//   - bool: 是否匹配
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

// matchCondition 判断交易是否匹配单个条件
// 根据条件字段（描述、金额、账户、分类、标签、日期等）提取交易字段值，再用运算符比较
// 特殊处理：标签字段需遍历所有标签逐个匹配；日期字段直接比较字符串
// 参数：
//   - condition: 规则条件
//   - txn: 待匹配的交易
// 返回：
//   - bool: 是否匹配
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
		return false
	case model.ConditionFieldBill:
		if txn.RecurringID != nil {
			fieldValue = fmt.Sprintf("%d", *txn.RecurringID)
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
// 参数：
//   - operator: 比较运算符
//   - fieldValue: 交易字段值
//   - conditionValue: 条件值
// 返回：
//   - bool: 是否匹配
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

// applyActions 对交易应用规则操作列表
// 修复：实现SetCategory、SetBudget、AddTag、RemoveTag操作
// - SetCategory: 按名称查找分类并设置CategoryID
// - SetBudget: 交易模型无BudgetID字段，记录警告日志
// - AddTag: 按名称查找标签并添加到交易
// - RemoveTag: 按名称查找标签并从交易中移除
func (s *RuleGroupService) applyActions(actions []model.RuleAction, txn *model.Transaction, userID uint64) bool {
	needsUpdate := false
	var tagIDs []uint64

	// 收集当前交易的标签ID列表
	for _, tag := range txn.Tags {
		tagIDs = append(tagIDs, tag.ID)
	}

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
		case model.ActionTypeClearRecurring:
			txn.RecurringID = nil
			needsUpdate = true
		case model.ActionTypeSetCategory:
			// 按名称查找分类并设置CategoryID
			categories, err := s.categoryRepo.List(userID)
			if err != nil {
				continue
			}
			for _, cat := range categories {
				if strings.EqualFold(cat.Name, action.Value) {
					catID := cat.ID
					txn.CategoryID = &catID
					needsUpdate = true
					break
				}
			}
		case model.ActionTypeSetBudget:
			// 交易模型没有BudgetID字段，无法直接设置预算关联
			// 记录警告日志提示此操作暂不支持
			log.Printf("[RuleGroupService] SetBudget action not supported: transaction model has no BudgetID field, value=%s", action.Value)
		case model.ActionTypeAddTag:
			// 按名称查找标签并添加到交易的标签列表
			tags, err := s.tagRepo.List(userID)
			if err != nil {
				continue
			}
			for _, tag := range tags {
				if strings.EqualFold(tag.Name, action.Value) {
					// 检查标签是否已存在，避免重复添加
					found := false
					for _, existingID := range tagIDs {
						if existingID == tag.ID {
							found = true
							break
						}
					}
					if !found {
						tagIDs = append(tagIDs, tag.ID)
						needsUpdate = true
					}
					break
				}
			}
		case model.ActionTypeRemoveTag:
			// 按名称查找标签并从交易的标签列表中移除
			tags, err := s.tagRepo.List(userID)
			if err != nil {
				continue
			}
			for _, tag := range tags {
				if strings.EqualFold(tag.Name, action.Value) {
					// 从标签ID列表中移除匹配的标签
					newTagIDs := make([]uint64, 0, len(tagIDs))
					for _, existingID := range tagIDs {
						if existingID != tag.ID {
							newTagIDs = append(newTagIDs, existingID)
						}
					}
					if len(newTagIDs) != len(tagIDs) {
						tagIDs = newTagIDs
						needsUpdate = true
					}
					break
				}
			}
		}
	}

	if needsUpdate {
		// 如果涉及标签变更，使用UpdateWithTags同时更新交易和标签关联
		if err := s.txnRepo.UpdateWithTags(txn, tagIDs); err != nil {
			return false
		}
	}
	return true
}

// toResp 将规则组模型转换为响应DTO
// 参数：
//   - g3: 规则组模型
//   - ruleCount: 规则组内的规则数量
// 返回：
//   - *response.RuleGroupResp: 规则组响应对象
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

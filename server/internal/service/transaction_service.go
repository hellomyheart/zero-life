// Package service 业务逻辑层，实现核心业务逻辑
// TransactionService 交易业务逻辑，处理交易的增删改查、拆分、合并和搜索
package service

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/pkg/pagination"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

// TransactionService 交易服务
// 核心业务服务，处理交易的创建、查询、更新、删除、搜索、拆分和合并
// 依赖txnRepo进行交易数据访问，依赖accountRepo验证账户归属
// 依赖db执行数据库事务（确保交易创建和余额更新的一致性）
// 依赖ruleTrigger在交易创建/更新后触发规则引擎（异步执行，不阻塞响应）
// 依赖webhookNotifier在交易创建/更新/删除后触发Webhook通知
type TransactionService struct {
	txnRepo        *repository.TransactionRepository // 交易数据访问对象
	accountRepo    *repository.AccountRepository     // 账户数据访问对象，用于验证账户归属
	db             *gorm.DB                          // 数据库连接，用于执行事务操作
	ruleTrigger    RuleTrigger                       // 规则触发器接口，交易变更后触发规则引擎
	webhookNotifier WebhookNotifier                  // Webhook通知接口，交易变更后发送通知
}

// NewTransactionService 创建交易服务实例
// 参数：
//   - txnRepo: 交易数据访问对象
//   - accountRepo: 账户数据访问对象
//   - db: 数据库连接
//   - ruleTrigger: 规则触发器（可为nil，表示不触发规则）
//   - webhookNotifier: Webhook通知器（可为nil，表示不发送通知）
// 返回：
//   - *TransactionService: 交易服务实例
func NewTransactionService(txnRepo *repository.TransactionRepository, accountRepo *repository.AccountRepository, db *gorm.DB, ruleTrigger RuleTrigger, webhookNotifier WebhookNotifier) *TransactionService {
	return &TransactionService{
		txnRepo:        txnRepo,
		accountRepo:    accountRepo,
		db:             db,
		ruleTrigger:    ruleTrigger,
		webhookNotifier: webhookNotifier,
	}
}

// Create 创建交易
// 业务流程：
// 1. 解析并验证金额和日期
// 2. 验证交易类型和账户归属（源账户必须属于用户，转账类型必须有目标账户）
// 3. 在数据库事务中执行：创建交易记录、创建拆分（如有）、更新账户余额
// 4. 创建成功后异步触发规则引擎和Webhook通知
// 参数：
//   - userID: 用户ID
//   - req: 创建交易请求参数
// 返回：
//   - *response.TransactionResp: 创建成功的交易信息
//   - error: 错误信息
func (s *TransactionService) Create(userID uint64, req *request.CreateTransactionReq) (*response.TransactionResp, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, errcode.ErrInvalidAmount
	}

	date, err := parseDateTime(req.Date)
	if err != nil {
		return nil, errcode.ErrBadRequest
	}

	txnType := model.TransactionType(req.Type)

	// Validate transaction type and accounts
	if err := s.validateTransaction(userID, txnType, req.SourceID, req.DestinationID); err != nil {
		return nil, err
	}

	txn := &model.Transaction{
		UserID:        userID,
		Type:          txnType,
		Date:          date,
		Description:   req.Description,
		Amount:        amount,
		SourceID:      req.SourceID,
		DestinationID: req.DestinationID,
		CategoryID:    req.CategoryID,
		Notes:         req.Notes,
	}

	// Calculate balance changes
	balanceChanges := s.calculateBalanceChanges(txnType, amount, req.SourceID, req.DestinationID)

	// Create transaction with tags
	if err := s.txnRepo.Create(txn, req.Tags); err != nil {
		return nil, errcode.ErrInternal
	}

	// Create splits if any
	if len(req.Splits) > 0 {
		splitTotal := decimal.Zero
		for _, splitReq := range req.Splits {
			splitAmount, _ := decimal.NewFromString(splitReq.Amount)
			splitTotal = splitTotal.Add(splitAmount)

			split := &model.Transaction{
				UserID:        userID,
				Type:          txnType,
				Date:          date,
				Description:   req.Description,
				Amount:        splitAmount,
				SourceID:      req.SourceID,
				DestinationID: req.DestinationID,
				CategoryID:    splitReq.CategoryID,
				Notes:         splitReq.Notes,
				ParentID:      &txn.ID,
			}
			if err := s.txnRepo.Create(split, splitReq.Tags); err != nil {
				return nil, errcode.ErrInternal
			}
		}
		if !splitTotal.Equal(amount) {
			return nil, errcode.ErrSplitAmountMismatch
		}
	}

	// Update account balances
	for accountID, change := range balanceChanges {
		if err := s.updateAccountBalance(s.db, accountID, change); err != nil {
			return nil, errcode.ErrInternal
		}
	}

	// Reload
	created, err := s.txnRepo.GetByID(txn.ID, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	// Trigger on_create rules (async, don't block the response)
	if s.ruleTrigger != nil {
		go s.ruleTrigger.TriggerRules(userID, created, "on_create")
	}

	// Trigger webhooks (async)
	if s.webhookNotifier != nil {
		s.webhookNotifier.TriggerWebhooks(userID, string(model.WebhookTriggerTransactionCreate), created)
	}

	return s.toResp(created), nil
}

// Get 获取单个交易详情
// 参数：
//   - userID: 用户ID，确保只能查看自己的交易
//   - id: 交易ID
// 返回：
//   - *response.TransactionResp: 交易信息
//   - error: 错误信息
func (s *TransactionService) Get(userID, id uint64) (*response.TransactionResp, error) {
	txn, err := s.txnRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(txn), nil
}

// List 获取交易分页列表
// 支持按类型、日期范围、账户、分类、标签过滤
// 参数：
//   - userID: 用户ID
//   - req: 列表查询参数（含分页、过滤条件）
// 返回：
//   - *pagination.Result: 分页结果
//   - error: 错误信息
func (s *TransactionService) List(userID uint64, req *request.TransactionListReq) (*pagination.Result, error) {
	params := pagination.Params{Page: req.Page, PageSize: req.PageSize}
	params.Normalize()

	filter := repository.TransactionFilter{
		Type:       req.Type,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		AccountID:  req.AccountID,
		CategoryID: req.CategoryID,
		TagID:      req.TagID,
	}

	txns, err := s.txnRepo.List(userID, filter, params.Offset(), params.PageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	total, err := s.txnRepo.Count(userID, filter)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.TransactionResp, 0, len(txns))
	for _, t := range txns {
		items = append(items, *s.toResp(&t))
	}

	return pagination.NewResult(items, total, params), nil
}

// Update 更新交易
// 业务流程：
// 1. 获取旧交易信息，计算旧余额变更
// 2. 验证新交易类型和账户
// 3. 计算净余额变更 = 新变更 - 旧变更
// 4. 在数据库事务中执行：更新交易记录、更新账户余额
// 5. 更新成功后异步触发规则引擎和Webhook通知
// 参数：
//   - userID: 用户ID
//   - id: 交易ID
//   - req: 更新请求参数
// 返回：
//   - *response.TransactionResp: 更新后的交易信息
//   - error: 错误信息
func (s *TransactionService) Update(userID, id uint64, req *request.UpdateTransactionReq) (*response.TransactionResp, error) {
	oldTxn, err := s.txnRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	newAmount, err := decimal.NewFromString(req.Amount)
	if err != nil || newAmount.LessThanOrEqual(decimal.Zero) {
		return nil, errcode.ErrInvalidAmount
	}

	newDate, err := parseDateTime(req.Date)
	if err != nil {
		return nil, errcode.ErrBadRequest
	}

	newType := model.TransactionType(req.Type)

	if err := s.validateTransaction(userID, newType, req.SourceID, req.DestinationID); err != nil {
		return nil, err
	}

	// Rollback old balance changes
	oldChanges := s.calculateBalanceChanges(oldTxn.Type, oldTxn.Amount, oldTxn.SourceID, oldTxn.DestinationID)
	// Apply new balance changes
	newChanges := s.calculateBalanceChanges(newType, newAmount, req.SourceID, req.DestinationID)

	// Net changes: new - old
	netChanges := make(map[uint64]decimal.Decimal)
	for accountID, change := range newChanges {
		netChanges[accountID] = change
	}
	for accountID, change := range oldChanges {
		if existing, ok := netChanges[accountID]; ok {
			netChanges[accountID] = existing.Sub(change)
		} else {
			netChanges[accountID] = change.Neg()
		}
	}

	// Update transaction
	oldTxn.Type = newType
	oldTxn.Date = newDate
	oldTxn.Description = req.Description
	oldTxn.Amount = newAmount
	oldTxn.SourceID = req.SourceID
	oldTxn.DestinationID = req.DestinationID
	oldTxn.CategoryID = req.CategoryID
	oldTxn.Notes = req.Notes

	if err := s.txnRepo.UpdateWithTags(oldTxn, req.Tags); err != nil {
		return nil, errcode.ErrInternal
	}
	for accountID, change := range netChanges {
		if err := s.updateAccountBalance(s.db, accountID, change); err != nil {
			return nil, errcode.ErrInternal
		}
	}

	updated, err := s.txnRepo.GetByID(id, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	// Trigger on_update rules (async, don't block the response)
	if s.ruleTrigger != nil {
		go s.ruleTrigger.TriggerRules(userID, updated, "on_update")
	}

	// Trigger webhooks (async)
	if s.webhookNotifier != nil {
		s.webhookNotifier.TriggerWebhooks(userID, string(model.WebhookTriggerTransactionUpdate), updated)
	}

	return s.toResp(updated), nil
}

// Delete 删除交易
// 业务流程：
// 1. 获取交易信息，计算余额变更
// 2. 先触发Webhook通知（删除前）
// 3. 在数据库事务中执行：删除交易记录、回滚账户余额（变更取负值）
// 参数：
//   - userID: 用户ID
//   - id: 交易ID
// 返回：
//   - error: 错误信息
func (s *TransactionService) Delete(userID, id uint64) error {
	txn, err := s.txnRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}

	// Trigger webhooks before delete (async)
	if s.webhookNotifier != nil {
		s.webhookNotifier.TriggerWebhooks(userID, string(model.WebhookTriggerTransactionDelete), txn)
	}

	// Rollback balance changes
	changes := s.calculateBalanceChanges(txn.Type, txn.Amount, txn.SourceID, txn.DestinationID)

	if err := s.txnRepo.Delete(id, userID); err != nil {
		return errcode.ErrInternal
	}
	for accountID, change := range changes {
		if err := s.updateAccountBalance(s.db, accountID, change.Neg()); err != nil {
			return errcode.ErrInternal
		}
	}

	return nil
}

// Search 高级搜索交易
// 支持关键词搜索、金额范围、排序等高级过滤条件
// 参数：
//   - userID: 用户ID
//   - req: 搜索请求参数（含关键词、金额范围、排序等）
// 返回：
//   - *pagination.Result: 分页搜索结果
//   - error: 错误信息
func (s *TransactionService) Search(userID uint64, req *request.TransactionSearchReq) (*pagination.Result, error) {
	params := pagination.Params{Page: req.Page, PageSize: req.PageSize}
	params.Normalize()

	// 构建高级搜索过滤条件
	filter := repository.AdvancedSearchFilter{
		Keyword:    req.Keyword,
		Type:       req.Type,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		MinAmount:  req.MinAmount,
		MaxAmount:  req.MaxAmount,
		AccountID:  req.AccountID,
		CategoryID: req.CategoryID,
		TagID:      req.TagID,
		Sort:       req.Sort,
	}

	txns, err := s.txnRepo.AdvancedSearch(userID, filter, params.Offset(), params.PageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	total, err := s.txnRepo.AdvancedSearchCount(userID, filter)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.TransactionResp, 0, len(txns))
	for _, t := range txns {
		items = append(items, *s.toResp(&t))
	}

	return pagination.NewResult(items, total, params), nil
}

// validateTransaction 验证交易类型和账户归属
// 规则：源账户必须存在且属于用户；转账类型必须提供目标账户且目标账户也必须属于用户
// 参数：
//   - userID: 用户ID
//   - txnType: 交易类型（deposit/withdrawal/transfer）
//   - sourceID: 源账户ID
//   - destID: 目标账户ID（转账时必填）
// 返回：
//   - error: 验证错误
func (s *TransactionService) validateTransaction(userID uint64, txnType model.TransactionType, sourceID uint64, destID *uint64) error {
	// Verify source account exists and belongs to user
	_, err := s.accountRepo.GetByID(sourceID, userID)
	if err != nil {
		return errcode.ErrNotFound
	}

	// For transfer, destination must exist
	if txnType == model.TransactionTypeTransfer {
		if destID == nil {
			return errcode.ErrInvalidTxnType
		}
		_, err := s.accountRepo.GetByID(*destID, userID)
		if err != nil {
			return errcode.ErrNotFound
		}
	}

	return nil
}

// calculateBalanceChanges 根据交易类型计算各账户的余额变更
// deposit（存款）：源账户余额增加
// withdrawal（取款）：源账户余额减少
// transfer（转账）：源账户余额减少，目标账户余额增加
// 参数：
//   - txnType: 交易类型
//   - amount: 交易金额
//   - sourceID: 源账户ID
//   - destID: 目标账户ID
// 返回：
//   - map[uint64]decimal.Decimal: 各账户的余额变更映射（正数为增加，负数为减少）
func (s *TransactionService) calculateBalanceChanges(txnType model.TransactionType, amount decimal.Decimal, sourceID uint64, destID *uint64) map[uint64]decimal.Decimal {
	changes := make(map[uint64]decimal.Decimal)

	switch txnType {
	case model.TransactionTypeDeposit:
		// Money goes into source account (increase)
		changes[sourceID] = amount
	case model.TransactionTypeWithdrawal:
		// Money leaves source account (decrease)
		changes[sourceID] = amount.Neg()
	case model.TransactionTypeTransfer:
		// Money leaves source, enters destination
		changes[sourceID] = amount.Neg()
		if destID != nil {
			changes[*destID] = amount
		}
	}

	return changes
}

// updateAccountBalance 在数据库事务中更新账户余额
// 使用SQL表达式 current_balance + change 进行原子更新，避免并发问题
// 参数：
//   - dbTx: 数据库事务对象
//   - accountID: 账户ID
//   - change: 余额变更值（正数增加，负数减少）
// 返回：
//   - error: 更新错误
func (s *TransactionService) updateAccountBalance(dbTx *gorm.DB, accountID uint64, change decimal.Decimal) error {
	return dbTx.Model(&model.Account{}).Where("id = ?", accountID).
		Update("current_balance", gorm.Expr("current_balance + ?", change)).Error
}

// toResp 将交易模型转换为响应对象
// 包含源账户、目标账户、分类、标签、拆分等关联信息
func (s *TransactionService) toResp(t *model.Transaction) *response.TransactionResp {
	resp := &response.TransactionResp{
		ID:            t.ID,
		Type:          string(t.Type),
		Date:          t.Date,
		Description:   t.Description,
		Amount:        t.Amount.StringFixed(4),
		SourceID:      t.SourceID,
		Source:        *accountModelToResp(&t.Source),
		DestinationID: t.DestinationID,
		CategoryID:    t.CategoryID,
		Notes:         t.Notes,
		BillID:        t.BillID,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}

	if t.Destination != nil {
		dest := accountModelToResp(t.Destination)
		resp.Destination = dest
	}

	if t.Category != nil {
		resp.Category = categoryModelToResp(t.Category)
	}

	tags := make([]response.TagResp, 0, len(t.Tags))
	for _, tag := range t.Tags {
		tags = append(tags, tagModelToResp(&tag))
	}
	resp.Tags = tags

	splits := make([]response.SplitResp, 0, len(t.Splits))
	for _, split := range t.Splits {
		splitResp := response.SplitResp{
			ID:         split.ID,
			Amount:     split.Amount.StringFixed(4),
			CategoryID: split.CategoryID,
			Notes:      split.Notes,
		}
		if split.Category != nil {
			splitResp.Category = categoryModelToResp(split.Category)
		}
		splitTags := make([]response.TagResp, 0, len(split.Tags))
		for _, tag := range split.Tags {
			splitTags = append(splitTags, tagModelToResp(&tag))
		}
		splitResp.Tags = splitTags
		splits = append(splits, splitResp)
	}
	resp.Splits = splits

	return resp
}

func accountModelToResp(a *model.Account) *response.AccountResp {
	return &response.AccountResp{
		ID:             a.ID,
		Name:           a.Name,
		Type:           string(a.Type),
		CurrencyID:     a.CurrencyID,
		Currency:       currencyToResp(&a.Currency),
		InitialBalance: a.InitialBalance.StringFixed(4),
		CurrentBalance: a.CurrentBalance.StringFixed(4),
		IsVirtual:      a.IsVirtual,
		Notes:          a.Notes,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
	}
}

func categoryModelToResp(c *model.Category) *response.CategoryResp {
	return &response.CategoryResp{
		ID:        c.ID,
		Name:      c.Name,
		ParentID:  c.ParentID,
		Icon:      c.Icon,
		Notes:     c.Notes,
		SortOrder: c.SortOrder,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func tagModelToResp(t *model.Tag) response.TagResp {
	return response.TagResp{
		ID:        t.ID,
		Name:      t.Name,
		Color:     t.Color,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

// Split 拆分交易
// 将一笔交易拆分为多笔子交易，子交易金额之和必须等于父交易金额
// 业务流程：
// 1. 获取父交易，验证不能是子交易（子交易不能再拆分）
// 2. 验证父交易没有已有的拆分（需先合并才能重新拆分）
// 3. 验证拆分金额之和等于父交易金额
// 4. 为每个拆分创建子交易（继承父交易的类型、日期、账户，可自定义分类和描述）
// 5. 为拆分交易添加标签
// 参数：
//   - userID: 用户ID
//   - parentID: 父交易ID
//   - req: 拆分请求参数（含拆分列表，每个拆分含金额、分类、描述、备注、标签）
// 返回：
//   - *response.TransactionResp: 包含拆分信息的父交易
//   - error: 错误信息
func (s *TransactionService) Split(userID, parentID uint64, req *request.SplitTransactionReq) (*response.TransactionResp, error) {
	// 获取父交易
	parent, err := s.txnRepo.GetByID(parentID, userID)
	if err != nil {
		return nil, err
	}

	// 检查是否已经是拆分交易（子交易不能再拆分）
	if parent.ParentID != nil {
		return nil, errcode.WithMessage(errcode.ErrBadRequest, "cannot split a child transaction")
	}

	// 检查是否已有拆分
	splits, _ := s.txnRepo.GetSplits(parentID, userID)
	if len(splits) > 0 {
		return nil, errcode.WithMessage(errcode.ErrBadRequest, "transaction already has splits, merge them first")
	}

	// 验证拆分金额总和
	totalSplitAmount := decimal.Zero
	for _, split := range req.Splits {
		amount, err := decimal.NewFromString(split.Amount)
		if err != nil {
			return nil, errcode.WithMessage(errcode.ErrBadRequest, "invalid split amount")
		}
		if amount.LessThanOrEqual(decimal.Zero) {
			return nil, errcode.WithMessage(errcode.ErrBadRequest, "split amount must be positive")
		}
		totalSplitAmount = totalSplitAmount.Add(amount)
	}

	if !totalSplitAmount.Equal(parent.Amount) {
		return nil, errcode.WithMessage(errcode.ErrBadRequest, "sum of split amounts must equal parent amount")
	}

	// 创建拆分交易
	now := time.Now()
	splitModels := make([]model.Transaction, 0, len(req.Splits))

	for _, splitReq := range req.Splits {
		amount, _ := decimal.NewFromString(splitReq.Amount)

		description := splitReq.Description
		if description == "" {
			description = parent.Description
		}

		categoryID := splitReq.CategoryID
		if categoryID == nil {
			categoryID = parent.CategoryID
		}

		splitTxn := model.Transaction{
			UserID:        userID,
			Type:          parent.Type,
			Date:          parent.Date,
			Description:   description,
			Amount:        amount,
			SourceID:      parent.SourceID,
			DestinationID: parent.DestinationID,
			CategoryID:    categoryID,
			Notes:         splitReq.Notes,
			ParentID:      &parentID,
			IsReconciled:  false,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		splitModels = append(splitModels, splitTxn)
	}

	// 批量保存拆分交易
	if err := s.txnRepo.CreateBatch(splitModels); err != nil {
		return nil, err
	}

	// 为拆分交易添加标签
	for i, splitReq := range req.Splits {
		if len(splitReq.Tags) > 0 {
			if err := s.txnRepo.AttachTags(splitModels[i].ID, splitReq.Tags); err != nil {
				// 记录错误但不回滚，标签不是关键数据
			}
		}
	}

	// 重新获取父交易（包含拆分）
	return s.Get(userID, parentID)
}

// GetSplits 获取拆分交易列表
// 参数：
//   - userID: 用户ID
//   - parentID: 父交易ID
// 返回：
//   - []*response.TransactionResp: 拆分交易列表
//   - error: 错误信息
func (s *TransactionService) GetSplits(userID uint64, parentID uint64) ([]*response.TransactionResp, error) {
	splits, err := s.txnRepo.GetSplits(parentID, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*response.TransactionResp, 0, len(splits))
	for _, split := range splits {
		result = append(result, s.toResp(&split))
	}

	return result, nil
}

// MergeSplits 合并拆分交易
// 删除所有子交易，保留父交易
// 业务流程：
// 1. 验证父交易存在且不是子交易
// 2. 获取所有拆分子交易
// 3. 批量删除所有子交易
// 参数：
//   - userID: 用户ID
//   - parentID: 父交易ID
// 返回：
//   - error: 错误信息
func (s *TransactionService) MergeSplits(userID, parentID uint64) error {
	// 验证父交易存在
	parent, err := s.txnRepo.GetByID(parentID, userID)
	if err != nil {
		return err
	}

	if parent.ParentID != nil {
		return errcode.WithMessage(errcode.ErrBadRequest, "cannot merge a child transaction")
	}

	// 获取所有拆分
	splits, err := s.txnRepo.GetSplits(parentID, userID)
	if err != nil {
		return err
	}

	if len(splits) == 0 {
		return errcode.WithMessage(errcode.ErrBadRequest, "transaction has no splits to merge")
	}

	// 删除所有拆分交易
	splitIDs := make([]uint64, 0, len(splits))
	for _, split := range splits {
		splitIDs = append(splitIDs, split.ID)
	}

	return s.txnRepo.DeleteBatch(splitIDs, userID)
}

// parseDateTime 解析日期时间字符串，支持两种格式：
// - "2006-01-02"（仅日期，时间为 00:00）
// - "2006-01-02 15:04"（日期+时分）
// - "2006-01-02T15:04"（ISO 格式）
// 参数 dateStr: 日期时间字符串
// 返回: 解析后的 time.Time（使用本地时区）和错误
func parseDateTime(dateStr string) (time.Time, error) {
	// 优先尝试带时分的格式
	for _, layout := range []string{
		"2006-01-02 15:04",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
	} {
		if t, err := time.ParseInLocation(layout, dateStr, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
}

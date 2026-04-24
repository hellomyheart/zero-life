// Package service 业务逻辑层，实现核心业务逻辑
// TransactionService 交易业务逻辑，处理交易的增删改查、拆分、合并和搜索

import (
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

type TransactionService struct {
	txnRepo        *repository.TransactionRepository
	accountRepo    *repository.AccountRepository
	db             *gorm.DB
	ruleTrigger    RuleTrigger
	webhookNotifier WebhookNotifier
}

func NewTransactionService(txnRepo *repository.TransactionRepository, accountRepo *repository.AccountRepository, db *gorm.DB, ruleTrigger RuleTrigger, webhookNotifier WebhookNotifier) *TransactionService {
	return &TransactionService{
		txnRepo:        txnRepo,
		accountRepo:    accountRepo,
		db:             db,
		ruleTrigger:    ruleTrigger,
		webhookNotifier: webhookNotifier,
	}
}

func (s *TransactionService) Create(userID uint64, req *request.CreateTransactionReq) (*response.TransactionResp, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, errcode.ErrInvalidAmount
	}

	date, err := time.Parse("2006-01-02", req.Date)
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

	// Execute in DB transaction
	err = s.db.Transaction(func(dbTx *gorm.DB) error {
		// Create transaction with tags
		if err := s.txnRepo.Create(txn, req.Tags); err != nil {
			return err
		}

		// Create splits if any
		if len(req.Splits) > 0 {
			splitTotal := decimal.Zero
			for _, splitReq := range req.Splits {
				splitAmount, _ := decimal.NewFromString(splitReq.Amount)
				splitTotal = splitTotal.Add(splitAmount)

				split := &model.Transaction{
					UserID:     userID,
					Type:        txnType,
					Date:        date,
					Description: req.Description,
					Amount:      splitAmount,
					SourceID:    req.SourceID,
					DestinationID: req.DestinationID,
					CategoryID:  splitReq.CategoryID,
					Notes:       splitReq.Notes,
					ParentID:    &txn.ID,
				}
				if err := s.txnRepo.Create(split, splitReq.Tags); err != nil {
					return err
				}
			}
			if !splitTotal.Equal(amount) {
				return errcode.ErrSplitAmountMismatch
			}
		}

		// Update account balances
		for accountID, change := range balanceChanges {
			if err := s.updateAccountBalance(dbTx, accountID, change); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		if e, ok := err.(*errcode.Error); ok {
			return nil, e
		}
		return nil, errcode.ErrInternal
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

	newDate, err := time.Parse("2006-01-02", req.Date)
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

	err = s.db.Transaction(func(dbTx *gorm.DB) error {
		if err := s.txnRepo.UpdateWithTags(oldTxn, req.Tags); err != nil {
			return err
		}
		for accountID, change := range netChanges {
			if err := s.updateAccountBalance(dbTx, accountID, change); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, errcode.ErrInternal
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

	err = s.db.Transaction(func(dbTx *gorm.DB) error {
		if err := s.txnRepo.Delete(id, userID); err != nil {
			return err
		}
		for accountID, change := range changes {
			if err := s.updateAccountBalance(dbTx, accountID, change.Neg()); err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

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

func (s *TransactionService) updateAccountBalance(dbTx *gorm.DB, accountID uint64, change decimal.Decimal) error {
	return dbTx.Model(&model.Account{}).Where("id = ?", accountID).
		Update("current_balance", gorm.Expr("current_balance + ?", change)).Error
}

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
// 将一笔交易拆分为多笔子交易，子交易金额之和必须等于父交易金额
func (s *TransactionService) Split(userID, parentID uint64, req *request.SplitTransactionReq) (*response.TransactionResp, error) {
	// 获取父交易	parent, err := s.txnRepo.GetByID(parentID, userID)
	if err != nil {
		return nil, err
	}

	// 妫€鏌ユ槸鍚﹀凡缁忔槸拆分交易
	if parent.ParentID != nil {
		return nil, errcode.WithMessage(errcode.ErrBadRequest, "cannot split a child transaction")
	}

	// 检查是否已有拆分	splits, _ := s.txnRepo.GetSplits(parentID, userID)
	if len(splits) > 0 {
		return nil, errcode.WithMessage(errcode.ErrBadRequest, "transaction already has splits, merge them first")
	}

	// 楠岃瘉鎷嗗垎閲戦鎬诲拰
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

	// 鍒涘缓拆分交易
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

	// 淇濆瓨拆分交易
	if err := s.txnRepo.CreateBatch(splitModels); err != nil {
		return nil, err
	}

	
		// 为拆分交易添加标签	for i, splitReq := range req.Splits {
		if len(splitReq.Tags) > 0 {
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
func (s *TransactionService) GetSplits(userID, parentID uint64) ([]*response.TransactionResp, error) {
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

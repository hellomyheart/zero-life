// Package service 业务逻辑层，实现核心业务逻辑
// TransactionBulkService 交易批量操作业务逻辑，支持批量编辑、删除、类型转换和克隆
package service

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

// TransactionBulkService 交易批量操作服务
// 负责处理交易的批量编辑、批量删除、类型转换和克隆操作
// 所有批量操作都在数据库事务中执行，确保数据一致性
// 依赖txnRepo进行交易数据访问，依赖accountRepo验证账户，依赖db执行事务
type TransactionBulkService struct {
	txnRepo     *repository.TransactionRepository // 交易数据访问对象
	accountRepo *repository.AccountRepository     // 账户数据访问对象
	db          *gorm.DB                          // 数据库连接，用于事务操作
}

// NewTransactionBulkService 创建交易批量操作服务实例
func NewTransactionBulkService(
	txnRepo *repository.TransactionRepository,
	accountRepo *repository.AccountRepository,
	db *gorm.DB,
) *TransactionBulkService {
	return &TransactionBulkService{
		txnRepo:     txnRepo,
		accountRepo: accountRepo,
		db:          db,
	}
}

// BulkEdit 批量编辑交易
// 在数据库事务中更新多笔交易的分类、备注和标签
// 不存在的交易ID会被跳过
// 参数：
//   - userID: 用户ID
//   - req: 批量编辑请求参数（交易ID列表、分类ID、备注、标签ID列表）
// 返回：
//   - error: 错误信息
func (s *TransactionBulkService) BulkEdit(userID uint64, req *request.BulkEditReq) error {
	for _, id := range req.IDs {
		txn, err := s.txnRepo.GetByID(id, userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			return err
		}

		updated := false
		if req.CategoryID != nil {
			txn.CategoryID = req.CategoryID
			updated = true
		}
		if req.Notes != "" {
			txn.Notes = req.Notes
			updated = true
		}

		if updated {
			if err := s.txnRepo.Update(txn); err != nil {
				return err
			}
		}

		if len(req.TagIDs) > 0 {
			if err := s.txnRepo.AddTags(txn.ID, req.TagIDs); err != nil {
				return err
			}
		}
	}
	return nil
}

// BulkDelete 批量删除交易
// 在数据库事务中逐笔删除交易并回滚对应的账户余额变更
// 参数：
//   - userID: 用户ID
//   - req: 批量删除请求参数（交易ID列表）
// 返回：
//   - error: 错误信息
func (s *TransactionBulkService) BulkDelete(userID uint64, req *request.BulkDeleteReq) error {
	type txnInfo struct {
		id      uint64
		changes map[uint64]decimal.Decimal
	}
	txns := make([]txnInfo, 0, len(req.IDs))

	for _, id := range req.IDs {
		txn, err := s.txnRepo.GetByID(id, userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			return err
		}
		changes := s.calculateBalanceChanges(txn.Type, txn.Amount, txn.SourceID, txn.DestinationID)
		txns = append(txns, txnInfo{id: id, changes: changes})
	}

	if len(txns) == 0 {
		return nil
	}

	ids := make([]uint64, len(txns))
	for i, t := range txns {
		ids[i] = t.id
	}

	return s.db.Transaction(func(dbTx *gorm.DB) error {
		if err := dbTx.Where("transaction_id IN ?", ids).Delete(&model.TransactionTag{}).Error; err != nil {
			return err
		}
		if err := dbTx.Unscoped().Where("parent_id IN ?", ids).Delete(&model.Transaction{}).Error; err != nil {
			return err
		}
		if err := dbTx.Unscoped().Where("id IN ? AND user_id = ?", ids, userID).Delete(&model.Transaction{}).Error; err != nil {
			return err
		}
		for _, t := range txns {
			for accountID, change := range t.changes {
				if err := dbTx.Model(&model.Account{}).Where("id = ?", accountID).
					Update("current_balance", gorm.Expr("current_balance - ?", change)).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// ConvertType 转换交易类型
// 支持以下转换：
//   - 支出 <-> 收入：可指定新的源账户
//   - 任意 -> 转账：必须提供源账户和目标账户
//   - 转账 -> 支出/收入：移除目标账户
// 转换后自动计算并更新账户余额的净变更
// 参数：
//   - userID: 用户ID
//   - id: 交易ID
//   - req: 类型转换请求参数（新类型、新源账户ID、新目标账户ID）
// 返回：
//   - *response.TransactionResp: 转换后的交易信息
//   - error: 错误信息
func (s *TransactionBulkService) ConvertType(userID, id uint64, req *request.ConvertReq) (*response.TransactionResp, error) {
	txn, err := s.txnRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	newType := model.TransactionType(req.Type)

	// 计算旧余额变更（转换前的余额影响）
	oldChanges := s.calculateBalanceChanges(txn.Type, txn.Amount, txn.SourceID, txn.DestinationID)

	// 根据转换类型更新交易字段
	switch {
	case txn.Type == model.TransactionTypeWithdrawal && newType == model.TransactionTypeDeposit:
		// 支出转收入：源账户变为目标账户，需要指定新的源账户（收入类账户）
		if req.SourceID != nil {
			txn.SourceID = *req.SourceID
		}
		txn.Type = newType

	case txn.Type == model.TransactionTypeDeposit && newType == model.TransactionTypeWithdrawal:
		// 收入转支出：源账户保持不变，仅更改类型
		if req.SourceID != nil {
			txn.SourceID = *req.SourceID
		}
		txn.Type = newType

	case newType == model.TransactionTypeTransfer:
		// 任意类型转转账：必须同时提供源账户和目标账户
		if req.SourceID == nil || req.DestinationID == nil {
			return nil, errcode.ErrInvalidTxnType
		}
		txn.SourceID = *req.SourceID
		txn.DestinationID = req.DestinationID
		txn.Type = newType

	case txn.Type == model.TransactionTypeTransfer && newType != model.TransactionTypeTransfer:
		// 转账转支出/收入：移除目标账户
		if req.SourceID != nil {
			txn.SourceID = *req.SourceID
		}
		txn.DestinationID = nil
		txn.Type = newType

	default:
		return nil, errcode.ErrInvalidTxnType
	}

	// 计算新余额变更（转换后的余额影响）
	newChanges := s.calculateBalanceChanges(txn.Type, txn.Amount, txn.SourceID, txn.DestinationID)

	// 计算净余额变更 = 新变更 - 旧变更
	// 对每个受影响的账户，净变更 = 新影响 - 旧影响
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

	if err := s.db.Transaction(func(dbTx *gorm.DB) error {
		if err := s.txnRepo.UpdateWithTagsAndDB(dbTx, txn, nil); err != nil {
			return err
		}
		for accountID, change := range netChanges {
			if err := dbTx.Model(&model.Account{}).Where("id = ?", accountID).
				Update("current_balance", gorm.Expr("current_balance + ?", change)).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, errcode.ErrInternal
	}

	updated, err := s.txnRepo.GetByID(id, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(updated), nil
}

// CloneTransaction 克隆交易
// 复制一笔交易，日期设为今天，同时复制标签和分类
// 克隆后自动更新相关账户余额
// 参数：
//   - userID: 用户ID
//   - id: 被克隆的交易ID
// 返回：
//   - *response.TransactionResp: 克隆后的新交易信息
//   - error: 错误信息
func (s *TransactionBulkService) CloneTransaction(userID, id uint64) (*response.TransactionResp, error) {
	txn, err := s.txnRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	// Get tag IDs
	tagIDs := make([]uint64, 0, len(txn.Tags))
	for _, tag := range txn.Tags {
		tagIDs = append(tagIDs, tag.ID)
	}

	clone := &model.Transaction{
		UserID:        userID,
		Type:          txn.Type,
		Date:          time.Now(),
		Description:   txn.Description,
		Amount:        txn.Amount,
		SourceID:      txn.SourceID,
		DestinationID: txn.DestinationID,
		CategoryID:    txn.CategoryID,
		Notes:         txn.Notes,
	}

	// Calculate balance changes
	changes := s.calculateBalanceChanges(clone.Type, clone.Amount, clone.SourceID, clone.DestinationID)

	if err := s.db.Transaction(func(dbTx *gorm.DB) error {
		if err := s.txnRepo.CreateWithDB(dbTx, clone, tagIDs); err != nil {
			return err
		}
		for accountID, change := range changes {
			if err := dbTx.Model(&model.Account{}).Where("id = ?", accountID).
				Update("current_balance", gorm.Expr("current_balance + ?", change)).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, errcode.ErrInternal
	}

	created, err := s.txnRepo.GetByID(clone.ID, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(created), nil
}

// calculateBalanceChanges 根据交易类型计算各账户的余额变更
// deposit：目标账户余额增加；withdrawal：源账户余额减少；transfer：源账户减少，目标账户增加
func (s *TransactionBulkService) calculateBalanceChanges(txnType model.TransactionType, amount decimal.Decimal, sourceID uint64, destID *uint64) map[uint64]decimal.Decimal {
	changes := make(map[uint64]decimal.Decimal)

	switch txnType {
	case model.TransactionTypeDeposit:
		if destID != nil {
			changes[*destID] = amount
		}
	case model.TransactionTypeWithdrawal:
		changes[sourceID] = amount.Neg()
	case model.TransactionTypeTransfer:
		changes[sourceID] = amount.Neg()
		if destID != nil {
			changes[*destID] = amount
		}
	}

	return changes
}

// toResp 将交易模型转换为响应对象
func (s *TransactionBulkService) toResp(t *model.Transaction) *response.TransactionResp {
	resp := &response.TransactionResp{
		ID:            t.ID,
		Type:          string(t.Type),
		Date:          t.Date,
		Description:   t.Description,
		Amount:        t.Amount.StringFixed(4),
		SourceID:      t.SourceID,
		DestinationID: t.DestinationID,
		CategoryID:    t.CategoryID,
		Notes:         t.Notes,
		BillID:        t.BillID,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}

	if t.Source.ID > 0 {
		resp.Source = *accountModelToResp(&t.Source)
	}
	if t.Destination != nil && t.Destination.ID > 0 {
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

	return resp
}

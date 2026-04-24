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

type TransactionBulkService struct {
	txnRepo     *repository.TransactionRepository
	accountRepo *repository.AccountRepository
	db          *gorm.DB
}

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

// BulkEdit updates category, notes, and/or tags for multiple transactions within a DB transaction.
func (s *TransactionBulkService) BulkEdit(userID uint64, req *request.BulkEditReq) error {
	return s.db.Transaction(func(dbTx *gorm.DB) error {
		for _, id := range req.IDs {
			txn, err := s.txnRepo.GetByID(id, userID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					continue // skip not found
				}
				return err
			}

			if req.CategoryID != nil {
				txn.CategoryID = req.CategoryID
			}
			if req.Notes != "" {
				txn.Notes = req.Notes
			}

			if err := s.txnRepo.UpdateWithTags(txn, req.TagIDs); err != nil {
				return err
			}
		}
		return nil
	})
}

// BulkDelete deletes multiple transactions within a DB transaction, rolling back balances for each.
func (s *TransactionBulkService) BulkDelete(userID uint64, req *request.BulkDeleteReq) error {
	return s.db.Transaction(func(dbTx *gorm.DB) error {
		for _, id := range req.IDs {
			txn, err := s.txnRepo.GetByID(id, userID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					continue
				}
				return err
			}

			// Rollback balance changes
			changes := s.calculateBalanceChanges(txn.Type, txn.Amount, txn.SourceID, txn.DestinationID)
			for accountID, change := range changes {
				if err := dbTx.Model(&model.Account{}).Where("id = ?", accountID).
					Update("current_balance", gorm.Expr("current_balance - ?", change)).Error; err != nil {
					return err
				}
			}

			if err := s.txnRepo.Delete(id, userID); err != nil {
				return err
			}
		}
		return nil
	})
}

// ConvertType changes a transaction's type and swaps source/destination as needed.
// For withdrawal <-> deposit: swap SourceID and DestinationID.
// For transfer: the new SourceID/DestinationID must be provided.
func (s *TransactionBulkService) ConvertType(userID, id uint64, req *request.ConvertReq) (*response.TransactionResp, error) {
	txn, err := s.txnRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	newType := model.TransactionType(req.Type)

	// Calculate old balance changes
	oldChanges := s.calculateBalanceChanges(txn.Type, txn.Amount, txn.SourceID, txn.DestinationID)

	// Apply type conversion
	switch {
	case txn.Type == model.TransactionTypeWithdrawal && newType == model.TransactionTypeDeposit:
		// withdrawal -> deposit: source becomes destination, need a new source (revenue account)
		if req.SourceID != nil {
			txn.SourceID = *req.SourceID
		}
		txn.Type = newType

	case txn.Type == model.TransactionTypeDeposit && newType == model.TransactionTypeWithdrawal:
		// deposit -> withdrawal: source stays, just change type
		if req.SourceID != nil {
			txn.SourceID = *req.SourceID
		}
		txn.Type = newType

	case newType == model.TransactionTypeTransfer:
		// any -> transfer: need both source and destination
		if req.SourceID == nil || req.DestinationID == nil {
			return nil, errcode.ErrInvalidTxnType
		}
		txn.SourceID = *req.SourceID
		txn.DestinationID = req.DestinationID
		txn.Type = newType

	case txn.Type == model.TransactionTypeTransfer && newType != model.TransactionTypeTransfer:
		// transfer -> deposit/withdrawal
		if req.SourceID != nil {
			txn.SourceID = *req.SourceID
		}
		txn.DestinationID = nil
		txn.Type = newType

	default:
		return nil, errcode.ErrInvalidTxnType
	}

	// Calculate new balance changes
	newChanges := s.calculateBalanceChanges(txn.Type, txn.Amount, txn.SourceID, txn.DestinationID)

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

	err = s.db.Transaction(func(dbTx *gorm.DB) error {
		if err := s.txnRepo.Update(txn); err != nil {
			return err
		}
		for accountID, change := range netChanges {
			if err := dbTx.Model(&model.Account{}).Where("id = ?", accountID).
				Update("current_balance", gorm.Expr("current_balance + ?", change)).Error; err != nil {
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

	return s.toResp(updated), nil
}

// CloneTransaction duplicates a transaction with today's date.
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

	err = s.db.Transaction(func(dbTx *gorm.DB) error {
		if err := s.txnRepo.Create(clone, tagIDs); err != nil {
			return err
		}
		for accountID, change := range changes {
			if err := dbTx.Model(&model.Account{}).Where("id = ?", accountID).
				Update("current_balance", gorm.Expr("current_balance + ?", change)).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, errcode.ErrInternal
	}

	created, err := s.txnRepo.GetByID(clone.ID, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(created), nil
}

func (s *TransactionBulkService) calculateBalanceChanges(txnType model.TransactionType, amount decimal.Decimal, sourceID uint64, destID *uint64) map[uint64]decimal.Decimal {
	changes := make(map[uint64]decimal.Decimal)

	switch txnType {
	case model.TransactionTypeDeposit:
		changes[sourceID] = amount
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

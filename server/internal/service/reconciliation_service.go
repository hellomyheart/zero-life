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

type ReconciliationService struct {
	reconRepo   *repository.ReconciliationRepository
	accountRepo *repository.AccountRepository
	txnRepo     *repository.TransactionRepository
	db          *gorm.DB
}

func NewReconciliationService(
	reconRepo *repository.ReconciliationRepository,
	accountRepo *repository.AccountRepository,
	txnRepo *repository.TransactionRepository,
	db *gorm.DB,
) *ReconciliationService {
	return &ReconciliationService{
		reconRepo:   reconRepo,
		accountRepo: accountRepo,
		txnRepo:     txnRepo,
		db:          db,
	}
}

// GetReconciliationData calculates start/end balances and returns transactions in the date range.
func (s *ReconciliationService) GetReconciliationData(userID, accountID uint64, req *request.GetReconciliationReq) (*response.ReconciliationDataResp, error) {
	// Verify account belongs to user
	_, err := s.accountRepo.GetByID(accountID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errcode.ErrBadRequest
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, errcode.ErrBadRequest
	}

	// Calculate start balance: account initial balance + sum of transactions before start_date
	startBalance, err := s.calculateBalanceAtDate(userID, accountID, startDate)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	// Calculate end balance: account initial balance + sum of transactions up to end_date
	endBalance, err := s.calculateBalanceAtDate(userID, accountID, endDate.AddDate(0, 0, 1))
	if err != nil {
		return nil, errcode.ErrInternal
	}

	// Get transactions in the date range
	filter := repository.TransactionFilter{
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		AccountID:  &accountID,
	}
	txns, err := s.txnRepo.List(userID, filter, 0, 1000)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	txnResps := make([]response.TransactionResp, 0, len(txns))
	for _, t := range txns {
		txnResps = append(txnResps, *s.txnToResp(&t))
	}

	return &response.ReconciliationDataResp{
		AccountID:    accountID,
		StartDate:    startDate,
		EndDate:      endDate,
		StartBalance: startBalance.StringFixed(4),
		EndBalance:   endBalance.StringFixed(4),
		Transactions: txnResps,
	}, nil
}

// SubmitReconciliation calculates the difference and creates a reconciliation record.
// If difference != 0, it creates an adjustment transaction.
// It also marks all transactions in the range as reconciled.
func (s *ReconciliationService) SubmitReconciliation(userID, accountID uint64, req *request.SubmitReconciliationReq) (*response.ReconciliationResp, error) {
	// Verify account belongs to user
	_, err := s.accountRepo.GetByID(accountID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errcode.ErrBadRequest
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, errcode.ErrBadRequest
	}

	submittedBalance, err := decimal.NewFromString(req.SubmittedBalance)
	if err != nil {
		return nil, errcode.ErrBadRequest
	}

	// Calculate start and end balances
	startBalance, err := s.calculateBalanceAtDate(userID, accountID, startDate)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	endBalance, err := s.calculateBalanceAtDate(userID, accountID, endDate.AddDate(0, 0, 1))
	if err != nil {
		return nil, errcode.ErrInternal
	}

	difference := submittedBalance.Sub(endBalance)

	// Execute in DB transaction
	var recon model.Reconciliation
	err = s.db.Transaction(func(dbTx *gorm.DB) error {
		// If difference != 0, create an adjustment transaction
		if !difference.IsZero() {
			adjustTxn := &model.Transaction{
				UserID:      userID,
				Type:        model.TransactionTypeDeposit,
				Date:        endDate,
				Description: "Reconciliation adjustment",
				Amount:      difference.Abs(),
				SourceID:    accountID,
				Notes:       "Auto-generated reconciliation adjustment",
				IsReconciled: true,
			}

			if difference.IsNegative() {
				adjustTxn.Type = model.TransactionTypeWithdrawal
			}

			if err := dbTx.Create(adjustTxn).Error; err != nil {
				return err
			}

			// Update account balance
			if difference.IsPositive() {
				if err := dbTx.Model(&model.Account{}).Where("id = ?", accountID).
					Update("current_balance", gorm.Expr("current_balance + ?", difference)).Error; err != nil {
					return err
				}
			} else {
				if err := dbTx.Model(&model.Account{}).Where("id = ?", accountID).
					Update("current_balance", gorm.Expr("current_balance - ?", difference.Abs())).Error; err != nil {
					return err
				}
			}
		}

		// Mark transactions in the range as reconciled
		if err := dbTx.Model(&model.Transaction{}).
			Where("user_id = ? AND (source_id = ? OR destination_id = ?) AND date >= ? AND date < ?",
				userID, accountID, accountID, startDate, endDate.AddDate(0, 0, 1)).
			Update("is_reconciled", true).Error; err != nil {
			return err
		}

		// Create reconciliation record
		recon = model.Reconciliation{
			UserID:           userID,
			AccountID:        accountID,
			StartDate:        startDate,
			EndDate:          endDate,
			StartBalance:     startBalance,
			EndBalance:       endBalance,
			SubmittedBalance: submittedBalance,
			Difference:       difference,
			CreatedAt:        time.Now(),
		}

		return dbTx.Create(&recon).Error
	})

	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(&recon), nil
}

// calculateBalanceAtDate returns the account balance at the start of the given date.
// It sums all transactions affecting the account before the given date.
func (s *ReconciliationService) calculateBalanceAtDate(userID, accountID uint64, date time.Time) (decimal.Decimal, error) {
	account, err := s.accountRepo.GetByID(accountID, userID)
	if err != nil {
		return decimal.Zero, err
	}

	balance := account.InitialBalance

	// Sum deposits to this account
	var depositSum decimal.Decimal
	if err := s.db.Model(&model.Transaction{}).
		Where("user_id = ? AND source_id = ? AND type = ? AND date < ?", userID, accountID, model.TransactionTypeDeposit, date).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&depositSum); err != nil {
		return decimal.Zero, err
	}

	// Sum withdrawals from this account
	var withdrawalSum decimal.Decimal
	if err := s.db.Model(&model.Transaction{}).
		Where("user_id = ? AND source_id = ? AND type = ? AND date < ?", userID, accountID, model.TransactionTypeWithdrawal, date).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&withdrawalSum); err != nil {
		return decimal.Zero, err
	}

	// Sum transfers to this account (destination)
	var transferInSum decimal.Decimal
	if err := s.db.Model(&model.Transaction{}).
		Where("user_id = ? AND destination_id = ? AND type = ? AND date < ?", userID, accountID, model.TransactionTypeTransfer, date).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&transferInSum); err != nil {
		return decimal.Zero, err
	}

	// Sum transfers from this account (source)
	var transferOutSum decimal.Decimal
	if err := s.db.Model(&model.Transaction{}).
		Where("user_id = ? AND source_id = ? AND type = ? AND date < ?", userID, accountID, model.TransactionTypeTransfer, date).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&transferOutSum); err != nil {
		return decimal.Zero, err
	}

	balance = balance.Add(depositSum).Sub(withdrawalSum).Add(transferInSum).Sub(transferOutSum)
	return balance, nil
}

func (s *ReconciliationService) toResp(r *model.Reconciliation) *response.ReconciliationResp {
	return &response.ReconciliationResp{
		ID:               r.ID,
		AccountID:        r.AccountID,
		StartDate:        r.StartDate,
		EndDate:          r.EndDate,
		StartBalance:     r.StartBalance.StringFixed(4),
		EndBalance:       r.EndBalance.StringFixed(4),
		SubmittedBalance: r.SubmittedBalance.StringFixed(4),
		Difference:       r.Difference.StringFixed(4),
		CreatedAt:        r.CreatedAt,
	}
}

// txnToResp converts a transaction model to response (simplified for reconciliation context).
func (s *ReconciliationService) txnToResp(t *model.Transaction) *response.TransactionResp {
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

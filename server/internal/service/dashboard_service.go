// Package service 业务逻辑层，实现核心业务逻辑
// DashboardService 仪表盘业务逻辑，汇总首页展示的收支、预算、账单等
package service

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/repository"
)

type DashboardService struct {
	txnRepo    *repository.TransactionRepository
	accountRepo *repository.AccountRepository
	budgetRepo *repository.BudgetRepository
	billRepo   *repository.BillRepository
}

func NewDashboardService(
	txnRepo *repository.TransactionRepository,
	accountRepo *repository.AccountRepository,
	budgetRepo *repository.BudgetRepository,
	billRepo *repository.BillRepository,
) *DashboardService {
	return &DashboardService{
		txnRepo:    txnRepo,
		accountRepo: accountRepo,
		budgetRepo: budgetRepo,
		billRepo:   billRepo,
	}
}

func (s *DashboardService) Get(userID uint64) (*response.DashboardResp, error) {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Monthly income & expense
	monthFilter := repository.TransactionFilter{
		StartDate: monthStart.Format("2006-01-02"),
		EndDate:   now.Format("2006-01-02"),
	}
	monthTxns, err := s.txnRepo.List(userID, monthFilter, 0, 10000)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	monthIncome := decimal.Zero
	monthExpense := decimal.Zero
	for _, txn := range monthTxns {
		switch txn.Type {
		case model.TransactionTypeDeposit:
			monthIncome = monthIncome.Add(txn.Amount)
		case model.TransactionTypeWithdrawal:
			monthExpense = monthExpense.Add(txn.Amount)
		}
	}

	// Total balance (all asset accounts)
	assetAccounts, err := s.accountRepo.List(userID, string(model.AccountTypeAsset), "", "name", 0, 1000)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	totalBalance := decimal.Zero
	for _, a := range assetAccounts {
		totalBalance = totalBalance.Add(a.CurrentBalance)
	}

	// Budget alerts
	budgets, err := s.budgetRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	budgetAlerts := make([]response.BudgetAlertResp, 0)
	for _, b := range budgets {
		if !b.IsEnabled {
			continue
		}
		// Calculate spent for this budget's categories
		spent := decimal.Zero
		for _, cat := range b.Categories {
			catFilter := repository.TransactionFilter{
				StartDate:  monthStart.Format("2006-01-02"),
				EndDate:    now.Format("2006-01-02"),
				CategoryID: &cat.ID,
			}
			catTxns, err := s.txnRepo.List(userID, catFilter, 0, 10000)
			if err != nil {
				continue
			}
			for _, txn := range catTxns {
				spent = spent.Add(txn.Amount)
			}
		}

		usageRate, _ := spent.Div(b.Amount).Float64()
		status := "normal"
		if usageRate >= 1.0 {
			status = "overspent"
		} else if usageRate >= 0.8 {
			status = "warning"
		}

		if status != "normal" {
			budgetAlerts = append(budgetAlerts, response.BudgetAlertResp{
				BudgetID:   b.ID,
				BudgetName: b.Name,
				Spent:      spent.StringFixed(4),
				Amount:     b.Amount.StringFixed(4),
				UsageRate:  usageRate,
				Status:     status,
			})
		}
	}

	// Bill reminders (upcoming 7 days)
	bills, err := s.billRepo.GetUpcoming(userID, 7)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	billReminders := make([]response.BillReminderResp, 0, len(bills))
	for _, b := range bills {
		billReminders = append(billReminders, response.BillReminderResp{
			BillID:   b.ID,
			BillName: b.Name,
			Amount:   b.Amount.StringFixed(4),
			NextDue:  b.NextDue.Format("2006-01-02"),
		})
	}

	// Recent transactions (last 5)
	recentTxns, err := s.txnRepo.List(userID, repository.TransactionFilter{}, 0, 5)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	recentTxnResps := make([]response.TransactionResp, 0, len(recentTxns))
	for _, t := range recentTxns {
		recentTxnResps = append(recentTxnResps, *transactionModelToResp(&t))
	}

	return &response.DashboardResp{
		MonthIncome:   monthIncome.StringFixed(4),
		MonthExpense:  monthExpense.StringFixed(4),
		NetIncome:     monthIncome.Sub(monthExpense).StringFixed(4),
		TotalBalance:  totalBalance.StringFixed(4),
		BudgetAlerts:  budgetAlerts,
		BillReminders: billReminders,
		RecentTxns:    recentTxnResps,
	}, nil
}

func transactionModelToResp(t *model.Transaction) *response.TransactionResp {
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

	if t.Destination != nil {
		resp.Destination = &response.AccountResp{
			ID:             t.Destination.ID,
			Name:           t.Destination.Name,
			Type:           string(t.Destination.Type),
			CurrencyID:     t.Destination.CurrencyID,
			InitialBalance: t.Destination.InitialBalance.StringFixed(4),
			CurrentBalance: t.Destination.CurrentBalance.StringFixed(4),
			IsVirtual:      t.Destination.IsVirtual,
			Notes:          t.Destination.Notes,
			CreatedAt:      t.Destination.CreatedAt,
			UpdatedAt:      t.Destination.UpdatedAt,
		}
	}

	if t.Category != nil {
		resp.Category = &response.CategoryResp{
			ID:        t.Category.ID,
			Name:      t.Category.Name,
			ParentID:  t.Category.ParentID,
			Icon:      t.Category.Icon,
			Notes:     t.Category.Notes,
			SortOrder: t.Category.SortOrder,
			CreatedAt: t.Category.CreatedAt,
			UpdatedAt: t.Category.UpdatedAt,
		}
	}

	tags := make([]response.TagResp, 0, len(t.Tags))
	for _, tag := range t.Tags {
		tags = append(tags, response.TagResp{
			ID:        tag.ID,
			Name:      tag.Name,
			Color:     tag.Color,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		})
	}
	resp.Tags = tags

	return resp
}

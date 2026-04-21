package service

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/zero-life/server/internal/dto/request"
	"github.com/zero-life/server/internal/dto/response"
	"github.com/zero-life/server/internal/model"
	"github.com/zero-life/server/internal/pkg/errcode"
	"github.com/zero-life/server/internal/repository"
)

type ReportService struct {
	txnRepo     *repository.TransactionRepository
	accountRepo *repository.AccountRepository
	budgetRepo  *repository.BudgetRepository
	categoryRepo *repository.CategoryRepository
}

func NewReportService(
	txnRepo *repository.TransactionRepository,
	accountRepo *repository.AccountRepository,
	budgetRepo *repository.BudgetRepository,
	categoryRepo *repository.CategoryRepository,
) *ReportService {
	return &ReportService{
		txnRepo:     txnRepo,
		accountRepo: accountRepo,
		budgetRepo:  budgetRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *ReportService) IncomeExpense(userID uint64, req *request.ReportReq) (*response.IncomeExpenseResp, error) {
	startDate, endDate, err := s.parseDateRange(req)
	if err != nil {
		return nil, err
	}

	filter := repository.TransactionFilter{
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),
	}

	txns, err := s.txnRepo.List(userID, filter, 0, 10000)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	totalIncome := decimal.Zero
	totalExpense := decimal.Zero

	for _, txn := range txns {
		switch txn.Type {
		case model.TransactionTypeDeposit:
			totalIncome = totalIncome.Add(txn.Amount)
		case model.TransactionTypeWithdrawal:
			totalExpense = totalExpense.Add(txn.Amount)
		}
	}

	return &response.IncomeExpenseResp{
		TotalIncome:  totalIncome.StringFixed(4),
		TotalExpense: totalExpense.StringFixed(4),
		NetIncome:    totalIncome.Sub(totalExpense).StringFixed(4),
	}, nil
}

func (s *ReportService) Category(userID uint64, req *request.ReportReq) (*response.CategoryReportResp, error) {
	startDate, endDate, err := s.parseDateRange(req)
	if err != nil {
		return nil, err
	}

	filter := repository.TransactionFilter{
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),
	}

	txns, err := s.txnRepo.List(userID, filter, 0, 10000)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	expenseByCategory := make(map[uint64]decimal.Decimal)
	incomeByCategory := make(map[uint64]decimal.Decimal)
	categoryNames := make(map[uint64]string)

	for _, txn := range txns {
		if txn.CategoryID != nil {
			catID := *txn.CategoryID
			if txn.Category != nil {
				categoryNames[catID] = txn.Category.Name
			}
			switch txn.Type {
			case model.TransactionTypeWithdrawal:
				expenseByCategory[catID] = expenseByCategory[catID].Add(txn.Amount)
			case model.TransactionTypeDeposit:
				incomeByCategory[catID] = incomeByCategory[catID].Add(txn.Amount)
			}
		}
	}

	totalExpense := decimal.Zero
	for _, v := range expenseByCategory {
		totalExpense = totalExpense.Add(v)
	}
	totalIncome := decimal.Zero
	for _, v := range incomeByCategory {
		totalIncome = totalIncome.Add(v)
	}

	expenseItems := make([]response.CategoryItemResp, 0)
	for catID, amount := range expenseByCategory {
		pct, _ := amount.Div(totalExpense).Mul(decimal.NewFromInt(100)).Float64()
		expenseItems = append(expenseItems, response.CategoryItemResp{
			CategoryID:   catID,
			CategoryName: categoryNames[catID],
			Amount:       amount.StringFixed(4),
			Percentage:   pct,
		})
	}

	incomeItems := make([]response.CategoryItemResp, 0)
	for catID, amount := range incomeByCategory {
		pct, _ := amount.Div(totalIncome).Mul(decimal.NewFromInt(100)).Float64()
		incomeItems = append(incomeItems, response.CategoryItemResp{
			CategoryID:   catID,
			CategoryName: categoryNames[catID],
			Amount:       amount.StringFixed(4),
			Percentage:   pct,
		})
	}

	return &response.CategoryReportResp{
		ExpenseByCategory: expenseItems,
		IncomeByCategory:  incomeItems,
	}, nil
}

func (s *ReportService) Budget(userID uint64, req *request.ReportReq) (*response.BudgetReportResp, error) {
	budgets, err := s.budgetRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.BudgetReportItemResp, 0, len(budgets))
	for _, b := range budgets {
		if !b.IsEnabled {
			continue
		}
		// Simplified: use budget amount as is
		usageRate, _ := decimal.Zero.Div(b.Amount).Float64()
		items = append(items, response.BudgetReportItemResp{
			BudgetID:   b.ID,
			BudgetName: b.Name,
			Amount:     b.Amount.StringFixed(4),
			Spent:      "0.0000",
			Remaining:  b.Amount.StringFixed(4),
			UsageRate:  usageRate,
		})
	}

	return &response.BudgetReportResp{Items: items}, nil
}

func (s *ReportService) NetWorth(userID uint64, req *request.ReportReq) (*response.NetWorthResp, error) {
	// Get all asset accounts
	accounts, err := s.accountRepo.List(userID, string(model.AccountTypeAsset), "", "name", 0, 1000)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	// Get all liability accounts
	liabilities, err := s.accountRepo.List(userID, string(model.AccountTypeLiability), "", "name", 0, 1000)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	totalAssets := decimal.Zero
	for _, a := range accounts {
		totalAssets = totalAssets.Add(a.CurrentBalance)
	}

	totalLiabilities := decimal.Zero
	for _, a := range liabilities {
		totalLiabilities = totalLiabilities.Add(a.CurrentBalance)
	}

	netWorth := totalAssets.Sub(totalLiabilities)

	return &response.NetWorthResp{
		TotalNetWorth: netWorth.StringFixed(4),
		Trend:         []response.NetWorthPointResp{},
	}, nil
}

func (s *ReportService) Trend(userID uint64, req *request.ReportReq) (*response.TrendResp, error) {
	startDate, endDate, err := s.parseDateRange(req)
	if err != nil {
		return nil, err
	}

	granularity := req.Granularity
	if granularity == "" {
		granularity = "month"
	}

	var periods []string
	var periodStarts []time.Time
	var periodEnds []time.Time

	switch granularity {
	case "month":
		current := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
		for current.Before(endDate) || current.Equal(endDate) {
			end := current.AddDate(0, 1, 0).Add(-time.Second)
			periods = append(periods, current.Format("2006-01"))
			periodStarts = append(periodStarts, current)
			periodEnds = append(periodEnds, end)
			current = current.AddDate(0, 1, 0)
		}
	default:
		periods = append(periods, startDate.Format("2006-01-02")+" ~ "+endDate.Format("2006-01-02"))
		periodStarts = append(periodStarts, startDate)
		periodEnds = append(periodEnds, endDate)
	}

	items := make([]response.TrendItemResp, 0, len(periods))
	for i, period := range periods {
		filter := repository.TransactionFilter{
			StartDate: periodStarts[i].Format("2006-01-02"),
			EndDate:   periodEnds[i].Format("2006-01-02"),
		}

		txns, err := s.txnRepo.List(userID, filter, 0, 10000)
		if err != nil {
			continue
		}

		income := decimal.Zero
		expense := decimal.Zero
		for _, txn := range txns {
			switch txn.Type {
			case model.TransactionTypeDeposit:
				income = income.Add(txn.Amount)
			case model.TransactionTypeWithdrawal:
				expense = expense.Add(txn.Amount)
			}
		}

		items = append(items, response.TrendItemResp{
			Period:  period,
			Income:  income.StringFixed(4),
			Expense: expense.StringFixed(4),
		})
	}

	return &response.TrendResp{Items: items}, nil
}

func (s *ReportService) parseDateRange(req *request.ReportReq) (time.Time, time.Time, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start_date: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end_date: %w", err)
	}
	return startDate, endDate, nil
}

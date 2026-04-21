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

type BudgetService struct {
	budgetRepo *repository.BudgetRepository
	txnRepo    *repository.TransactionRepository
}

func NewBudgetService(budgetRepo *repository.BudgetRepository, txnRepo *repository.TransactionRepository) *BudgetService {
	return &BudgetService{
		budgetRepo: budgetRepo,
		txnRepo:    txnRepo,
	}
}

func (s *BudgetService) Create(userID uint64, req *request.CreateBudgetReq) (*response.BudgetResp, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, errcode.ErrBudgetAmountInvalid
	}

	budget := &model.Budget{
		UserID:    userID,
		Name:      req.Name,
		Amount:    amount,
		Period:    model.BudgetPeriod(req.Period),
		IsEnabled: true,
	}

	if err := s.budgetRepo.Create(budget, req.CategoryIDs); err != nil {
		return nil, errcode.ErrInternal
	}

	created, err := s.budgetRepo.GetByID(budget.ID, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toRespWithUsage(created, userID)
}

func (s *BudgetService) Get(userID, id uint64) (*response.BudgetResp, error) {
	budget, err := s.budgetRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	return s.toRespWithUsage(budget, userID)
}

func (s *BudgetService) List(userID uint64) ([]response.BudgetResp, error) {
	budgets, err := s.budgetRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.BudgetResp, 0, len(budgets))
	for _, b := range budgets {
		resp, err := s.toRespWithUsage(&b, userID)
		if err != nil {
			return nil, err
		}
		items = append(items, *resp)
	}

	return items, nil
}

func (s *BudgetService) Update(userID, id uint64, req *request.UpdateBudgetReq) (*response.BudgetResp, error) {
	budget, err := s.budgetRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Name != "" {
		budget.Name = req.Name
	}
	if req.Amount != "" {
		amount, err := decimal.NewFromString(req.Amount)
		if err != nil || amount.LessThanOrEqual(decimal.Zero) {
			return nil, errcode.ErrBudgetAmountInvalid
		}
		budget.Amount = amount
	}
	if req.Period != "" {
		budget.Period = model.BudgetPeriod(req.Period)
	}
	if req.IsEnabled != nil {
		budget.IsEnabled = *req.IsEnabled
	}

	categoryIDs := req.CategoryIDs
	// If categoryIDs not provided, keep existing
	if categoryIDs == nil {
		existingIDs := make([]uint64, 0, len(budget.Categories))
		for _, c := range budget.Categories {
			existingIDs = append(existingIDs, c.ID)
		}
		categoryIDs = existingIDs
	}

	if err := s.budgetRepo.Update(budget, categoryIDs); err != nil {
		return nil, errcode.ErrInternal
	}

	updated, err := s.budgetRepo.GetByID(id, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toRespWithUsage(updated, userID)
}

func (s *BudgetService) Delete(userID, id uint64) error {
	_, err := s.budgetRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}

	return s.budgetRepo.Delete(id, userID)
}

func (s *BudgetService) GetHistory(userID, id uint64) ([]response.BudgetHistoryResp, error) {
	_, err := s.budgetRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	history, err := s.budgetRepo.GetHistory(id)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.BudgetHistoryResp, 0, len(history))
	for _, h := range history {
		usageRate, _ := h.Spent.Div(h.Amount).Float64()
		items = append(items, response.BudgetHistoryResp{
			ID:          h.ID,
			PeriodStart: h.PeriodStart,
			PeriodEnd:   h.PeriodEnd,
			Amount:      h.Amount.StringFixed(4),
			Spent:       h.Spent.StringFixed(4),
			UsageRate:   usageRate,
			CreatedAt:   h.CreatedAt,
		})
	}

	return items, nil
}

func (s *BudgetService) calculateSpent(budget *model.Budget, userID uint64) decimal.Decimal {
	now := time.Now()
	var start, end time.Time

	if budget.Period == model.BudgetPeriodMonthly {
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, 0).Add(-time.Second)
	} else {
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(1, 0, 0).Add(-time.Second)
	}

	// Sum transactions for the budget's categories in the period
	var total decimal.Decimal
	for _, cat := range budget.Categories {
		filter := repository.TransactionFilter{
			StartDate:  start.Format("2006-01-02"),
			EndDate:    end.Format("2006-01-02"),
			CategoryID: &cat.ID,
		}
		txns, err := s.txnRepo.List(userID, filter, 0, 10000)
		if err != nil {
			continue
		}
		for _, txn := range txns {
			total = total.Add(txn.Amount)
		}
	}

	return total
}

func (s *BudgetService) toRespWithUsage(budget *model.Budget, userID uint64) (*response.BudgetResp, error) {
	spent := s.calculateSpent(budget, userID)
	remaining := budget.Amount.Sub(spent)
	usageRate, _ := spent.Div(budget.Amount).Float64()

	status := "normal"
	if usageRate >= 1.0 {
		status = "overspent"
	} else if usageRate >= 0.8 {
		status = "warning"
	}

	categories := make([]response.CategoryResp, 0, len(budget.Categories))
	for _, c := range budget.Categories {
		categories = append(categories, *categoryModelToResp(&c))
	}

	return &response.BudgetResp{
		ID:         budget.ID,
		Name:       budget.Name,
		Amount:     budget.Amount.StringFixed(4),
		Period:     string(budget.Period),
		IsEnabled:  budget.IsEnabled,
		Categories: categories,
		Spent:      spent.StringFixed(4),
		Remaining:  remaining.StringFixed(4),
		UsageRate:  usageRate,
		Status:     status,
		CreatedAt:  budget.CreatedAt,
		UpdatedAt:  budget.UpdatedAt,
	}, nil
}

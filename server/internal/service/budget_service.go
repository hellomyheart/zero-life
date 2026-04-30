package service

import (
	"fmt"
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
	budgetRepo   *repository.BudgetRepository
	txnRepo      *repository.TransactionRepository
	categoryRepo *repository.CategoryRepository
}

func NewBudgetService(budgetRepo *repository.BudgetRepository, txnRepo *repository.TransactionRepository, categoryRepo *repository.CategoryRepository) *BudgetService {
	return &BudgetService{
		budgetRepo:   budgetRepo,
		txnRepo:      txnRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *BudgetService) validateCategoryIDs(categoryIDs []uint64, userID uint64) error {
	if len(categoryIDs) == 0 {
		return errcode.ErrBudgetCategoryRequired
	}
	count, err := s.categoryRepo.CountByIDsAndUserID(categoryIDs, userID)
	if err != nil {
		return errcode.ErrInternal
	}
	if int(count) != len(categoryIDs) {
		return errcode.ErrBudgetCategoryInvalid
	}
	return nil
}

func (s *BudgetService) Create(userID uint64, req *request.CreateBudgetReq) (*response.BudgetResp, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, errcode.ErrBudgetAmountInvalid
	}

	if err := s.validateCategoryIDs(req.CategoryIDs, userID); err != nil {
		return nil, err
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

	created, err := s.budgetRepo.GetByIDWithDB(s.budgetRepo.WriteDB(), budget.ID, userID)
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
	if categoryIDs == nil {
		existingIDs := make([]uint64, 0, len(budget.Categories))
		for _, c := range budget.Categories {
			existingIDs = append(existingIDs, c.ID)
		}
		categoryIDs = existingIDs
	} else {
		if len(categoryIDs) == 0 {
			return nil, errcode.ErrBudgetCategoryRequired
		}
		if err := s.validateCategoryIDs(categoryIDs, userID); err != nil {
			return nil, err
		}
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
		var usageRate float64
		if !h.Amount.IsZero() {
			usageRate, _ = h.Spent.Div(h.Amount).Float64()
		}
		items = append(items, response.BudgetHistoryResp{
			ID:          h.ID,
			PeriodStart: h.PeriodStart,
			PeriodEnd:   h.PeriodEnd,
			Amount:      h.Amount.StringFixed(4),
			Spent:       h.Spent.StringFixed(4),
			UsageRate:   usageRate,
			CreatedAt:   h.CreatedAt,
			UpdatedAt:   h.UpdatedAt,
		})
	}

	return items, nil
}

func (s *BudgetService) SnapshotHistory() (int, []string) {
	budgets, err := s.budgetRepo.ListAllEnabled()
	if err != nil {
		return 0, []string{err.Error()}
	}

	now := time.Now()
	twoYearsAgo := now.AddDate(-2, 0, 0)
	created := 0
	var errs []string

	for i := range budgets {
		b := &budgets[i]

		currentStart, _ := budgetPeriodRange(b.Period, now)

		cursor := now
		for {
			prevStart, prevEnd := budgetPeriodRange(b.Period, cursor)

			if prevStart.Before(twoYearsAgo) {
				break
			}

			if !prevStart.Before(currentStart) {
				cursor = budgetPeriodPrev(b.Period, cursor)
				continue
			}

			spent, spentErr := s.calculateSpentInRangeWithError(b, b.UserID, prevStart, prevEnd)
			if spentErr != nil {
				errs = append(errs, fmt.Sprintf("budget %d period %s: %v", b.ID, prevStart.Format("2006-01-02"), spentErr))
				cursor = budgetPeriodPrev(b.Period, cursor)
				continue
			}

			if err := s.budgetRepo.UpsertHistory(&model.BudgetHistory{
				BudgetID:    b.ID,
				PeriodStart: prevStart,
				PeriodEnd:   prevEnd,
				Amount:      b.Amount,
				Spent:       spent,
				CreatedAt:   now,
				UpdatedAt:   now,
			}); err != nil {
				errs = append(errs, fmt.Sprintf("budget %d upsert history: %v", b.ID, err))
			} else {
				created++
			}

			cursor = budgetPeriodPrev(b.Period, cursor)
		}
	}

	return created, errs
}

func budgetPeriodRange(period model.BudgetPeriod, t time.Time) (time.Time, time.Time) {
	loc := t.Location()
	switch period {
	case model.BudgetPeriodDaily:
		start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
		return start, start.AddDate(0, 0, 1).Add(-time.Second)
	case model.BudgetPeriodWeekly:
		weekday := int(t.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start := time.Date(t.Year(), t.Month(), t.Day()-weekday+1, 0, 0, 0, 0, loc)
		return start, start.AddDate(0, 0, 7).Add(-time.Second)
	case model.BudgetPeriodMonthly:
		start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, loc)
		return start, start.AddDate(0, 1, 0).Add(-time.Second)
	case model.BudgetPeriodQuarterly:
		quarter := (int(t.Month())-1)/3 + 1
		startMonth := time.Month((quarter-1)*3 + 1)
		start := time.Date(t.Year(), startMonth, 1, 0, 0, 0, 0, loc)
		return start, start.AddDate(0, 3, 0).Add(-time.Second)
	default:
		start := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, loc)
		return start, start.AddDate(1, 0, 0).Add(-time.Second)
	}
}

func budgetPeriodPrev(period model.BudgetPeriod, t time.Time) time.Time {
	switch period {
	case model.BudgetPeriodDaily:
		return t.AddDate(0, 0, -1)
	case model.BudgetPeriodWeekly:
		return t.AddDate(0, 0, -7)
	case model.BudgetPeriodMonthly:
		return t.AddDate(0, -1, 0)
	case model.BudgetPeriodQuarterly:
		return t.AddDate(0, -3, 0)
	default:
		return t.AddDate(-1, 0, 0)
	}
}

func (s *BudgetService) calculateSpent(budget *model.Budget, userID uint64) decimal.Decimal {
	if !budget.IsEnabled {
		return decimal.Zero
	}
	now := time.Now()
	start, end := budgetPeriodRange(budget.Period, now)
	spent, _ := s.calculateSpentInRangeWithError(budget, userID, start, end)
	return spent
}

func (s *BudgetService) calculateSpentInRangeWithError(budget *model.Budget, userID uint64, start, end time.Time) (decimal.Decimal, error) {
	if !budget.IsEnabled {
		return decimal.Zero, nil
	}

	var allCategoryIDs []uint64
	for _, cat := range budget.Categories {
		allCategoryIDs = append(allCategoryIDs, cat.ID)
	}
	expandedIDs, err := s.categoryRepo.GetDescendantIDs(allCategoryIDs, userID)
	if err != nil {
		return decimal.Zero, err
	}
	if len(expandedIDs) == 0 {
		return decimal.Zero, nil
	}

	filter := repository.TransactionFilter{
		Type:        string(model.TransactionTypeWithdrawal),
		StartDate:   start.Format("2006-01-02"),
		EndDate:     end.Format("2006-01-02"),
		CategoryIDs: expandedIDs,
	}
	txns, err := s.txnRepo.ListAll(userID, filter)
	if err != nil {
		return decimal.Zero, err
	}

	var total decimal.Decimal
	for _, txn := range txns {
		total = total.Add(txn.Amount)
	}
	return total, nil
}

func (s *BudgetService) toRespWithUsage(budget *model.Budget, userID uint64) (*response.BudgetResp, error) {
	spent := s.calculateSpent(budget, userID)
	remaining := budget.Amount.Sub(spent)

	var usageRate float64
	usageRateDecimal := decimal.Zero
	if !budget.Amount.IsZero() {
		usageRateDecimal = spent.Div(budget.Amount)
		usageRate, _ = usageRateDecimal.Float64()
	}

	status := "normal"
	if usageRateDecimal.GreaterThanOrEqual(decimal.NewFromInt(1)) {
		status = "overspent"
	} else if usageRateDecimal.GreaterThanOrEqual(decimal.NewFromFloat(0.8)) {
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

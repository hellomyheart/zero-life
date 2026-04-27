// Package service 业务逻辑层，实现核心业务逻辑
// BudgetService 预算业务逻辑，处理预算的增删改查及使用率计算
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

// BudgetService 预算服务
// 负责处理预算相关的业务逻辑，包括预算的创建、查询、更新、删除及使用率计算
// 依赖budgetRepo进行预算数据访问，依赖txnRepo查询交易以计算预算已花费金额
// 依赖categoryRepo展开分类的子孙节点ID（预算关联父分类时自动包含子分类的支出）
type BudgetService struct {
	budgetRepo   *repository.BudgetRepository       // 预算数据访问对象
	txnRepo      *repository.TransactionRepository   // 交易数据访问对象，用于计算预算周期内的支出
	categoryRepo *repository.CategoryRepository      // 分类数据访问对象，用于展开子孙分类ID
}

// NewBudgetService 创建预算服务实例
func NewBudgetService(budgetRepo *repository.BudgetRepository, txnRepo *repository.TransactionRepository, categoryRepo *repository.CategoryRepository) *BudgetService {
	return &BudgetService{
		budgetRepo:   budgetRepo,
		txnRepo:      txnRepo,
		categoryRepo: categoryRepo,
	}
}

// Create 创建预算
// 解析并验证金额，创建预算记录并关联分类，返回带使用率的预算信息
// 参数：
//   - userID: 用户ID
//   - req: 创建预算请求参数（名称、金额、周期、分类ID列表）
// 返回：
//   - *response.BudgetResp: 创建成功的预算信息（含已花费金额、剩余金额、使用率、状态）
//   - error: 错误信息
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

// Get 获取单个预算详情（含使用率计算）
// 参数：
//   - userID: 用户ID
//   - id: 预算ID
// 返回：
//   - *response.BudgetResp: 预算信息（含已花费金额、剩余金额、使用率、状态）
//   - error: 错误信息
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

// List 获取用户所有预算列表（每个预算含使用率计算）
// 参数：
//   - userID: 用户ID
// 返回：
//   - []response.BudgetResp: 预算列表
//   - error: 错误信息
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

// Update 更新预算信息
// 支持部分更新：名称、金额、周期、启用状态、关联分类
// 如果未提供分类ID列表，则保留现有分类关联
// 参数：
//   - userID: 用户ID
//   - id: 预算ID
//   - req: 更新请求参数
// 返回：
//   - *response.BudgetResp: 更新后的预算信息（含使用率）
//   - error: 错误信息
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

// Delete 删除预算
// 参数：
//   - userID: 用户ID
//   - id: 预算ID
// 返回：
//   - error: 错误信息
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

// GetHistory 获取预算的历史记录
// 返回每个历史周期的金额、已花费金额和使用率
// 参数：
//   - userID: 用户ID
//   - id: 预算ID
// 返回：
//   - []response.BudgetHistoryResp: 预算历史列表
//   - error: 错误信息
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
		// 修复除以零：当预算金额为零时，使用率为0
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
		})
	}

	return items, nil
}

// SnapshotCurrentPeriod 为所有已启用的预算生成当前周期的快照
// 遍历所有 is_enabled=true 的预算，检查当前周期是否已有快照，没有则计算支出并写入 BudgetHistory
// 返回生成的快照数量和错误列表
func (s *BudgetService) SnapshotCurrentPeriod() (int, []string) {
	budgets, err := s.budgetRepo.ListAllEnabled()
	if err != nil {
		return 0, []string{err.Error()}
	}

	now := time.Now()
	created := 0
	var errs []string

	for i := range budgets {
		b := &budgets[i]

		var periodStart, periodEnd time.Time
		periodStart, periodEnd = budgetPeriodRange(b.Period, now)

		exists, err := s.budgetRepo.HasHistory(b.ID, periodStart)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		if exists {
			continue
		}

		spent := s.calculateSpent(b, b.UserID)

		history := &model.BudgetHistory{
			BudgetID:    b.ID,
			PeriodStart: periodStart,
			PeriodEnd:   periodEnd,
			Amount:      b.Amount,
			Spent:       spent,
			CreatedAt:   now,
		}

		if err := s.budgetRepo.CreateHistory(history); err != nil {
			errs = append(errs, err.Error())
			continue
		}
		created++
	}

	return created, errs
}

func budgetPeriodRange(period model.BudgetPeriod, now time.Time) (time.Time, time.Time) {
	loc := now.Location()
	switch period {
	case model.BudgetPeriodDaily:
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		return start, start.AddDate(0, 0, 1).Add(-time.Second)
	case model.BudgetPeriodWeekly:
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start := time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, loc)
		return start, start.AddDate(0, 0, 7).Add(-time.Second)
	case model.BudgetPeriodMonthly:
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		return start, start.AddDate(0, 1, 0).Add(-time.Second)
	case model.BudgetPeriodQuarterly:
		quarter := (int(now.Month())-1)/3 + 1
		startMonth := time.Month((quarter-1)*3 + 1)
		start := time.Date(now.Year(), startMonth, 1, 0, 0, 0, 0, loc)
		return start, start.AddDate(0, 3, 0).Add(-time.Second)
	default:
		start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, loc)
		return start, start.AddDate(1, 0, 0).Add(-time.Second)
	}
}

// calculateSpent 计算预算在当前周期内已花费金额
// 业务流程：
// 1. 根据预算周期（月度/年度）计算当前周期的起止时间
// 2. 只统计支出类型(withdrawal)交易，预算追踪的是支出而非收入
// 3. 遍历预算关联的所有分类，查询每个分类下的支出交易并累加金额
// 修复：添加交易类型过滤，只统计支出类型(withdrawal)交易，避免将收入交易计入预算支出
func (s *BudgetService) calculateSpent(budget *model.Budget, userID uint64) decimal.Decimal {
	if !budget.IsEnabled {
		return decimal.Zero
	}

	now := time.Now()
	start, end := budgetPeriodRange(budget.Period, now)

	// 只统计支出类型交易，预算追踪的是支出而非收入
	withdrawalType := string(model.TransactionTypeWithdrawal)

	// 收集所有关联分类ID，并展开子孙分类（选中父分类时自动包含子分类的支出）
	var allCategoryIDs []uint64
	for _, cat := range budget.Categories {
		allCategoryIDs = append(allCategoryIDs, cat.ID)
	}
	expandedIDs, err := s.categoryRepo.GetDescendantIDs(allCategoryIDs, userID)
	if err != nil {
		return decimal.Zero
	}

	filter := repository.TransactionFilter{
		Type:        withdrawalType,
		StartDate:   start.Format("2006-01-02"),
		EndDate:     end.Format("2006-01-02"),
		CategoryIDs: expandedIDs,
	}
	txns, err := s.txnRepo.List(userID, filter, 0, 10000)
	if err != nil {
		return decimal.Zero
	}

	var total decimal.Decimal
	for _, txn := range txns {
		total = total.Add(txn.Amount)
	}

	return total
}

// toRespWithUsage 将预算模型转换为带使用率的响应对象
// 自动计算已花费金额、剩余金额、使用率和状态（normal/warning/overspent）
// 参数：
//   - budget: 预算模型
//   - userID: 用户ID（用于计算已花费金额）
// 返回：
//   - *response.BudgetResp: 预算响应对象
//   - error: 错误信息
func (s *BudgetService) toRespWithUsage(budget *model.Budget, userID uint64) (*response.BudgetResp, error) {
	spent := s.calculateSpent(budget, userID)
	remaining := budget.Amount.Sub(spent)
	// 修复除以零：当预算金额为零时，使用率为0
	var usageRate float64
	if !budget.Amount.IsZero() {
		usageRate, _ = spent.Div(budget.Amount).Float64()
	}

// 计算预算状态：使用率>=100%为超支(overspent)，>=80%为预警(warning)，否则正常(normal)
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

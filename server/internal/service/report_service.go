// Package service 业务逻辑层，实现核心业务逻辑
// ReportService 报表业务逻辑，生成收支、预算、净值等报表数据
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

// ReportService 报表服务
// 负责生成各类报表数据，包括收支报表、分类报表、预算报表、净值报表、趋势报表、标签报表和审计报表
// 依赖txnRepo查询交易数据，依赖accountRepo查询账户余额，依赖budgetRepo查询预算数据
// 依赖categoryRepo查询分类信息，依赖tagRepo查询标签信息
type ReportService struct {
	txnRepo     *repository.TransactionRepository  // 交易数据访问对象
	accountRepo *repository.AccountRepository      // 账户数据访问对象
	budgetRepo  *repository.BudgetRepository       // 预算数据访问对象
	categoryRepo *repository.CategoryRepository    // 分类数据访问对象
	tagRepo     *repository.TagRepository          // 标签数据访问对象
}

// NewReportService 创建报表服务实例
func NewReportService(
	txnRepo *repository.TransactionRepository,
	accountRepo *repository.AccountRepository,
	budgetRepo *repository.BudgetRepository,
	categoryRepo *repository.CategoryRepository,
	tagRepo *repository.TagRepository,
) *ReportService {
	return &ReportService{
		txnRepo:     txnRepo,
		accountRepo: accountRepo,
		budgetRepo:  budgetRepo,
		categoryRepo: categoryRepo,
		tagRepo:     tagRepo,
	}
}

// IncomeExpense 生成收支报表
// 统计指定时间范围内的总收入、总支出和净收入
// 参数：
//   - userID: 用户ID
//   - req: 报表请求参数（含日期范围）
// 返回：
//   - *response.IncomeExpenseResp: 收支报表数据
//   - error: 错误信息
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

// Category 生成分类报表
// 按分类统计支出和收入金额及百分比
// 参数：
//   - userID: 用户ID
//   - req: 报表请求参数（含日期范围）
// 返回：
//   - *response.CategoryReportResp: 分类报表数据
//   - error: 错误信息
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
		var pct float64
		// 修复除以零：当总支出为零时，百分比为0
		if !totalExpense.IsZero() {
			pct, _ = amount.Div(totalExpense).Mul(decimal.NewFromInt(100)).Float64()
		}
		expenseItems = append(expenseItems, response.CategoryItemResp{
			CategoryID:   catID,
			CategoryName: categoryNames[catID],
			Amount:       amount.StringFixed(4),
			Percentage:   pct,
		})
	}

	incomeItems := make([]response.CategoryItemResp, 0)
	for catID, amount := range incomeByCategory {
		var pct float64
		// 修复除以零：当总收入为零时，百分比为0
		if !totalIncome.IsZero() {
			pct, _ = amount.Div(totalIncome).Mul(decimal.NewFromInt(100)).Float64()
		}
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

// Budget 生成预算报表
// 统计每个启用预算的当前周期内已花费金额、剩余金额和使用率
// 使用预算自身的周期计算时间范围，展开分类后代确保子分类支出被计入
// 参数：
//   - userID: 用户ID
//   - req: 报表请求参数（含日期范围，用于筛选预算周期内的交易）
// 返回：
//   - *response.BudgetReportResp: 预算报表数据
//   - error: 错误信息
func (s *ReportService) Budget(userID uint64, req *request.ReportReq) (*response.BudgetReportResp, error) {
	budgets, err := s.budgetRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	now := time.Now()
	items := make([]response.BudgetReportItemResp, 0, len(budgets))
	for _, b := range budgets {
		if !b.IsEnabled {
			continue
		}

		// 根据预算周期计算当前周期范围
		periodStart, periodEnd := budgetPeriodRange(b.Period, now)

		// 展开分类后代：选择父分类时自动包含所有子分类的交易
		catIDs := make([]uint64, 0, len(b.Categories))
		for _, cat := range b.Categories {
			catIDs = append(catIDs, cat.ID)
		}
		expandedIDs, err := s.categoryRepo.GetDescendantIDs(catIDs, userID)
		if err != nil {
			continue
		}

		filter := repository.TransactionFilter{
			Type:        string(model.TransactionTypeWithdrawal),
			StartDate:   periodStart.Format("2006-01-02"),
			EndDate:     periodEnd.Format("2006-01-02"),
			CategoryIDs: expandedIDs,
		}

		txns, err := s.txnRepo.List(userID, filter, 0, 10000)
		if err != nil {
			continue
		}

		spent := decimal.Zero
		for _, txn := range txns {
			spent = spent.Add(txn.Amount)
		}

		remaining := b.Amount.Sub(spent)
		var usageRate float64
		if !b.Amount.IsZero() {
			usageRate, _ = spent.Div(b.Amount).Float64()
		}

		items = append(items, response.BudgetReportItemResp{
			BudgetID:   b.ID,
			BudgetName: b.Name,
			Amount:     b.Amount.StringFixed(4),
			Spent:      spent.StringFixed(4),
			Remaining:  remaining.StringFixed(4),
			UsageRate:  usageRate,
		})
	}

	return &response.BudgetReportResp{Items: items}, nil
}

// NetWorth 生成净值报表
// 业务流程：
// 1. 计算总资产（所有资产账户的当前余额之和）
// 2. 计算总负债（所有负债账户的当前余额之和）
// 3. 净值 = 总资产 - 总负债
// 4. 按月生成净值趋势：从初始余额开始，逐月累加交易影响
//    - 存款增加资产，取款减少资产，转账不影响净值
// 参数：
//   - userID: 用户ID
//   - req: 报表请求参数（含日期范围）
// 返回：
//   - *response.NetWorthResp: 净值报表数据（含趋势）
//   - error: 错误信息
func (s *ReportService) NetWorth(userID uint64, req *request.ReportReq) (*response.NetWorthResp, error) {
	accounts, err := s.accountRepo.List(userID, string(model.AccountTypeAsset), "", "name", 0, 1000)
	if err != nil {
		return nil, errcode.ErrInternal
	}

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

	// Calculate trend by month
	// 修正逻辑：每月净值应从初始余额开始，逐月累加从报告起始到该月末的交易影响
	// 这样每月的净值反映的是"从开始到该月末的累计净值"，而非"仅该月的净值"
	startDate, endDate, err := s.parseDateRange(req)
	if err != nil {
		return &response.NetWorthResp{
			TotalNetWorth: netWorth.StringFixed(4),
			Trend:         []response.NetWorthPointResp{},
		}, nil
	}

	trend := make([]response.NetWorthPointResp, 0)

	// 先计算初始余额作为基准
	initialAssets := decimal.Zero
	for _, a := range accounts {
		initialAssets = initialAssets.Add(a.InitialBalance)
	}
	initialLiabilities := decimal.Zero
	for _, a := range liabilities {
		initialLiabilities = initialLiabilities.Add(a.InitialBalance);
	}

	// 获取整个报告时间范围内的所有交易，一次性查询避免重复查询
	allFilter := repository.TransactionFilter{
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),
	}
	allTxns, err := s.txnRepo.List(userID, allFilter, 0, 10000)
	if err != nil {
		allTxns = nil
	}

	// 按月分组交易
	monthlyTxnEffects := make(map[string]decimal.Decimal)
	for _, txn := range allTxns {
		monthKey := txn.Date.Format("2006-01")
		effect := decimal.Zero
		switch txn.Type {
		case model.TransactionTypeDeposit:
			effect = txn.Amount
		case model.TransactionTypeWithdrawal:
			effect = txn.Amount.Neg()
		case model.TransactionTypeTransfer:
			// 转账不影响净值
		}
		monthlyTxnEffects[monthKey] = monthlyTxnEffects[monthKey].Add(effect)
	}

	// 逐月累加净值趋势
	current := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
	cumulativeNetWorth := initialAssets.Sub(initialLiabilities)
	for current.Before(endDate) || current.Equal(endDate) {
		monthKey := current.Format("2006-01")
		// 累加该月的交易影响
		if effect, ok := monthlyTxnEffects[monthKey]; ok {
			cumulativeNetWorth = cumulativeNetWorth.Add(effect)
		}

		trend = append(trend, response.NetWorthPointResp{
			Date:     monthKey,
			NetWorth: cumulativeNetWorth.StringFixed(4),
		})

		current = current.AddDate(0, 1, 0)
	}

	return &response.NetWorthResp{
		TotalNetWorth: netWorth.StringFixed(4),
		Trend:         trend,
	}, nil
}

// Trend 生成趋势报表
// 按月统计收入和支出趋势，支持自定义粒度
// 参数：
//   - userID: 用户ID
//   - req: 报表请求参数（含日期范围、粒度）
// 返回：
//   - *response.TrendResp: 趋势报表数据
//   - error: 错误信息
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

// Tag 生成标签报表
// 按标签统计收入和支出金额
// 参数：
//   - userID: 用户ID
//   - req: 报表请求参数（含日期范围）
// 返回：
//   - *response.TagReportResp: 标签报表数据
//   - error: 错误信息
func (s *ReportService) Tag(userID uint64, req *request.ReportReq) (*response.TagReportResp, error) {
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

	type tagStats struct {
		income  decimal.Decimal
		expense decimal.Decimal
	}
	tagData := make(map[uint64]*tagStats)
	tagNames := make(map[uint64]string)

	for _, txn := range txns {
		for _, t := range txn.Tags {
			if _, ok := tagData[t.ID]; !ok {
				tagData[t.ID] = &tagStats{income: decimal.Zero, expense: decimal.Zero}
				tagNames[t.ID] = t.Name
			}
			switch txn.Type {
			case model.TransactionTypeDeposit:
				tagData[t.ID].income = tagData[t.ID].income.Add(txn.Amount)
			case model.TransactionTypeWithdrawal:
				tagData[t.ID].expense = tagData[t.ID].expense.Add(txn.Amount)
			}
		}
	}

	items := make([]response.TagReportItemResp, 0, len(tagData))
	for tagID, stats := range tagData {
		items = append(items, response.TagReportItemResp{
			TagID:   tagID,
			TagName: tagNames[tagID],
			Income:  stats.income.StringFixed(4),
			Expense: stats.expense.StringFixed(4),
		})
	}

	return &response.TagReportResp{Items: items}, nil
}

// AuditReport 生成审计报表
// 按账户生成交易流水，包含每笔交易后的运行余额（running balance）
// 参数：
//   - userID: 用户ID
//   - accountID: 账户ID
//   - startDate: 开始日期
//   - endDate: 结束日期
//   - reconciled: 是否只显示已对账的交易（nil表示不过滤）
// 返回：
//   - *response.AuditReportResp: 审计报表数据
//   - error: 错误信息
func (s *ReportService) AuditReport(userID uint64, accountID uint64, startDate, endDate string, reconciled *bool) (*response.AuditReportResp, error) {
	// Verify account belongs to user
	account, err := s.accountRepo.GetByID(accountID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	txns, err := s.txnRepo.GetForAudit(userID, accountID, startDate, endDate, reconciled)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	runningBalance := account.InitialBalance
	items := make([]response.AuditReportItemResp, 0, len(txns))

	for _, txn := range txns {
		switch txn.Type {
		case model.TransactionTypeDeposit:
			if txn.DestinationID != nil && *txn.DestinationID == accountID {
				runningBalance = runningBalance.Add(txn.Amount)
			}
		case model.TransactionTypeWithdrawal:
			if txn.SourceID == accountID {
				runningBalance = runningBalance.Sub(txn.Amount)
			}
		case model.TransactionTypeTransfer:
			if txn.SourceID == accountID {
				runningBalance = runningBalance.Sub(txn.Amount)
			}
			if txn.DestinationID != nil && *txn.DestinationID == accountID {
				runningBalance = runningBalance.Add(txn.Amount)
			}
		}

		item := response.AuditReportItemResp{
			TransactionID: txn.ID,
			Date:          txn.Date,
			Description:   txn.Description,
			Type:          string(txn.Type),
			Amount:        txn.Amount.StringFixed(4),
			RunningBalance: runningBalance.StringFixed(4),
			IsReconciled:  txn.IsReconciled,
		}

		if txn.Category != nil {
			item.CategoryName = txn.Category.Name
		}

		items = append(items, item)
	}

	return &response.AuditReportResp{
		AccountID:      accountID,
		AccountName:    account.Name,
		InitialBalance: account.InitialBalance.StringFixed(4),
		FinalBalance:   runningBalance.StringFixed(4),
		StartDate:      startDate,
		EndDate:        endDate,
		Items:          items,
	}, nil
}

// parseDateRange 解析报表请求中的日期范围
// 参数：
//   - req: 报表请求参数
// 返回：
//   - time.Time: 开始日期
//   - time.Time: 结束日期
//   - error: 解析错误
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

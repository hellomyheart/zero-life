// Package service 业务逻辑层，实现核心业务逻辑
// ChartService 图表业务逻辑，提供各类图表数据的计算和聚合
package service

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/repository"
)

// ChartService 图表服务
// 提供各类图表数据的计算和聚合，包括账户余额趋势、预算支出、分类分布、标签分布和交易趋势
// 依赖txnRepo查询交易数据，依赖accountRepo查询账户信息，依赖budgetRepo查询预算数据
// 依赖categoryRepo查询分类名称，依赖tagRepo查询标签名称
type ChartService struct {
	txnRepo     *repository.TransactionRepository // 交易数据访问对象
	accountRepo *repository.AccountRepository     // 账户数据访问对象
	budgetRepo  *repository.BudgetRepository      // 预算数据访问对象
	categoryRepo *repository.CategoryRepository   // 分类数据访问对象
	tagRepo     *repository.TagRepository         // 标签数据访问对象
}

// NewChartService 创建图表服务实例
// 参数：
//   - txnRepo: 交易数据访问对象
//   - accountRepo: 账户数据访问对象
//   - budgetRepo: 预算数据访问对象
//   - categoryRepo: 分类数据访问对象
//   - tagRepo: 标签数据访问对象
// 返回：
//   - *ChartService: 图表服务实例
func NewChartService(
	txnRepo *repository.TransactionRepository,
	accountRepo *repository.AccountRepository,
	budgetRepo *repository.BudgetRepository,
	categoryRepo *repository.CategoryRepository,
	tagRepo *repository.TagRepository,
) *ChartService {
	return &ChartService{
		txnRepo:     txnRepo,
		accountRepo: accountRepo,
		budgetRepo:  budgetRepo,
		categoryRepo: categoryRepo,
		tagRepo:     tagRepo,
	}
}

// ChartDataPoint 图表数据点
type ChartDataPoint struct {
	Date  string          `json:"date"`  // 日期
	Value decimal.Decimal `json:"value"` // 数值
}

// AccountBalanceData 账户余额图表数据
type AccountBalanceData struct {
	AccountID   uint64            `json:"account_id"`
	AccountName string            `json:"account_name"`
	Points      []ChartDataPoint  `json:"points"`
}

// AccountBalance 获取账户余额变化趋势
// 业务流程：
// 1. 获取账户信息（含初始余额）
// 2. 获取时间范围内的所有交易
// 3. 按日期计算每笔交易对账户余额的影响（支出/转出减少，收入/转入增加）
// 4. 从初始余额开始，逐日累加生成余额趋势数据点
// 参数：
//   - userID: 用户ID
//   - accountID: 账户ID
//   - startDate: 开始日期
//   - endDate: 结束日期
// 返回：
//   - *AccountBalanceData: 账户余额趋势数据
//   - error: 错误信息
func (s *ChartService) AccountBalance(userID, accountID uint64, startDate, endDate time.Time) (*AccountBalanceData, error) {
	// 获取账户信息
	account, err := s.accountRepo.GetByID(accountID, userID)
	if err != nil {
		return nil, err
	}

	// 获取时间范围内的交易
	txns, err := s.txnRepo.GetByAccountAndDateRange(userID, accountID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 计算每日余额
	points := make([]ChartDataPoint, 0)
	balance := account.InitialBalance

	// 按日期分组计算
	dateMap := make(map[string]decimal.Decimal)
	for _, txn := range txns {
		dateStr := txn.Date.Format("2006-01-02")
		if _, exists := dateMap[dateStr]; !exists {
			dateMap[dateStr] = decimal.Zero
		}

		// 根据交易类型调整余额
		if txn.SourceID == accountID {
			// 支出或转出
			dateMap[dateStr] = dateMap[dateStr].Sub(txn.Amount)
		}
		if txn.DestinationID != nil && *txn.DestinationID == accountID {
			// 收入或转入
			dateMap[dateStr] = dateMap[dateStr].Add(txn.Amount)
		}
	}

	// 生成数据点
	currentDate := startDate
	for !currentDate.After(endDate) {
		dateStr := currentDate.Format("2006-01-02")
		if delta, exists := dateMap[dateStr]; exists {
			balance = balance.Add(delta)
		}
		points = append(points, ChartDataPoint{
			Date:  dateStr,
			Value: balance,
		})
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return &AccountBalanceData{
		AccountID:   accountID,
		AccountName: account.Name,
		Points:      points,
	}, nil
}

// BudgetSpendingData 预算支出图表数据
type BudgetSpendingData struct {
	BudgetID   uint64           `json:"budget_id"`
	BudgetName string           `json:"budget_name"`
	Amount     decimal.Decimal  `json:"amount"`
	Spent      decimal.Decimal  `json:"spent"`
	Remaining  decimal.Decimal  `json:"remaining"`
	Percentage decimal.Decimal  `json:"percentage"`
}

// BudgetSpending 获取预算支出数据
// 业务流程：
// 1. 获取预算信息（含预算金额）
// 2. 计算已花费金额
// 3. 计算剩余金额和支出百分比
// 参数：
//   - userID: 用户ID
//   - budgetID: 预算ID
// 返回：
//   - *BudgetSpendingData: 预算支出数据
//   - error: 错误信息
func (s *ChartService) BudgetSpending(userID, budgetID uint64) (*BudgetSpendingData, error) {
	budget, err := s.budgetRepo.GetByID(budgetID, userID)
	if err != nil {
		return nil, err
	}

	// 计算已花费金额（含分类后代展开）
	var spent decimal.Decimal
	if budget.IsEnabled {
		now := time.Now()
		start, end := budgetPeriodRange(budget.Period, now)

		var allCategoryIDs []uint64
		for _, cat := range budget.Categories {
			allCategoryIDs = append(allCategoryIDs, cat.ID)
		}
		expandedIDs, err := s.categoryRepo.GetDescendantIDs(allCategoryIDs, userID)
		if err == nil && len(expandedIDs) > 0 {
			filter := repository.TransactionFilter{
				Type:        string(model.TransactionTypeWithdrawal),
				StartDate:   start.Format("2006-01-02"),
				EndDate:     end.Format("2006-01-02"),
				CategoryIDs: expandedIDs,
			}
			txns, err := s.txnRepo.List(userID, filter, 0, 10000)
			if err == nil {
				for _, txn := range txns {
					spent = spent.Add(txn.Amount)
				}
			}
		}
	}

	remaining := budget.Amount.Sub(spent)
	percentage := decimal.Zero
	if !budget.Amount.IsZero() {
		percentage = spent.Div(budget.Amount).Mul(decimal.NewFromInt(100))
	}

	return &BudgetSpendingData{
		BudgetID:   budgetID,
		BudgetName: budget.Name,
		Amount:     budget.Amount,
		Spent:      spent,
		Remaining:  remaining,
		Percentage: percentage,
	}, nil
}

// CategoryDistributionData 分类分布图表数据
type CategoryDistributionData struct {
	CategoryID   uint64          `json:"category_id"`
	CategoryName string          `json:"category_name"`
	Amount       decimal.Decimal `json:"amount"`
	Percentage   decimal.Decimal `json:"percentage"`
}

// CategoryDistribution 获取分类支出/收入分布
// 业务流程：
// 1. 获取时间范围内的所有交易
// 2. 根据chartType过滤交易类型（expense=支出，income=收入）
// 3. 按分类汇总金额和计算百分比
// 参数：
//   - userID: 用户ID
//   - startDate: 开始日期
//   - endDate: 结束日期
//   - chartType: 图表类型（expense=支出分布，income=收入分布）
// 返回：
//   - []CategoryDistributionData: 分类分布数据列表
//   - error: 错误信息
func (s *ChartService) CategoryDistribution(userID uint64, startDate, endDate time.Time, chartType string) ([]CategoryDistributionData, error) {
	// 获取时间范围内的交易
	txns, err := s.txnRepo.GetByDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 按分类汇总
	categoryMap := make(map[uint64]decimal.Decimal)
	total := decimal.Zero

	for _, txn := range txns {
		// 根据类型筛选
		if chartType == "expense" && txn.Type != "withdrawal" {
			continue
		}
		if chartType == "income" && txn.Type != "deposit" {
			continue
		}

		if txn.CategoryID != nil {
			categoryMap[*txn.CategoryID] = categoryMap[*txn.CategoryID].Add(txn.Amount)
			total = total.Add(txn.Amount)
		}
	}

	// 生成结果
	result := make([]CategoryDistributionData, 0)
	for categoryID, amount := range categoryMap {
		category, err := s.categoryRepo.GetByID(categoryID, userID)
		if err != nil {
			continue
		}

		percentage := decimal.Zero
		if !total.IsZero() {
			percentage = amount.Div(total).Mul(decimal.NewFromInt(100))
		}

		result = append(result, CategoryDistributionData{
			CategoryID:   categoryID,
			CategoryName: category.Name,
			Amount:       amount,
			Percentage:   percentage,
		})
	}

	return result, nil
}

// TagDistributionData 标签分布图表数据
type TagDistributionData struct {
	TagID      uint64          `json:"tag_id"`
	TagName    string          `json:"tag_name"`
	Amount     decimal.Decimal `json:"amount"`
	Percentage decimal.Decimal `json:"percentage"`
}

// TagDistribution 获取标签分布
// 业务流程：
// 1. 获取时间范围内的所有交易
// 2. 按标签汇总金额和计算百分比
// 参数：
//   - userID: 用户ID
//   - startDate: 开始日期
//   - endDate: 结束日期
// 返回：
//   - []TagDistributionData: 标签分布数据列表
//   - error: 错误信息
func (s *ChartService) TagDistribution(userID uint64, startDate, endDate time.Time) ([]TagDistributionData, error) {
	// 获取时间范围内的交易
	txns, err := s.txnRepo.GetByDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 按标签汇总
	tagMap := make(map[uint64]decimal.Decimal)
	total := decimal.Zero

	for _, txn := range txns {
		for _, tag := range txn.Tags {
			tagMap[tag.ID] = tagMap[tag.ID].Add(txn.Amount)
			total = total.Add(txn.Amount)
		}
	}

	// 生成结果
	result := make([]TagDistributionData, 0)
	for tagID, amount := range tagMap {
		tag, err := s.tagRepo.GetByID(tagID, userID)
		if err != nil {
			continue
		}

		percentage := decimal.Zero
		if !total.IsZero() {
			percentage = amount.Div(total).Mul(decimal.NewFromInt(100))
		}

		result = append(result, TagDistributionData{
			TagID:      tagID,
			TagName:    tag.Name,
			Amount:     amount,
			Percentage: percentage,
		})
	}

	return result, nil
}

// TransactionTrendData 交易趋势图表数据
type TransactionTrendData struct {
	Income     []ChartDataPoint `json:"income"`
	Expense    []ChartDataPoint `json:"expense"`
	Transfer   []ChartDataPoint `json:"transfer"`
}

// TransactionTrend 获取交易趋势数据
// 业务流程：
// 1. 获取时间范围内的所有交易
// 2. 按日期和交易类型（收入/支出/转账）汇总金额
// 3. 为每个日期生成三种类型的数据点（无交易的日期值为0）
// 参数：
//   - userID: 用户ID
//   - startDate: 开始日期
//   - endDate: 结束日期
// 返回：
//   - *TransactionTrendData: 交易趋势数据（含收入、支出、转账三条曲线）
//   - error: 错误信息
func (s *ChartService) TransactionTrend(userID uint64, startDate, endDate time.Time) (*TransactionTrendData, error) {
	// 获取时间范围内的交易
	txns, err := s.txnRepo.GetByDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 按日期和类型汇总
	incomeMap := make(map[string]decimal.Decimal)
	expenseMap := make(map[string]decimal.Decimal)
	transferMap := make(map[string]decimal.Decimal)

	for _, txn := range txns {
		dateStr := txn.Date.Format("2006-01-02")
		switch txn.Type {
		case "deposit":
			incomeMap[dateStr] = incomeMap[dateStr].Add(txn.Amount)
		case "withdrawal":
			expenseMap[dateStr] = expenseMap[dateStr].Add(txn.Amount)
		case "transfer":
			transferMap[dateStr] = transferMap[dateStr].Add(txn.Amount)
		}
	}

	// 生成数据点
	incomePoints := make([]ChartDataPoint, 0)
	expensePoints := make([]ChartDataPoint, 0)
	transferPoints := make([]ChartDataPoint, 0)

	currentDate := startDate
	for !currentDate.After(endDate) {
		dateStr := currentDate.Format("2006-01-02")
		incomePoints = append(incomePoints, ChartDataPoint{
			Date:  dateStr,
			Value: incomeMap[dateStr],
		})
		expensePoints = append(expensePoints, ChartDataPoint{
			Date:  dateStr,
			Value: expenseMap[dateStr],
		})
		transferPoints = append(transferPoints, ChartDataPoint{
			Date:  dateStr,
			Value: transferMap[dateStr],
		})
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return &TransactionTrendData{
		Income:   incomePoints,
		Expense:  expensePoints,
		Transfer: transferPoints,
	}, nil
}

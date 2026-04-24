// Package service 业务逻辑层
package service

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/hellomyheart/zero-life/server/internal/repository"
)

// ChartService 图表服务
// 提供各类图表数据的计算和聚合
type ChartService struct {
	txnRepo     *repository.TransactionRepository
	accountRepo *repository.AccountRepository
	budgetRepo  *repository.BudgetRepository
	categoryRepo *repository.CategoryRepository
	tagRepo     *repository.TagRepository
}

// NewChartService 创建图表服务实例
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
// 计算指定时间范围内账户余额的变化
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
func (s *ChartService) BudgetSpending(userID, budgetID uint64) (*BudgetSpendingData, error) {
	budget, err := s.budgetRepo.GetByID(budgetID, userID)
	if err != nil {
		return nil, err
	}

	// 计算已花费金额
	spent, err := s.budgetRepo.GetSpentAmount(budgetID, userID)
	if err != nil {
		spent = decimal.Zero
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

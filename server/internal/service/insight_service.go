// Package service 业务逻辑层，实现核心业务逻辑
// InsightService 数据洞察业务逻辑，提供深度数据分析功能
package service

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/hellomyheart/zero-life/server/internal/repository"
)

// InsightService 数据洞察服务
// 提供深度数据分析功能，包括支出洞察、收入洞察和转账洞察
// 依赖txnRepo查询交易数据，依赖accountRepo查询账户名称，依赖categoryRepo查询分类名称
type InsightService struct {
	txnRepo     *repository.TransactionRepository // 交易数据访问对象
	accountRepo *repository.AccountRepository     // 账户数据访问对象
	categoryRepo *repository.CategoryRepository   // 分类数据访问对象
}

// NewInsightService 创建数据洞察服务实例
// 参数：
//   - txnRepo: 交易数据访问对象
//   - accountRepo: 账户数据访问对象
//   - categoryRepo: 分类数据访问对象
// 返回：
//   - *InsightService: 数据洞察服务实例
func NewInsightService(
	txnRepo *repository.TransactionRepository,
	accountRepo *repository.AccountRepository,
	categoryRepo *repository.CategoryRepository,
) *InsightService {
	return &InsightService{
		txnRepo:     txnRepo,
		accountRepo: accountRepo,
		categoryRepo: categoryRepo,
	}
}

// ExpenseInsightData 支出洞察数据
type ExpenseInsightData struct {
	Total          decimal.Decimal              `json:"total"`           // 总支出
	Average        decimal.Decimal              `json:"average"`         // 平均支出
	Max            decimal.Decimal              `json:"max"`             // 最大单笔支出
	Min            decimal.Decimal              `json:"min"`             // 最小单笔支出
	Count          int                          `json:"count"`           // 交易笔数
	ByCategory     []CategoryInsight            `json:"by_category"`     // 按分类统计
	ByAccount      []AccountInsight             `json:"by_account"`      // 按账户统计
	ByDate         []DateInsight                `json:"by_date"`         // 按日期统计
	Trend          string                       `json:"trend"`           // 趋势（up/down/stable）
}

// IncomeInsightData 收入洞察数据
type IncomeInsightData struct {
	Total          decimal.Decimal              `json:"total"`           // 总收入
	Average        decimal.Decimal              `json:"average"`         // 平均收入
	Max            decimal.Decimal              `json:"max"`             // 最大单笔收入
	Min            decimal.Decimal              `json:"min"`             // 最小单笔收入
	Count          int                          `json:"count"`           // 交易笔数
	ByCategory     []CategoryInsight            `json:"by_category"`     // 按分类统计
	ByAccount      []AccountInsight             `json:"by_account"`      // 按账户统计
	ByDate         []DateInsight                `json:"by_date"`         // 按日期统计
	Trend          string                       `json:"trend"`           // 趋势
}

// TransferInsightData 转账洞察数据
type TransferInsightData struct {
	Total          decimal.Decimal              `json:"total"`           // 总转账金额
	Count          int                          `json:"count"`           // 转账笔数
	ByAccount      []AccountInsight             `json:"by_account"`      // 按账户统计
	FrequentPairs  []TransferPair               `json:"frequent_pairs"`  // 频繁转账对
}

// CategoryInsight 分类洞察
type CategoryInsight struct {
	CategoryID   uint64          `json:"category_id"`
	CategoryName string          `json:"category_name"`
	Amount       decimal.Decimal `json:"amount"`
	Count        int             `json:"count"`
	Percentage   decimal.Decimal `json:"percentage"`
}

// AccountInsight 账户洞察
type AccountInsight struct {
	AccountID   uint64          `json:"account_id"`
	AccountName string          `json:"account_name"`
	Amount      decimal.Decimal `json:"amount"`
	Count       int             `json:"count"`
	Percentage  decimal.Decimal `json:"percentage"`
}

// DateInsight 日期洞察
type DateInsight struct {
	Date   string          `json:"date"`
	Amount decimal.Decimal `json:"amount"`
	Count  int             `json:"count"`
}

// TransferPair 转账对
type TransferPair struct {
	SourceAccountID   uint64          `json:"source_account_id"`
	SourceAccountName string          `json:"source_account_name"`
	DestAccountID     uint64          `json:"dest_account_id"`
	DestAccountName   string          `json:"dest_account_name"`
	Count             int             `json:"count"`
	TotalAmount       decimal.Decimal `json:"total_amount"`
}

// ExpenseInsight 支出洞察分析
// 业务流程：
// 1. 获取指定时间范围内的支出交易
// 2. 计算基础统计：总支出、平均支出、最大/最小单笔支出、交易笔数
// 3. 按分类汇总：每个分类的支出金额、笔数、百分比
// 4. 按账户汇总：每个支出账户的支出金额、笔数、百分比
// 5. 按日期汇总：每天的支出金额和笔数
// 参数：
//   - userID: 用户ID
//   - startDate: 开始日期
//   - endDate: 结束日期
// 返回：
//   - *ExpenseInsightData: 支出洞察数据
//   - error: 错误信息
func (s *InsightService) ExpenseInsight(userID uint64, startDate, endDate time.Time) (*ExpenseInsightData, error) {
	// 获取支出交易
	txns, err := s.txnRepo.GetByTypeAndDateRange(userID, "withdrawal", startDate, endDate)
	if err != nil {
		return nil, err
	}

	if len(txns) == 0 {
		return &ExpenseInsightData{}, nil
	}

	// 计算基础统计
	total := decimal.Zero
	max := decimal.Zero
	min := decimal.Zero
	count := len(txns)

	categoryMap := make(map[uint64]*CategoryInsight)
	accountMap := make(map[uint64]*AccountInsight)
	dateMap := make(map[string]*DateInsight)

	for i, txn := range txns {
		total = total.Add(txn.Amount)
		if i == 0 || txn.Amount.GreaterThan(max) {
			max = txn.Amount
		}
		if i == 0 || txn.Amount.LessThan(min) {
			min = txn.Amount
		}

		// 按分类统计
		if txn.CategoryID != nil {
			catID := *txn.CategoryID
			if _, exists := categoryMap[catID]; !exists {
				category, _ := s.categoryRepo.GetByID(catID, userID)
				categoryMap[catID] = &CategoryInsight{
					CategoryID:   catID,
					CategoryName: category.Name,
				}
			}
			categoryMap[catID].Amount = categoryMap[catID].Amount.Add(txn.Amount)
			categoryMap[catID].Count++
		}

		// 按账户统计
		accID := txn.SourceID
		if _, exists := accountMap[accID]; !exists {
			account, _ := s.accountRepo.GetByID(accID, userID)
			accountMap[accID] = &AccountInsight{
				AccountID:   accID,
				AccountName: account.Name,
			}
		}
		accountMap[accID].Amount = accountMap[accID].Amount.Add(txn.Amount)
		accountMap[accID].Count++

		// 按日期统计
		dateStr := txn.Date.Format("2006-01-02")
		if _, exists := dateMap[dateStr]; !exists {
			dateMap[dateStr] = &DateInsight{Date: dateStr}
		}
		dateMap[dateStr].Amount = dateMap[dateStr].Amount.Add(txn.Amount)
		dateMap[dateStr].Count++
	}

	// 计算平均值和百分比
	average := total.Div(decimal.NewFromInt(int64(count)))

	byCategory := make([]CategoryInsight, 0)
	for _, cat := range categoryMap {
		if !total.IsZero() {
			cat.Percentage = cat.Amount.Div(total).Mul(decimal.NewFromInt(100))
		}
		byCategory = append(byCategory, *cat)
	}

	byAccount := make([]AccountInsight, 0)
	for _, acc := range accountMap {
		if !total.IsZero() {
			acc.Percentage = acc.Amount.Div(total).Mul(decimal.NewFromInt(100))
		}
		byAccount = append(byAccount, *acc)
	}

	byDate := make([]DateInsight, 0)
	for _, date := range dateMap {
		byDate = append(byDate, *date)
	}

	return &ExpenseInsightData{
		Total:      total,
		Average:    average,
		Max:        max,
		Min:        min,
		Count:      count,
		ByCategory: byCategory,
		ByAccount:  byAccount,
		ByDate:     byDate,
		Trend:      calculateTrend(byDate),
	}, nil
}

// IncomeInsight 收入洞察分析
// 业务流程：
// 1. 获取指定时间范围内的收入交易
// 2. 计算基础统计：总收入、平均收入、最大/最小单笔收入、交易笔数
// 3. 按分类汇总：每个分类的收入金额、笔数、百分比
// 4. 按账户汇总：每个收入账户的收入金额、笔数、百分比
// 5. 按日期汇总：每天的收入金额和笔数
// 6. 计算趋势：比较前半段和后半段的收入总额判断趋势方向
// 参数：
//   - userID: 用户ID
//   - startDate: 开始日期
//   - endDate: 结束日期
// 返回：
//   - *IncomeInsightData: 收入洞察数据
//   - error: 错误信息
func (s *InsightService) IncomeInsight(userID uint64, startDate, endDate time.Time) (*IncomeInsightData, error) {
	// 获取收入交易
	txns, err := s.txnRepo.GetByTypeAndDateRange(userID, "deposit", startDate, endDate)
	if err != nil {
		return nil, err
	}

	if len(txns) == 0 {
		return &IncomeInsightData{}, nil
	}

	// 计算基础统计（逻辑与支出类似）
	total := decimal.Zero
	max := decimal.Zero
	min := decimal.Zero
	count := len(txns)

	categoryMap := make(map[uint64]*CategoryInsight)
	accountMap := make(map[uint64]*AccountInsight)
	dateMap := make(map[string]*DateInsight)

	for i, txn := range txns {
		total = total.Add(txn.Amount)
		if i == 0 || txn.Amount.GreaterThan(max) {
			max = txn.Amount
		}
		if i == 0 || txn.Amount.LessThan(min) {
			min = txn.Amount
		}

		// 按分类统计
		if txn.CategoryID != nil {
			catID := *txn.CategoryID
			if _, exists := categoryMap[catID]; !exists {
				category, _ := s.categoryRepo.GetByID(catID, userID)
				categoryMap[catID] = &CategoryInsight{
					CategoryID:   catID,
					CategoryName: category.Name,
				}
			}
			categoryMap[catID].Amount = categoryMap[catID].Amount.Add(txn.Amount)
			categoryMap[catID].Count++
		}

		// 按账户统计（收入交易的目标账户）
		if txn.DestinationID != nil {
			accID := *txn.DestinationID
			if _, exists := accountMap[accID]; !exists {
				account, _ := s.accountRepo.GetByID(accID, userID)
				accountMap[accID] = &AccountInsight{
					AccountID:   accID,
					AccountName: account.Name,
				}
			}
			accountMap[accID].Amount = accountMap[accID].Amount.Add(txn.Amount)
			accountMap[accID].Count++
		}

		// 按日期统计
		dateStr := txn.Date.Format("2006-01-02")
		if _, exists := dateMap[dateStr]; !exists {
			dateMap[dateStr] = &DateInsight{Date: dateStr}
		}
		dateMap[dateStr].Amount = dateMap[dateStr].Amount.Add(txn.Amount)
		dateMap[dateStr].Count++
	}

	// 计算平均值和百分比
	average := total.Div(decimal.NewFromInt(int64(count)))

	byCategory := make([]CategoryInsight, 0)
	for _, cat := range categoryMap {
		if !total.IsZero() {
			cat.Percentage = cat.Amount.Div(total).Mul(decimal.NewFromInt(100))
		}
		byCategory = append(byCategory, *cat)
	}

	byAccount := make([]AccountInsight, 0)
	for _, acc := range accountMap {
		if !total.IsZero() {
			acc.Percentage = acc.Amount.Div(total).Mul(decimal.NewFromInt(100))
		}
		byAccount = append(byAccount, *acc)
	}

	byDate := make([]DateInsight, 0)
	for _, date := range dateMap {
		byDate = append(byDate, *date)
	}

	return &IncomeInsightData{
		Total:      total,
		Average:    average,
		Max:       max,
		Min:       min,
		Count:     count,
		ByCategory: byCategory,
		ByAccount:  byAccount,
		ByDate:     byDate,
		Trend:      calculateTrend(byDate),
	}, nil
}

// TransferInsight 转账洞察分析
// 业务流程：
// 1. 获取指定时间范围内的转账交易
// 2. 计算总转账金额和笔数
// 3. 统计频繁转账对：按"源账户->目标账户"分组，记录每对的转账次数和总金额
// 参数：
//   - userID: 用户ID
//   - startDate: 开始日期
//   - endDate: 结束日期
// 返回：
//   - *TransferInsightData: 转账洞察数据
//   - error: 错误信息
func (s *InsightService) TransferInsight(userID uint64, startDate, endDate time.Time) (*TransferInsightData, error) {
	// 获取转账交易
	txns, err := s.txnRepo.GetByTypeAndDateRange(userID, "transfer", startDate, endDate)
	if err != nil {
		return nil, err
	}

	if len(txns) == 0 {
		return &TransferInsightData{}, nil
	}

	total := decimal.Zero
	count := len(txns)

	// 统计转账对
	pairMap := make(map[string]*TransferPair)

	for _, txn := range txns {
		total = total.Add(txn.Amount)

		if txn.DestinationID != nil {
			// 创建转账对键
			pairKey := ""
			srcAccount, _ := s.accountRepo.GetByID(txn.SourceID, userID)
			destAccount, _ := s.accountRepo.GetByID(*txn.DestinationID, userID)
			pairKey = srcAccount.Name + "->" + destAccount.Name

			if _, exists := pairMap[pairKey]; !exists {
				pairMap[pairKey] = &TransferPair{
					SourceAccountID:   txn.SourceID,
					SourceAccountName: srcAccount.Name,
					DestAccountID:     *txn.DestinationID,
					DestAccountName:   destAccount.Name,
				}
			}
			pairMap[pairKey].Count++
			pairMap[pairKey].TotalAmount = pairMap[pairKey].TotalAmount.Add(txn.Amount)
		}
	}

	frequentPairs := make([]TransferPair, 0)
	for _, pair := range pairMap {
		frequentPairs = append(frequentPairs, *pair)
	}

	return &TransferInsightData{
		Total:         total,
		Count:         count,
		FrequentPairs: frequentPairs,
	}, nil
}

// calculateTrend 根据按日期的数据计算趋势方向
// 将日期数据分成前后两半，比较后半段总额与前半段总额
// 如果后半段 > 前半段 → "up"，后半段 < 前半段 → "down"，相等 → "stable"
// 参数：
//   - byDate: 按日期统计的数据列表
// 返回：
//   - string: 趋势方向（"up"/"down"/"stable"）
func calculateTrend(byDate []DateInsight) string {
	if len(byDate) < 2 {
		return "stable"
	}

	// 将数据分成前后两半
	mid := len(byDate) / 2
	firstHalf := decimal.Zero
	secondHalf := decimal.Zero

	for i, d := range byDate {
		if i < mid {
			firstHalf = firstHalf.Add(d.Amount)
		} else {
			secondHalf = secondHalf.Add(d.Amount)
		}
	}

	if secondHalf.GreaterThan(firstHalf) {
		return "up"
	}
	if secondHalf.LessThan(firstHalf) {
		return "down"
	}
	return "stable"
}

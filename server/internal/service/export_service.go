// Package service 业务逻辑层，实现核心业务逻辑
// ExportService 数据导出业务逻辑，支持导出为CSV和JSON格式
package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/repository"
)

// ExportService 数据导出服务
// 负责将各类数据导出为CSV或JSON格式，支持导出交易、账户、循环交易、预算、分类、标签、存钱罐和规则
// 依赖各repository获取数据，不依赖service层避免循环依赖
type ExportService struct {
	txnRepo       *repository.TransactionRepository
	accountRepo   *repository.AccountRepository
	rtRepo        *repository.RecurringTransactionRepository
	budgetRepo    *repository.BudgetRepository
	categoryRepo  *repository.CategoryRepository
	tagRepo       *repository.TagRepository
	piggyBankRepo *repository.PiggyBankRepository
	ruleRepo      *repository.RuleRepository
}

func NewExportService(
	txnRepo *repository.TransactionRepository,
	accountRepo *repository.AccountRepository,
	rtRepo *repository.RecurringTransactionRepository,
	budgetRepo *repository.BudgetRepository,
	categoryRepo *repository.CategoryRepository,
	tagRepo *repository.TagRepository,
	piggyBankRepo *repository.PiggyBankRepository,
	ruleRepo *repository.RuleRepository,
) *ExportService {
	return &ExportService{
		txnRepo:       txnRepo,
		accountRepo:   accountRepo,
		rtRepo:        rtRepo,
		budgetRepo:    budgetRepo,
		categoryRepo:  categoryRepo,
		tagRepo:       tagRepo,
		piggyBankRepo: piggyBankRepo,
		ruleRepo:      ruleRepo,
	}
}

// ExportTransactions 导出交易数据
// 参数：
//   - userID: 用户ID
//   - startDate: 开始日期（格式：2006-01-02）
//   - endDate: 结束日期
//   - format: 导出格式（csv/json）
// 返回：
//   - []byte: 导出数据
//   - string: 文件名
//   - error: 错误信息
func (s *ExportService) ExportTransactions(userID uint64, startDate, endDate, format string) ([]byte, string, error) {
	filter := repository.TransactionFilter{
		StartDate: startDate,
		EndDate:   endDate,
	}

	txns, err := s.txnRepo.ListAll(userID, filter)
	if err != nil {
		return nil, "", err
	}

	switch format {
	case "csv":
		return s.exportTransactionsCSV(txns)
	case "json":
		return s.exportTransactionsJSON(txns)
	default:
		return s.exportTransactionsCSV(txns)
	}
}

// ExportAccounts 导出账户数据
// 参数：
//   - userID: 用户ID
//   - format: 导出格式（csv/json）
// 返回：
//   - []byte: 导出数据
//   - string: 文件名
//   - error: 错误信息
func (s *ExportService) ExportAccounts(userID uint64, format string) ([]byte, string, error) {
	accounts, err := s.accountRepo.List(userID, "", "", "name", 0, 1000)
	if err != nil {
		return nil, "", err
	}

	switch format {
	case "csv":
		return s.exportAccountsCSV(accounts)
	case "json":
		return s.exportAccountsJSON(accounts)
	default:
		return s.exportAccountsCSV(accounts)
	}
}

// exportTransactionsCSV 将交易列表导出为CSV格式
func (s *ExportService) exportTransactionsCSV(txns []model.Transaction) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"date", "type", "description", "amount", "source_account", "destination_account", "category", "tags", "notes"}
	writer.Write(header)

	for _, txn := range txns {
		sourceName := ""
		if txn.Source.ID != 0 {
			sourceName = txn.Source.Name
		}
		destName := ""
		if txn.Destination != nil && txn.Destination.ID != 0 {
			destName = txn.Destination.Name
		}
		categoryName := ""
		if txn.CategoryID != nil && txn.Category != nil {
			categoryName = txn.Category.Name
		}
		tags := ""
		for i, t := range txn.Tags {
			if i > 0 {
				tags += ","
			}
			tags += t.Name
		}

		row := []string{
			txn.Date.Format("2006-01-02"),
			string(txn.Type),
			txn.Description,
			txn.Amount.StringFixed(4),
			sourceName,
			destName,
			categoryName,
			tags,
			txn.Notes,
		}
		writer.Write(row)
	}

	writer.Flush()
	filename := fmt.Sprintf("transactions_%s.csv", time.Now().Format("20060102"))
	return buf.Bytes(), filename, nil
}

// exportTransactionsJSON 将交易列表导出为JSON格式
func (s *ExportService) exportTransactionsJSON(txns []model.Transaction) ([]byte, string, error) {
	type txnExport struct {
		Date              string `json:"date"`
		Type              string `json:"type"`
		Description       string `json:"description"`
		Amount            string `json:"amount"`
		SourceAccount     string `json:"source_account"`
		DestinationAccount string `json:"destination_account"`
		Category          string `json:"category"`
		Tags              string `json:"tags"`
		Notes             string `json:"notes"`
	}

	items := make([]txnExport, 0, len(txns))
	for _, txn := range txns {
		sourceName := ""
		if txn.Source.ID != 0 {
			sourceName = txn.Source.Name
		}
		destName := ""
		if txn.Destination != nil && txn.Destination.ID != 0 {
			destName = txn.Destination.Name
		}
		categoryName := ""
		if txn.CategoryID != nil && txn.Category != nil {
			categoryName = txn.Category.Name
		}
		tags := ""
		for i, t := range txn.Tags {
			if i > 0 {
				tags += ","
			}
			tags += t.Name
		}

		items = append(items, txnExport{
			Date:              txn.Date.Format("2006-01-02"),
			Type:              string(txn.Type),
			Description:       txn.Description,
			Amount:            txn.Amount.StringFixed(4),
			SourceAccount:     sourceName,
			DestinationAccount: destName,
			Category:          categoryName,
			Tags:              tags,
			Notes:             txn.Notes,
		})
	}

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return nil, "", err
	}
	filename := fmt.Sprintf("transactions_%s.json", time.Now().Format("20060102"))
	return data, filename, nil
}

// exportAccountsCSV 将账户列表导出为CSV格式
func (s *ExportService) exportAccountsCSV(accounts []model.Account) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"name", "type", "currency", "initial_balance", "current_balance", "is_virtual", "notes"}
	writer.Write(header)

	for _, a := range accounts {
		currencyCode := ""
		if a.Currency.ID != 0 {
			currencyCode = a.Currency.Code
		}
		row := []string{
			a.Name,
			string(a.Type),
			currencyCode,
			a.InitialBalance.StringFixed(4),
			a.CurrentBalance.StringFixed(4),
			fmt.Sprintf("%v", a.IsVirtual),
			a.Notes,
		}
		writer.Write(row)
	}

	writer.Flush()
	filename := fmt.Sprintf("accounts_%s.csv", time.Now().Format("20060102"))
	return buf.Bytes(), filename, nil
}

// exportAccountsJSON 将账户列表导出为JSON格式
func (s *ExportService) exportAccountsJSON(accounts []model.Account) ([]byte, string, error) {
	type accountExport struct {
		Name           string `json:"name"`
		Type           string `json:"type"`
		Currency       string `json:"currency"`
		InitialBalance string `json:"initial_balance"`
		CurrentBalance string `json:"current_balance"`
		IsVirtual      bool   `json:"is_virtual"`
		Notes          string `json:"notes"`
	}

	items := make([]accountExport, 0, len(accounts))
	for _, a := range accounts {
		currencyCode := ""
		if a.Currency.ID != 0 {
			currencyCode = a.Currency.Code
		}
		items = append(items, accountExport{
			Name:           a.Name,
			Type:           string(a.Type),
			Currency:       currencyCode,
			InitialBalance: a.InitialBalance.StringFixed(4),
			CurrentBalance: a.CurrentBalance.StringFixed(4),
			IsVirtual:      a.IsVirtual,
			Notes:          a.Notes,
		})
	}

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return nil, "", err
	}
	filename := fmt.Sprintf("accounts_%s.json", time.Now().Format("20060102"))
	return data, filename, nil
}

// ExportRecurringTransactions 导出循环交易数据
// 参数：
//   - userID: 用户ID
//   - format: 导出格式（csv/json）
// 返回：
//   - []byte: 导出数据
//   - string: 文件名
//   - error: 错误信息
func (s *ExportService) ExportRecurringTransactions(userID uint64, format string) ([]byte, string, error) {
	rts, err := s.rtRepo.List(userID, nil, 0, 10000)
	if err != nil {
		return nil, "", err
	}

	switch format {
	case "csv":
		return s.exportRecurringTxnsCSV(rts)
	case "json":
		return s.exportRecurringTxnsJSON(rts)
	default:
		return s.exportRecurringTxnsCSV(rts)
	}
}

// ExportBudgets 导出预算数据
// 参数：
//   - userID: 用户ID
//   - format: 导出格式（csv/json）
// 返回：
//   - []byte: 导出数据
//   - string: 文件名
//   - error: 错误信息
func (s *ExportService) ExportBudgets(userID uint64, format string) ([]byte, string, error) {
	budgets, err := s.budgetRepo.List(userID)
	if err != nil {
		return nil, "", err
	}

	switch format {
	case "csv":
		return s.exportBudgetsCSV(budgets)
	case "json":
		return s.exportBudgetsJSON(budgets)
	default:
		return s.exportBudgetsCSV(budgets)
	}
}

// ExportCategories 导出分类数据
// 参数：
//   - userID: 用户ID
//   - format: 导出格式（csv/json）
// 返回：
//   - []byte: 导出数据
//   - string: 文件名
//   - error: 错误信息
func (s *ExportService) ExportCategories(userID uint64, format string) ([]byte, string, error) {
	categories, err := s.categoryRepo.List(userID)
	if err != nil {
		return nil, "", err
	}

	switch format {
	case "csv":
		return s.exportCategoriesCSV(categories)
	case "json":
		return s.exportCategoriesJSON(categories)
	default:
		return s.exportCategoriesCSV(categories)
	}
}

// ExportTags 导出标签数据
// 参数：
//   - userID: 用户ID
//   - format: 导出格式（csv/json）
// 返回：
//   - []byte: 导出数据
//   - string: 文件名
//   - error: 错误信息
func (s *ExportService) ExportTags(userID uint64, format string) ([]byte, string, error) {
	tags, err := s.tagRepo.List(userID)
	if err != nil {
		return nil, "", err
	}

	switch format {
	case "csv":
		return s.exportTagsCSV(tags)
	case "json":
		return s.exportTagsJSON(tags)
	default:
		return s.exportTagsCSV(tags)
	}
}

// ExportPiggyBanks 导出存钱罐数据
// 参数：
//   - userID: 用户ID
//   - format: 导出格式（csv/json）
// 返回：
//   - []byte: 导出数据
//   - string: 文件名
//   - error: 错误信息
func (s *ExportService) ExportPiggyBanks(userID uint64, format string) ([]byte, string, error) {
	piggyBanks, err := s.piggyBankRepo.List(userID)
	if err != nil {
		return nil, "", err
	}

	switch format {
	case "csv":
		return s.exportPiggyBanksCSV(piggyBanks)
	case "json":
		return s.exportPiggyBanksJSON(piggyBanks)
	default:
		return s.exportPiggyBanksCSV(piggyBanks)
	}
}

// ExportRules 导出规则数据
// 参数：
//   - userID: 用户ID
//   - format: 导出格式（csv/json）
// 返回：
//   - []byte: 导出数据
//   - string: 文件名
//   - error: 错误信息
func (s *ExportService) ExportRules(userID uint64, format string) ([]byte, string, error) {
	rules, err := s.ruleRepo.List(userID)
	if err != nil {
		return nil, "", err
	}

	switch format {
	case "csv":
		return s.exportRulesCSV(rules)
	case "json":
		return s.exportRulesJSON(rules)
	default:
		return s.exportRulesCSV(rules)
	}
}

// Recurring transactions export helpers
func (s *ExportService) exportRecurringTxnsCSV(rts []model.RecurringTransaction) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Write([]string{"description", "amount", "recurrence_type", "repeat_every", "next_occurrence", "is_active", "notes"})
	for _, rt := range rts {
		writer.Write([]string{rt.Description, rt.Amount.StringFixed(4), string(rt.RecurrenceType), fmt.Sprintf("%d", rt.RepeatEvery), rt.NextOccurrence.Format("2006-01-02"), fmt.Sprintf("%v", rt.IsActive), rt.Notes})
	}
	writer.Flush()
	return buf.Bytes(), fmt.Sprintf("recurring_transactions_%s.csv", time.Now().Format("20060102")), nil
}

func (s *ExportService) exportRecurringTxnsJSON(rts []model.RecurringTransaction) ([]byte, string, error) {
	type rtExport struct {
		Description    string `json:"description"`
		Amount         string `json:"amount"`
		RecurrenceType string `json:"recurrence_type"`
		RepeatEvery    int    `json:"repeat_every"`
		NextOccurrence string `json:"next_occurrence"`
		IsActive       bool   `json:"is_active"`
		Notes          string `json:"notes"`
	}
	items := make([]rtExport, 0, len(rts))
	for _, rt := range rts {
		items = append(items, rtExport{Description: rt.Description, Amount: rt.Amount.StringFixed(4), RecurrenceType: string(rt.RecurrenceType), RepeatEvery: rt.RepeatEvery, NextOccurrence: rt.NextOccurrence.Format("2006-01-02"), IsActive: rt.IsActive, Notes: rt.Notes})
	}
	data, _ := json.MarshalIndent(items, "", "  ")
	return data, fmt.Sprintf("recurring_transactions_%s.json", time.Now().Format("20060102")), nil
}

// Budgets export helpers
func (s *ExportService) exportBudgetsCSV(budgets []model.Budget) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Write([]string{"name", "amount", "period", "is_enabled"})
	for _, b := range budgets {
		writer.Write([]string{b.Name, b.Amount.StringFixed(4), string(b.Period), fmt.Sprintf("%v", b.IsEnabled)})
	}
	writer.Flush()
	return buf.Bytes(), fmt.Sprintf("budgets_%s.csv", time.Now().Format("20060102")), nil
}

func (s *ExportService) exportBudgetsJSON(budgets []model.Budget) ([]byte, string, error) {
	type budgetExport struct {
		Name      string `json:"name"`
		Amount    string `json:"amount"`
		Period    string `json:"period"`
		IsEnabled bool   `json:"is_enabled"`
	}
	items := make([]budgetExport, 0, len(budgets))
	for _, b := range budgets {
		items = append(items, budgetExport{Name: b.Name, Amount: b.Amount.StringFixed(4), Period: string(b.Period), IsEnabled: b.IsEnabled})
	}
	data, _ := json.MarshalIndent(items, "", "  ")
	return data, fmt.Sprintf("budgets_%s.json", time.Now().Format("20060102")), nil
}

// Categories export helpers
func (s *ExportService) exportCategoriesCSV(categories []model.Category) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Write([]string{"name", "icon", "notes"})
	for _, c := range categories {
		writer.Write([]string{c.Name, c.Icon, c.Notes})
	}
	writer.Flush()
	return buf.Bytes(), fmt.Sprintf("categories_%s.csv", time.Now().Format("20060102")), nil
}

func (s *ExportService) exportCategoriesJSON(categories []model.Category) ([]byte, string, error) {
	type categoryExport struct {
		Name  string `json:"name"`
		Icon  string `json:"icon"`
		Notes string `json:"notes"`
	}
	items := make([]categoryExport, 0, len(categories))
	for _, c := range categories {
		items = append(items, categoryExport{Name: c.Name, Icon: c.Icon, Notes: c.Notes})
	}
	data, _ := json.MarshalIndent(items, "", "  ")
	return data, fmt.Sprintf("categories_%s.json", time.Now().Format("20060102")), nil
}

// Tags export helpers
func (s *ExportService) exportTagsCSV(tags []model.Tag) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Write([]string{"name", "color"})
	for _, t := range tags {
		writer.Write([]string{t.Name, t.Color})
	}
	writer.Flush()
	return buf.Bytes(), fmt.Sprintf("tags_%s.csv", time.Now().Format("20060102")), nil
}

func (s *ExportService) exportTagsJSON(tags []model.Tag) ([]byte, string, error) {
	type tagExport struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	items := make([]tagExport, 0, len(tags))
	for _, t := range tags {
		items = append(items, tagExport{Name: t.Name, Color: t.Color})
	}
	data, _ := json.MarshalIndent(items, "", "  ")
	return data, fmt.Sprintf("tags_%s.json", time.Now().Format("20060102")), nil
}

// Piggy Banks export helpers
func (s *ExportService) exportPiggyBanksCSV(piggyBanks []model.PiggyBank) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Write([]string{"name", "target_amount", "current_amount", "notes"})
	for _, p := range piggyBanks {
		writer.Write([]string{p.Name, p.TargetAmount.StringFixed(4), p.CurrentAmount.StringFixed(4), p.Notes})
	}
	writer.Flush()
	return buf.Bytes(), fmt.Sprintf("piggy_banks_%s.csv", time.Now().Format("20060102")), nil
}

func (s *ExportService) exportPiggyBanksJSON(piggyBanks []model.PiggyBank) ([]byte, string, error) {
	type piggyBankExport struct {
		Name          string `json:"name"`
		TargetAmount  string `json:"target_amount"`
		CurrentAmount string `json:"current_amount"`
		Notes         string `json:"notes"`
	}
	items := make([]piggyBankExport, 0, len(piggyBanks))
	for _, p := range piggyBanks {
		items = append(items, piggyBankExport{Name: p.Name, TargetAmount: p.TargetAmount.StringFixed(4), CurrentAmount: p.CurrentAmount.StringFixed(4), Notes: p.Notes})
	}
	data, _ := json.MarshalIndent(items, "", "  ")
	return data, fmt.Sprintf("piggy_banks_%s.json", time.Now().Format("20060102")), nil
}

// Rules export helpers
func (s *ExportService) exportRulesCSV(rules []model.Rule) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Write([]string{"name", "priority", "is_enabled", "logic_type", "trigger"})
	for _, r := range rules {
		writer.Write([]string{r.Name, fmt.Sprintf("%d", r.Priority), fmt.Sprintf("%v", r.IsEnabled), string(r.LogicType), string(r.Trigger)})
	}
	writer.Flush()
	return buf.Bytes(), fmt.Sprintf("rules_%s.csv", time.Now().Format("20060102")), nil
}

func (s *ExportService) exportRulesJSON(rules []model.Rule) ([]byte, string, error) {
	type ruleExport struct {
		Name      string `json:"name"`
		Priority  int    `json:"priority"`
		IsEnabled bool   `json:"is_enabled"`
		LogicType string `json:"logic_type"`
		Trigger   string `json:"trigger"`
	}
	items := make([]ruleExport, 0, len(rules))
	for _, r := range rules {
		items = append(items, ruleExport{Name: r.Name, Priority: r.Priority, IsEnabled: r.IsEnabled, LogicType: string(r.LogicType), Trigger: string(r.Trigger)})
	}
	data, _ := json.MarshalIndent(items, "", "  ")
	return data, fmt.Sprintf("rules_%s.json", time.Now().Format("20060102")), nil
}

// Package service 业务逻辑层，实现核心业务逻辑
// CronService 定时任务业务逻辑，处理账单提醒和循环交易
package service

import (
	"time"

	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/repository"
)

// CronService 定时任务服务
// 依赖recurrenceService执行循环交易，依赖billService/billRepo处理到期账单
// 依赖recurrenceRepo获取到期循环交易，依赖txnService创建交易
type CronService struct {
	recurrenceService *RecurrenceService
	billService       *BillService
	billRepo          *repository.BillRepository
	recurrenceRepo    *repository.RecurrenceRepository
	txnService        *TransactionService
	budgetService     *BudgetService
}

// NewCronService 创建定时任务服务实例
// 参数：
//   - recurrenceService: 循环交易服务，用于执行到期的循环交易
//   - billService: 账单服务，用于处理到期账单
//   - billRepo: 账单数据访问对象，用于获取所有用户的到期账单
//   - recurrenceRepo: 循环交易数据访问对象，用于获取到期循环交易
//   - txnService: 交易服务，用于从账单创建交易（确保余额更新）
// 返回：
//   - *CronService: 定时任务服务实例
func NewCronService(
	recurrenceService *RecurrenceService,
	billService *BillService,
	billRepo *repository.BillRepository,
	recurrenceRepo *repository.RecurrenceRepository,
	txnService *TransactionService,
	budgetService *BudgetService,
) *CronService {
	return &CronService{
		recurrenceService: recurrenceService,
		billService:       billService,
		billRepo:          billRepo,
		recurrenceRepo:    recurrenceRepo,
		txnService:        txnService,
		budgetService:     budgetService,
	}
}

// CronResult 定时任务执行结果
type CronResult struct {
	RecurrencesExecuted int      `json:"recurrences_executed"`
	RecurrenceErrors    []string `json:"recurrence_errors,omitempty"`
	BillsExecuted       int      `json:"bills_executed"`
	BillErrors          []string `json:"bill_errors,omitempty"`
	BudgetSnapshots     int      `json:"budget_snapshots"`
	BudgetErrors        []string `json:"budget_errors,omitempty"`
}

// CronRun 执行所有到期的循环交易和账单
// 修复：使用GetAllDueBills()替代GetUpcoming(0,0)，避免userID=0导致的数据隔离问题
func (s *CronService) CronRun() (*CronResult, error) {
	result := &CronResult{
		RecurrenceErrors: make([]string, 0),
		BillErrors:       make([]string, 0),
		BudgetErrors:     make([]string, 0),
	}

	// 执行到期的循环交易
	today := time.Now()
	dueRecurrences, err := s.recurrenceRepo.GetDueRecurrences(today)
	if err != nil {
		return nil, err
	}

	for _, rec := range dueRecurrences {
		if err := s.recurrenceService.ExecuteRecurrence(&rec); err != nil {
			result.RecurrenceErrors = append(result.RecurrenceErrors, err.Error())
		} else {
			result.RecurrencesExecuted++
		}
	}

	// 执行到期的账单：使用GetAllDueBills获取所有用户的到期账单
	// 原代码使用GetUpcoming(0,0)传入userID=0，由于repo的WHERE条件会过滤user_id=0，
	// 导致无法获取任何用户的账单。GetAllDueBills不按用户过滤，适合定时任务批量处理
	dueBills, err := s.billRepo.GetAllDueBills()
	if err != nil {
		result.BillErrors = append(result.BillErrors, err.Error())
	} else {
		for i := range dueBills {
			if err := s.createTransactionFromBill(&dueBills[i]); err != nil {
				result.BillErrors = append(result.BillErrors, err.Error())
			} else {
				result.BillsExecuted++
			}
		}
	}

	// 生成预算快照：为所有已启用预算写入当前周期的 BudgetHistory
	snapshots, budgetErrs := s.budgetService.SnapshotCurrentPeriod()
	result.BudgetSnapshots = snapshots
	result.BudgetErrors = budgetErrs

	return result, nil
}

// createTransactionFromBill 从账单创建交易并更新账单的下次到期日期
// 只处理到期日期在今天或之前的账单，通过txnService.Create()确保账户余额更新
func (s *CronService) createTransactionFromBill(bill *model.Bill) error {
	// 只处理到期日期在今天或之前的账单
	if bill.NextDue.After(time.Now()) {
		return nil
	}

	// 如果账单指定了支出账户，创建对应的支出交易
	// 通过txnService.Create()确保：1)账户余额正确更新 2)规则触发 3)Webhook通知
	if bill.SourceID != nil {
		txnReq := &request.CreateTransactionReq{
			Type:        string(model.TransactionTypeWithdrawal),
			Date:        bill.NextDue.Format("2006-01-02"),
			Description: bill.Name,
			Amount:      bill.Amount.StringFixed(4),
			SourceID:    *bill.SourceID,
			CategoryID:  bill.CategoryID,
			Notes:       bill.Notes,
		}

		_, err := s.txnService.Create(bill.UserID, txnReq)
		if err != nil {
			return err
		}
	}

	// 根据重复规则更新下次到期日期
	nextDue := s.calculateBillNextDue(bill.NextDue, bill.RepeatRule)
	bill.NextDue = nextDue

	return s.billRepo.Update(bill)
}

// calculateBillNextDue 根据重复规则计算账单的下次到期日期
func (s *CronService) calculateBillNextDue(currentDue time.Time, rule model.RepeatRule) time.Time {
	switch rule {
	case model.RepeatRuleDaily:
		return currentDue.AddDate(0, 0, 1)
	case model.RepeatRuleWeekly:
		return currentDue.AddDate(0, 0, 7)
	case model.RepeatRuleMonthly:
		return currentDue.AddDate(0, 1, 0)
	case model.RepeatRuleYearly:
		return currentDue.AddDate(1, 0, 0)
	default:
		return currentDue.AddDate(0, 1, 0)
	}
}

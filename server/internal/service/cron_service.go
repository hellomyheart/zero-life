// Package service 业务逻辑层，实现核心业务逻辑
// CronService 定时任务业务逻辑，处理账单提醒和循环交易
package service

import (
	"time"

	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/repository"
)

type CronService struct {
	recurrenceService *RecurrenceService
	billService       *BillService
	billRepo          *repository.BillRepository
	recurrenceRepo    *repository.RecurrenceRepository
	txnService        *TransactionService
}

func NewCronService(
	recurrenceService *RecurrenceService,
	billService *BillService,
	billRepo *repository.BillRepository,
	recurrenceRepo *repository.RecurrenceRepository,
	txnService *TransactionService,
) *CronService {
	return &CronService{
		recurrenceService: recurrenceService,
		billService:       billService,
		billRepo:          billRepo,
		recurrenceRepo:    recurrenceRepo,
		txnService:        txnService,
	}
}

type CronResult struct {
	RecurrencesExecuted int      `json:"recurrences_executed"`
	RecurrenceErrors    []string `json:"recurrence_errors,omitempty"`
	BillsExecuted       int      `json:"bills_executed"`
	BillErrors          []string `json:"bill_errors,omitempty"`
}

// CronRun executes all due recurrences and bills.
func (s *CronService) CronRun() (*CronResult, error) {
	result := &CronResult{
		RecurrenceErrors: make([]string, 0),
		BillErrors:       make([]string, 0),
	}

	// Execute due recurrences
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

	// Execute due bills (get bills due within 0 days = today)
	dueBills, err := s.billRepo.GetUpcoming(0, 0)
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

	return result, nil
}

// createTransactionFromBill creates a transaction from a bill and updates the bill's next due date.
func (s *CronService) createTransactionFromBill(bill *model.Bill) error {
	// Only process bills that are due today or earlier
	if bill.NextDue.After(time.Now()) {
		return nil
	}

	// Create transaction from bill
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

	// Update bill's next due date
	nextDue := s.calculateBillNextDue(bill.NextDue, bill.RepeatRule)
	bill.NextDue = nextDue

	return s.billRepo.Update(bill)
}

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

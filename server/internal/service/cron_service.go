package service

import (
	"time"

	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type CronTask struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Schedule    string `json:"schedule"`
}

type CronTaskResult struct {
	TaskID  string   `json:"task_id"`
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

type CronService struct {
	recurrenceService *RecurrenceService
	billService       *BillService
	billRepo          *repository.BillRepository
	recurrenceRepo    *repository.RecurrenceRepository
	txnService        *TransactionService
	budgetService     *BudgetService
	cron              *cron.Cron
	logger            *zap.Logger
	tasks             []CronTask
}

func NewCronService(
	recurrenceService *RecurrenceService,
	billService *BillService,
	billRepo *repository.BillRepository,
	recurrenceRepo *repository.RecurrenceRepository,
	txnService *TransactionService,
	budgetService *BudgetService,
	logger *zap.Logger,
) *CronService {
	s := &CronService{
		recurrenceService: recurrenceService,
		billService:       billService,
		billRepo:          billRepo,
		recurrenceRepo:    recurrenceRepo,
		txnService:        txnService,
		budgetService:     budgetService,
		cron:              cron.New(cron.WithSeconds()),
		logger:            logger,
	}
	s.tasks = []CronTask{
		{ID: "recurrences_bills", Name: "循环交易与到期账单", Description: "执行到期的循环交易和到期账单，创建对应交易记录", Schedule: "每天 00:00"},
		{ID: "budget_snapshot", Name: "预算历史快照", Description: "为所有启用预算生成近2年内已结束周期的历史快照（支持重复跑）", Schedule: "每天 09:00"},
	}
	return s
}

func (s *CronService) ListTasks() []CronTask {
	return s.tasks
}

func (s *CronService) RunTask(taskID string) (*CronTaskResult, error) {
	switch taskID {
	case "recurrences_bills":
		return s.runRecurrencesAndBillsTask(), nil
	case "budget_snapshot":
		return s.runBudgetSnapshotTask(), nil
	default:
		return nil, ErrCronTaskNotFound
	}
}

func (s *CronService) runRecurrencesAndBillsTask() *CronTaskResult {
	result := &CronTaskResult{
		TaskID: "recurrences_bills",
		Errors: make([]string, 0),
	}

	today := time.Now()
	dueRecurrences, err := s.recurrenceRepo.GetDueRecurrences(today)
	if err != nil {
		result.Errors = append(result.Errors, err.Error())
	} else {
		for _, rec := range dueRecurrences {
			if err := s.recurrenceService.ExecuteRecurrence(&rec); err != nil {
				result.Errors = append(result.Errors, err.Error())
			}
		}
	}

	dueBills, err := s.billRepo.GetAllDueBills()
	if err != nil {
		result.Errors = append(result.Errors, err.Error())
	} else {
		for i := range dueBills {
			if err := s.createTransactionFromBill(&dueBills[i]); err != nil {
				result.Errors = append(result.Errors, err.Error())
			}
		}
	}

	result.Success = len(result.Errors) == 0
	if result.Success {
		result.Message = "执行成功"
	} else {
		result.Message = "执行完成，部分错误"
	}
	return result
}

func (s *CronService) runBudgetSnapshotTask() *CronTaskResult {
	result := &CronTaskResult{
		TaskID: "budget_snapshot",
		Errors: make([]string, 0),
	}

	created, errs := s.budgetService.SnapshotHistory()
	result.Errors = errs
	result.Success = len(errs) == 0
	if result.Success {
		result.Message = "执行成功"
	} else {
		result.Message = "执行完成，部分错误"
	}
	_ = created
	return result
}

var ErrCronTaskNotFound = &cronTaskNotFoundError{}

type cronTaskNotFoundError struct{}

func (e *cronTaskNotFoundError) Error() string { return "cron task not found" }

func (s *CronService) StartScheduler() {
	if _, err := s.cron.AddFunc("0 0 0 * * *", func() {
		s.logger.Info("cron: daily job started - recurrences and bills")
		s.runRecurrencesAndBillsTask()
		s.logger.Info("cron: daily job finished")
	}); err != nil {
		s.logger.Error("cron: failed to register daily job", zap.Error(err))
	}

	if _, err := s.cron.AddFunc("0 0 9 * * *", func() {
		s.logger.Info("cron: budget history snapshot job started")
		s.runBudgetSnapshotTask()
		s.logger.Info("cron: budget history snapshot job finished")
	}); err != nil {
		s.logger.Error("cron: failed to register budget snapshot job", zap.Error(err))
	}

	s.cron.Start()
	s.logger.Info("cron: scheduler started - daily 00:00 (recurrences+bills), daily 09:00 (budget history snapshot)")
}

func (s *CronService) StopScheduler() {
	if s.cron != nil {
		ctx := s.cron.Stop()
		<-ctx.Done()
		s.logger.Info("cron: scheduler stopped")
	}
}

func (s *CronService) createTransactionFromBill(bill *model.Bill) error {
	if bill.NextDue.After(time.Now()) {
		return nil
	}

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

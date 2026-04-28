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
	rtRepo        *repository.RecurringTransactionRepository
	rtService     *RecurringTransactionService
	budgetService *BudgetService
	cron          *cron.Cron
	logger        *zap.Logger
	tasks         []CronTask
}

func NewCronService(
	rtRepo *repository.RecurringTransactionRepository,
	rtService *RecurringTransactionService,
	budgetService *BudgetService,
	logger *zap.Logger,
) *CronService {
	s := &CronService{
		rtRepo:        rtRepo,
		rtService:     rtService,
		budgetService: budgetService,
		cron:          cron.New(cron.WithSeconds()),
		logger:        logger,
	}
	s.tasks = []CronTask{
		{ID: "recurring_transactions", Name: "到期循环交易", Description: "执行到期的循环交易，创建对应交易记录", Schedule: "每天 00:00"},
		{ID: "budget_snapshot", Name: "预算历史快照", Description: "为所有启用预算生成近2年内已结束周期的历史快照（支持重复跑）", Schedule: "每天 09:00"},
	}
	return s
}

func (s *CronService) ListTasks() []CronTask {
	return s.tasks
}

func (s *CronService) RunTask(taskID string) (*CronTaskResult, error) {
	switch taskID {
	case "recurring_transactions":
		return s.runRecurringTransactionsTask(), nil
	case "budget_snapshot":
		return s.runBudgetSnapshotTask(), nil
	default:
		return nil, ErrCronTaskNotFound
	}
}

func (s *CronService) runRecurringTransactionsTask() *CronTaskResult {
	result := &CronTaskResult{
		TaskID: "recurring_transactions",
		Errors: make([]string, 0),
	}

	dueRTs, err := s.rtRepo.GetAllDue()
	if err != nil {
		result.Errors = append(result.Errors, err.Error())
	} else {
		for i := range dueRTs {
			if err := s.createTransactionFromRecurring(&dueRTs[i]); err != nil {
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

func (s *CronService) createTransactionFromRecurring(rt *model.RecurringTransaction) error {
	if rt.NextOccurrence.After(time.Now()) {
		return nil
	}

	txnType := determineTransactionType(rt)
	txnReq := &request.CreateTransactionReq{
		Type:          string(txnType),
		Date:          rt.NextOccurrence.Format("2006-01-02"),
		Description:   rt.Description,
		Amount:        rt.Amount.StringFixed(4),
		SourceID:      rt.SourceID,
		DestinationID: rt.DestinationID,
		CategoryID:    rt.CategoryID,
		Notes:         rt.Notes,
		RecurringID:   &rt.ID,
	}

	txnResp, err := s.rtService.txnService.Create(rt.UserID, txnReq)
	if err != nil {
		return err
	}

	log := &model.RecurringTransactionLog{
		RecurringTransactionID: rt.ID,
		TransactionID:          txnResp.ID,
		OccurrenceDate:         rt.NextOccurrence,
	}
	if err := s.rtRepo.CreateLog(log); err != nil {
		return err
	}

	rt.NextOccurrence = s.calculateNextOccurrence(rt.NextOccurrence, rt.RecurrenceType, rt.RepeatEvery)
	if rt.EndDate != nil && rt.NextOccurrence.After(*rt.EndDate) {
		rt.IsActive = false
	}

	return s.rtRepo.Update(rt)
}

func (s *CronService) calculateNextOccurrence(from time.Time, recurrenceType model.RecurrenceType, repeatEvery int) time.Time {
	switch recurrenceType {
	case model.RecurrenceTypeDaily:
		return from.AddDate(0, 0, repeatEvery)
	case model.RecurrenceTypeWeekly:
		return from.AddDate(0, 0, 7*repeatEvery)
	case model.RecurrenceTypeMonthly:
		return from.AddDate(0, repeatEvery, 0)
	case model.RecurrenceTypeYearly:
		return from.AddDate(repeatEvery, 0, 0)
	default:
		return from.AddDate(0, 1, 0)
	}
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
		s.logger.Info("cron: daily job started - recurring transactions")
		s.runRecurringTransactionsTask()
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
	s.logger.Info("cron: scheduler started - daily 00:00 (recurring transactions), daily 09:00 (budget history snapshot)")
}

func (s *CronService) StopScheduler() {
	if s.cron != nil {
		ctx := s.cron.Stop()
		<-ctx.Done()
		s.logger.Info("cron: scheduler stopped")
	}
}

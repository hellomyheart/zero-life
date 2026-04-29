package service

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	rtService     *RecurringTransactionService
	budgetService *BudgetService
	writeDB       *gorm.DB
	cron          *cron.Cron
	logger        *zap.Logger
	tasks         []CronTask
}

func NewCronService(
	rtService *RecurringTransactionService,
	budgetService *BudgetService,
	writeDB *gorm.DB,
	logger *zap.Logger,
) *CronService {
	s := &CronService{
		rtService:     rtService,
		budgetService: budgetService,
		writeDB:       writeDB,
		cron:          cron.New(cron.WithSeconds()),
		logger:        logger,
	}
	s.tasks = []CronTask{
		{ID: "recurring_transactions", Name: "到期循环交易", Description: "执行到期的循环交易，创建对应交易记录", Schedule: "每天 00:00"},
		{ID: "budget_snapshot", Name: "预算历史快照", Description: "为所有启用预算生成近2年内已结束周期的历史快照（支持重复跑）", Schedule: "每天 09:00"},
		{ID: "wal_checkpoint", Name: "WAL检查点", Description: "执行 SQLite WAL 检查点，将 WAL 文件合并回主数据库，控制文件大小", Schedule: "每6小时"},
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
	case "wal_checkpoint":
		return s.runWALCheckpointTask(), nil
	default:
		return nil, ErrCronTaskNotFound
	}
}

func (s *CronService) runRecurringTransactionsTask() *CronTaskResult {
	result := &CronTaskResult{
		TaskID: "recurring_transactions",
		Errors: make([]string, 0),
	}

	errs := s.rtService.ProcessAllDue()
	result.Errors = errs
	result.Success = len(errs) == 0
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

func (s *CronService) runWALCheckpointTask() *CronTaskResult {
	result := &CronTaskResult{
		TaskID: "wal_checkpoint",
		Errors: make([]string, 0),
	}

	if err := s.writeDB.Exec("PRAGMA wal_checkpoint(TRUNCATE)").Error; err != nil {
		result.Success = false
		result.Message = "WAL检查点执行失败"
		result.Errors = append(result.Errors, err.Error())
	} else {
		result.Success = true
		result.Message = "WAL检查点执行成功"
	}
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

	// 每6小时执行一次 WAL 检查点，将 WAL 文件合并回主数据库
	if _, err := s.cron.AddFunc("0 0 */6 * * *", func() {
		s.logger.Info("cron: WAL checkpoint job started")
		s.runWALCheckpointTask()
		s.logger.Info("cron: WAL checkpoint job finished")
	}); err != nil {
		s.logger.Error("cron: failed to register WAL checkpoint job", zap.Error(err))
	}

	s.cron.Start()
	s.logger.Info("cron: scheduler started - daily 00:00 (recurring transactions), daily 09:00 (budget history snapshot), every 6h (WAL checkpoint)")
}

func (s *CronService) StopScheduler() {
	if s.cron != nil {
		ctx := s.cron.Stop()
		<-ctx.Done()
		s.logger.Info("cron: scheduler stopped")
	}
}
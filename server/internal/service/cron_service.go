package service

import (
	"time"

	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type CronService struct {
	recurrenceService *RecurrenceService
	billService       *BillService
	billRepo          *repository.BillRepository
	recurrenceRepo    *repository.RecurrenceRepository
	txnService        *TransactionService
	budgetService     *BudgetService
	cron              *cron.Cron
	logger            *zap.Logger
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
	return &CronService{
		recurrenceService: recurrenceService,
		billService:       billService,
		billRepo:          billRepo,
		recurrenceRepo:    recurrenceRepo,
		txnService:        txnService,
		budgetService:     budgetService,
		cron:              cron.New(cron.WithSeconds()),
		logger:            logger,
	}
}

func (s *CronService) StartScheduler() {
	if _, err := s.cron.AddFunc("0 0 0 * * *", func() {
		s.logger.Info("cron: daily job started - recurrences and bills")
		s.runRecurrencesAndBills()
		s.logger.Info("cron: daily job finished")
	}); err != nil {
		s.logger.Error("cron: failed to register daily job", zap.Error(err))
	}

	if _, err := s.cron.AddFunc("0 0 9 1 * *", func() {
		s.logger.Info("cron: monthly budget snapshot job started")
		created, errs := s.budgetService.SnapshotCurrentPeriod()
		s.logger.Info("cron: monthly budget snapshot job finished",
			zap.Int("created", created),
			zap.Int("errors", len(errs)),
		)
		for _, e := range errs {
			s.logger.Warn("cron: budget snapshot error", zap.String("error", e))
		}
	}); err != nil {
		s.logger.Error("cron: failed to register budget snapshot job", zap.Error(err))
	}

	s.cron.Start()
	s.logger.Info("cron: scheduler started - daily 00:00 (recurrences+bills), monthly 1st 09:00 (budget snapshot)")
}

func (s *CronService) StopScheduler() {
	if s.cron != nil {
		ctx := s.cron.Stop()
		<-ctx.Done()
		s.logger.Info("cron: scheduler stopped")
	}
}

func (s *CronService) runRecurrencesAndBills() {
	today := time.Now()
	dueRecurrences, err := s.recurrenceRepo.GetDueRecurrences(today)
	if err != nil {
		s.logger.Error("cron: failed to get due recurrences", zap.Error(err))
	} else {
		for _, rec := range dueRecurrences {
			if err := s.recurrenceService.ExecuteRecurrence(&rec); err != nil {
				s.logger.Warn("cron: recurrence execution failed", zap.Uint64("id", rec.ID), zap.Error(err))
			} else {
				s.logger.Info("cron: recurrence executed", zap.Uint64("id", rec.ID))
			}
		}
	}

	dueBills, err := s.billRepo.GetAllDueBills()
	if err != nil {
		s.logger.Error("cron: failed to get due bills", zap.Error(err))
	} else {
		for i := range dueBills {
			if err := s.createTransactionFromBill(&dueBills[i]); err != nil {
				s.logger.Warn("cron: bill execution failed", zap.Uint64("id", dueBills[i].ID), zap.Error(err))
			} else {
				s.logger.Info("cron: bill executed", zap.Uint64("id", dueBills[i].ID))
			}
		}
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

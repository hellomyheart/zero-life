package service

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/pkg/pagination"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

type RecurringTransactionService struct {
	rtRepo    *repository.RecurringTransactionRepository
	txnRepo   *repository.TransactionRepository
	accountRepo *repository.AccountRepository
	db        *gorm.DB
}

func NewRecurringTransactionService(
	rtRepo *repository.RecurringTransactionRepository,
	txnRepo *repository.TransactionRepository,
	accountRepo *repository.AccountRepository,
	db *gorm.DB,
) *RecurringTransactionService {
	return &RecurringTransactionService{rtRepo: rtRepo, txnRepo: txnRepo, accountRepo: accountRepo, db: db}
}

func (s *RecurringTransactionService) Create(userID uint64, req *request.CreateRecurringTransactionReq) (*response.RecurringTransactionResp, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, errcode.ErrInvalidAmount
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errcode.ErrBadRequest
	}

	var endDate *time.Time
	if req.EndDate != nil {
		ed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, errcode.ErrBadRequest
		}
		endDate = &ed
	}

	repeatEvery := req.RepeatEvery
	if repeatEvery <= 0 {
		repeatEvery = 1
	}

	nextOccurrence := s.calculateNextOccurrence(startDate, model.RecurrenceType(req.RecurrenceType), repeatEvery)

	rt := &model.RecurringTransaction{
		UserID:         userID,
		Description:    req.Description,
		Amount:         amount,
		SourceID:       req.SourceID,
		DestinationID:  req.DestinationID,
		CategoryID:     req.CategoryID,
		Notes:          req.Notes,
		RecurrenceType: model.RecurrenceType(req.RecurrenceType),
		RepeatEvery:    repeatEvery,
		StartDate:      startDate,
		EndDate:        endDate,
		NextOccurrence: nextOccurrence,
		IsActive:       true,
	}

	if err := s.rtRepo.Create(rt); err != nil {
		return nil, errcode.ErrInternal
	}

	created, err := s.rtRepo.GetByID(rt.ID, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(created), nil
}

func (s *RecurringTransactionService) Get(userID, id uint64) (*response.RecurringTransactionResp, error) {
	rt, err := s.rtRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(rt), nil
}

func (s *RecurringTransactionService) List(userID uint64, req *request.RecurringTransactionListReq) (*pagination.Result, error) {
	params := pagination.Params{Page: req.Page, PageSize: req.PageSize}
	params.Normalize()

	rts, err := s.rtRepo.List(userID, req.Active, params.Offset(), params.PageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	total, err := s.rtRepo.Count(userID, req.Active)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.RecurringTransactionResp, 0, len(rts))
	for _, rt := range rts {
		items = append(items, *s.toResp(&rt))
	}

	return pagination.NewResult(items, total, params), nil
}

func (s *RecurringTransactionService) Update(userID, id uint64, req *request.UpdateRecurringTransactionReq) (*response.RecurringTransactionResp, error) {
	rt, err := s.rtRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Description != "" {
		rt.Description = req.Description
	}
	if req.Amount != "" {
		amount, err := decimal.NewFromString(req.Amount)
		if err != nil || amount.LessThanOrEqual(decimal.Zero) {
			return nil, errcode.ErrInvalidAmount
		}
		rt.Amount = amount
	}
	if req.SourceID != nil {
		rt.SourceID = *req.SourceID
	}
	if req.DestinationID != nil {
		rt.DestinationID = req.DestinationID
	}
	if req.CategoryID != nil {
		rt.CategoryID = req.CategoryID
	}
	if req.Notes != "" {
		rt.Notes = req.Notes
	}
	if req.RecurrenceType != "" {
		rt.RecurrenceType = model.RecurrenceType(req.RecurrenceType)
	}
	if req.RepeatEvery != nil {
		rt.RepeatEvery = *req.RepeatEvery
	}
	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, errcode.ErrBadRequest
		}
		rt.StartDate = startDate
	}
	if req.EndDate != nil {
		ed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, errcode.ErrBadRequest
		}
		rt.EndDate = &ed
	}
	if req.IsActive != nil {
		rt.IsActive = *req.IsActive
	}

	rt.NextOccurrence = s.calculateNextOccurrence(rt.StartDate, rt.RecurrenceType, rt.RepeatEvery)

	if err := s.rtRepo.Update(rt); err != nil {
		return nil, errcode.ErrInternal
	}

	updated, err := s.rtRepo.GetByID(id, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(updated), nil
}

func (s *RecurringTransactionService) Delete(userID, id uint64) error {
	_, err := s.rtRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}
	return s.rtRepo.Delete(id, userID)
}

func (s *RecurringTransactionService) ProcessDue(userID uint64) (int, error) {
	rts, err := s.rtRepo.GetDueRecurringTransactions(userID)
	if err != nil {
		return 0, errcode.ErrInternal
	}

	created := 0
	for _, rt := range rts {
		if err := s.createTransactionFromRecurring(&rt); err != nil {
			continue
		}
		created++

		rt.NextOccurrence = s.calculateNextOccurrence(rt.NextOccurrence, rt.RecurrenceType, rt.RepeatEvery)
		if rt.EndDate != nil && rt.NextOccurrence.After(*rt.EndDate) {
			rt.IsActive = false
		}
		s.rtRepo.Update(&rt)
	}

	return created, nil
}

func (s *RecurringTransactionService) createTransactionFromRecurring(rt *model.RecurringTransaction) error {
	txn := &model.Transaction{
		UserID:        rt.UserID,
		Type:          model.TransactionTypeWithdrawal,
		Date:          rt.NextOccurrence,
		Description:   rt.Description,
		Amount:        rt.Amount,
		SourceID:      rt.SourceID,
		DestinationID: rt.DestinationID,
		CategoryID:    rt.CategoryID,
		Notes:         rt.Notes,
	}

	if err := s.txnRepo.Create(txn, nil); err != nil {
		return err
	}

	log := &model.RecurringTransactionLog{
		RecurringTransactionID: rt.ID,
		TransactionID:          txn.ID,
		OccurrenceDate:         rt.NextOccurrence,
	}
	return s.rtRepo.CreateLog(log)
}

func (s *RecurringTransactionService) calculateNextOccurrence(from time.Time, recurrenceType model.RecurrenceType, repeatEvery int) time.Time {
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

func (s *RecurringTransactionService) toResp(rt *model.RecurringTransaction) *response.RecurringTransactionResp {
	return &response.RecurringTransactionResp{
		ID:             rt.ID,
		Description:    rt.Description,
		Amount:         rt.Amount.StringFixed(4),
		SourceID:       rt.SourceID,
		DestinationID:  rt.DestinationID,
		CategoryID:     rt.CategoryID,
		Notes:          rt.Notes,
		RecurrenceType: string(rt.RecurrenceType),
		RepeatEvery:    rt.RepeatEvery,
		StartDate:      rt.StartDate,
		EndDate:        rt.EndDate,
		NextOccurrence: rt.NextOccurrence,
		IsActive:       rt.IsActive,
		CreatedAt:      rt.CreatedAt,
		UpdatedAt:      rt.UpdatedAt,
	}
}
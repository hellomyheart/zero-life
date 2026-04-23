package service

import (
	"strings"
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

type RecurrenceService struct {
	recurrenceRepo *repository.RecurrenceRepository
	txnService     *TransactionService
	accountRepo    *repository.AccountRepository
}

func NewRecurrenceService(recurrenceRepo *repository.RecurrenceRepository, txnService *TransactionService, accountRepo *repository.AccountRepository) *RecurrenceService {
	return &RecurrenceService{
		recurrenceRepo: recurrenceRepo,
		txnService:     txnService,
		accountRepo:    accountRepo,
	}
}

func (s *RecurrenceService) Create(userID uint64, req *request.CreateRecurrenceReq) (*response.RecurrenceResp, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, errcode.ErrRecurrenceAmountInvalid
	}

	nextDate, err := time.Parse("2006-01-02", req.NextDate)
	if err != nil {
		return nil, errcode.ErrBadRequest
	}

	rec := &model.Recurrence{
		UserID:         userID,
		Title:          req.Title,
		Type:           model.TransactionType(req.Type),
		Amount:         amount,
		SourceID:       req.SourceID,
		DestinationID:  req.DestinationID,
		CategoryID:     req.CategoryID,
		Description:    req.Description,
		Notes:          req.Notes,
		TagNames:       req.TagNames,
		RepeatFreq:     model.RepeatFreq(req.RepeatFreq),
		RepeatInterval: req.RepeatInterval,
		NextDate:       nextDate,
		MaxRepetitions: req.MaxRepetitions,
		IsActive:       true,
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, errcode.ErrBadRequest
		}
		rec.EndDate = &endDate
	}

	if err := s.recurrenceRepo.Create(rec); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(rec), nil
}

func (s *RecurrenceService) Get(userID, id uint64) (*response.RecurrenceResp, error) {
	rec, err := s.recurrenceRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(rec), nil
}

func (s *RecurrenceService) List(userID uint64, req *request.RecurrenceListReq) (*pagination.Result, error) {
	params := pagination.Params{Page: req.Page, PageSize: req.PageSize}
	params.Normalize()

	recs, err := s.recurrenceRepo.List(userID, params.Offset(), params.PageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	total, err := s.recurrenceRepo.Count(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.RecurrenceResp, 0, len(recs))
	for _, r := range recs {
		items = append(items, *s.toResp(&r))
	}

	return pagination.NewResult(items, total, params), nil
}

func (s *RecurrenceService) Update(userID, id uint64, req *request.UpdateRecurrenceReq) (*response.RecurrenceResp, error) {
	rec, err := s.recurrenceRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Title != "" {
		rec.Title = req.Title
	}
	if req.Type != "" {
		rec.Type = model.TransactionType(req.Type)
	}
	if req.Amount != "" {
		amount, err := decimal.NewFromString(req.Amount)
		if err != nil || amount.LessThanOrEqual(decimal.Zero) {
			return nil, errcode.ErrRecurrenceAmountInvalid
		}
		rec.Amount = amount
	}
	if req.SourceID != nil {
		rec.SourceID = *req.SourceID
	}
	if req.DestinationID != nil {
		rec.DestinationID = req.DestinationID
	}
	if req.CategoryID != nil {
		rec.CategoryID = req.CategoryID
	}
	if req.Description != "" {
		rec.Description = req.Description
	}
	if req.Notes != "" {
		rec.Notes = req.Notes
	}
	if req.TagNames != "" {
		rec.TagNames = req.TagNames
	}
	if req.RepeatFreq != "" {
		rec.RepeatFreq = model.RepeatFreq(req.RepeatFreq)
	}
	if req.RepeatInterval != nil {
		if *req.RepeatInterval < 1 {
			return nil, errcode.ErrBadRequest
		}
		rec.RepeatInterval = *req.RepeatInterval
	}
	if req.NextDate != "" {
		nextDate, err := time.Parse("2006-01-02", req.NextDate)
		if err != nil {
			return nil, errcode.ErrBadRequest
		}
		rec.NextDate = nextDate
	}
	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, errcode.ErrBadRequest
		}
		rec.EndDate = &endDate
	}
	if req.MaxRepetitions != nil {
		rec.MaxRepetitions = req.MaxRepetitions
	}
	if req.IsActive != nil {
		rec.IsActive = *req.IsActive
	}

	if err := s.recurrenceRepo.Update(rec); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(rec), nil
}

func (s *RecurrenceService) Delete(userID, id uint64) error {
	_, err := s.recurrenceRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}

	return s.recurrenceRepo.Delete(id, userID)
}

// Trigger manually triggers a recurrence to create a transaction without updating NextDate.
func (s *RecurrenceService) Trigger(userID, id uint64) (*response.TransactionResp, error) {
	rec, err := s.recurrenceRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if !rec.IsActive {
		return nil, errcode.ErrRecurrenceInactive
	}

	return s.createTransactionFromRecurrence(rec)
}

// ExecuteRecurrence creates a transaction from the recurrence template, updates NextDate,
// increments Repetitions, and checks end conditions.
func (s *RecurrenceService) ExecuteRecurrence(rec *model.Recurrence) error {
	if !rec.IsActive {
		return nil
	}

	// Create transaction from recurrence template
	_, err := s.createTransactionFromRecurrence(rec)
	if err != nil {
		return err
	}

	// Increment repetitions
	rec.Repetitions++

	// Calculate next date
	nextDate := s.calculateNextDate(rec.NextDate, rec.RepeatFreq, rec.RepeatInterval)
	rec.NextDate = nextDate

	// Check end conditions
	if rec.EndDate != nil && nextDate.After(*rec.EndDate) {
		rec.IsActive = false
	}
	if rec.MaxRepetitions != nil && rec.Repetitions >= *rec.MaxRepetitions {
		rec.IsActive = false
	}

	// Update recurrence
	if err := s.recurrenceRepo.Update(rec); err != nil {
		return err
	}

	return nil
}

func (s *RecurrenceService) createTransactionFromRecurrence(rec *model.Recurrence) (*response.TransactionResp, error) {
	// Parse tag names to tag IDs (empty for now, tags are stored by name)
	var tagIDs []uint64
	if rec.TagNames != "" {
		_ = strings.Split(rec.TagNames, ",")
		// Tag names are stored as comma-separated; we pass empty tagIDs
		// since the transaction service expects tag IDs
	}

	dateStr := rec.NextDate.Format("2006-01-02")

	txnReq := &request.CreateTransactionReq{
		Type:          string(rec.Type),
		Date:          dateStr,
		Description:   rec.Description,
		Amount:        rec.Amount.StringFixed(4),
		SourceID:      rec.SourceID,
		DestinationID: rec.DestinationID,
		CategoryID:    rec.CategoryID,
		Notes:         rec.Notes,
		Tags:          tagIDs,
	}

	return s.txnService.Create(rec.UserID, txnReq)
}

// calculateNextDate calculates the next occurrence date based on frequency and interval.
func (s *RecurrenceService) calculateNextDate(currentDate time.Time, freq model.RepeatFreq, interval int) time.Time {
	switch freq {
	case model.RepeatFreqDaily:
		return currentDate.AddDate(0, 0, interval)
	case model.RepeatFreqWeekly:
		return currentDate.AddDate(0, 0, interval*7)
	case model.RepeatFreqMonthly:
		return currentDate.AddDate(0, interval, 0)
	case model.RepeatFreqYearly:
		return currentDate.AddDate(interval, 0, 0)
	default:
		return currentDate.AddDate(0, 0, interval)
	}
}

func (s *RecurrenceService) toResp(r *model.Recurrence) *response.RecurrenceResp {
	resp := &response.RecurrenceResp{
		ID:             r.ID,
		UserID:         r.UserID,
		Title:          r.Title,
		Type:           string(r.Type),
		Amount:         r.Amount.StringFixed(4),
		SourceID:       r.SourceID,
		DestinationID:  r.DestinationID,
		CategoryID:     r.CategoryID,
		Description:    r.Description,
		Notes:          r.Notes,
		TagNames:       r.TagNames,
		RepeatFreq:     string(r.RepeatFreq),
		RepeatInterval: r.RepeatInterval,
		NextDate:       r.NextDate,
		EndDate:        r.EndDate,
		Repetitions:    r.Repetitions,
		MaxRepetitions: r.MaxRepetitions,
		IsActive:       r.IsActive,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}

	// Load source account name
	source, err := s.accountRepo.GetByID(r.SourceID, r.UserID)
	if err == nil {
		resp.SourceName = source.Name
	}

	// Load destination account name
	if r.DestinationID != nil {
		dest, err := s.accountRepo.GetByID(*r.DestinationID, r.UserID)
		if err == nil {
			resp.DestinationName = dest.Name
		}
	}

	return resp
}

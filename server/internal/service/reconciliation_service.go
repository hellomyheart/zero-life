// Package service 业务逻辑层，实现核心业务逻辑
// ReconciliationService 对账业务逻辑，核对账户余额与实际余额
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

type ReconciliationService struct {
	recRepo     *repository.ReconciliationRepository
	accountRepo *repository.AccountRepository
	txnRepo     *repository.TransactionRepository
}

func NewReconciliationService(
	recRepo *repository.ReconciliationRepository,
	accountRepo *repository.AccountRepository,
	txnRepo *repository.TransactionRepository,
) *ReconciliationService {
	return &ReconciliationService{recRepo: recRepo, accountRepo: accountRepo, txnRepo: txnRepo}
}

func (s *ReconciliationService) Create(userID uint64, req *request.CreateReconciliationReq) (*response.ReconciliationResp, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errcode.ErrBadRequest
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, errcode.ErrBadRequest
	}

	startingBalance, err := decimal.NewFromString(req.StartingBalance)
	if err != nil {
		return nil, errcode.ErrBadRequest
	}
	endingBalance, err := decimal.NewFromString(req.EndingBalance)
	if err != nil {
		return nil, errcode.ErrBadRequest
	}

	account, err := s.accountRepo.GetByID(req.AccountID, userID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}

	bookBalance := account.CurrentBalance
	difference := endingBalance.Sub(bookBalance)

	rec := &model.TransactionReconciliation{
		AccountID:       req.AccountID,
		StartDate:       startDate,
		EndDate:         endDate,
		StartingBalance: startingBalance.StringFixed(4),
		EndingBalance:   endingBalance.StringFixed(4),
		BookBalance:     bookBalance.StringFixed(4),
		Difference:      difference.StringFixed(4),
		Status:          "open",
	}

	if err := s.recRepo.Create(rec); err != nil {
		return nil, errcode.ErrInternal
	}

	created, err := s.recRepo.GetByID(rec.ID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(created), nil
}

func (s *ReconciliationService) Get(id uint64) (*response.ReconciliationResp, error) {
	rec, err := s.recRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(rec), nil
}

func (s *ReconciliationService) List(req *request.ReconciliationListReq) (*pagination.Result, error) {
	params := pagination.Params{Page: req.Page, PageSize: req.PageSize}
	params.Normalize()

	recs, err := s.recRepo.List(req.AccountID, params.Offset(), params.PageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	total, err := s.recRepo.Count(req.AccountID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.ReconciliationResp, 0, len(recs))
	for _, r := range recs {
		items = append(items, *s.toResp(&r))
	}

	return pagination.NewResult(items, total, params), nil
}

func (s *ReconciliationService) Update(id uint64, req *request.UpdateReconciliationReq) (*response.ReconciliationResp, error) {
	rec, err := s.recRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.EndingBalance != "" {
		endingBalance, err := decimal.NewFromString(req.EndingBalance)
		if err != nil {
			return nil, errcode.ErrBadRequest
		}
		rec.EndingBalance = endingBalance.StringFixed(4)
		bookBalance, _ := decimal.NewFromString(rec.BookBalance)
		rec.Difference = endingBalance.Sub(bookBalance).StringFixed(4)
	}
	if req.Status != "" {
		rec.Status = req.Status
	}

	if err := s.recRepo.Update(rec); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(rec), nil
}

func (s *ReconciliationService) Delete(id uint64) error {
	_, err := s.recRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}
	return s.recRepo.Delete(id)
}

func (s *ReconciliationService) toResp(rec *model.TransactionReconciliation) *response.ReconciliationResp {
	return &response.ReconciliationResp{
		ID:              rec.ID,
		AccountID:       rec.AccountID,
		StartDate:       rec.StartDate,
		EndDate:         rec.EndDate,
		StartingBalance: rec.StartingBalance,
		EndingBalance:   rec.EndingBalance,
		BookBalance:     rec.BookBalance,
		Difference:      rec.Difference,
		Status:          rec.Status,
		CreatedAt:       rec.CreatedAt,
		UpdatedAt:       rec.UpdatedAt,
	}
}
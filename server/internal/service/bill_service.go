package service

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/zero-life/server/internal/dto/request"
	"github.com/zero-life/server/internal/dto/response"
	"github.com/zero-life/server/internal/model"
	"github.com/zero-life/server/internal/pkg/errcode"
	"github.com/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

type BillService struct {
	billRepo *repository.BillRepository
}

func NewBillService(billRepo *repository.BillRepository) *BillService {
	return &BillService{billRepo: billRepo}
}

func (s *BillService) Create(userID uint64, req *request.CreateBillReq) (*response.BillResp, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, errcode.ErrBillAmountInvalid
	}

	nextDue, err := time.Parse("2006-01-02", req.NextDue)
	if err != nil {
		return nil, errcode.ErrBadRequest
	}

	bill := &model.Bill{
		UserID:     userID,
		Name:       req.Name,
		Amount:     amount,
		RepeatRule: model.RepeatRule(req.RepeatRule),
		NextDue:    nextDue,
		SourceID:   req.SourceID,
		CategoryID: req.CategoryID,
		Notes:      req.Notes,
	}

	if err := s.billRepo.Create(bill); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(bill), nil
}

func (s *BillService) Get(userID, id uint64) (*response.BillResp, error) {
	bill, err := s.billRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(bill), nil
}

func (s *BillService) List(userID uint64) ([]response.BillResp, error) {
	bills, err := s.billRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.BillResp, 0, len(bills))
	for _, b := range bills {
		items = append(items, *s.toResp(&b))
	}

	return items, nil
}

func (s *BillService) Update(userID, id uint64, req *request.UpdateBillReq) (*response.BillResp, error) {
	bill, err := s.billRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Name != "" {
		bill.Name = req.Name
	}
	if req.Amount != "" {
		amount, err := decimal.NewFromString(req.Amount)
		if err != nil || amount.LessThanOrEqual(decimal.Zero) {
			return nil, errcode.ErrBillAmountInvalid
		}
		bill.Amount = amount
	}
	if req.RepeatRule != "" {
		bill.RepeatRule = model.RepeatRule(req.RepeatRule)
	}
	if req.NextDue != "" {
		nextDue, err := time.Parse("2006-01-02", req.NextDue)
		if err != nil {
			return nil, errcode.ErrBadRequest
		}
		bill.NextDue = nextDue
	}
	if req.SourceID != nil {
		bill.SourceID = req.SourceID
	}
	if req.CategoryID != nil {
		bill.CategoryID = req.CategoryID
	}
	if req.Notes != "" {
		bill.Notes = req.Notes
	}

	if err := s.billRepo.Update(bill); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(bill), nil
}

func (s *BillService) Delete(userID, id uint64) error {
	_, err := s.billRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}

	return s.billRepo.Delete(id, userID)
}

func (s *BillService) toResp(b *model.Bill) *response.BillResp {
	return &response.BillResp{
		ID:         b.ID,
		Name:       b.Name,
		Amount:     b.Amount.StringFixed(4),
		RepeatRule: string(b.RepeatRule),
		NextDue:    b.NextDue,
		SourceID:   b.SourceID,
		CategoryID: b.CategoryID,
		Notes:      b.Notes,
		CreatedAt:  b.CreatedAt,
		UpdatedAt:  b.UpdatedAt,
	}
}

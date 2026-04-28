package service

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

type BillService struct {
	billRepo   *repository.BillRepository
	txnRepo    *repository.TransactionRepository
	txnService *TransactionService
	db         *gorm.DB
}

func NewBillService(billRepo *repository.BillRepository, txnRepo *repository.TransactionRepository, txnService *TransactionService, db *gorm.DB) *BillService {
	return &BillService{billRepo: billRepo, txnRepo: txnRepo, txnService: txnService, db: db}
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
	if req.ClearSourceID {
		bill.SourceID = nil
	}
	if req.CategoryID != nil {
		bill.CategoryID = req.CategoryID
	}
	if req.ClearCategoryID {
		bill.CategoryID = nil
	}
	if req.Notes != "" {
		bill.Notes = req.Notes
	}
	if req.ClearNotes {
		bill.Notes = ""
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

	if err := s.db.Transaction(func(dbTx *gorm.DB) error {
		if err := s.billRepo.Delete(id, userID); err != nil {
			return err
		}
		if err := dbTx.Model(&model.Transaction{}).
			Where("bill_id = ? AND user_id = ?", id, userID).
			Update("bill_id", nil).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return errcode.ErrInternal
	}

	return nil
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

func (s *BillService) calculateNextDue(currentDue time.Time, rule model.RepeatRule) time.Time {
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

func (s *BillService) CreateTransactionFromBill(userID, billID uint64) (*response.BillResp, error) {
	bill, err := s.billRepo.GetByID(billID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if bill.NextDue.After(time.Now()) {
		return s.toResp(bill), nil
	}

	var createdTxn *model.Transaction

	if err := s.db.Transaction(func(dbTx *gorm.DB) error {
		if bill.SourceID != nil {
			txnReq := &request.CreateTransactionReq{
				Type:        string(model.TransactionTypeWithdrawal),
				Date:        bill.NextDue.Format("2006-01-02"),
				Description: bill.Name,
				Amount:      bill.Amount.StringFixed(4),
				SourceID:    *bill.SourceID,
				CategoryID:  bill.CategoryID,
				Notes:       bill.Notes,
				BillID:      &bill.ID,
			}

			txn, err := s.txnService.CreateWithDB(dbTx, bill.UserID, txnReq)
			if err != nil {
				return err
			}
			createdTxn = txn
		}

		nextDue := s.calculateNextDue(bill.NextDue, bill.RepeatRule)
		bill.NextDue = nextDue

		return s.billRepo.UpdateWithDB(dbTx, bill)
	}); err != nil {
		return nil, errcode.ErrInternal
	}

	if createdTxn != nil {
		created, err := s.txnRepo.GetByID(createdTxn.ID, bill.UserID)
		if err == nil {
			s.txnService.TriggerPostCreate(bill.UserID, created)
		}
	}

	return s.toResp(bill), nil
}

func (s *BillService) GetDueBills(userID uint64, days int) ([]response.BillResp, error) {
	if days <= 0 {
		days = 7
	}

	bills, err := s.billRepo.GetUpcoming(userID, days)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.BillResp, 0, len(bills))
	for _, b := range bills {
		items = append(items, *s.toResp(&b))
	}
	return items, nil
}
// Package service 业务逻辑层，实现核心业务逻辑
// PiggyBankService 存钱罐业务逻辑，处理储蓄目标、存取和进度计算
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

type PiggyBankService struct {
	piggyBankRepo *repository.PiggyBankRepository
	accountRepo   *repository.AccountRepository
}

func NewPiggyBankService(piggyBankRepo *repository.PiggyBankRepository, accountRepo *repository.AccountRepository) *PiggyBankService {
	return &PiggyBankService{piggyBankRepo: piggyBankRepo, accountRepo: accountRepo}
}

func (s *PiggyBankService) Create(userID uint64, req *request.CreatePiggyBankReq) (*response.PiggyBankResp, error) {
	targetAmount, err := decimal.NewFromString(req.TargetAmount)
	if err != nil || !targetAmount.IsPositive() {
		return nil, errcode.ErrBadRequest
	}

	piggyBank := &model.PiggyBank{
		UserID:        userID,
		Name:          req.Name,
		TargetAmount:  targetAmount,
		CurrentAmount: decimal.Zero,
		AccountID:     req.AccountID,
		Notes:         req.Notes,
	}

	if req.TargetDate != nil {
		td, err := time.Parse("2006-01-02", *req.TargetDate)
		if err != nil {
			return nil, errcode.ErrBadRequest
		}
		piggyBank.TargetDate = &td
	}

	if err := s.piggyBankRepo.Create(piggyBank); err != nil {
		return nil, errcode.ErrInternal
	}

	created, err := s.piggyBankRepo.GetByID(piggyBank.ID, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(created), nil
}

func (s *PiggyBankService) Get(userID, id uint64) (*response.PiggyBankResp, error) {
	piggyBank, err := s.piggyBankRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(piggyBank), nil
}

func (s *PiggyBankService) List(userID uint64) ([]response.PiggyBankResp, error) {
	piggyBanks, err := s.piggyBankRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.PiggyBankResp, 0, len(piggyBanks))
	for _, pb := range piggyBanks {
		items = append(items, *s.toResp(&pb))
	}
	return items, nil
}

func (s *PiggyBankService) Update(userID, id uint64, req *request.UpdatePiggyBankReq) (*response.PiggyBankResp, error) {
	piggyBank, err := s.piggyBankRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Name != "" {
		piggyBank.Name = req.Name
	}
	if req.TargetAmount != nil {
		ta, err := decimal.NewFromString(*req.TargetAmount)
		if err != nil || !ta.IsPositive() {
			return nil, errcode.ErrBadRequest
		}
		piggyBank.TargetAmount = ta
	}
	if req.TargetDate != nil {
		td, err := time.Parse("2006-01-02", *req.TargetDate)
		if err != nil {
			return nil, errcode.ErrBadRequest
		}
		piggyBank.TargetDate = &td
	}
	if req.Notes != "" {
		piggyBank.Notes = req.Notes
	}

	if err := s.piggyBankRepo.Update(piggyBank); err != nil {
		return nil, errcode.ErrInternal
	}

	updated, err := s.piggyBankRepo.GetByID(id, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	return s.toResp(updated), nil
}

func (s *PiggyBankService) Delete(userID, id uint64) error {
	_, err := s.piggyBankRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}
	return s.piggyBankRepo.Delete(id, userID)
}

func (s *PiggyBankService) AddAmount(userID, id uint64, req *request.AddAmountReq) (*response.PiggyBankResp, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || !amount.IsPositive() {
		return nil, errcode.ErrBadRequest
	}

	piggyBank, err := s.piggyBankRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	newAmount := piggyBank.CurrentAmount.Add(amount)
	if newAmount.GreaterThan(piggyBank.TargetAmount) {
		return nil, errcode.ErrBadRequest
	}

	piggyBank.CurrentAmount = newAmount
	if err := s.piggyBankRepo.Update(piggyBank); err != nil {
		return nil, errcode.ErrInternal
	}

	event := &model.PiggyEvent{
		PiggyBankID: id,
		Amount:      amount,
		Note:        req.Note,
	}
	if err := s.piggyBankRepo.CreateEvent(event); err != nil {
		return nil, errcode.ErrInternal
	}

	updated, _ := s.piggyBankRepo.GetByID(id, userID)
	return s.toResp(updated), nil
}

func (s *PiggyBankService) RemoveAmount(userID, id uint64, req *request.RemoveAmountReq) (*response.PiggyBankResp, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || !amount.IsPositive() {
		return nil, errcode.ErrBadRequest
	}

	piggyBank, err := s.piggyBankRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if amount.GreaterThan(piggyBank.CurrentAmount) {
		return nil, errcode.ErrBadRequest
	}

	piggyBank.CurrentAmount = piggyBank.CurrentAmount.Sub(amount)
	if err := s.piggyBankRepo.Update(piggyBank); err != nil {
		return nil, errcode.ErrInternal
	}

	event := &model.PiggyEvent{
		PiggyBankID: id,
		Amount:      amount.Neg(),
		Note:        req.Note,
	}
	if err := s.piggyBankRepo.CreateEvent(event); err != nil {
		return nil, errcode.ErrInternal
	}

	updated, _ := s.piggyBankRepo.GetByID(id, userID)
	return s.toResp(updated), nil
}

func (s *PiggyBankService) GetEvents(userID, id uint64) ([]response.PiggyEventResp, error) {
	_, err := s.piggyBankRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	events, err := s.piggyBankRepo.ListEvents(id)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.PiggyEventResp, 0, len(events))
	for _, e := range events {
		items = append(items, response.PiggyEventResp{
			ID:            e.ID,
			PiggyBankID:   e.PiggyBankID,
			Amount:        e.Amount.StringFixed(4),
			TransactionID: e.TransactionID,
			Note:          e.Note,
			CreatedAt:     e.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return items, nil
}

func (s *PiggyBankService) toResp(pb *model.PiggyBank) *response.PiggyBankResp {
	percentage := decimal.Zero
	if pb.TargetAmount.IsPositive() {
		percentage = pb.CurrentAmount.Div(pb.TargetAmount).Mul(decimal.NewFromInt(100))
		if percentage.GreaterThan(decimal.NewFromInt(100)) {
			percentage = decimal.NewFromInt(100)
		}
	}
	pct, _ := percentage.Float64()

	var targetDate *string
	if pb.TargetDate != nil {
		td := pb.TargetDate.Format("2006-01-02")
		targetDate = &td
	}

	return &response.PiggyBankResp{
		ID:            pb.ID,
		Name:          pb.Name,
		TargetAmount:  pb.TargetAmount.StringFixed(4),
		CurrentAmount: pb.CurrentAmount.StringFixed(4),
		AccountID:     pb.AccountID,
		TargetDate:    targetDate,
		Notes:         pb.Notes,
		Percentage:    pct,
		CreatedAt:     pb.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     pb.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// Reorder updates the order of multiple piggy banks.
func (s *PiggyBankService) Reorder(userID uint64, orders map[uint64]int) error {
	return s.piggyBankRepo.Reorder(userID, orders)
}

// ResetHistory deletes all events for a piggy bank and resets current amount to zero.
func (s *PiggyBankService) ResetHistory(userID, id uint64) (*response.PiggyBankResp, error) {
	piggyBank, err := s.piggyBankRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	// Delete all events
	if err := s.piggyBankRepo.DeleteEvents(id); err != nil {
		return nil, errcode.ErrInternal
	}

	// Reset current amount to zero
	piggyBank.CurrentAmount = decimal.Zero
	if err := s.piggyBankRepo.Update(piggyBank); err != nil {
		return nil, errcode.ErrInternal
	}

	updated, err := s.piggyBankRepo.GetByID(id, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	return s.toResp(updated), nil
}

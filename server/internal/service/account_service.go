package service

import (
	"github.com/shopspring/decimal"
	"github.com/zero-life/server/internal/dto/request"
	"github.com/zero-life/server/internal/dto/response"
	"github.com/zero-life/server/internal/model"
	"github.com/zero-life/server/internal/pkg/errcode"
	"github.com/zero-life/server/internal/pkg/pagination"
	"github.com/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

type AccountService struct {
	accountRepo *repository.AccountRepository
}

func NewAccountService(accountRepo *repository.AccountRepository) *AccountService {
	return &AccountService{accountRepo: accountRepo}
}

func (s *AccountService) Create(userID uint64, req *request.CreateAccountReq) (*response.AccountResp, error) {
	// Check name uniqueness under same type
	accounts, err := s.accountRepo.List(userID, req.Type, "", "name", 0, 0)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	for _, a := range accounts {
		if a.Name == req.Name {
			return nil, errcode.ErrAccountNameExists
		}
	}

	initialBalance, err := decimal.NewFromString(req.InitialBalance)
	if err != nil {
		return nil, errcode.ErrBadRequest
	}

	account := &model.Account{
		UserID:         userID,
		Name:           req.Name,
		Type:           model.AccountType(req.Type),
		CurrencyID:     req.CurrencyID,
		InitialBalance: initialBalance,
		CurrentBalance: initialBalance,
		IsVirtual:      req.IsVirtual,
		Notes:          req.Notes,
	}

	if err := s.accountRepo.Create(account); err != nil {
		return nil, errcode.ErrInternal
	}

	// Reload with currency
	created, err := s.accountRepo.GetByID(account.ID, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(created), nil
}

func (s *AccountService) Get(userID, id uint64) (*response.AccountResp, error) {
	account, err := s.accountRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(account), nil
}

func (s *AccountService) List(userID uint64, req *request.AccountListReq) (*pagination.Result, error) {
	params := pagination.Params{Page: req.Page, PageSize: req.PageSize}
	params.Normalize()

	accounts, err := s.accountRepo.List(userID, req.Type, req.Search, req.Sort, params.Offset(), params.PageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	total, err := s.accountRepo.Count(userID, req.Type, req.Search)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.AccountResp, 0, len(accounts))
	for _, a := range accounts {
		items = append(items, *s.toResp(&a))
	}

	return pagination.NewResult(items, total, params), nil
}

func (s *AccountService) Update(userID, id uint64, req *request.UpdateAccountReq) (*response.AccountResp, error) {
	account, err := s.accountRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Name != "" {
		account.Name = req.Name
	}
	if req.Notes != "" {
		account.Notes = req.Notes
	}
	if req.IsVirtual != nil {
		account.IsVirtual = *req.IsVirtual
	}

	if err := s.accountRepo.Update(account); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(account), nil
}

func (s *AccountService) Delete(userID, id uint64) error {
	hasTxns, err := s.accountRepo.HasTransactions(id, userID)
	if err != nil {
		return errcode.ErrInternal
	}
	if hasTxns {
		return errcode.ErrAccountHasTxns
	}

	return s.accountRepo.Delete(id, userID)
}

func (s *AccountService) toResp(a *model.Account) *response.AccountResp {
	return &response.AccountResp{
		ID:             a.ID,
		Name:           a.Name,
		Type:           string(a.Type),
		CurrencyID:     a.CurrencyID,
		Currency:       currencyToResp(&a.Currency),
		InitialBalance: a.InitialBalance.StringFixed(4),
		CurrentBalance: a.CurrentBalance.StringFixed(4),
		IsVirtual:      a.IsVirtual,
		Notes:          a.Notes,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
	}
}

func currencyToResp(c *model.Currency) response.CurrencyResp {
	return response.CurrencyResp{
		ID:            c.ID,
		Code:          c.Code,
		Name:          c.Name,
		Symbol:        c.Symbol,
		DecimalPlaces: c.DecimalPlaces,
		IsEnabled:     c.IsEnabled,
		IsDefault:     c.IsDefault,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}

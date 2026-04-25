// Package service 业务逻辑层，实现核心业务逻辑
// AccountService 账户业务逻辑，处理账户的增删改查及资产账户余额计算
package service

import (
	"github.com/shopspring/decimal"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/pkg/pagination"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

// AccountService 账户业务服务
// 负责处理账户相关的业务逻辑，包括账户的创建、查询、更新、删除等操作
type AccountService struct {
	accountRepo *repository.AccountRepository // 账户数据访问对象
}

// NewAccountService 创建账户服务实例
// 参数：
//   - accountRepo: 账户数据访问对象
// 返回：
//   - *AccountService: 账户服务实例
func NewAccountService(accountRepo *repository.AccountRepository) *AccountService {
	return &AccountService{accountRepo: accountRepo}
}

// Create 创建新账户
// 业务流程：
// 1. 检查同类型下账户名称是否重复
// 2. 解析并验证初始余额
// 3. 创建账户记录，初始余额作为当前余额
// 4. 重新加载账户信息（包含关联的货币信息）
// 参数：
//   - userID: 用户ID
//   - req: 创建账户请求参数
// 返回：
//   - *response.AccountResp: 创建成功的账户信息
//   - error: 错误信息（如名称重复、参数错误等）
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

	// 解析初始余额字符串为decimal类型，避免浮点精度问题
	initialBalance, err := decimal.NewFromString(req.InitialBalance)
	if err != nil {
		return nil, errcode.ErrBadRequest
	}

	// 创建账户时，初始余额即为当前余额（后续交易会更新当前余额）
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

// Get 获取账户详情
// 根据账户ID和用户ID查询账户信息
// 参数：
//   - userID: 用户ID
//   - id: 账户ID
// 返回：
//   - *response.AccountResp: 账户信息
//   - error: 错误信息（如账户不存在）
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

// List 获取账户列表
// 支持按类型过滤、搜索、排序和分页
// 参数：
//   - userID: 用户ID
//   - req: 列表查询参数（包含类型、搜索关键词、排序、分页等）
// 返回：
//   - *pagination.Result: 分页结果（包含账户列表和总数）
//   - error: 错误信息
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

// Update 更新账户信息
// 仅更新请求中提供的字段（部分更新）
// 参数：
//   - userID: 用户ID
//   - id: 账户ID
//   - req: 更新请求参数
// 返回：
//   - *response.AccountResp: 更新后的账户信息
//   - error: 错误信息（如账户不存在）
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

// Delete 删除账户
// 业务规则：如果账户下存在交易记录，则不允许删除
// 参数：
//   - userID: 用户ID
//   - id: 账户ID
// 返回：
//   - error: 错误信息（如账户不存在、账户有交易记录）
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

// toResp 将账户模型转换为响应对象
// 参数：
//   - a: 账户模型
// 返回：
//   - *response.AccountResp: 账户响应对象
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

// currencyToResp 将货币模型转换为响应对象
// 参数：
//   - c: 货币模型
// 返回：
//   - response.CurrencyResp: 货币响应对象
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

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

// ReconciliationService 对账业务服务
// 负责处理对账相关的业务逻辑，包括对账记录的创建、查询、更新和删除
// 所有操作均需传入用户ID，确保用户只能操作自己的对账数据
type ReconciliationService struct {
	recRepo     *repository.ReconciliationRepository
	accountRepo *repository.AccountRepository
	txnRepo     *repository.TransactionRepository
}

// NewReconciliationService 创建对账服务实例
// 参数：
//   - recRepo: 对账数据访问对象
//   - accountRepo: 账户数据访问对象
//   - txnRepo: 交易数据访问对象
// 返回：
//   - *ReconciliationService: 对账服务实例
func NewReconciliationService(
	recRepo *repository.ReconciliationRepository,
	accountRepo *repository.AccountRepository,
	txnRepo *repository.TransactionRepository,
) *ReconciliationService {
	return &ReconciliationService{recRepo: recRepo, accountRepo: accountRepo, txnRepo: txnRepo}
}

// Create 创建对账记录
// 业务流程：
// 1. 解析并验证开始日期、结束日期
// 2. 解析并验证期初余额、期末余额
// 3. 根据账户ID和用户ID查询账户信息（确保账户属于当前用户）
// 4. 计算账面余额和差额
// 5. 创建对账记录并保存
// 参数：
//   - userID: 用户ID，用于验证账户归属
//   - req: 创建对账请求参数
// 返回：
//   - *response.ReconciliationResp: 创建成功的对账信息
//   - error: 错误信息
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

	created, err := s.recRepo.GetByID(rec.ID, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(created), nil
}

// Get 获取对账记录详情
// 根据对账ID和用户ID查询对账记录，通过关联账户表确保用户只能访问自己的对账数据
// 参数：
//   - userID: 用户ID，用于验证数据归属
//   - id: 对账记录ID
// 返回：
//   - *response.ReconciliationResp: 对账信息
//   - error: 错误信息（如对账记录不存在或不属于当前用户）
func (s *ReconciliationService) Get(userID, id uint64) (*response.ReconciliationResp, error) {
	rec, err := s.recRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(rec), nil
}

// List 获取对账记录列表
// 根据用户ID和查询条件获取对账记录列表，通过关联账户表确保用户只能查看自己的对账数据
// 参数：
//   - userID: 用户ID，用于过滤数据归属
//   - req: 列表查询参数（包含账户ID、分页等）
// 返回：
//   - *pagination.Result: 分页结果
//   - error: 错误信息
func (s *ReconciliationService) List(userID uint64, req *request.ReconciliationListReq) (*pagination.Result, error) {
	params := pagination.Params{Page: req.Page, PageSize: req.PageSize}
	params.Normalize()

	recs, err := s.recRepo.List(userID, req.AccountID, params.Offset(), params.PageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	total, err := s.recRepo.Count(userID, req.AccountID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.ReconciliationResp, 0, len(recs))
	for _, r := range recs {
		items = append(items, *s.toResp(&r))
	}

	return pagination.NewResult(items, total, params), nil
}

// Update 更新对账记录
// 根据对账ID和用户ID查询对账记录，确保用户只能更新自己的对账数据
// 支持更新期末余额和对账状态，更新期末余额时会自动重新计算差额
// 参数：
//   - userID: 用户ID，用于验证数据归属
//   - id: 对账记录ID
//   - req: 更新请求参数
// 返回：
//   - *response.ReconciliationResp: 更新后的对账信息
//   - error: 错误信息
func (s *ReconciliationService) Update(userID, id uint64, req *request.UpdateReconciliationReq) (*response.ReconciliationResp, error) {
	rec, err := s.recRepo.GetByID(id, userID)
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

// Delete 删除对账记录
// 根据对账ID和用户ID查询对账记录，确保用户只能删除自己的对账数据
// 同时会删除该对账记录下的所有条目
// 参数：
//   - userID: 用户ID，用于验证数据归属
//   - id: 对账记录ID
// 返回：
//   - error: 错误信息
func (s *ReconciliationService) Delete(userID, id uint64) error {
	_, err := s.recRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}
	return s.recRepo.Delete(id)
}

// toResp 将对账模型转换为响应对象
// 参数：
//   - rec: 对账模型
// 返回：
//   - *response.ReconciliationResp: 对账响应对象
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

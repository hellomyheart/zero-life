// Package service 业务逻辑层，实现核心业务逻辑
// BillService 账单业务逻辑，处理账单的增删改查和定期计算
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

// BillService 账单服务
// 依赖billRepo进行账单数据访问，依赖txnService创建关联交易
type BillService struct {
	billRepo  *repository.BillRepository
	txnService *TransactionService
}

// NewBillService 创建账单服务实例
// billRepo: 账单数据访问层
// txnService: 交易服务，用于从账单创建交易时更新账户余额、触发规则和Webhook
func NewBillService(billRepo *repository.BillRepository, txnService *TransactionService) *BillService {
	return &BillService{billRepo: billRepo, txnService: txnService}
}

// Create 创建账单
// 业务流程：
// 1. 解析并验证金额（必须大于0）
// 2. 解析并验证下次到期日期
// 3. 创建账单记录
// 参数：
//   - userID: 用户ID
//   - req: 创建账单请求参数（名称、金额、重复规则、下次到期日、源账户、分类、备注）
// 返回：
//   - *response.BillResp: 创建成功的账单信息
//   - error: 错误信息（如金额无效、日期格式错误）
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

// Get 获取单个账单详情
// 参数：
//   - userID: 用户ID
//   - id: 账单ID
// 返回：
//   - *response.BillResp: 账单信息
//   - error: 错误信息（如账单不存在）
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

// List 获取用户所有账单列表
// 参数：
//   - userID: 用户ID
// 返回：
//   - []response.BillResp: 账单列表
//   - error: 错误信息
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

// Update 更新账单信息
// 支持部分更新：名称、金额、重复规则、下次到期日、源账户、分类、备注
// 更新金额时会重新验证金额有效性
// 参数：
//   - userID: 用户ID
//   - id: 账单ID
//   - req: 更新请求参数
// 返回：
//   - *response.BillResp: 更新后的账单信息
//   - error: 错误信息
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

// Delete 删除账单
// 参数：
//   - userID: 用户ID
//   - id: 账单ID
// 返回：
//   - error: 错误信息（如账单不存在）
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

// toResp 将账单模型转换为响应DTO
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

// CreateTransactionFromBill 从账单创建交易并推进下次到期日期
// 修复：原方法只推进到期日期而不创建交易，现在通过txnService.Create()创建支出交易，
// 确保账户余额更新、规则触发和Webhook通知正常执行
func (s *BillService) CreateTransactionFromBill(userID, billID uint64) (*response.BillResp, error) {
	bill, err := s.billRepo.GetByID(billID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	// 如果账单指定了支出账户，则创建对应的支出交易
	// 通过txnService.Create()确保：1)账户余额正确更新 2)规则触发 3)Webhook通知
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

		if _, err := s.txnService.Create(bill.UserID, txnReq); err != nil {
			return nil, errcode.ErrInternal
		}
	}

	// 根据重复规则推进下次到期日期
	switch bill.RepeatRule {
	case model.RepeatRuleDaily:
		bill.NextDue = bill.NextDue.AddDate(0, 0, 1)
	case model.RepeatRuleWeekly:
		bill.NextDue = bill.NextDue.AddDate(0, 0, 7)
	case model.RepeatRuleMonthly:
		bill.NextDue = bill.NextDue.AddDate(0, 1, 0)
	case model.RepeatRuleYearly:
		bill.NextDue = bill.NextDue.AddDate(1, 0, 0)
	}

	if err := s.billRepo.Update(bill); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(bill), nil
}

// GetDueBills 获取指定天数内到期的账单列表
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

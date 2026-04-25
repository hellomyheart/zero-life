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

// PiggyBankService 存钱罐服务
// 负责处理储蓄目标的创建、存取款、进度计算和事件记录
// 依赖piggyBankRepo进行存钱罐数据访问，依赖accountRepo验证关联账户
type PiggyBankService struct {
	piggyBankRepo *repository.PiggyBankRepository // 存钱罐数据访问对象
	accountRepo   *repository.AccountRepository   // 账户数据访问对象
}

// NewPiggyBankService 创建存钱罐服务实例
func NewPiggyBankService(piggyBankRepo *repository.PiggyBankRepository, accountRepo *repository.AccountRepository) *PiggyBankService {
	return &PiggyBankService{piggyBankRepo: piggyBankRepo, accountRepo: accountRepo}
}

// Create 创建存钱罐
// 设置目标金额和可选的目标日期，初始当前金额为0
// 参数：
//   - userID: 用户ID
//   - req: 创建请求参数（名称、目标金额、目标日期、关联账户、备注）
// 返回：
//   - *response.PiggyBankResp: 创建成功的存钱罐信息（含完成百分比）
//   - error: 错误信息
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

// Get 获取单个存钱罐详情
// 参数：
//   - userID: 用户ID
//   - id: 存钱罐ID
// 返回：
//   - *response.PiggyBankResp: 存钱罐信息（含完成百分比）
//   - error: 错误信息
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

// List 获取用户所有存钱罐列表
// 参数：
//   - userID: 用户ID
// 返回：
//   - []response.PiggyBankResp: 存钱罐列表
//   - error: 错误信息
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

// Update 更新存钱罐信息
// 支持更新名称、目标金额、目标日期、备注
// 参数：
//   - userID: 用户ID
//   - id: 存钱罐ID
//   - req: 更新请求参数
// 返回：
//   - *response.PiggyBankResp: 更新后的存钱罐信息
//   - error: 错误信息
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

// Delete 删除存钱罐
// 参数：
//   - userID: 用户ID
//   - id: 存钱罐ID
// 返回：
//   - error: 错误信息
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

// AddAmount 向存钱罐存入金额
// 业务规则：存入后当前金额不能超过目标金额
// 同时创建一条存入事件记录
// 参数：
//   - userID: 用户ID
//   - id: 存钱罐ID
//   - req: 存入请求参数（金额、备注）
// 返回：
//   - *response.PiggyBankResp: 更新后的存钱罐信息
//   - error: 错误信息
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

	// 业务规则：存入后当前金额不能超过目标金额（防止超额存入）
	newAmount := piggyBank.CurrentAmount.Add(amount)
	if newAmount.GreaterThan(piggyBank.TargetAmount) {
		return nil, errcode.ErrBadRequest
	}

	piggyBank.CurrentAmount = newAmount
	if err := s.piggyBankRepo.Update(piggyBank); err != nil {
		return nil, errcode.ErrInternal
	}

	// 创建存入事件记录（正数金额）
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

// RemoveAmount 从存钱罐取出金额
// 业务规则：取出金额不能超过当前已存金额
// 同时创建一条取出事件记录（金额为负值）
// 参数：
//   - userID: 用户ID
//   - id: 存钱罐ID
//   - req: 取出请求参数（金额、备注）
// 返回：
//   - *response.PiggyBankResp: 更新后的存钱罐信息
//   - error: 错误信息
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

	// 业务规则：取出金额不能超过当前已存金额（防止超额取出）
	// 业务规则：取出金额不能超过当前已存金额（防止超额取出）
	if amount.GreaterThan(piggyBank.CurrentAmount) {
		return nil, errcode.ErrBadRequest
	}

	piggyBank.CurrentAmount = piggyBank.CurrentAmount.Sub(amount)
	if err := s.piggyBankRepo.Update(piggyBank); err != nil {
		return nil, errcode.ErrInternal
	}

	// 创建取出事件记录（负数金额，使用Neg()将正数转为负数）
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

// GetEvents 获取存钱罐的事件记录列表（存入/取出历史）
// 参数：
//   - userID: 用户ID
//   - id: 存钱罐ID
// 返回：
//   - []response.PiggyEventResp: 事件列表
//   - error: 错误信息
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

// toResp 将存钱罐模型转换为响应对象
// 自动计算完成百分比（当前金额/目标金额*100，上限100%）
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

// Reorder 批量更新存钱罐的排序顺序
// 参数：
//   - userID: 用户ID
//   - orders: 存钱罐ID到排序值的映射
// 返回：
//   - error: 错误信息
func (s *PiggyBankService) Reorder(userID uint64, orders map[uint64]int) error {
	return s.piggyBankRepo.Reorder(userID, orders)
}

// ResetHistory 重置存钱罐历史
// 删除所有存取事件记录，将当前金额重置为0
// 参数：
//   - userID: 用户ID
//   - id: 存钱罐ID
// 返回：
//   - *response.PiggyBankResp: 重置后的存钱罐信息
//   - error: 错误信息
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

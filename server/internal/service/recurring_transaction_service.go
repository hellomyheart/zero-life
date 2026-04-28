// Package service 业务逻辑层，实现核心业务逻辑
// RecurringTransactionService 循环交易业务逻辑，自动创建循环交易
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

// RecurringTransactionService 循环交易服务
// 依赖rtRepo进行循环交易数据访问，依赖txnService创建交易（确保余额更新、规则触发、Webhook通知）
// 依赖accountRepo进行账户验证，依赖db进行数据库操作
type RecurringTransactionService struct {
	rtRepo      *repository.RecurringTransactionRepository
	txnRepo     *repository.TransactionRepository
	txnService  *TransactionService
	accountRepo *repository.AccountRepository
	db          *gorm.DB
}

// NewRecurringTransactionService 创建循环交易服务实例
// rtRepo: 循环交易数据访问层
// txnRepo: 交易数据访问层（仅用于查询，创建交易走txnService）
// txnService: 交易服务，用于创建交易时更新账户余额、触发规则和Webhook
// accountRepo: 账户数据访问层
// db: 数据库连接
func NewRecurringTransactionService(
	rtRepo *repository.RecurringTransactionRepository,
	txnRepo *repository.TransactionRepository,
	txnService *TransactionService,
	accountRepo *repository.AccountRepository,
	db *gorm.DB,
) *RecurringTransactionService {
	return &RecurringTransactionService{
		rtRepo:      rtRepo,
		txnRepo:     txnRepo,
		txnService:  txnService,
		accountRepo: accountRepo,
		db:          db,
	}
}

// Create 创建循环交易
// 业务流程：
// 1. 解析并验证金额（必须大于0）
// 2. 解析并验证开始日期和可选的结束日期
// 3. 设置重复间隔（默认为1）
// 4. 计算下次执行日期
// 5. 创建循环交易记录
// 参数：
//   - userID: 用户ID
//   - req: 创建请求参数（描述、金额、账户、分类、重复类型、间隔、日期等）
// 返回：
//   - *response.RecurringTransactionResp: 创建成功的循环交易信息
//   - error: 错误信息
func (s *RecurringTransactionService) Create(userID uint64, req *request.CreateRecurringTransactionReq) (*response.RecurringTransactionResp, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, errcode.ErrInvalidAmount
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errcode.ErrBadRequest
	}

	var endDate *time.Time
	if req.EndDate != nil {
		ed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, errcode.ErrBadRequest
		}
		endDate = &ed
	}

	repeatEvery := req.RepeatEvery
	if repeatEvery <= 0 {
		repeatEvery = 1
	}

	nextOccurrence := calculateNextOccurrence(startDate, model.RecurrenceType(req.RecurrenceType), repeatEvery)

	rt := &model.RecurringTransaction{
		UserID:         userID,
		Description:    req.Description,
		Amount:         amount,
		SourceID:       req.SourceID,
		DestinationID:  req.DestinationID,
		CategoryID:     req.CategoryID,
		Notes:          req.Notes,
		RecurrenceType: model.RecurrenceType(req.RecurrenceType),
		RepeatEvery:    repeatEvery,
		StartDate:      startDate,
		EndDate:        endDate,
		NextOccurrence: nextOccurrence,
		IsActive:       true,
		ReminderDays:   req.ReminderDays,
	}

	if err := s.rtRepo.Create(rt); err != nil {
		return nil, errcode.ErrInternal
	}

	created, err := s.rtRepo.GetByID(rt.ID, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(created), nil
}

// Get 获取单个循环交易详情
// 参数：
//   - userID: 用户ID
//   - id: 循环交易ID
// 返回：
//   - *response.RecurringTransactionResp: 循环交易信息
//   - error: 错误信息
func (s *RecurringTransactionService) Get(userID, id uint64) (*response.RecurringTransactionResp, error) {
	rt, err := s.rtRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(rt), nil
}

// List 获取循环交易分页列表
// 支持按活跃状态过滤
// 参数：
//   - userID: 用户ID
//   - req: 列表查询参数（含分页、活跃状态过滤）
// 返回：
//   - *pagination.Result: 分页结果
//   - error: 错误信息
func (s *RecurringTransactionService) List(userID uint64, req *request.RecurringTransactionListReq) (*pagination.Result, error) {
	params := pagination.Params{Page: req.Page, PageSize: req.PageSize}
	params.Normalize()

	rts, err := s.rtRepo.List(userID, req.Active, params.Offset(), params.PageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	total, err := s.rtRepo.Count(userID, req.Active)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.RecurringTransactionResp, 0, len(rts))
	for _, rt := range rts {
		items = append(items, *s.toResp(&rt))
	}

	return pagination.NewResult(items, total, params), nil
}

// Update 更新循环交易信息
// 支持更新描述、金额、账户、分类、重复类型、间隔、日期、结束条件等
// 更新后会重新计算下次执行日期
// 参数：
//   - userID: 用户ID
//   - id: 循环交易ID
//   - req: 更新请求参数
// 返回：
//   - *response.RecurringTransactionResp: 更新后的循环交易信息
//   - error: 错误信息
func (s *RecurringTransactionService) Update(userID, id uint64, req *request.UpdateRecurringTransactionReq) (*response.RecurringTransactionResp, error) {
	rt, err := s.rtRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Description != "" {
		rt.Description = req.Description
	}
	if req.Amount != "" {
		amount, err := decimal.NewFromString(req.Amount)
		if err != nil || amount.LessThanOrEqual(decimal.Zero) {
			return nil, errcode.ErrInvalidAmount
		}
		rt.Amount = amount
	}
	if req.SourceID != nil {
		rt.SourceID = *req.SourceID
	}
	if req.DestinationID != nil {
		rt.DestinationID = req.DestinationID
	}
	if req.CategoryID != nil {
		rt.CategoryID = req.CategoryID
	}
	if req.Notes != "" {
		rt.Notes = req.Notes
	}
	if req.RecurrenceType != "" {
		rt.RecurrenceType = model.RecurrenceType(req.RecurrenceType)
	}
	if req.RepeatEvery != nil {
		rt.RepeatEvery = *req.RepeatEvery
	}
	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, errcode.ErrBadRequest
		}
		rt.StartDate = startDate
	}
	if req.EndDate != nil {
		ed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, errcode.ErrBadRequest
		}
		rt.EndDate = &ed
	}
	if req.IsActive != nil {
		rt.IsActive = *req.IsActive
	}
	if req.ReminderDays != nil {
		rt.ReminderDays = *req.ReminderDays
	}

	// 仅在循环规则相关字段变更时才重算 NextOccurrence
	if req.RecurrenceType != "" || req.RepeatEvery != nil || req.StartDate != "" {
		rt.NextOccurrence = calculateNextOccurrence(rt.StartDate, rt.RecurrenceType, rt.RepeatEvery)
	}

	if err := s.rtRepo.Update(rt); err != nil {
		return nil, errcode.ErrInternal
	}

	updated, err := s.rtRepo.GetByID(id, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(updated), nil
}

// Delete 删除循环交易
// 参数：
//   - userID: 用户ID
//   - id: 循环交易ID
// 返回：
//   - error: 错误信息
func (s *RecurringTransactionService) Delete(userID, id uint64) error {
	_, err := s.rtRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}
	return s.rtRepo.Delete(id, userID)
}

// ProcessDue 处理指定用户所有到期的循环交易
func (s *RecurringTransactionService) ProcessDue(userID uint64) (int, error) {
	rts, err := s.rtRepo.GetDueRecurringTransactions(userID)
	if err != nil {
		return 0, errcode.ErrInternal
	}

	created := 0
	for _, rt := range rts {
		if err := s.processOneRecurring(&rt); err != nil {
			continue
		}
		created++
	}

	return created, nil
}

// ProcessAllDue 处理所有用户到期的循环交易，供 CronService 定时任务调用。
// 每条循环交易在事务中完成：创建交易 + 写日志 + 推进日期。
func (s *RecurringTransactionService) ProcessAllDue() []string {
	rts, err := s.rtRepo.GetAllDue()
	if err != nil {
		return []string{err.Error()}
	}

	var errs []string
	for i := range rts {
		if err := s.processOneRecurring(&rts[i]); err != nil {
			errs = append(errs, err.Error())
		}
	}
	return errs
}

func determineTransactionType(rt *model.RecurringTransaction) model.TransactionType {
	if rt.SourceID != 0 && rt.DestinationID != nil && *rt.DestinationID != 0 {
		return model.TransactionTypeTransfer
	}
	if rt.DestinationID != nil && *rt.DestinationID != 0 {
		return model.TransactionTypeDeposit
	}
	return model.TransactionTypeWithdrawal
}

func calculateNextOccurrence(from time.Time, recurrenceType model.RecurrenceType, repeatEvery int) time.Time {
	switch recurrenceType {
	case model.RecurrenceTypeDaily:
		return from.AddDate(0, 0, repeatEvery)
	case model.RecurrenceTypeWeekly:
		return from.AddDate(0, 0, 7*repeatEvery)
	case model.RecurrenceTypeMonthly:
		return from.AddDate(0, repeatEvery, 0)
	case model.RecurrenceTypeYearly:
		return from.AddDate(repeatEvery, 0, 0)
	default:
		return from.AddDate(0, 1, 0)
	}
}

// processOneRecurring 处理单条到期循环交易：创建交易、写日志、推进下次日期。
// 使用数据库事务确保原子性。
func (s *RecurringTransactionService) processOneRecurring(rt *model.RecurringTransaction) error {
	if rt.NextOccurrence.After(time.Now()) {
		return nil
	}

	txnType := determineTransactionType(rt)
	txnReq := &request.CreateTransactionReq{
		Type:          string(txnType),
		Date:          rt.NextOccurrence.Format("2006-01-02"),
		Description:   rt.Description,
		Amount:        rt.Amount.StringFixed(4),
		SourceID:      rt.SourceID,
		DestinationID: rt.DestinationID,
		CategoryID:    rt.CategoryID,
		Notes:         rt.Notes,
		RecurringID:   &rt.ID,
	}

	var txn *model.Transaction
	if err := s.db.Transaction(func(dbTx *gorm.DB) error {
		var err error
		txn, err = s.txnService.CreateWithDB(dbTx, rt.UserID, txnReq)
		if err != nil {
			return err
		}

		log := &model.RecurringTransactionLog{
			RecurringTransactionID: rt.ID,
			TransactionID:          txn.ID,
			OccurrenceDate:         rt.NextOccurrence,
		}
		if err := s.rtRepo.CreateLogWithDB(dbTx, log); err != nil {
			return err
		}

		rt.NextOccurrence = calculateNextOccurrence(rt.NextOccurrence, rt.RecurrenceType, rt.RepeatEvery)
		if rt.EndDate != nil && rt.NextOccurrence.After(*rt.EndDate) {
			rt.IsActive = false
		}
		return s.rtRepo.UpdateWithDB(dbTx, rt)
	}); err != nil {
		return err
	}

	s.txnService.TriggerPostCreate(rt.UserID, txn)
	return nil
}

// toResp 将循环交易模型转换为响应DTO
func (s *RecurringTransactionService) toResp(rt *model.RecurringTransaction) *response.RecurringTransactionResp {
	return &response.RecurringTransactionResp{
		ID:             rt.ID,
		Description:    rt.Description,
		Amount:         rt.Amount.StringFixed(4),
		SourceID:       rt.SourceID,
		DestinationID:  rt.DestinationID,
		CategoryID:     rt.CategoryID,
		Notes:          rt.Notes,
		RecurrenceType: string(rt.RecurrenceType),
		RepeatEvery:    rt.RepeatEvery,
		StartDate:      rt.StartDate,
		EndDate:        rt.EndDate,
		NextOccurrence: rt.NextOccurrence,
		IsActive:       rt.IsActive,
		ReminderDays:   rt.ReminderDays,
		CreatedAt:      rt.CreatedAt,
		UpdatedAt:      rt.UpdatedAt,
	}
}

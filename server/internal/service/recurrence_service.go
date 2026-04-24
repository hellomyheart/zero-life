package service

import (
	"strings"
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

// RecurrenceService 循环交易服务
// 负责处理循环交易的增删改查、手动触发和自动执行
// 循环交易按指定频率（日/周/月/年）自动创建交易，支持设置结束条件和最大重复次数
// 依赖recurrenceRepo进行循环交易数据访问，依赖txnService创建交易（确保余额更新）
// 依赖accountRepo验证账户归属和加载账户名称
type RecurrenceService struct {
	recurrenceRepo *repository.RecurrenceRepository // 循环交易数据访问对象
	txnService     *TransactionService              // 交易服务，用于创建循环产生的交易
	accountRepo    *repository.AccountRepository    // 账户数据访问对象
}

// NewRecurrenceService 创建循环交易服务实例
func NewRecurrenceService(recurrenceRepo *repository.RecurrenceRepository, txnService *TransactionService, accountRepo *repository.AccountRepository) *RecurrenceService {
	return &RecurrenceService{
		recurrenceRepo: recurrenceRepo,
		txnService:     txnService,
		accountRepo:    accountRepo,
	}
}

// Create 创建循环交易
// 设置下次执行日期、重复频率、间隔和结束条件
// 参数：
//   - userID: 用户ID
//   - req: 创建请求参数
// 返回：
//   - *response.RecurrenceResp: 创建成功的循环交易信息
//   - error: 错误信息
func (s *RecurrenceService) Create(userID uint64, req *request.CreateRecurrenceReq) (*response.RecurrenceResp, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, errcode.ErrRecurrenceAmountInvalid
	}

	nextDate, err := time.Parse("2006-01-02", req.NextDate)
	if err != nil {
		return nil, errcode.ErrBadRequest
	}

	rec := &model.Recurrence{
		UserID:         userID,
		Title:          req.Title,
		Type:           model.TransactionType(req.Type),
		Amount:         amount,
		SourceID:       req.SourceID,
		DestinationID:  req.DestinationID,
		CategoryID:     req.CategoryID,
		Description:    req.Description,
		Notes:          req.Notes,
		TagNames:       req.TagNames,
		RepeatFreq:     model.RepeatFreq(req.RepeatFreq),
		RepeatInterval: req.RepeatInterval,
		NextDate:       nextDate,
		MaxRepetitions: req.MaxRepetitions,
		IsActive:       true,
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, errcode.ErrBadRequest
		}
		rec.EndDate = &endDate
	}

	if err := s.recurrenceRepo.Create(rec); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(rec), nil
}

// Get 获取单个循环交易详情
// 参数：
//   - userID: 用户ID
//   - id: 循环交易ID
// 返回：
//   - *response.RecurrenceResp: 循环交易信息
//   - error: 错误信息
func (s *RecurrenceService) Get(userID, id uint64) (*response.RecurrenceResp, error) {
	rec, err := s.recurrenceRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(rec), nil
}

// List 获取循环交易分页列表
// 参数：
//   - userID: 用户ID
//   - req: 列表查询参数（含分页）
// 返回：
//   - *pagination.Result: 分页结果
//   - error: 错误信息
func (s *RecurrenceService) List(userID uint64, req *request.RecurrenceListReq) (*pagination.Result, error) {
	params := pagination.Params{Page: req.Page, PageSize: req.PageSize}
	params.Normalize()

	recs, err := s.recurrenceRepo.List(userID, params.Offset(), params.PageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	total, err := s.recurrenceRepo.Count(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.RecurrenceResp, 0, len(recs))
	for _, r := range recs {
		items = append(items, *s.toResp(&r))
	}

	return pagination.NewResult(items, total, params), nil
}

// Update 更新循环交易信息
// 支持更新标题、类型、金额、账户、分类、重复频率、间隔、日期、结束条件等
// 参数：
//   - userID: 用户ID
//   - id: 循环交易ID
//   - req: 更新请求参数
// 返回：
//   - *response.RecurrenceResp: 更新后的循环交易信息
//   - error: 错误信息
func (s *RecurrenceService) Update(userID, id uint64, req *request.UpdateRecurrenceReq) (*response.RecurrenceResp, error) {
	rec, err := s.recurrenceRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Title != "" {
		rec.Title = req.Title
	}
	if req.Type != "" {
		rec.Type = model.TransactionType(req.Type)
	}
	if req.Amount != "" {
		amount, err := decimal.NewFromString(req.Amount)
		if err != nil || amount.LessThanOrEqual(decimal.Zero) {
			return nil, errcode.ErrRecurrenceAmountInvalid
		}
		rec.Amount = amount
	}
	if req.SourceID != nil {
		rec.SourceID = *req.SourceID
	}
	if req.DestinationID != nil {
		rec.DestinationID = req.DestinationID
	}
	if req.CategoryID != nil {
		rec.CategoryID = req.CategoryID
	}
	if req.Description != "" {
		rec.Description = req.Description
	}
	if req.Notes != "" {
		rec.Notes = req.Notes
	}
	if req.TagNames != "" {
		rec.TagNames = req.TagNames
	}
	if req.RepeatFreq != "" {
		rec.RepeatFreq = model.RepeatFreq(req.RepeatFreq)
	}
	if req.RepeatInterval != nil {
		if *req.RepeatInterval < 1 {
			return nil, errcode.ErrBadRequest
		}
		rec.RepeatInterval = *req.RepeatInterval
	}
	if req.NextDate != "" {
		nextDate, err := time.Parse("2006-01-02", req.NextDate)
		if err != nil {
			return nil, errcode.ErrBadRequest
		}
		rec.NextDate = nextDate
	}
	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, errcode.ErrBadRequest
		}
		rec.EndDate = &endDate
	}
	if req.MaxRepetitions != nil {
		rec.MaxRepetitions = req.MaxRepetitions
	}
	if req.IsActive != nil {
		rec.IsActive = *req.IsActive
	}

	if err := s.recurrenceRepo.Update(rec); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(rec), nil
}

// Delete 删除循环交易
// 参数：
//   - userID: 用户ID
//   - id: 循环交易ID
// 返回：
//   - error: 错误信息
func (s *RecurrenceService) Delete(userID, id uint64) error {
	_, err := s.recurrenceRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}

	return s.recurrenceRepo.Delete(id, userID)
}

// Trigger 手动触发循环交易
// 立即创建一笔交易，但不更新NextDate（不推进下次执行日期）
// 参数：
//   - userID: 用户ID
//   - id: 循环交易ID
// 返回：
//   - *response.TransactionResp: 创建的交易信息
//   - error: 错误信息
func (s *RecurrenceService) Trigger(userID, id uint64) (*response.TransactionResp, error) {
	rec, err := s.recurrenceRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if !rec.IsActive {
		return nil, errcode.ErrRecurrenceInactive
	}

	return s.createTransactionFromRecurrence(rec)
}

// ExecuteRecurrence 自动执行循环交易
// 由定时任务调用，创建交易后推进NextDate，递增重复次数，检查结束条件
// 如果超过结束日期或达到最大重复次数，自动停用循环交易
// 参数：
//   - rec: 循环交易模型
// 返回：
//   - error: 错误信息
func (s *RecurrenceService) ExecuteRecurrence(rec *model.Recurrence) error {
	if !rec.IsActive {
		return nil
	}

	// Create transaction from recurrence template
	_, err := s.createTransactionFromRecurrence(rec)
	if err != nil {
		return err
	}

	// Increment repetitions
	rec.Repetitions++

	// Calculate next date
	nextDate := s.calculateNextDate(rec.NextDate, rec.RepeatFreq, rec.RepeatInterval)
	rec.NextDate = nextDate

	// Check end conditions
	if rec.EndDate != nil && nextDate.After(*rec.EndDate) {
		rec.IsActive = false
	}
	if rec.MaxRepetitions != nil && rec.Repetitions >= *rec.MaxRepetitions {
		rec.IsActive = false
	}

	// Update recurrence
	if err := s.recurrenceRepo.Update(rec); err != nil {
		return err
	}

	return nil
}

// createTransactionFromRecurrence 根据循环交易模板创建实际交易
// 通过txnService.Create()确保账户余额更新、规则触发和Webhook通知
func (s *RecurrenceService) createTransactionFromRecurrence(rec *model.Recurrence) (*response.TransactionResp, error) {
	// Parse tag names to tag IDs (empty for now, tags are stored by name)
	var tagIDs []uint64
	if rec.TagNames != "" {
		_ = strings.Split(rec.TagNames, ",")
		// Tag names are stored as comma-separated; we pass empty tagIDs
		// since the transaction service expects tag IDs
	}

	dateStr := rec.NextDate.Format("2006-01-02")

	txnReq := &request.CreateTransactionReq{
		Type:          string(rec.Type),
		Date:          dateStr,
		Description:   rec.Description,
		Amount:        rec.Amount.StringFixed(4),
		SourceID:      rec.SourceID,
		DestinationID: rec.DestinationID,
		CategoryID:    rec.CategoryID,
		Notes:         rec.Notes,
		Tags:          tagIDs,
	}

	return s.txnService.Create(rec.UserID, txnReq)
}

// calculateNextDate 根据频率和间隔计算下次执行日期
// 支持日/周/月/年频率，间隔表示跳过几个周期
func (s *RecurrenceService) calculateNextDate(currentDate time.Time, freq model.RepeatFreq, interval int) time.Time {
	switch freq {
	case model.RepeatFreqDaily:
		return currentDate.AddDate(0, 0, interval)
	case model.RepeatFreqWeekly:
		return currentDate.AddDate(0, 0, interval*7)
	case model.RepeatFreqMonthly:
		return currentDate.AddDate(0, interval, 0)
	case model.RepeatFreqYearly:
		return currentDate.AddDate(interval, 0, 0)
	default:
		return currentDate.AddDate(0, 0, interval)
	}
}

// toResp 将循环交易模型转换为响应对象
// 加载源账户和目标账户的名称
func (s *RecurrenceService) toResp(r *model.Recurrence) *response.RecurrenceResp {
	resp := &response.RecurrenceResp{
		ID:             r.ID,
		UserID:         r.UserID,
		Title:          r.Title,
		Type:           string(r.Type),
		Amount:         r.Amount.StringFixed(4),
		SourceID:       r.SourceID,
		DestinationID:  r.DestinationID,
		CategoryID:     r.CategoryID,
		Description:    r.Description,
		Notes:          r.Notes,
		TagNames:       r.TagNames,
		RepeatFreq:     string(r.RepeatFreq),
		RepeatInterval: r.RepeatInterval,
		NextDate:       r.NextDate,
		EndDate:        r.EndDate,
		Repetitions:    r.Repetitions,
		MaxRepetitions: r.MaxRepetitions,
		IsActive:       r.IsActive,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}

	// Load source account name
	source, err := s.accountRepo.GetByID(r.SourceID, r.UserID)
	if err == nil {
		resp.SourceName = source.Name
	}

	// Load destination account name
	if r.DestinationID != nil {
		dest, err := s.accountRepo.GetByID(*r.DestinationID, r.UserID)
		if err == nil {
			resp.DestinationName = dest.Name
		}
	}

	return resp
}

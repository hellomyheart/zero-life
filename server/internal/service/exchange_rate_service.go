// Package service 业务逻辑层，实现核心业务逻辑
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

// ExchangeRateService 汇率业务服务
// 负责处理汇率相关的业务逻辑，包括汇率的创建、查询、更新、删除等操作
type ExchangeRateService struct {
	rateRepo *repository.ExchangeRateRepository
}

// NewExchangeRateService 创建汇率服务实例
// 参数：
//   - rateRepo: 汇率数据访问对象
// 返回：
//   - *ExchangeRateService: 汇率服务实例
func NewExchangeRateService(rateRepo *repository.ExchangeRateRepository) *ExchangeRateService {
	return &ExchangeRateService{rateRepo: rateRepo}
}

// Create 创建汇率记录
// 业务流程：
// 1. 检查是否已存在相同货币对和日期的汇率
// 2. 解析并验证汇率值
// 3. 创建汇率记录
// 参数：
//   - userID: 用户ID
//   - req: 创建汇率请求参数
// 返回：
//   - *response.ExchangeRateResp: 创建成功的汇率信息
//   - error: 错误信息
func (s *ExchangeRateService) Create(userID uint64, req *request.CreateExchangeRateReq) (*response.ExchangeRateResp, error) {
	// 检查是否已存在
	existing, err := s.rateRepo.GetSpecificRateOnDate(userID, req.FromCurrencyID, req.ToCurrencyID, req.Date)
	if err == nil && existing.ID > 0 {
		return nil, errcode.ErrDuplicate
	}

	rate, err := decimal.NewFromString(req.Rate)
	if err != nil {
		return nil, errcode.ErrBadRequest
	}

	exchangeRate := &model.ExchangeRate{
		UserID:         userID,
		FromCurrencyID: req.FromCurrencyID,
		ToCurrencyID:   req.ToCurrencyID,
		Date:           req.Date,
		Rate:           rate,
	}

	if req.UserRate != "" {
		userRate, err := decimal.NewFromString(req.UserRate)
		if err != nil {
			return nil, errcode.ErrBadRequest
		}
		exchangeRate.UserRate = userRate
	}

	if err := s.rateRepo.Create(exchangeRate); err != nil {
		return nil, errcode.ErrInternal
	}

	// 重新加载关联数据
	created, err := s.rateRepo.GetByID(exchangeRate.ID, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(created), nil
}

// Get 获取汇率详情
// 参数：
//   - userID: 用户ID
//   - id: 汇率ID
// 返回：
//   - *response.ExchangeRateResp: 汇率信息
//   - error: 错误信息
func (s *ExchangeRateService) Get(userID, id uint64) (*response.ExchangeRateResp, error) {
	rate, err := s.rateRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(rate), nil
}

// List 获取汇率列表
// 支持按货币对、日期范围过滤和分页
// 参数：
//   - userID: 用户ID
//   - req: 列表查询参数
// 返回：
//   - *pagination.Result: 分页结果
//   - error: 错误信息
func (s *ExchangeRateService) List(userID uint64, req *request.ExchangeRateListReq) (*pagination.Result, error) {
	params := pagination.Params{Page: req.Page, PageSize: req.PageSize}
	params.Normalize()

	var startDate, endDate *time.Time
	if req.StartDate != "" {
		t, err := time.Parse("2006-01-02", req.StartDate)
		if err == nil {
			startDate = &t
		}
	}
	if req.EndDate != "" {
		t, err := time.Parse("2006-01-02", req.EndDate)
		if err == nil {
			endDate = &t
		}
	}

	rates, err := s.rateRepo.List(userID, req.FromCurrencyID, req.ToCurrencyID, startDate, endDate, params.Offset(), params.PageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	total, err := s.rateRepo.Count(userID, req.FromCurrencyID, req.ToCurrencyID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.ExchangeRateResp, 0, len(rates))
	for _, r := range rates {
		items = append(items, *s.toResp(&r))
	}

	return pagination.NewResult(items, total, params), nil
}

// Update 更新汇率信息
// 参数：
//   - userID: 用户ID
//   - id: 汇率ID
//   - req: 更新请求参数
// 返回：
//   - *response.ExchangeRateResp: 更新后的汇率信息
//   - error: 错误信息
func (s *ExchangeRateService) Update(userID, id uint64, req *request.UpdateExchangeRateReq) (*response.ExchangeRateResp, error) {
	rate, err := s.rateRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Rate != "" {
		r, err := decimal.NewFromString(req.Rate)
		if err != nil {
			return nil, errcode.ErrBadRequest
		}
		rate.Rate = r
	}

	if req.UserRate != "" {
		userRate, err := decimal.NewFromString(req.UserRate)
		if err != nil {
			return nil, errcode.ErrBadRequest
		}
		rate.UserRate = userRate
	}

	if err := s.rateRepo.Update(rate); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(rate), nil
}

// Delete 删除汇率记录
// 参数：
//   - userID: 用户ID
//   - id: 汇率ID
// 返回：
//   - error: 错误信息
func (s *ExchangeRateService) Delete(userID, id uint64) error {
	return s.rateRepo.Delete(id, userID)
}

// GetLatest 获取最新汇率
// 参数：
//   - userID: 用户ID
//   - fromCurrencyID: 源货币ID
//   - toCurrencyID: 目标货币ID
// 返回：
//   - *response.ExchangeRateResp: 汇率信息
//   - error: 错误信息
func (s *ExchangeRateService) GetLatest(userID, fromCurrencyID, toCurrencyID uint64) (*response.ExchangeRateResp, error) {
	rate, err := s.rateRepo.GetLatestRate(userID, fromCurrencyID, toCurrencyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(rate), nil
}

// Convert 货币转换
// 根据汇率将金额从源货币转换为目标货币
// 参数：
//   - userID: 用户ID
//   - fromCurrencyID: 源货币ID
//   - toCurrencyID: 目标货币ID
//   - amount: 金额
//   - date: 日期（可选，不指定则使用最新汇率）
// 返回：
//   - string: 转换后的金额
//   - error: 错误信息
func (s *ExchangeRateService) Convert(userID, fromCurrencyID, toCurrencyID uint64, amount string, date *time.Time) (string, error) {
	amt, err := decimal.NewFromString(amount)
	if err != nil {
		return "", errcode.ErrBadRequest
	}

	var rate *model.ExchangeRate
	if date != nil {
		rate, err = s.rateRepo.GetSpecificRateOnDate(userID, fromCurrencyID, toCurrencyID, *date)
	} else {
		rate, err = s.rateRepo.GetLatestRate(userID, fromCurrencyID, toCurrencyID)
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errcode.ErrNotFound
		}
		return "", errcode.ErrInternal
	}

	// 优先使用用户自定义汇率
	rateValue := rate.Rate
	if !rate.UserRate.IsZero() {
		rateValue = rate.UserRate
	}

	result := amt.Mul(rateValue)
	return result.StringFixed(4), nil
}

// toResp 将汇率模型转换为响应对象
// 参数：
//   - rate: 汇率模型
// 返回：
//   - *response.ExchangeRateResp: 汇率响应对象
func (s *ExchangeRateService) toResp(rate *model.ExchangeRate) *response.ExchangeRateResp {
	resp := &response.ExchangeRateResp{
		ID:             rate.ID,
		UserID:         rate.UserID,
		FromCurrencyID: rate.FromCurrencyID,
		ToCurrencyID:   rate.ToCurrencyID,
		Date:           rate.Date.Format("2006-01-02"),
		Rate:           rate.Rate.StringFixed(12),
		CreatedAt:      rate.CreatedAt,
		UpdatedAt:      rate.UpdatedAt,
	}

	if !rate.UserRate.IsZero() {
		resp.UserRate = rate.UserRate.StringFixed(12)
	}

	if rate.FromCurrency.ID > 0 {
		resp.FromCurrency = response.CurrencyResp{
			ID:            rate.FromCurrency.ID,
			Code:          rate.FromCurrency.Code,
			Name:          rate.FromCurrency.Name,
			Symbol:        rate.FromCurrency.Symbol,
			DecimalPlaces: rate.FromCurrency.DecimalPlaces,
		}
	}

	if rate.ToCurrency.ID > 0 {
		resp.ToCurrency = response.CurrencyResp{
			ID:            rate.ToCurrency.ID,
			Code:          rate.ToCurrency.Code,
			Name:          rate.ToCurrency.Name,
			Symbol:        rate.ToCurrency.Symbol,
			DecimalPlaces: rate.ToCurrency.DecimalPlaces,
		}
	}

	return resp
}

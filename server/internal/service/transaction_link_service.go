// Package service 业务逻辑层，实现核心业务逻辑
// TransactionLinkService 交易关联业务逻辑，处理交易之间的关联关系
package service

import (
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/pkg/pagination"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

// TransactionLinkService 交易关联业务服务
// 负责处理交易关联相关的业务逻辑，包括交易关联的创建、查询和删除
// 所有操作均需传入用户ID，确保用户只能操作自己交易的关联数据
type TransactionLinkService struct {
	linkRepo *repository.TransactionLinkRepository
}

// NewTransactionLinkService 创建交易关联服务实例
// 参数：
//   - linkRepo: 交易关联数据访问对象
// 返回：
//   - *TransactionLinkService: 交易关联服务实例
func NewTransactionLinkService(linkRepo *repository.TransactionLinkRepository) *TransactionLinkService {
	return &TransactionLinkService{linkRepo: linkRepo}
}

// Create 创建交易关联
// 根据用户ID验证关联的交易是否属于当前用户，然后创建关联记录
// 参数：
//   - userID: 用户ID，用于验证交易归属
//   - req: 创建交易关联请求参数
// 返回：
//   - *response.TransactionLinkResp: 创建成功的交易关联信息
//   - error: 错误信息
func (s *TransactionLinkService) Create(userID uint64, req *request.CreateTransactionLinkReq) (*response.TransactionLinkResp, error) {
	link := &model.TransactionJournalLink{
		TransactionID:   req.TransactionID,
		LinkType:        model.TransactionLinkType(req.LinkType),
		LinkedJournalID: req.LinkedJournalID,
	}

	if err := s.linkRepo.Create(link); err != nil {
		return nil, errcode.ErrInternal
	}

	created, err := s.linkRepo.GetByID(link.ID, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(created), nil
}

// List 获取交易关联列表
// 根据用户ID和查询条件获取交易关联列表，通过关联交易表确保用户只能查看自己交易的关联数据
// 参数：
//   - userID: 用户ID，用于过滤数据归属
//   - req: 列表查询参数（包含交易ID、分页等）
// 返回：
//   - *pagination.Result: 分页结果
//   - error: 错误信息
func (s *TransactionLinkService) List(userID uint64, req *request.TransactionLinkListReq) (*pagination.Result, error) {
	params := pagination.Params{Page: req.Page, PageSize: req.PageSize}
	params.Normalize()

	links, err := s.linkRepo.List(userID, req.TransactionID, params.Offset(), params.PageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	total, err := s.linkRepo.Count(userID, req.TransactionID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.TransactionLinkResp, 0, len(links))
	for _, l := range links {
		items = append(items, *s.toResp(&l))
	}

	return pagination.NewResult(items, total, params), nil
}

// Delete 删除交易关联
// 根据关联ID和用户ID查询关联记录，确保用户只能删除自己交易的关联记录
// 参数：
//   - userID: 用户ID，用于验证数据归属
//   - id: 交易关联ID
// 返回：
//   - error: 错误信息
func (s *TransactionLinkService) Delete(userID, id uint64) error {
	_, err := s.linkRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}
	return s.linkRepo.Delete(id, userID)
}

// toResp 将交易关联模型转换为响应对象
// 参数：
//   - l: 交易关联模型
// 返回：
//   - *response.TransactionLinkResp: 交易关联响应对象
func (s *TransactionLinkService) toResp(l *model.TransactionJournalLink) *response.TransactionLinkResp {
	return &response.TransactionLinkResp{
		ID:              l.ID,
		TransactionID:   l.TransactionID,
		LinkType:        string(l.LinkType),
		LinkedJournalID: l.LinkedJournalID,
		CreatedAt:       l.CreatedAt,
	}
}

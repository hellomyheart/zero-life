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

type TransactionLinkService struct {
	linkRepo *repository.TransactionLinkRepository
}

func NewTransactionLinkService(linkRepo *repository.TransactionLinkRepository) *TransactionLinkService {
	return &TransactionLinkService{linkRepo: linkRepo}
}

func (s *TransactionLinkService) Create(req *request.CreateTransactionLinkReq) (*response.TransactionLinkResp, error) {
	link := &model.TransactionJournalLink{
		TransactionID:   req.TransactionID,
		LinkType:        model.TransactionLinkType(req.LinkType),
		LinkedJournalID: req.LinkedJournalID,
	}

	if err := s.linkRepo.Create(link); err != nil {
		return nil, errcode.ErrInternal
	}

	created, err := s.linkRepo.GetByID(link.ID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(created), nil
}

func (s *TransactionLinkService) List(req *request.TransactionLinkListReq) (*pagination.Result, error) {
	params := pagination.Params{Page: req.Page, PageSize: req.PageSize}
	params.Normalize()

	links, err := s.linkRepo.List(req.TransactionID, params.Offset(), params.PageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	total, err := s.linkRepo.Count(req.TransactionID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.TransactionLinkResp, 0, len(links))
	for _, l := range links {
		items = append(items, *s.toResp(&l))
	}

	return pagination.NewResult(items, total, params), nil
}

func (s *TransactionLinkService) Delete(id uint64) error {
	_, err := s.linkRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}
	return s.linkRepo.Delete(id)
}

func (s *TransactionLinkService) toResp(l *model.TransactionJournalLink) *response.TransactionLinkResp {
	return &response.TransactionLinkResp{
		ID:              l.ID,
		TransactionID:   l.TransactionID,
		LinkType:        string(l.LinkType),
		LinkedJournalID: l.LinkedJournalID,
		CreatedAt:       l.CreatedAt,
	}
}
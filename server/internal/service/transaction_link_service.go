package service

import (
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

type TransactionLinkService struct {
	txnLinkRepo *repository.TransactionLinkRepository
}

func NewTransactionLinkService(txnLinkRepo *repository.TransactionLinkRepository) *TransactionLinkService {
	return &TransactionLinkService{txnLinkRepo: txnLinkRepo}
}

func (s *TransactionLinkService) Create(req *request.CreateTransactionLinkReq) (*response.TransactionLinkResp, error) {
	// Check for duplicate link
	exists, err := s.txnLinkRepo.Exists(req.LinkTypeID, req.SourceID, req.DestinationID)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	if exists {
		return nil, errcode.ErrTxnLinkDuplicate
	}

	link := &model.TransactionLink{
		LinkTypeID:    req.LinkTypeID,
		SourceID:      req.SourceID,
		DestinationID: req.DestinationID,
		Comment:       req.Comment,
	}

	if err := s.txnLinkRepo.Create(link); err != nil {
		return nil, errcode.ErrInternal
	}

	// Reload with LinkType
	created, err := s.txnLinkRepo.GetByID(link.ID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.toResp(created), nil
}

func (s *TransactionLinkService) List(transactionID uint64) ([]response.TransactionLinkResp, error) {
	links, err := s.txnLinkRepo.ListByTransactionID(transactionID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.TransactionLinkResp, 0, len(links))
	for _, l := range links {
		items = append(items, *s.toResp(&l))
	}

	return items, nil
}

func (s *TransactionLinkService) Delete(id uint64) error {
	_, err := s.txnLinkRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}
	return s.txnLinkRepo.Delete(id)
}

func (s *TransactionLinkService) toResp(l *model.TransactionLink) *response.TransactionLinkResp {
	resp := &response.TransactionLinkResp{
		ID:            l.ID,
		LinkTypeID:    l.LinkTypeID,
		SourceID:      l.SourceID,
		DestinationID: l.DestinationID,
		Comment:       l.Comment,
		CreatedAt:     l.CreatedAt,
	}
	if l.LinkType.ID > 0 {
		resp.LinkType = response.LinkTypeResp{
			ID:            l.LinkType.ID,
			Name:          l.LinkType.Name,
			Outward:       l.LinkType.Outward,
			Inward:        l.LinkType.Inward,
			IsDirectional: l.LinkType.IsDirectional,
			CreatedAt:     l.LinkType.CreatedAt,
			UpdatedAt:     l.LinkType.UpdatedAt,
		}
	}
	return resp
}

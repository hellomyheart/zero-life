package service

import (
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

type PreferenceService struct {
	prefRepo *repository.PreferenceRepository
}

func NewPreferenceService(prefRepo *repository.PreferenceRepository) *PreferenceService {
	return &PreferenceService{prefRepo: prefRepo}
}

func (s *PreferenceService) Get(userID uint64, key string) (*response.PreferenceResp, error) {
	pref, err := s.prefRepo.Get(userID, key)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return &response.PreferenceResp{Key: pref.Key, Value: pref.Value}, nil
}

func (s *PreferenceService) Set(userID uint64, req *request.SetPreferenceReq) (*response.PreferenceResp, error) {
	if err := s.prefRepo.Set(userID, req.Key, req.Value); err != nil {
		return nil, errcode.ErrInternal
	}
	return &response.PreferenceResp{Key: req.Key, Value: req.Value}, nil
}

func (s *PreferenceService) List(userID uint64) ([]response.PreferenceResp, error) {
	prefs, err := s.prefRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.PreferenceResp, 0, len(prefs))
	for _, p := range prefs {
		items = append(items, response.PreferenceResp{Key: p.Key, Value: p.Value})
	}
	return items, nil
}

func (s *PreferenceService) Delete(userID uint64, key string) error {
	return s.prefRepo.Delete(userID, key)
}
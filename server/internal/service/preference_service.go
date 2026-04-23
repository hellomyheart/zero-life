package service

import (
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

// Predefined preference keys and their default values
var defaultPreferences = map[string]string{
	"default_currency":    "CNY",
	"fiscal_year_start":  "1",
	"list_page_size":     "20",
	"frontpage_accounts": "",
}

type PreferenceService struct {
	prefRepo *repository.PreferenceRepository
}

func NewPreferenceService(prefRepo *repository.PreferenceRepository) *PreferenceService {
	return &PreferenceService{prefRepo: prefRepo}
}

func (s *PreferenceService) Get(userID uint64, name string) (*response.PreferenceResp, error) {
	pref, err := s.prefRepo.Get(userID, name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return default value if preference not set
			if defaultVal, ok := defaultPreferences[name]; ok {
				return &response.PreferenceResp{
					Name:  name,
					Value: defaultVal,
				}, nil
			}
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return s.toResp(pref), nil
}

func (s *PreferenceService) GetWithDefault(userID uint64, name string) string {
	pref, err := s.prefRepo.Get(userID, name)
	if err != nil {
		if defaultVal, ok := defaultPreferences[name]; ok {
			return defaultVal
		}
		return ""
	}
	return pref.Value
}

func (s *PreferenceService) Set(userID uint64, name, value string) (*response.PreferenceResp, error) {
	if err := s.prefRepo.Set(userID, name, value); err != nil {
		return nil, errcode.ErrInternal
	}

	pref, err := s.prefRepo.Get(userID, name)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	return s.toResp(pref), nil
}

func (s *PreferenceService) List(userID uint64) ([]response.PreferenceResp, error) {
	prefs, err := s.prefRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.PreferenceResp, 0, len(prefs))
	for _, p := range prefs {
		items = append(items, *s.toResp(&p))
	}

	return items, nil
}

func (s *PreferenceService) toResp(p *model.Preference) *response.PreferenceResp {
	return &response.PreferenceResp{
		ID:        p.ID,
		Name:      p.Name,
		Value:     p.Value,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

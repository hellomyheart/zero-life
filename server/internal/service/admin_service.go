package service

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/email"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/pkg/hash"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

type AdminService struct {
	userRepo  *repository.UserRepository
	configRepo *repository.ConfigurationRepository
}

func NewAdminService(
	userRepo *repository.UserRepository,
	configRepo *repository.ConfigurationRepository,
) *AdminService {
	return &AdminService{
		userRepo:  userRepo,
		configRepo: configRepo,
	}
}

// ListUsers returns a paginated list of users.
func (s *AdminService) ListUsers(req *request.AdminListUsersReq) ([]response.AdminUserResp, int64, error) {
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	users, total, err := s.userRepo.List(req.Search, offset, pageSize)
	if err != nil {
		return nil, 0, errcode.ErrInternal
	}

	items := make([]response.AdminUserResp, 0, len(users))
	for _, u := range users {
		items = append(items, s.toUserResp(&u))
	}
	return items, total, nil
}

// UpdateUser updates a user's profile (admin operation).
func (s *AdminService) UpdateUser(userID uint64, req *request.AdminUpdateUserReq) (*response.AdminUserResp, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Language != "" {
		user.Language = req.Language
	}
	if req.Timezone != "" {
		user.Timezone = req.Timezone
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errcode.ErrInternal
	}

	return &response.AdminUserResp{
		ID:        user.ID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Role:      user.Role,
		Language:  user.Language,
		Timezone:  user.Timezone,
		MFAEnabled: user.MFAEnabled,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// DeleteUser soft-deletes a user.
func (s *AdminService) DeleteUser(userID uint64) error {
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}
	return s.userRepo.Delete(userID)
}

// InviteUser creates a new user with a random password (they must reset it).
func (s *AdminService) InviteUser(req *request.AdminInviteUserReq) (*response.AdminUserResp, error) {
	// Check email uniqueness
	existing, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existing != nil {
		return nil, errcode.ErrEmailExists
	}

	// Generate a random temporary password
	tempPassword, err := hash.HashPassword("TempPass_" + randomHex(16))
	if err != nil {
		return nil, errcode.ErrInternal
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	user := &model.User{
		Email:    req.Email,
		Password: tempPassword,
		Nickname: req.Nickname,
		Role:     role,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, errcode.ErrInternal
	}

	return &response.AdminUserResp{
		ID:        user.ID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Role:      user.Role,
		Language:  user.Language,
		Timezone:  user.Timezone,
		MFAEnabled: user.MFAEnabled,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetConfiguration returns a configuration value by name.
func (s *AdminService) GetConfiguration(name string) (*response.AdminConfigurationResp, error) {
	cfg, err := s.configRepo.Get(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}
	return &response.AdminConfigurationResp{
		ID:        cfg.ID,
		Name:      cfg.Name,
		Value:     cfg.Value,
		CreatedAt: cfg.CreatedAt,
		UpdatedAt: cfg.UpdatedAt,
	}, nil
}

// ListConfigurations returns all configurations.
func (s *AdminService) ListConfigurations() ([]response.AdminConfigurationResp, error) {
	configs, err := s.configRepo.List()
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.AdminConfigurationResp, 0, len(configs))
	for _, cfg := range configs {
		items = append(items, response.AdminConfigurationResp{
			ID:        cfg.ID,
			Name:      cfg.Name,
			Value:     cfg.Value,
			CreatedAt: cfg.CreatedAt,
			UpdatedAt: cfg.UpdatedAt,
		})
	}
	return items, nil
}

// UpdateConfiguration sets a configuration value.
func (s *AdminService) UpdateConfiguration(name string, req *request.AdminUpdateConfigurationReq) (*response.AdminConfigurationResp, error) {
	cfg := &model.Configuration{
		Name:  name,
		Value: req.Value,
	}
	if err := s.configRepo.Set(cfg); err != nil {
		return nil, errcode.ErrInternal
	}

	return s.GetConfiguration(name)
}

// TestEmail sends a test email to verify SMTP configuration.
func (s *AdminService) TestEmail(to string) (*response.AdminTestEmailResp, error) {
	err := email.Send(to, "Test Email from ZeroLife", `
	<html>
	<body>
		<h2>Test Email</h2>
		<p>This is a test email from ZeroLife to verify SMTP configuration.</p>
	</body>
	</html>`)
	if err != nil {
		return &response.AdminTestEmailResp{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	return &response.AdminTestEmailResp{
		Success: true,
		Message: "Test email sent successfully",
	}, nil
}

func (s *AdminService) toUserResp(u *model.User) response.AdminUserResp {
	return response.AdminUserResp{
		ID:        u.ID,
		Email:     u.Email,
		Nickname:  u.Nickname,
		Role:      u.Role,
		Language:  u.Language,
		Timezone:  u.Timezone,
		MFAEnabled: u.MFAEnabled,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

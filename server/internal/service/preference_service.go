// Package service 业务逻辑层，实现核心业务逻辑
// PreferenceService 偏好设置业务逻辑，管理用户个性化配置
package service

import (
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

// PreferenceService 偏好设置服务
// 管理用户个性化配置，如默认货币、主题、语言等键值对设置
// 依赖prefRepo进行偏好设置数据访问
type PreferenceService struct {
	prefRepo *repository.PreferenceRepository // 偏好设置数据访问对象
}

// NewPreferenceService 创建偏好设置服务实例
func NewPreferenceService(prefRepo *repository.PreferenceRepository) *PreferenceService {
	return &PreferenceService{prefRepo: prefRepo}
}

// Get 获取单个偏好设置
// 参数：
//   - userID: 用户ID
//   - key: 偏好设置键名
// 返回：
//   - *response.PreferenceResp: 偏好设置信息
//   - error: 错误信息
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

// Set 设置偏好设置（不存在则创建，存在则更新）
// 参数：
//   - userID: 用户ID
//   - req: 设置请求参数（键名、值）
// 返回：
//   - *response.PreferenceResp: 偏好设置信息
//   - error: 错误信息
func (s *PreferenceService) Set(userID uint64, req *request.SetPreferenceReq) (*response.PreferenceResp, error) {
	if err := s.prefRepo.Set(userID, req.Key, req.Value); err != nil {
		return nil, errcode.ErrInternal
	}
	return &response.PreferenceResp{Key: req.Key, Value: req.Value}, nil
}

// List 获取用户所有偏好设置
// 参数：
//   - userID: 用户ID
// 返回：
//   - []response.PreferenceResp: 偏好设置列表
//   - error: 错误信息
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

// Delete 删除偏好设置
// 参数：
//   - userID: 用户ID
//   - key: 偏好设置键名
// 返回：
//   - error: 错误信息
func (s *PreferenceService) Delete(userID uint64, key string) error {
	return s.prefRepo.Delete(userID, key)
}
// Package service 业务逻辑层，实现核心业务逻辑
// TagService 标签业务逻辑，处理标签的增删改查
package service

import (
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"gorm.io/gorm"
)

// TagService 标签服务
// 负责处理标签的增删改查，标签用于对交易进行分类标记
// 依赖tagRepo进行标签数据访问
type TagService struct {
	tagRepo *repository.TagRepository // 标签数据访问对象
}

// NewTagService 创建标签服务实例
func NewTagService(tagRepo *repository.TagRepository) *TagService {
	return &TagService{tagRepo: tagRepo}
}

// Create 创建标签
// 检查名称唯一性，如果未指定颜色则默认使用 #409EFF
// 参数：
//   - userID: 用户ID
//   - req: 创建请求参数（名称、颜色）
// 返回：
//   - *response.TagResp: 创建成功的标签信息（含关联交易数）
//   - error: 错误信息（如名称重复）
func (s *TagService) Create(userID uint64, req *request.CreateTagReq) (*response.TagResp, error) {
	// Check name uniqueness
	tags, err := s.tagRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	for _, t := range tags {
		if t.Name == req.Name {
			return nil, errcode.ErrTagNameExists
		}
	}

	tag := &model.Tag{
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
	}
	if tag.Color == "" {
		tag.Color = "#409EFF"
	}

	if err := s.tagRepo.Create(tag); err != nil {
		return nil, errcode.ErrInternal
	}

	count, _ := s.tagRepo.CountTransactions(tag.ID)

	return &response.TagResp{
		ID:               tag.ID,
		Name:             tag.Name,
		Color:            tag.Color,
		TransactionCount: count,
		CreatedAt:        tag.CreatedAt,
		UpdatedAt:        tag.UpdatedAt,
	}, nil
}

// List 获取用户所有标签列表（含每个标签的关联交易数）
// 参数：
//   - userID: 用户ID
// 返回：
//   - []response.TagResp: 标签列表
//   - error: 错误信息
func (s *TagService) List(userID uint64) ([]response.TagResp, error) {
	tags, err := s.tagRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	items := make([]response.TagResp, 0, len(tags))
	for _, t := range tags {
		count, _ := s.tagRepo.CountTransactions(t.ID)
		items = append(items, response.TagResp{
			ID:               t.ID,
			Name:             t.Name,
			Color:            t.Color,
			TransactionCount: count,
			CreatedAt:        t.CreatedAt,
			UpdatedAt:        t.UpdatedAt,
		})
	}

	return items, nil
}

// Update 更新标签信息（名称和颜色）
// 参数：
//   - userID: 用户ID
//   - id: 标签ID
//   - req: 更新请求参数
// 返回：
//   - *response.TagResp: 更新后的标签信息（含关联交易数）
//   - error: 错误信息
func (s *TagService) Update(userID, id uint64, req *request.UpdateTagReq) (*response.TagResp, error) {
	tag, err := s.tagRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	if req.Name != "" {
		tag.Name = req.Name
	}
	if req.Color != "" {
		tag.Color = req.Color
	}

	if err := s.tagRepo.Update(tag); err != nil {
		return nil, errcode.ErrInternal
	}

	count, _ := s.tagRepo.CountTransactions(tag.ID)

	return &response.TagResp{
		ID:               tag.ID,
		Name:             tag.Name,
		Color:            tag.Color,
		TransactionCount: count,
		CreatedAt:        tag.CreatedAt,
		UpdatedAt:        tag.UpdatedAt,
	}, nil
}

// Delete 删除标签
// 参数：
//   - userID: 用户ID
//   - id: 标签ID
// 返回：
//   - error: 错误信息
func (s *TagService) Delete(userID, id uint64) error {
	_, err := s.tagRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrInternal
	}

	return s.tagRepo.Delete(id, userID)
}

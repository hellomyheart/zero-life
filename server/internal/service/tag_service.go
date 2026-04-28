// Package service 业务逻辑层，实现核心业务逻辑
// TagService 标签业务逻辑，处理标签的增删改查（支持树形结构）
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
// 支持最多5级树形结构（通过 parent_id 形成层级关系），与分类类似
// 依赖tagRepo进行标签数据访问，依赖db执行删除事务
type TagService struct {
	tagRepo *repository.TagRepository // 标签数据访问对象
	db      *gorm.DB                  // 数据库连接，用于删除标签时清理子标签
}

// NewTagService 创建标签服务实例
func NewTagService(tagRepo *repository.TagRepository, db *gorm.DB) *TagService {
	return &TagService{tagRepo: tagRepo, db: db}
}

// Create 创建标签
// 检查名称唯一性和5级深度限制，如果未指定颜色则默认使用 #409EFF
// 参数：
//   - userID: 用户ID
//   - req: 创建请求参数（名称、颜色、父标签ID）
// 返回：
//   - *response.TagResp: 创建成功的标签信息（含关联交易数）
//   - error: 错误信息（如名称重复、层级过深）
func (s *TagService) Create(userID uint64, req *request.CreateTagReq) (*response.TagResp, error) {
	// 检查名称唯一性
	tags, err := s.tagRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	for _, t := range tags {
		if t.Name == req.Name {
			return nil, errcode.ErrTagNameExists
		}
	}

	// 检查5级深度限制：沿 parent 链向上计算深度
	if req.ParentID != nil {
		depth := 1
		pid := *req.ParentID
		for pid != 0 {
			parent, err := s.tagRepo.GetByID(pid, userID)
			if err != nil {
				return nil, errcode.ErrNotFound
			}
			if parent.ParentID == nil {
				break
			}
			depth++
			pid = *parent.ParentID
		}
		if depth >= 5 {
			return nil, errcode.ErrTagTooDeep
		}
	}

	tag := &model.Tag{
		UserID:   userID,
		Name:     req.Name,
		Color:    req.Color,
		ParentID: req.ParentID,
	}
	if tag.Color == "" {
		tag.Color = "#409EFF"
	}

	if err := s.tagRepo.Create(tag); err != nil {
		return nil, errcode.ErrInternal
	}

	count, _ := s.tagRepo.CountTransactions(tag.ID)

	return s.toResp(tag, count), nil
}

// List 获取用户所有标签的树形列表
// 返回两级树形结构的标签数据，含每个标签的关联交易数
// 参数：
//   - userID: 用户ID
// 返回：
//   - []response.TagResp: 标签树形列表
//   - error: 错误信息
func (s *TagService) List(userID uint64) ([]response.TagResp, error) {
	tags, err := s.tagRepo.List(userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	return s.buildTree(tags), nil
}

// Get 获取单个标签详情
// 参数：
//   - userID: 用户ID
//   - id: 标签ID
// 返回：
//   - *response.TagResp: 标签信息（含关联交易数）
//   - error: 错误信息
func (s *TagService) Get(userID, id uint64) (*response.TagResp, error) {
	tag, err := s.tagRepo.GetByID(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrInternal
	}

	count, _ := s.tagRepo.CountTransactions(tag.ID)
	return s.toResp(tag, count), nil
}

// Update 更新标签信息（名称、颜色、父标签ID）
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
	// 更新父标签ID时需要检查5级深度限制
	if req.ParentID != nil {
		// 不能将自己设为自己的子标签
		if *req.ParentID == id {
			return nil, errcode.ErrTagTooDeep
		}
		// 检查是否会形成循环：新父标签不能是当前标签的后代
		descendantIDs, _ := s.tagRepo.GetDescendantIDs([]uint64{id}, userID)
		for _, did := range descendantIDs {
			if did == *req.ParentID {
				return nil, errcode.ErrTagTooDeep
			}
		}
		// 计算新父标签的深度
		newParentDepth := 0
		pid := *req.ParentID
		for pid != 0 {
			parent, err := s.tagRepo.GetByID(pid, userID)
			if err != nil {
				return nil, errcode.ErrNotFound
			}
			newParentDepth++
			if parent.ParentID == nil {
				break
			}
			pid = *parent.ParentID
		}
		// 当前标签作为子树的最大深度
		subtreeDepth := s.getSubtreeDepth(id, userID)
		// 新父标签深度 + 子树深度 不能超过5
		if newParentDepth+subtreeDepth > 5 {
			return nil, errcode.ErrTagTooDeep
		}
		tag.ParentID = req.ParentID
	}

	if err := s.tagRepo.Update(tag); err != nil {
		return nil, errcode.ErrInternal
	}

	count, _ := s.tagRepo.CountTransactions(tag.ID)
	return s.toResp(tag, count), nil
}

// Delete 删除标签
// 级联删除所有子孙标签，以及关联的 transaction_tags 记录
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

	// 获取所有后代标签ID（包含自身）
	descendantIDs, err := s.tagRepo.GetDescendantIDs([]uint64{id}, userID)
	if err != nil {
		return errcode.ErrInternal
	}

	// 按从深到浅的顺序删除：先删后代，再删自身
	// GetDescendantIDs 返回顺序为 BFS（自身→子→孙），反序删除即可
	for i := len(descendantIDs) - 1; i >= 0; i-- {
		if err := s.tagRepo.Delete(descendantIDs[i], userID); err != nil {
			return errcode.ErrInternal
		}
	}

	return nil
}

// toResp 将标签模型转换为响应对象
func (s *TagService) toResp(t *model.Tag, transactionCount int64) *response.TagResp {
	return &response.TagResp{
		ID:               t.ID,
		Name:             t.Name,
		Color:            t.Color,
		ParentID:         t.ParentID,
		TransactionCount: transactionCount,
		CreatedAt:        t.CreatedAt,
		UpdatedAt:        t.UpdatedAt,
	}
}

// buildTree 将扁平的标签列表构建为树形结构
// 使用map快速查找，将子标签挂载到父标签的Children字段
// 使用指针切片避免值语义导致深层子节点丢失
func (s *TagService) buildTree(tags []model.Tag) []response.TagResp {
	// 预计算每个标签的关联交易数
	countMap := make(map[uint64]int64, len(tags))
	for _, t := range tags {
		count, _ := s.tagRepo.CountTransactions(t.ID)
		countMap[t.ID] = count
	}

	// 创建节点映射（使用指针）
	nodeMap := make(map[uint64]*response.TagResp, len(tags))
	for _, t := range tags {
		nodeMap[t.ID] = &response.TagResp{
			ID:               t.ID,
			Name:             t.Name,
			Color:            t.Color,
			ParentID:         t.ParentID,
			TransactionCount: countMap[t.ID],
			CreatedAt:        t.CreatedAt,
			UpdatedAt:        t.UpdatedAt,
		}
	}

	// 将子标签指针挂载到父标签的Children字段
	for _, t := range tags {
		if t.ParentID != nil {
			if parent, ok := nodeMap[*t.ParentID]; ok {
				parent.Children = append(parent.Children, nodeMap[t.ID])
			}
		}
	}

	// 收集根节点（parent_id 为 nil 的标签）
	var roots []response.TagResp
	for _, t := range tags {
		if t.ParentID == nil {
			roots = append(roots, *nodeMap[t.ID])
		}
	}

	return roots
}

// getSubtreeDepth 计算以指定标签为根的子树深度（含自身）
// 用于更新父标签时检查5级深度限制
func (s *TagService) getSubtreeDepth(id uint64, userID uint64) int {
	maxDepth := 1
	var currentLevel []uint64
	currentLevel = append(currentLevel, id)

	for len(currentLevel) > 0 {
		var nextLevel []uint64
		for _, pid := range currentLevel {
			subTags, _ := s.tagRepo.GetSubTags(pid, userID)
			for _, sub := range subTags {
				nextLevel = append(nextLevel, sub.ID)
			}
		}
		if len(nextLevel) > 0 {
			maxDepth++
			currentLevel = nextLevel
		} else {
			break
		}
	}

	return maxDepth
}

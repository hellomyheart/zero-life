// Package repository 提供数据访问层，封装所有与数据库的交互逻辑。
// 每个 Repository 结构体对应一个数据库表，提供增删改查等基本操作。
// 所有 Repository 都依赖 GORM 作为 ORM 框架，通过 *gorm.DB 执行数据库操作。
package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// AttachmentRepository 附件仓库，负责附件的数据访问。
// 附件（Attachment）用于为各种对象（如交易、账单等）添加文件附件，如收据照片、合同扫描件等。
// 通过 attachable_type 和 attachable_id 实现多态关联，即一个附件可以关联不同类型的对象。
// 例如：attachable_type="Transaction" + attachable_id=123 表示该附件关联 ID 为 123 的交易。
// 所有方法都通过 user_id 参数进行数据隔离，确保用户只能访问自己的附件。
type AttachmentRepository struct {
	db *gorm.DB
}

// NewAttachmentRepository 创建附件仓库实例。
// 参数 db: GORM 数据库连接实例。
func NewAttachmentRepository(db *gorm.DB) *AttachmentRepository {
	return &AttachmentRepository{db: db}
}

// Create 创建一条新的附件记录。
// 执行 SQL: INSERT INTO attachments (...)
// 参数 attachment: 要创建的附件对象，GORM 会自动填充 ID、CreatedAt 等字段。
// 返回: 创建失败时返回错误。
func (r *AttachmentRepository) Create(attachment *model.Attachment) error {
	return r.db.Create(attachment).Error
}

// GetByID 根据 ID 和用户 ID 获取单条附件。
// 同时验证该附件属于指定用户（权限校验），确保数据隔离。
// 执行 SQL: SELECT * FROM attachments WHERE id = ? AND user_id = ? LIMIT 1
// 参数 id: 附件 ID。
// 参数 userID: 当前登录用户 ID，用于权限校验，确保用户只能查看自己的附件。
// 返回: 找到的附件对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *AttachmentRepository) GetByID(id, userID uint64) (*model.Attachment, error) {
	var attachment model.Attachment
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&attachment).Error; err != nil {
		return nil, err
	}
	return &attachment, nil
}

// List 获取指定用户的附件列表，可选按关联对象类型和 ID 过滤，按创建时间倒序排列。
// 执行 SQL: SELECT * FROM attachments WHERE user_id = ? [AND attachable_type = ?] [AND attachable_id = ?] ORDER BY created_at DESC
// 参数 userID: 用户 ID，用于数据隔离。
// 参数 attachableType: 可选的关联对象类型过滤，如 "Transaction"、"Bill" 等，为空时不过滤。
// 参数 attachableID: 可选的关联对象 ID 过滤，为 0 时不过滤。通常与 attachableType 一起使用，定位某个具体对象的附件。
// 返回: 附件列表，按创建时间从新到旧排列。
func (r *AttachmentRepository) List(userID uint64, attachableType string, attachableID uint64) ([]model.Attachment, error) {
	var attachments []model.Attachment
	// 基于 user_id 过滤，确保数据隔离：用户只能查看自己的附件
	query := r.db.Where("user_id = ?", userID)
	// 按关联对象类型过滤（多态关联中的类型字段）
	if attachableType != "" {
		query = query.Where("attachable_type = ?", attachableType)
	}
	// 按关联对象 ID 过滤（多态关联中的 ID 字段）
	if attachableID != 0 {
		query = query.Where("attachable_id = ?", attachableID)
	}
	if err := query.Order("created_at DESC").Find(&attachments).Error; err != nil {
		return nil, err
	}
	return attachments, nil
}

// Delete 根据 ID 和用户 ID 删除附件，同时验证用户权限。
// 执行 SQL: DELETE FROM attachments WHERE id = ? AND user_id = ?
// 参数 id: 附件 ID。
// 参数 userID: 当前登录用户 ID，确保只能删除自己的附件。
// 返回: 删除失败时返回错误。
func (r *AttachmentRepository) Delete(id, userID uint64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Attachment{}).Error
}

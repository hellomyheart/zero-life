// Package repository 提供数据访问层，封装所有与数据库的交互逻辑。
// BackupCodeRepository 负责MFA备用验证码的数据访问。
// 备用验证码是在用户无法使用TOTP验证器时，用于登录验证的一次性代码。
// 每个备用码只能使用一次，使用后标记为已使用（used_at不为空）。
package repository

import (
	"time"

	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// BackupCodeRepository MFA备用验证码仓库，负责备用验证码的数据访问。
// 备用验证码是用户启用MFA后生成的恢复代码，用于在TOTP设备不可用时验证身份。
// 每个备用码只能使用一次，使用后会被标记为已使用。
type BackupCodeRepository struct {
	db *gorm.DB
}

// NewBackupCodeRepository 创建备用验证码仓库实例。
// 参数 db: GORM 数据库连接实例。
// 返回: 初始化后的 BackupCodeRepository 指针。
func NewBackupCodeRepository(db *gorm.DB) *BackupCodeRepository {
	return &BackupCodeRepository{db: db}
}

// Create 批量创建备用验证码。
// 在用户启用MFA时调用，一次性生成多个备用码存入数据库。
// 执行 SQL: INSERT INTO backup_codes (user_id, code) VALUES (?, ?), (?, ?), ...
// 参数 codes: 要创建的备用验证码切片，GORM 会自动填充 ID、CreatedAt 等字段。
// 返回: 创建失败时返回错误。
func (r *BackupCodeRepository) Create(codes []model.BackupCode) error {
	return r.db.Create(&codes).Error
}

// ListByUserID 获取指定用户的所有备用验证码，按 ID 升序排列。
// 包含已使用和未使用的备用码，前端可根据 used_at 字段区分状态。
// 执行 SQL: SELECT * FROM backup_codes WHERE user_id = ? ORDER BY id ASC
// 参数 userID: 用户 ID。
// 返回: 备用验证码列表。
func (r *BackupCodeRepository) ListByUserID(userID uint64) ([]model.BackupCode, error) {
	var codes []model.BackupCode
	if err := r.db.Where("user_id = ?", userID).Order("id ASC").Find(&codes).Error; err != nil {
		return nil, err
	}
	return codes, nil
}

// FindByCode 根据用户ID和验证码查找未使用的备用码。
// 用于MFA验证时校验用户输入的备用码是否有效。
// 条件 used_at IS NULL 确保只能使用未使用过的备用码。
// 执行 SQL: SELECT * FROM backup_codes WHERE user_id = ? AND code = ? AND used_at IS NULL LIMIT 1
// 参数 userID: 用户 ID。
// 参数 code: 用户输入的备用验证码。
// 返回: 找到的备用验证码对象；未找到或已使用时返回 gorm.ErrRecordNotFound 错误。
func (r *BackupCodeRepository) FindByCode(userID uint64, code string) (*model.BackupCode, error) {
	var bc model.BackupCode
	if err := r.db.Where("user_id = ? AND code = ? AND used_at IS NULL", userID, code).First(&bc).Error; err != nil {
		return nil, err
	}
	return &bc, nil
}

// MarkUsed 将指定的备用验证码标记为已使用。
// 记录使用时间（used_at），确保每个备用码只能使用一次。
// 执行 SQL: UPDATE backup_codes SET used_at = ? WHERE id = ?
// 参数 id: 备用验证码 ID。
// 返回: 更新失败时返回错误。
func (r *BackupCodeRepository) MarkUsed(id uint64) error {
	now := time.Now()
	return r.db.Model(&model.BackupCode{}).Where("id = ?", id).Update("used_at", now).Error
}

// DeleteByUserID 删除指定用户的所有备用验证码。
// 在用户重新生成备用码或禁用MFA时调用。
// 执行 SQL: DELETE FROM backup_codes WHERE user_id = ?
// 参数 userID: 用户 ID。
// 返回: 删除失败时返回错误。
func (r *BackupCodeRepository) DeleteByUserID(userID uint64) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.BackupCode{}).Error
}

// Package model 数据模型层，对应数据库结构
// 包含用户认证、MFA、备用码等安全相关模型
package model

import "time"

// BackupCode MFA 备用码模型，对应 backup_codes 表
// 
// 功能说明：
// - 存储两步验证（MFA）的备用恢复码
// - 当用户无法使用认证器 App 时，可用备用码登录
// - 每个备用码只能使用一次，使用后失效
//
// 使用场景：
// - 用户手机丢失，无法使用 Google Authenticator
// - 用户认证器 App 数据丢失
// - 用户设备损坏或更换
//
// 备用码特性：
// - 一次性：每个备用码只能使用一次
// - 随机性：使用加密安全的随机数生成
// - 时效性：建议设置有效期（如 90 天）
// - 数量限制：通常生成 10-20 个备用码
//
// 生成策略：
// 1. 用户启用 MFA 时，生成 10 个备用码
// 2. 备用码格式：8-10 位字母数字组合（如：A7K9-M2P5-Q8R3）
// 3. 使用 bcrypt 或 SHA256 哈希后存储
// 4. 明文展示给用户一次，之后不再显示
//
// 验证流程：
// 1. 用户输入备用码
// 2. 系统哈希后与数据库比对
// 3. 找到匹配且未使用的备用码
// 4. 标记 UsedAt 为当前时间（表示已使用）
// 5. 允许用户登录
//
// 安全建议：
// - 备用码应该离线保存（打印或存储在安全位置）
// - 不要将备用码存储在云端或密码管理器中
// - 使用后立即标记为已使用，防止重放攻击
// - 定期生成新的备用码，替换旧的
//
// 示例：
//   用户启用 MFA -> 生成 10 个备用码
//   备用码 1: A7K9-M2P5-Q8R3 (未使用)
//   备用码 2: B8L0-N3Q6-R9S4 (已使用，UsedAt=2026-04-20)
//   备用码 3: C9M1-P4R7-S0T5 (未使用)
type BackupCode struct {
	// ID 备用码唯一标识，主键自增
	// gorm:"primaryKey" 表示这是主键
	// gorm:"autoIncrement" 表示主键自动递增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	
	// UserID 所属用户 ID，关联 users 表
	// 每个用户可以有多个备用码
	// gorm:"not null" 表示数据库字段不允许为空
	// gorm:"index" 为此字段创建索引，加速按用户查询备用码
	UserID uint64 `gorm:"not null;index" json:"user_id"`
	
	// Code 备用码字符串
	// 通常是随机生成的 8-10 位字母数字组合
	// 格式示例：A7K9-M2P5-Q8R3（带分隔符便于阅读）
	// gorm:"size:20" 限制数据库字段最大长度为 20 个字符
	// gorm:"index" 为此字段创建索引，登录验证时快速查找备用码
	// 
	// 存储方式：
	// - 方案 1：明文存储（不推荐，安全性低）
	// - 方案 2：bcrypt 哈希（推荐，与密码处理相同）
	// - 方案 3：SHA256 哈希（推荐，计算速度快）
	Code string `gorm:"not null;size:20;index" json:"code"`
	
	// UsedAt 备用码使用时间，指针类型表示可为 nil（未使用）
	// nil: 备用码尚未使用，可以使用
	// 非 nil: 备用码已使用，值为使用时间
	// 
	// 验证逻辑：
	//   if UsedAt == nil {
	//     // 备用码可用，允许登录
	//     UsedAt = time.Now() // 标记为已使用
	//   } else {
	//     // 备用码已使用，拒绝登录
	//     return error
	//   }
	UsedAt *time.Time `gorm:"index" json:"used_at"`
	
	// CreatedAt 备用码创建时间
	// 即生成备用码的时间
	// 可用于判断备用码是否过期（如超过 90 天）
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
}

// TableName 指定 BackupCode 模型对应的数据库表名为 backup_codes
func (BackupCode) TableName() string { return "backup_codes" }

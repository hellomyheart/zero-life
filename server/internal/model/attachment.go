// Package model 定义了系统的数据模型，对应数据库表结构
// 包含所有业务实体：账户、交易、分类、标签、预算、账单等
package model

import (
	"time"

	"gorm.io/gorm"
)

// Attachment 附件模型，对应 attachments 表
// 
// 功能说明：
// - 支持为各类业务对象附加文件（图片、PDF、文档等）
// - 采用多态关联，可关联到交易、账户、账单等
// - 记录文件的元数据（文件名、MIME 类型、大小、路径）
//
// 使用场景：
// - 交易凭证：发票、收据、小票照片
// - 账户证明：银行卡照片、开户证明
// - 账单附件：电子账单 PDF、合同扫描件
// - 报销凭证：报销单、审批单
//
// 多态关联设计：
// - AttachableType: 关联对象类型（如："transaction", "account", "bill"）
// - AttachableID: 关联对象的 ID
// - 通过这两个字段可以关联到任何表
//
// 支持的文件类型：
// - 图片：image/jpeg, image/png, image/gif
// - 文档：application/pdf, application/msword
// - 表格：application/vnd.ms-excel, application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
//
// 存储策略：
// - 文件存储在服务器文件系统或对象存储（如 AWS S3）
// - Path 字段存储相对路径或 URL
// - 建议按日期分目录存储：/attachments/2026/04/25/xxx.jpg
//
// 示例：
//   附件：发票照片
//   AttachableType: "transaction"
//   AttachableID: 123
//   Filename: "invoice_20260425.jpg"
//   Mime: "image/jpeg"
//   Size: 524288 (512KB)
//   Path: "/uploads/attachments/2026/04/25/abc123.jpg"
type Attachment struct {
	// ID 附件唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	
	// UserID 所属用户 ID
	// 每个用户有自己的附件集合，用户间数据隔离
	// gorm:"index" 加速按用户查询
	UserID uint64 `gorm:"not null;index" json:"user_id"`
	
	// AttachableType 关联对象类型（多态）
	// 标识附件关联到哪类业务对象
	// 可选值：
	//   - "transaction": 关联交易
	//   - "account": 关联账户
	//   - "recurring_transaction": 关联循环交易
	//   - "budget": 关联预算
	//   - "category": 关联分类
	//   - "piggy_bank": 关联储蓄罐
	// gorm:"size:50" 限制最大长度为 50 个字符
	// gorm:"index" 加速按类型查询
	AttachableType string `gorm:"not null;size:50;index" json:"attachable_type"`
	
	// AttachableID 关联对象 ID（多态）
	// 与 AttachableType 配合，定位到具体的业务对象
	// 例如：AttachableType="transaction", AttachableID=123
	// gorm:"index" 加速按 ID 查询
	AttachableID uint64 `gorm:"not null;index" json:"attachable_id"`
	
	// Filename 原始文件名
	// 用户上传文件时的原始文件名
	// 如："invoice.jpg", "receipt.pdf"
	// gorm:"size:255" 限制最大长度为 255 个字符
	Filename string `gorm:"not null;size:255" json:"filename"`
	
	// Mime 文件 MIME 类型
	// 用于浏览器正确识别文件类型
	// 如："image/jpeg", "application/pdf"
	// gorm:"size:255" 限制最大长度为 255 个字符
	Mime string `gorm:"not null;size:255" json:"mime"`
	
	// Size 文件大小（字节）
	// 记录文件的字节数
	// 用于限制上传大小和显示文件信息
	Size int64 `gorm:"not null" json:"size"`
	
	// Path 文件存储路径
	// 可以是：
	//   - 本地文件系统路径：/var/www/uploads/xxx.jpg
	//   - Web 访问路径：/uploads/attachments/xxx.jpg
	//   - 对象存储 URL: https://s3.amazonaws.com/bucket/xxx.jpg
	// gorm:"size:500" 限制最大长度为 500 个字符
	Path string `gorm:"not null;size:500" json:"path"`
	
	// CreatedAt 附件上传时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	
	// UpdatedAt 附件最后更新时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
	
	// DeletedAt 软删除时间
	// GORM 软删除字段，记录删除时间而非真正删除
	// gorm:"index" 创建索引，GORM 查询时自动过滤已删除记录
	// json:"-" JSON 序列化时忽略此字段
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定 Attachment 模型对应的数据库表名为 attachments
func (Attachment) TableName() string { return "attachments" }

// Package model 定义了系统的数据模型，对应数据库表结构
// 包含所有业务实体：账户、交易、分类、标签、预算、账单等
package model

import (
	"time"

	"gorm.io/gorm"
)

// Tag 标签模型，对应 tags 表
//
// 功能说明：
// - 用于对交易进行灵活标记和分组
// - 支持自定义颜色，便于视觉区分
// - 一个交易可以有多个标签（多对多关系）
// - 标签可用于报表筛选和规则条件
// - 支持最多5级树形结构（通过 parent_id 形成层级关系）
//
// 使用场景：
// - 标记交易用途：如"出差"、"旅游"、"装修"
// - 标记交易来源：如"网购"、"实体店"
// - 标记特殊状态：如"待报销"、"AA 制"
// - 配合规则引擎自动打标签
//
// 与分类的区别：
// - 分类和标签都支持最多5级树形结构
// - 一个交易只能有一个分类，但可以有多个标签
// - 分类用于预算控制，标签用于灵活标记
//
// 示例：
//   交易：星巴克咖啡 -50 元
//   分类：餐饮
//   标签：["网购", "待报销", "出差"]
type Tag struct {
	// ID 标签唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	
	// UserID 所属用户 ID
	// 每个用户有自己的标签集合，用户间数据隔离
	// gorm:"index" 创建索引加速按用户查询
	UserID uint64 `gorm:"not null;index" json:"user_id"`
	
	// Name 标签名称
	// 如："网购"、"出差"、"待报销"
	// gorm:"size:100" 限制最大长度为 100 个字符
	Name string `gorm:"not null;size:100" json:"name"`
	
	// Color 标签颜色，十六进制格式
	// 用于前端展示时的视觉区分
	// 默认值：#409EFF（Element Plus 主色）
	// gorm:"size:7" 限制长度为 7 个字符（#RRGGBB）
	Color string `gorm:"size:7;default:#409EFF" json:"color"`

	// ParentID 父标签ID，支持最多5级树形结构
	// 为nil时表示顶级标签
	ParentID *uint64 `gorm:"index" json:"parent_id"`
	
	// CreatedAt 标签创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	
	// UpdatedAt 标签最后更新时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
	
	// DeletedAt 软删除时间
	// GORM 软删除字段，记录删除时间而非真正删除
	// gorm:"index" 创建索引，GORM 查询时自动过滤已删除记录
	// json:"-" JSON 序列化时忽略此字段
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定 Tag 模型对应的数据库表名为 tags
func (Tag) TableName() string { return "tags" }

// TransactionTag 交易 - 标签关联模型，对应 transaction_tags 表
// 
// 功能说明：
// - 实现交易与标签的多对多关联
// - 一个交易可以有多个标签
// - 一个标签可以关联到多个交易
//
// 数据结构：
// - 复合主键：(TransactionID, TagID)
// - 无额外字段，纯关联表
//
// 示例：
//   交易 ID=1 关联标签 ID=[1, 3, 5]
//   交易 ID=2 关联标签 ID=[2, 4]
//   标签 ID=1 关联交易 ID=[1, 5, 8]
type TransactionTag struct {
	// TransactionID 交易 ID，复合主键的一部分
	// 关联到 transactions 表的 ID 字段
	TransactionID uint64 `gorm:"primaryKey" json:"transaction_id"`
	
	// TagID 标签 ID，复合主键的另一部分
	// 关联到 tags 表的 ID 字段
	TagID uint64 `gorm:"primaryKey" json:"tag_id"`
}

// TableName 指定 TransactionTag 模型对应的数据库表名为 transaction_tags
func (TransactionTag) TableName() string { return "transaction_tags" }

// Package model 定义了系统的数据模型，对应数据库表结构
// 包含所有业务实体：账户、交易、分类、标签、预算、账单等
package model

import (
	"time"

	"gorm.io/gorm"
)

// Category 分类模型，对应 categories 表
//
// 功能说明：
// - 支持两级分类结构（父分类 + 子分类），用于对交易进行分类管理
// - 每笔交易只能属于一个分类（与标签的多对多不同）
// - 分类用于预算控制和报表统计
//
// 两级分类结构：
// - 一级分类（父分类）：ParentID 为 nil
//   如："餐饮"、"交通"、"居住"
// - 二级分类（子分类）：ParentID 指向父分类
//   如："餐饮"下的"外卖"、"堂食"、"咖啡"
//
// 与标签的区别：
// - 分类是树形结构（两级），标签是扁平结构
// - 一个交易只能有一个分类，但可以有多个标签
// - 分类用于预算控制，标签用于灵活标记
//
// 与其他模型的关系：
// - Budget: 多对多关系，一个预算可覆盖多个分类
// - Transaction: 一对多关系，每笔交易有一个分类
// - Bill: 一对多关系，账单可关联分类
//
// 示例：
//   一级分类：餐饮 (ParentID=nil)
//     二级分类：外卖 (ParentID=1)
//     二级分类：堂食 (ParentID=1)
//     二级分类：咖啡 (ParentID=1)
//   一级分类：交通 (ParentID=nil)
//     二级分类：公交 (ParentID=4)
//     二级分类：打车 (ParentID=4)
type Category struct {
	// ID 分类唯一标识，主键自增
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// UserID 所属用户 ID
	// 每个用户有自己独立的分类集合，用户间数据隔离
	// gorm:"index" 创建索引，加速按用户查询分类
	UserID uint64 `gorm:"not null;index" json:"user_id"`

	// Name 分类名称
	// 如："餐饮"、"外卖"、"交通"、"打车"
	// gorm:"size:100" 限制最大长度为 100 个字符
	Name string `gorm:"not null;size:100" json:"name"`

	// ParentID 父分类 ID
	// nil: 表示这是一级分类（顶级分类）
	// 非 nil: 表示这是二级分类，ParentID 指向其父分类
	// gorm:"index" 创建索引，加速按父分类查询子分类
	ParentID *uint64 `gorm:"index" json:"parent_id"`

	// Icon 图标标识
	// 用于前端展示分类图标
	// 可以是图标名称（如 "utensils"）或 emoji
	// gorm:"size:50" 限制最大长度为 50 个字符
	Icon string `gorm:"size:50" json:"icon"`

	// Notes 备注信息
	// 记录分类的附加说明
	// gorm:"type:text" 支持较长文本
	Notes string `gorm:"type:text" json:"notes"`

	// SortOrder 排序序号
	// 数值越小越靠前，用于控制分类在列表中的显示顺序
	// 默认值：0
	SortOrder int `gorm:"default:0" json:"sort_order"`

	// CreatedAt 分类创建时间
	CreatedAt time.Time `gorm:"not null" json:"created_at"`

	// UpdatedAt 分类最后更新时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	// DeletedAt 软删除时间
	// GORM 软删除字段，记录删除时间而非真正删除记录
	// gorm:"index" 创建索引，GORM 查询时自动过滤已删除记录
	// json:"-" JSON 序列化时忽略此字段
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Children 子分类列表
	// 通过 ParentID 外键关联，获取当前分类下的所有子分类
	// 仅一级分类会有子分类，二级分类的 Children 为空
	// json:"children,omitempty" 当字段为零值时 JSON 序列化时忽略
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// TableName 指定 Category 模型对应的数据库表名为 categories
func (Category) TableName() string { return "categories" }

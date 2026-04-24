package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// ConfigurationRepository 系统配置仓库，负责全局系统配置的数据访问。
// 系统配置（Configuration）以键值对形式存储系统级别的配置项，
// 如系统名称、是否允许注册、维护模式等。与 Preference（用户偏好）不同，
// Configuration 是全局的，不区分用户。
type ConfigurationRepository struct {
	db *gorm.DB
}

// NewConfigurationRepository 创建系统配置仓库实例。
// 参数 db: GORM 数据库连接实例。
func NewConfigurationRepository(db *gorm.DB) *ConfigurationRepository {
	return &ConfigurationRepository{db: db}
}

// Get 获取指定名称的系统配置项。
// 执行 SQL: SELECT * FROM configurations WHERE name = ? LIMIT 1
// 参数 name: 配置项名称，如 "allow_registration"、"maintenance_mode" 等。
// 返回: 找到的配置对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *ConfigurationRepository) Get(name string) (*model.Configuration, error) {
	var cfg model.Configuration
	if err := r.db.Where("name = ?", name).First(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Set 设置系统配置项，实现"存在则更新，不存在则创建"的逻辑（Upsert）。
// 先查询该配置名是否存在，如果不存在则创建新记录，否则更新已有记录的值。
// 执行 SQL:
//   查询: SELECT * FROM configurations WHERE name = ? LIMIT 1
//   创建: INSERT INTO configurations (name, value) VALUES (?, ?)
//   更新: UPDATE configurations SET value = ? WHERE id = ?
// 参数 cfg: 配置对象，必须包含 Name 和 Value 字段。
// 返回: 操作失败时返回错误。
func (r *ConfigurationRepository) Set(cfg *model.Configuration) error {
	var existing model.Configuration
	err := r.db.Where("name = ?", cfg.Name).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.Create(cfg).Error
	}
	if err != nil {
		return err
	}
	existing.Value = cfg.Value
	return r.db.Save(&existing).Error
}

// List 获取所有系统配置项，按名称升序排列。
// 执行 SQL: SELECT * FROM configurations ORDER BY name ASC
// 返回: 配置项列表。
func (r *ConfigurationRepository) List() ([]model.Configuration, error) {
	var configs []model.Configuration
	if err := r.db.Order("name ASC").Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

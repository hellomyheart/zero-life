package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// PiggyBankRepository 存钱罐仓库，负责存钱罐及其事件的数据访问。
// 存钱罐（PiggyBank）是用户设置的储蓄目标，如"旅行基金"、"应急储备"等。
// 每次存入或取出都会创建一条事件记录（PiggyEvent），追踪储蓄进度。
// 存钱罐关联一个账户（Account），资金变动通过该账户进行。
type PiggyBankRepository struct {
	db *gorm.DB
}

// NewPiggyBankRepository 创建存钱罐仓库实例。
// 参数 db: GORM 数据库连接实例。
func NewPiggyBankRepository(db *gorm.DB) *PiggyBankRepository {
	return &PiggyBankRepository{db: db}
}

// Create 创建一条新的存钱罐记录。
// 执行 SQL: INSERT INTO piggy_banks (...)
// 参数 piggyBank: 要创建的存钱罐对象。
// 返回: 创建失败时返回错误。
func (r *PiggyBankRepository) Create(piggyBank *model.PiggyBank) error {
	return r.db.Create(piggyBank).Error
}

// GetByID 根据 ID 和用户 ID 获取单条存钱罐，并预加载关联账户和事件。
// 事件按创建时间倒序排列（最新的事件排在前面）。
// 执行 SQL:
//   主查询: SELECT * FROM piggy_banks WHERE id = ? AND user_id = ? LIMIT 1
//   预加载: SELECT * FROM accounts WHERE id IN (...)  (Account)
//           SELECT * FROM piggy_events WHERE piggy_bank_id IN (...) ORDER BY created_at DESC  (Events)
// 参数 id: 存钱罐 ID。
// 参数 userID: 当前登录用户 ID，用于权限校验。
// 返回: 包含关联账户和事件的存钱罐对象。
func (r *PiggyBankRepository) GetByID(id, userID uint64) (*model.PiggyBank, error) {
	var piggyBank model.PiggyBank
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).
		Preload("Account").
		Preload("Events", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		First(&piggyBank).Error; err != nil {
		return nil, err
	}
	return &piggyBank, nil
}

// List 获取指定用户的所有存钱罐，按名称升序排列，并预加载关联账户。
// 执行 SQL: SELECT * FROM piggy_banks WHERE user_id = ? ORDER BY name ASC
// 参数 userID: 用户 ID。
// 返回: 存钱罐列表（含预加载的账户信息）。
func (r *PiggyBankRepository) List(userID uint64) ([]model.PiggyBank, error) {
	var piggyBanks []model.PiggyBank
	if err := r.db.Where("user_id = ?", userID).
		Preload("Account").
		Order("name ASC").Find(&piggyBanks).Error; err != nil {
		return nil, err
	}
	return piggyBanks, nil
}

// Update 更新存钱罐。GORM 的 Save 方法会更新所有字段。
// 执行 SQL: UPDATE piggy_banks SET ... WHERE id = ?
// 参数 piggyBank: 要更新的存钱罐对象（必须包含 ID 字段）。
// 返回: 更新失败时返回错误。
func (r *PiggyBankRepository) Update(piggyBank *model.PiggyBank) error {
	return r.db.Save(piggyBank).Error
}

// Delete 删除存钱罐及其关联的所有事件记录。使用数据库事务确保原子性。
// 执行 SQL（事务内）:
//   1. DELETE FROM piggy_events WHERE piggy_bank_id = ?
//   2. DELETE FROM piggy_banks WHERE id = ? AND user_id = ?
// 参数 id: 存钱罐 ID。
// 参数 userID: 当前登录用户 ID。
// 返回: 删除失败时返回错误，事务回滚。
func (r *PiggyBankRepository) Delete(id, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("piggy_bank_id = ?", id).Delete(&model.PiggyEvent{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&model.PiggyBank{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// CreateEvent 创建一条存钱罐事件记录，记录存入或取出操作。
// 执行 SQL: INSERT INTO piggy_events (...)
// 参数 event: 要创建的事件对象。
// 返回: 创建失败时返回错误。
func (r *PiggyBankRepository) CreateEvent(event *model.PiggyEvent) error {
	return r.db.Create(event).Error
}

// ListEvents 获取指定存钱罐的所有事件记录，按创建时间倒序排列。
// 执行 SQL: SELECT * FROM piggy_events WHERE piggy_bank_id = ? ORDER BY created_at DESC
// 参数 piggyBankID: 存钱罐 ID。
// 返回: 事件记录列表。
func (r *PiggyBankRepository) ListEvents(piggyBankID uint64) ([]model.PiggyEvent, error) {
	var events []model.PiggyEvent
	if err := r.db.Where("piggy_bank_id = ?", piggyBankID).
		Order("created_at DESC").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// Reorder 批量更新存钱罐的排序。使用数据库事务确保原子性。
// 遍历排序映射，逐条更新每个存钱罐的 order 字段。
// 执行 SQL（事务内，循环执行）: UPDATE piggy_banks SET `order` = ? WHERE id = ? AND user_id = ?
// 参数 userID: 用户 ID，确保只能修改自己的存钱罐排序。
// 参数 orders: 排序映射，key 为存钱罐 ID，value 为新的排序值。
// 返回: 更新失败时返回错误，事务回滚。
func (r *PiggyBankRepository) Reorder(userID uint64, orders map[uint64]int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for id, order := range orders {
			if err := tx.Model(&model.PiggyBank{}).Where("id = ? AND user_id = ?", id, userID).Update("order", order).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteEvents 删除指定存钱罐的所有事件记录，用于重置存钱罐时清空历史记录。
// 执行 SQL: DELETE FROM piggy_events WHERE piggy_bank_id = ?
// 参数 piggyBankID: 存钱罐 ID。
// 返回: 删除失败时返回错误。
func (r *PiggyBankRepository) DeleteEvents(piggyBankID uint64) error {
	return r.db.Where("piggy_bank_id = ?", piggyBankID).Delete(&model.PiggyEvent{}).Error
}

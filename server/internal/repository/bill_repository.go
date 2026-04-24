package repository

import (
	"time"

	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

// BillRepository 账单仓库，负责定期账单的数据访问。
// 账单（Bill）表示周期性需要支付的账单，如房租、水电费、订阅费等。
// 每条账单包含下次到期日（next_due），系统会根据此字段自动提醒用户。
type BillRepository struct {
	db *gorm.DB
}

// NewBillRepository 创建账单仓库实例。
// 参数 db: GORM 数据库连接实例。
func NewBillRepository(db *gorm.DB) *BillRepository {
	return &BillRepository{db: db}
}

// Create 创建一条新的账单记录。
// 执行 SQL: INSERT INTO bills (...)
// 参数 bill: 要创建的账单对象。
// 返回: 创建失败时返回错误。
func (r *BillRepository) Create(bill *model.Bill) error {
	return r.db.Create(bill).Error
}

// GetByID 根据 ID 和用户 ID 获取单条账单。
// 同时验证该账单属于指定用户（权限校验）。
// 执行 SQL: SELECT * FROM bills WHERE id = ? AND user_id = ? LIMIT 1
// 参数 id: 账单 ID。
// 参数 userID: 当前登录用户 ID，用于权限校验。
// 返回: 找到的账单对象；未找到时返回 gorm.ErrRecordNotFound 错误。
func (r *BillRepository) GetByID(id, userID uint64) (*model.Bill, error) {
	var bill model.Bill
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&bill).Error; err != nil {
		return nil, err
	}
	return &bill, nil
}

// List 获取指定用户的所有账单，按到期日升序排列（最快到期的排在前面）。
// 执行 SQL: SELECT * FROM bills WHERE user_id = ? ORDER BY next_due ASC
// 参数 userID: 用户 ID。
// 返回: 账单列表。
func (r *BillRepository) List(userID uint64) ([]model.Bill, error) {
	var bills []model.Bill
	if err := r.db.Where("user_id = ?", userID).Order("next_due ASC").Find(&bills).Error; err != nil {
		return nil, err
	}
	return bills, nil
}

// Update 更新账单记录。GORM 的 Save 方法会更新所有字段。
// 执行 SQL: UPDATE bills SET ... WHERE id = ?
// 参数 bill: 要更新的账单对象（必须包含 ID 字段）。
// 返回: 更新失败时返回错误。
func (r *BillRepository) Update(bill *model.Bill) error {
	return r.db.Save(bill).Error
}

// Delete 根据 ID 和用户 ID 删除账单，同时验证用户权限。
// 执行 SQL: DELETE FROM bills WHERE id = ? AND user_id = ?
// 参数 id: 账单 ID。
// 参数 userID: 当前登录用户 ID，确保只能删除自己的账单。
// 返回: 删除失败时返回错误。
func (r *BillRepository) Delete(id, userID uint64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Bill{}).Error
}

// GetAllDueBills 获取所有用户的到期账单（不按用户过滤），用于定时任务批量处理。
// 此方法不限制用户，因为定时任务需要处理所有用户的到期账单。
// 执行 SQL: SELECT * FROM bills WHERE next_due <= NOW() ORDER BY next_due ASC
// 返回: 所有已到期的账单列表。
func (r *BillRepository) GetAllDueBills() ([]model.Bill, error) {
	var bills []model.Bill
	now := time.Now()
	if err := r.db.Where("next_due <= ?", now).
		Order("next_due ASC").Find(&bills).Error; err != nil {
		return nil, err
	}
	return bills, nil
}

// GetUpcoming 获取指定用户未来若干天内到期的账单，用于提醒功能。
// 执行 SQL: SELECT * FROM bills WHERE user_id = ? AND next_due BETWEEN NOW() AND (NOW() + days天) ORDER BY next_due ASC
// 参数 userID: 用户 ID。
// 参数 days: 查询未来多少天内的账单。
// 返回: 即将到期的账单列表。
func (r *BillRepository) GetUpcoming(userID uint64, days int) ([]model.Bill, error) {
	var bills []model.Bill
	now := time.Now()
	endDate := now.AddDate(0, 0, days)
	if err := r.db.Where("user_id = ? AND next_due BETWEEN ? AND ?", userID, now, endDate).
		Order("next_due ASC").Find(&bills).Error; err != nil {
		return nil, err
	}
	return bills, nil
}

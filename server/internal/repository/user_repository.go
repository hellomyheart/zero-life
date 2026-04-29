package repository

import (
	"github.com/hellomyheart/zero-life/server/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewUserRepository(readDB, writeDB *gorm.DB) *UserRepository {
	return &UserRepository{readDB: readDB, writeDB: writeDB}
}

func (r *UserRepository) List(search string, offset, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.readDB.Model(&model.User{})
	if search != "" {
		query = query.Where("email LIKE ? OR nickname LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("id ASC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepository) GetByID(id uint64) (*model.User, error) {
	var user model.User
	if err := r.readDB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	return r.writeDB.Save(user).Error
}

func (r *UserRepository) Delete(id uint64) error {
	return r.writeDB.Delete(&model.User{}, id).Error
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.readDB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *model.User) error {
	return r.writeDB.Create(user).Error
}

func (r *UserRepository) Count() (int64, error) {
	var count int64
	err := r.readDB.Model(&model.User{}).Count(&count).Error
	return count, err
}

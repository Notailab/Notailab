package dao

import (
	"gorm.io/gorm"

	"github.com/Notailab/Notailab/internal/model"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (dao *UserDAO) CreateUser(user *model.User) error {
	return dao.db.Create(user).Error
}

func (dao *UserDAO) CheckUsername(username string) bool {
	var count int64
	err := dao.db.Model(&model.User{}).
		Where("username = ?", username).
		Count(&count).Error
	if err != nil || count != 0 {
		return false
	}
	return true
}

func (dao *UserDAO) CheckEmail(email string) bool {
	var count int64
	err := dao.db.Model(&model.User{}).
		Where("email = ?", email).
		Count(&count).Error
	if err != nil || count != 0 {
		return false
	}
	return true
}

func (dao *UserDAO) GetUserByID(user_id uint) (*model.User, error) {
	var user model.User
	err := dao.db.Where("user_id = ?", user_id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAO) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := dao.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAO) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := dao.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAO) UpdateUser(user *model.User) error {
	return dao.db.Save(user).Error
}

func (dao *UserDAO) DeleteUser(user *model.User) error {
	return dao.db.Delete(user).Error
}

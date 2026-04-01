package dao

import (
	"gorm.io/gorm"

	"github.com/Notailab/Notailab/internal/model"
)

type UserSettingDAO struct {
	db *gorm.DB
}

func NewUserSettingDAO(db *gorm.DB) *UserSettingDAO {
	return &UserSettingDAO{db: db}
}

func (dao *UserSettingDAO) GetByUserID(userID uint) (*model.UserSetting, error) {
	var setting model.UserSetting
	if err := dao.db.Where("user_id = ?", userID).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

func (dao *UserSettingDAO) Save(setting *model.UserSetting) error {
	return dao.db.Save(setting).Error
}

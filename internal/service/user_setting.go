package service

import (
	"errors"
	"strings"

	"github.com/Notailab/Notailab/internal/dao"
	"github.com/Notailab/Notailab/internal/model"
	"gorm.io/gorm"
)

type UserSettingService struct {
	dao *dao.UserSettingDAO
}

func NewUserSettingService(dao *dao.UserSettingDAO) *UserSettingService {
	return &UserSettingService{dao: dao}
}

func (s *UserSettingService) GetByUserID(userID uint) (*model.UserSetting, error) {
	setting, err := s.dao.GetByUserID(userID)
	if err == nil {
		return setting, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &model.UserSetting{UserID: userID, Temperature: 0.7}, nil
	}
	return nil, err
}

func (s *UserSettingService) Update(userID uint, setting *model.UserSetting) error {
	setting.UserID = userID
	setting.LLMBaseURL = strings.TrimSpace(setting.LLMBaseURL)
	setting.LLMAPIKey = strings.TrimSpace(setting.LLMAPIKey)
	setting.LLMModel = strings.TrimSpace(setting.LLMModel)
	if setting.Temperature == 0 {
		setting.Temperature = 0.7
	}
	return s.dao.Save(setting)
}

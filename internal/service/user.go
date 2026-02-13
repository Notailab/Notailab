package service

import (
	"fmt"
	"time"

	"github.com/Notailab/Notailab/internal/dao"
	"github.com/Notailab/Notailab/internal/model"
	"github.com/Notailab/Notailab/pkg/avatar"
	"github.com/Notailab/Notailab/pkg/encrypt"
)

type UserService struct {
	dao *dao.UserDAO
}

func NewUserService(dao *dao.UserDAO) *UserService {
	return &UserService{dao: dao}
}

func (s *UserService) UserLogin(username, password string) (bool, error) {
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, fmt.Errorf("User not found")
	}
	isValid, err := encrypt.VerifyPassword(user.Password, password)
	if err != nil {
		return false, err
	}
	return isValid, nil
}

func (s *UserService) CreateUser(username, password, email string) error {
	user, err := s.GetUserByUsername(username)
	if err == nil && user != nil {
		return fmt.Errorf("User already exists")
	}
	hashPassword, err := encrypt.HashPassword(password)
	if err != nil {
		return fmt.Errorf("Failed to hash password: %v", err)
	}
	avatar := avatar.GenerateTextAvatar(username)

	user = &model.User{
		Username:   username,
		Password:   hashPassword,
		Email:      email,
		Avatar:     avatar,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	return s.dao.CreateUser(user)
}

func (s *UserService) CheckUsername(username string) bool {
	return s.dao.CheckUsername(username)
}

func (s *UserService) CheckEmail(email string) bool {
	return s.dao.CheckEmail(email)
}

func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	return s.dao.GetUserByUsername(username)
}

func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	return s.dao.GetUserByUsername(email)
}

func (s *UserService) UpdateUser(user *model.User) error {
	return s.dao.UpdateUser(user)
}

func (s *UserService) DeleteUser(user *model.User) error {
	return s.dao.DeleteUser(user)
}


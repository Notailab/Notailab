package service

import (
	"fmt"

	"github.com/Notailab/Notailab/internal/dao"
	"github.com/Notailab/Notailab/internal/model"
)

type UserService struct {
	dao *dao.UserDAO
}

func NewUserService(dao *dao.UserDAO) *UserService {
	return &UserService{dao: dao}
}

func (s *UserService) CreateUser(username, password string) error {
	user, err := s.GetUserByUsername(username)
	if err == nil && user != nil {
		return fmt.Errorf("User already exists")
	}
	user = &model.User{
		Username: username,
		Password: password,
	}
	return s.dao.CreateUser(user)
}

func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	return s.dao.GetUserByUsername(username)
}

func (s *UserService) UpdateUser(user *model.User) error {
	return s.dao.UpdateUser(user)
}

func (s *UserService) DeleteUser(user *model.User) error {
	return s.dao.DeleteUser(user)
}
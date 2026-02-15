package service

import (
	"github.com/Notailab/Notailab/internal/dao"
	"github.com/Notailab/Notailab/internal/model"
)

type ProjectService struct {
	dao *dao.ProjectDAO
}

func NewProjectService(dao *dao.ProjectDAO) *ProjectService {
	return &ProjectService{dao: dao}
}

func (s *ProjectService) CreateProject(project *model.Project) error {
	return s.dao.CreateProject(project)
}

func (s *ProjectService) CheckTitleByUserId(user_id uint, title string) bool {
	return s.dao.CheckTitleByUserId(user_id, title)
}

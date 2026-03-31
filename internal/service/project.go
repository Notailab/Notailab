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

func (s *ProjectService) GetAllProjectTitlesByUserId(user_id uint) ([]string, error) {
	projects, err := s.dao.GetAllProjectByUserID(user_id)
	if err != nil {
		return nil, err
	}

	titles := make([]string, 0, len(projects))
	for _, p := range projects {
		titles = append(titles, p.Title)
	}

	return titles, nil
}

func (s *ProjectService) GetAllProjectByUserId(user_id uint) ([]model.Project, error) {
	projects, err := s.dao.GetAllProjectByUserID(user_id)
	if err != nil {
		return nil, err
	}

	return projects, nil
}
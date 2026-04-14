package dao

import (
	"errors"

	"gorm.io/gorm"

	"github.com/Notailab/Notailab/internal/model"
)

type ProjectDAO struct {
	db *gorm.DB
}

func NewProjectDAO(db *gorm.DB) *ProjectDAO {
	return &ProjectDAO{db: db}
}

func (dao *ProjectDAO) CreateProject(pro *model.Project) error {
	return dao.db.Create(pro).Error
}

func (dao *ProjectDAO) UpdateProjectByIDAndUserID(project *model.Project) error {
	return dao.db.Model(&model.Project{}).
		Where("project_id = ? AND user_id = ?", project.ProjectID, project.UserID).
		Updates(map[string]interface{}{
			"title":       project.Title,
			"description": project.Description,
			"start_date":  project.StartDate,
			"end_date":    project.EndDate,
			"updated_at":  gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}

func (dao *ProjectDAO) CheckTitleByUserId(user_id uint, title string) bool {
	var count int64
	err := dao.db.Model(&model.Project{}).
		Where("user_id = ? AND title = ?", user_id, title).
		Count(&count).Error
	if err != nil || count != 0 {
		return false
	}
	return true
}

func (dao *ProjectDAO) GetAllProjectByUserID(user_id uint) ([]model.Project, error) {
	var projects []model.Project

	if err := dao.db.Where("user_id = ?", user_id).Find(&projects).Error; err != nil {
		return nil, err
	}

	if len(projects) == 0 {
		return nil, errors.New("no projects found for the specified user")
	}

	return projects, nil
}

func (dao *ProjectDAO) GetProjectByIDAndUserID(projectID, userID uint) (*model.Project, error) {
	var project model.Project
	if err := dao.db.Where("project_id = ? AND user_id = ?", projectID, userID).First(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (dao *ProjectDAO) GetProjectByTitleAndUserID(title string, userID uint) (*model.Project, error) {
	var project model.Project
	if err := dao.db.Where("title = ? AND user_id = ?", title, userID).First(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

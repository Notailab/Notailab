package dao

import (
	"strings"

	"gorm.io/gorm"

	"github.com/Notailab/Notailab/internal/model"
)

type FileDAO struct {
	db *gorm.DB
}

func NewFileDAO(db *gorm.DB) *FileDAO {
	return &FileDAO{db: db}
}

func (dao *FileDAO) CreateFile(file *model.File) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		file.IsHidden = file.IsHidden || strings.HasPrefix(file.Name, ".")
		if err := tx.Create(file).Error; err != nil {
			return err
		}

		return tx.Model(&model.Project{}).
			Where("project_id = ?", file.ProjectID).
			Update("updated_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
	})
}

func (dao *FileDAO) GetFilesByProjectID(projectID uint) ([]model.File, error) {
	return dao.getFilesByProjectID(projectID, false)
}

func (dao *FileDAO) GetAllFilesByProjectID(projectID uint) ([]model.File, error) {
	return dao.getFilesByProjectID(projectID, true)
}

func (dao *FileDAO) getFilesByProjectID(projectID uint, includeHidden bool) ([]model.File, error) {
	var files []model.File

	query := dao.db.Where("project_id = ?", projectID)
	if !includeHidden {
		query = query.Where("is_hidden = ?", false)
	}
	if err := query.Find(&files).Error; err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return make([]model.File, 0), nil
	}

	return files, nil
}

func (dao *FileDAO) UpdateFileContent(file_id uint, content string) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		var file model.File
		if err := tx.Select("file_id", "project_id").Where("file_id = ?", file_id).First(&file).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.File{}).Where("file_id = ?", file_id).Update("content", content).Error; err != nil {
			return err
		}

		return tx.Model(&model.Project{}).
			Where("project_id = ?", file.ProjectID).
			Update("updated_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
	})
}

func (dao *FileDAO) DeleteFile(file_id uint) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		var file model.File
		if err := tx.Select("file_id", "project_id").Where("file_id = ?", file_id).First(&file).Error; err != nil {
			return err
		}

		if err := tx.Delete(&model.File{}, file_id).Error; err != nil {
			return err
		}

		return tx.Model(&model.Project{}).
			Where("project_id = ?", file.ProjectID).
			Update("updated_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
	})
}

func (dao *FileDAO) GetFileByIDAndProjectID(fileID, projectID uint) (*model.File, error) {
	var file model.File
	if err := dao.db.Where("file_id = ? AND project_id = ?", fileID, projectID).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

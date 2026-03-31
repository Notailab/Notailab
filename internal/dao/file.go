package dao

import (
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
	return dao.db.Create(file).Error
}

func (dao *FileDAO) GetFilesByProjectID(project_id uint) ([]model.File, error) {
	var files []model.File

	if err := dao.db.Where("project_id = ?", project_id).Find(&files).Error; err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return make([]model.File, 0), nil
	}

	return files, nil
}

func (dao *FileDAO) UpdateFileContent(file_id uint, content string) error {
	return dao.db.Model(&model.File{}).Where("file_id = ?", file_id).Update("content", content).Error
}

func (dao *FileDAO) DeleteFile(file_id uint) error {
	return dao.db.Delete(&model.File{}, file_id).Error
}
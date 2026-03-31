package service

import (
	"github.com/Notailab/Notailab/internal/dao"
	"github.com/Notailab/Notailab/internal/model"
)

type FileService struct {
	dao *dao.FileDAO
}

func NewFileService(dao *dao.FileDAO) *FileService {
	return &FileService{dao: dao}
}

func (s *FileService) CreateFile(file *model.File) error {
	return s.dao.CreateFile(file)
}

func (s *FileService) GetFilesByProjectID(project_id uint) ([]model.File, error) {
	files, err := s.dao.GetFilesByProjectID(project_id)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (s *FileService) UpdateFileContent(file_id uint, content string) error {
	return s.dao.UpdateFileContent(file_id, content)
}

func (s *FileService) DeleteFile(file_id uint) error {
	return s.dao.DeleteFile(file_id)
}
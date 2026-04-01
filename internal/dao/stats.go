package dao

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/Notailab/Notailab/internal/model"
)

type StatsDAO struct {
	db *gorm.DB
}

func NewStatsDAO(db *gorm.DB) *StatsDAO {
	return &StatsDAO{db: db}
}

func (dao *StatsDAO) CountProjectsByUserID(userID uint) (int64, error) {
	var count int64
	err := dao.db.Model(&model.Project{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (dao *StatsDAO) CountFilesByUserID(userID uint) (int64, error) {
	var count int64
	err := dao.db.Table("files").
		Select("count(files.file_id)").
		Joins("inner join projects on projects.project_id = files.project_id").
		Where("projects.user_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (dao *StatsDAO) CountAgentConversationsByUserID(userID uint) (int64, error) {
	var count int64
	err := dao.db.Model(&model.AgentConversation{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (dao *StatsDAO) CountRecentProjectsByUserID(userID uint, since time.Time) (int64, error) {
	var count int64
	err := dao.db.Model(&model.Project{}).
		Where("user_id = ? AND updated_at >= ?", userID, since).
		Count(&count).Error
	return count, err
}

func (dao *StatsDAO) CountRecentFilesByUserID(userID uint, since time.Time) (int64, error) {
	var count int64
	err := dao.db.Table("files").
		Select("count(files.file_id)").
		Joins("inner join projects on projects.project_id = files.project_id").
		Where("projects.user_id = ? AND files.updated_at >= ?", userID, since).
		Count(&count).Error
	return count, err
}

func (dao *StatsDAO) ListProjectsByUserID(userID uint) ([]model.Project, error) {
	var projects []model.Project
	err := dao.db.Where("user_id = ?", userID).Find(&projects).Error
	return projects, err
}

func (dao *StatsDAO) ListFilesByUserID(userID uint) ([]model.File, error) {
	var files []model.File
	err := dao.db.Table("files").
		Select("files.*").
		Joins("inner join projects on projects.project_id = files.project_id").
		Where("projects.user_id = ?", userID).
		Find(&files).Error
	return files, err
}

type MonthCount struct {
	Month string `gorm:"column:month"`
	Count int    `gorm:"column:count"`
}

func (dao *StatsDAO) MonthlyFileActivityByUserID(userID uint, months int) ([]MonthCount, error) {
	var records []MonthCount
	if months < 1 {
		months = 1
	}
	sql := fmt.Sprintf(`
		select to_char(date_trunc('month', files.updated_at), 'Mon') as month,
		       count(files.file_id)::int as count
		from files
		inner join projects on projects.project_id = files.project_id
		where projects.user_id = ?
		  and files.updated_at >= date_trunc('month', current_date) - interval '%d months'
		group by date_trunc('month', files.updated_at)
		order by date_trunc('month', files.updated_at)
	`, months-1)
	err := dao.db.Raw(sql, userID).Scan(&records).Error
	return records, err
}

type ProjectVelocityRow struct {
	ProjectID uint      `gorm:"column:project_id" json:"project_id"`
	Title     string    `gorm:"column:title" json:"title"`
	Progress  int       `gorm:"column:progress" json:"progress"`
	FileCount int       `gorm:"column:file_count" json:"file_count"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (dao *StatsDAO) TopProjectsByFileCount(userID uint, limit int) ([]ProjectVelocityRow, error) {
	var rows []ProjectVelocityRow
	err := dao.db.Raw(`
		select p.project_id,
		       p.title,
		       p.progress,
		       p.updated_at,
		       count(f.file_id)::int as file_count
		from projects p
		left join files f on f.project_id = p.project_id
		where p.user_id = ?
		group by p.project_id, p.title, p.progress, p.updated_at
		order by count(f.file_id) desc, p.updated_at desc
		limit ?
	`, userID, limit).Scan(&rows).Error
	return rows, err
}

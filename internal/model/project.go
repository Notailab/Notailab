package model

import "time"

type Project struct {
	ProjectID   uint      `gorm:"primaryKey;                 column:project_id"  json:"project_id"`
	UserID      uint      `gorm:"not null;                   column:user_id"     json:"user_id"`
	Title       string    `gorm:"type:varchar(255);not null; column:title"       json:"title"`
	Description string    `gorm:"type:text;                  column:description" json:"description"`
	Status      string    `gorm:"type:varchar(50);           column:status"      json:"status"`
	StartDate   time.Time `gorm:"default:current_timestamp;  column:start_date"  json:"start_date"`
	EndDate     time.Time `gorm:"default:null;               column:end_date"    json:"end_date"`
	Progress    int       `gorm:"default:0;                  column:progress"    json:"progress"`
	CreatedAt   time.Time `gorm:"default:current_timestamp;  column:created_at"  json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:current_timestamp;  column:updated_at"  json:"updated_at"`
}

func (p *Project) TableName() string {
	return "projects"
}

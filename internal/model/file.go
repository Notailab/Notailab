package model

import "time"

// CREATE TABLE files (
//     id BIGSERIAL PRIMARY KEY,
//     project_id BIGINT NOT NULL,
//     name VARCHAR(30) NOT NULL,
//     content TEXT NOT NULL,
//     is_hidden BOOLEAN NOT NULL DEFAULT FALSE,
//     suffix VARCHAR(10) NOT NULL,
//     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
//     updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
//     version INT NOT NULL DEFAULT 1,
//     FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
//     UNIQUE(project_id, name)
// );

type File struct {
	FileID    uint      `gorm:"primaryKey;                 column:file_id"    json:"file_id"`
	ProjectID uint      `gorm:"not null;                   column:project_id" json:"project_id"`
	Name      string    `gorm:"type:varchar(30);not null;  column:name"       json:"name"`
	Content   string    `gorm:"type:text;                  column:content"    json:"content"`
	IsHidden  bool      `gorm:"not null;default:false;     column:is_hidden"   json:"is_hidden"`
	Version   int       `gorm:"default:1;                  column:version"    json:"version"`
	CreatedAt time.Time `gorm:"default:current_timestamp;  column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp;  column:updated_at" json:"updated_at"`
}

func (f *File) TableName() string {
	return "files"
}

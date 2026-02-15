package model

import "time"

type User struct {
	UserID     uint      `gorm:"primaryKey;                column:user_id"     json:"user_id"`
	Username   string    `gorm:"unique;not null;size=50;   column:username"    json:"username"`
	Password   string    `gorm:"not null;size=100;         column:password"    json:"password"`
	Email      string    `gorm:"unique;not null;sizae=100; column:email"       json:"email"`
	Avatar     string    `gorm:"not null;                  column:avatar"      json:"avatar"`
	CreateTime time.Time `gorm:"not null;                  column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"not null;                  column:update_time" json:"update_time"`
}

func (u *User) TableName() string {
	return "users"
}

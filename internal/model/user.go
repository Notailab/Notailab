package model

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey;column:id" json:"id"`
	Username   string    `gorm:"column:username;unique;not null;size=50" json:"username"`
	Password   string    `gorm:"column:password;not null;size=100" json:"password"`
	Email      string    `gorm:"column:email;unique;not null;sizae=100" json:"email"`
	Avatar     string    `gorm:"column:avatar;not null;" json:"avatar"`
	CreateTime time.Time `gorm:"column:create_time;not null" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;not null" json:"update_time"`
}

func (u *User) TableName() string {
	return "users"
}

package model

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
}

func (u *User) TableName() string {
	return "users"
}
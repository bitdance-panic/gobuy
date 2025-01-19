package model

import "time"

type Base struct {
	ID        int `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// 贫血模型
type User struct {
	Base
	Email          string `gorm:"unique"`
	PasswordHashed string
	Username       string
}

// 设置这个po对应的表名
func (u User) TableName() string {
	return "user"
}

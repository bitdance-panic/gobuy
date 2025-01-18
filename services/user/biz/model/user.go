package model

import (
	"common/model"
	"context"

	"gorm.io/gorm"
)

type User struct {
	model.Base
	Email          string `gorm:"unique"`
	PasswordHashed string
}

// 设置这个po对应的表名
func (u User) TableName() string {
	return "user"
}

// func GetByID(db *gorm.DB, ctx context.Context, userID string) (*User, error) {
// 	err = db.WithContext(ctx).Model(&User{}).Where(&User{ID: userID}).First(&user).Error
// 	return
// }

func GetByEmailAndPass(db *gorm.DB, ctx context.Context, email string, password string) (*User, error) {
	userPO := &User{}
	err := db.WithContext(ctx).Model(&User{}).Where(&User{Email: email, PasswordHashed: password}).First(userPO).Error
	return userPO, err
}

func Create(db *gorm.DB, ctx context.Context, user *User) error {
	return db.WithContext(ctx).Create(user).Error
}

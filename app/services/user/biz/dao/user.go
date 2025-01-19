package dao

import (
	"context"

	"github.com/bitdance-panic/gobuy/app/models"

	"gorm.io/gorm"
)

type User = models.User

// func GetByID(db *gorm.DB, ctx context.Context, userID string) (*User, error) {
// 	err = db.WithContext(ctx).Model(&User{}).Where(&User{ID: userID}).First(&user).Error
// 	return
// }

func GetUserByEmailAndPass(db *gorm.DB, ctx context.Context, email string, password string) (*User, error) {
	userPO := &User{}
	err := db.WithContext(ctx).Model(&User{}).Where(&User{Email: email, PasswordHashed: password}).First(userPO).Error
	return userPO, err
}

func CreateUser(db *gorm.DB, ctx context.Context, user *User) error {
	return db.WithContext(ctx).Create(user).Error
}

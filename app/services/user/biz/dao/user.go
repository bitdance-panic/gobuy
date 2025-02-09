package dao

import (
	"context"

	"github.com/bitdance-panic/gobuy/app/models"

	"gorm.io/gorm"
)

type User = models.User
type Base = models.Base

// func GetByID(db *gorm.DB, ctx context.Context, userID string) (*User, error) {
// 	err = db.WithContext(ctx).Model(&User{}).Where(&User{ID: userID}).First(&user).Error
// 	return
// }

func GetUserByIDAndPass(db *gorm.DB, ctx context.Context, userID int, password string) (*User, error) {
	userPO := &User{}
	err := db.WithContext(ctx).
		Model(&User{}).
		Where(&User{
			Base: Base{
				ID: userID,
			},
			PasswordHashed: password,
		}).
		First(userPO).
		Error

	return userPO, err
}

func GetUserByEmailAndPass(db *gorm.DB, ctx context.Context, email string, password string) (*User, error) {
	userPO := &User{}
	err := db.WithContext(ctx).Model(&User{}).Where(&User{Email: email, PasswordHashed: password}).First(userPO).Error
	// newUser := &User{
	// 	Email:          email,
	// 	Username:       "qwer",
	// 	PasswordHashed: password,
	// }
	// db.WithContext(ctx).Create(newUser)
	return userPO, err
}

func CreateUser(db *gorm.DB, ctx context.Context, user *User) error {
	return db.WithContext(ctx).Create(user).Error
}

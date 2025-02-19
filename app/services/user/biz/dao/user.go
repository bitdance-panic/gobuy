package dao

import (
	"context"

	"github.com/bitdance-panic/gobuy/app/models"

	"gorm.io/gorm"
)

type User = models.User

//type Base = models.Base

func RegisterUser(db *gorm.DB, ctx context.Context, username, password, email string) (*User, error) {
	hashedPassword := password

	user := &User{
		Username:       username,
		PasswordHashed: hashedPassword,
		Email:          email,
	}

	// 插入新用户
	err := db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUsers 查询所有用户信息，并进行分页
func GetUsers(db *gorm.DB, ctx context.Context, page int, pageSize int) ([]User, error) {
	var users []User
	offset := (page - 1) * pageSize
	err := db.WithContext(ctx).Limit(pageSize).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

//func GetUserByIDAndPass(db *gorm.DB, ctx context.Context, userID int, password string) (*User, error) {
//	userPO := &User{}
//	err := db.WithContext(ctx).
//		Model(&User{}).
//		Where(&User{
//			Base: Base{
//				ID: userID,
//			},
//			PasswordHashed: password,
//		}).
//		First(userPO).
//		Error
//
//	return userPO, err
//}

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

// GetUserByID 根据用户 ID 查询用户信息
func GetUserByID(db *gorm.DB, ctx context.Context, userID int) (*User, error) {
	user := &User{}

	err := db.WithContext(ctx).
		Where("id = ?", userID).
		First(user).
		Error

	return user, err
}

// 更新
func UpdateUserByID(db *gorm.DB, ctx context.Context, userID int, username, email string) error {
	// if db.DB == nil {
	// 	return errors.New("database connection is nil")
	// }

	return db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"username": username,
			"email":    email,
		}).Error
}

// 封禁
func DeleteUserByID(db *gorm.DB, ctx context.Context, userID int) error {
	return db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", userID).
		Update("is_deleted", 1).Error
}

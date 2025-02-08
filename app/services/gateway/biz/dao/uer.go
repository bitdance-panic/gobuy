package dao

import (
	"errors"

	"github.com/bitdance-panic/gobuy/app/models"
	"gorm.io/gorm"
)

type User = models.User

// UpdateRefreshToken 更新刷新令牌
func UpdateRefreshToken(db *gorm.DB, userID int, refreshToken string) error {
	return db.Model(&User{}).Where("id = ?", userID).
		Update("refresh_token", refreshToken).
		Error
}

// 根据用户ID查询 Refresh Token
func GetRefreshTokenByUserID(db *gorm.DB, userID int) (string, error) {
	var user User
	result := db.Model(&User{}).Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return "", errors.New("refresh token not found")
	}
	return user.RefreshToken, nil
}

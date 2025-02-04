package dao

import (
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal/mysql"
)

// GetUserByUsername 通过用户名获取用户
func GetUserByUsername(username string) (*User, error) {
	var user User
	err := mysql.DB().Where("Username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail 通过用户名获取用户
func GetUserByEmail(email string) (*User, error) {
	var user User
	err := mysql.DB().Where("Email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 通过用户ID获取用户
func GetUserByID(userID int) (*User, error) {
	var user User
	err := mysql.DB().Where("Id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建一个新的用户
func CreateUser(user *User) error {
	err := mysql.DB().Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateRefreshToken 更新刷新令牌
func UpdateRefreshToken(userID int, refreshToken string) error {
	err := mysql.DB().Model(&User{}).Where("Id = ?", userID).Update("Refresh_token", refreshToken).Error
	if err != nil {
		return err
	}
	return nil
}

// ListProducts 获取商品列表
func ListProducts(page, pageSize int) ([]Product, int, error) {
	var products []Product
	var total int64

	// 查询商品总数
	err := mysql.DB().Model(&Product{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询商品列表
	err = mysql.DB().Offset(page * pageSize).Limit(pageSize).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, int(total), nil
}

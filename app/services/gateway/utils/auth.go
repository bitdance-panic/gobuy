package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 哈希密码
func HashPassword(password string) (string, error) {
	// 使用 bcrypt 算法对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword 验证密码
func VerifyPassword(providedPassword, storedHashedPassword string) bool {
	// 使用 bcrypt 校验密码是否与存储的哈希密码匹配
	err := bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(providedPassword))
	return err == nil
}

// GenerateTokens 生成双令牌
func GenerateTokens(userID int, secret string) (accessToken, refreshToken string, err error) {
	// Access Token (12小时)
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(12 * time.Hour).Unix(),
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(secret))
	if err != nil {
		return "", "", errors.New("生成 Access Token 时出错: " + err.Error())
	}

	// Refresh Token (7天)
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secret))
	if err != nil {
		return "", "", errors.New("生成 Refresh Token 时出错: " + err.Error())
	}

	return
}

// VerifyToken 验证JWT令牌
func VerifyToken(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保算法为HS256
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("不支持的签名方法")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// 如果验证通过，返回声明（claims）
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("无效的令牌")
}

// RefreshToken 刷新Token（可选：用于获取新的访问令牌）
func RefreshToken(refreshToken string, secret string) (string, error) {
	// 验证刷新令牌
	claims, err := VerifyToken(refreshToken, secret)
	if err != nil {
		return "", err
	}

	// 获取用户ID
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return "", errors.New("无法获取用户ID")
	}

	// 生成新的访问令牌
	newAccessToken, _, err := GenerateTokens(int(userID), secret)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}

package dao

import (
	"github.com/bitdance-panic/gobuy/app/models"
	"gorm.io/gorm"
)

type User = models.User
type Product = models.Product

func HasPermission(db *gorm.DB, userID uint, permissionName string) bool {
	var count int64
	db.Model(&User{}).
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Joins("JOIN role_permissions ON role_permissions.role_id = roles.id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("users.id = ? AND permissions.name = ?", userID, permissionName).
		Count(&count)
	return count > 0
}

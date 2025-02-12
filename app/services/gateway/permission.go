package main

import (
	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/gateway/casbin"
	"gorm.io/gorm"
)

type PermissionService struct {
	db *gorm.DB
}

func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{db: db}
}

// UpdateUserRoles 更新用户角色
func (s *PermissionService) UpdateUserRoles(userID uint, roleIDs []uint) error {
	// 删除旧角色
	if err := s.db.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
		return err
	}

	// 添加新角色
	var userRoles []models.UserRole
	for _, rid := range roleIDs {
		userRoles = append(userRoles, models.UserRole{
			UserID: userID,
			RoleID: rid,
		})
	}

	if err := s.db.Create(&userRoles).Error; err != nil {
		return err
	}

	// 更新策略
	casbin.Enforcer.LoadPolicy()
	return nil
}

// UpdateRolePermissions 更新角色权限
func (s *PermissionService) UpdateRolePermissions(roleID uint, permissions []uint) error {
	// 类似用户角色更新逻辑
	// ...

	// 更新策略
	casbin.Enforcer.LoadPolicy()
	return nil
}

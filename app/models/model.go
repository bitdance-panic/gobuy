package models

import (
	"log"
	"os"
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        int `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted bool
}

// 贫血模型
type User struct {
	Base
	Email          string `gorm:"unique"`
	PasswordHashed string
	Username       string
	RefreshToken   string

	Roles  []Role `gorm:"many2many:user_role;"` // 用户和角色的多对多关系
	Orders []Order
}

type Role struct {
	Base
	Name        string       `gorm:"unique;not null"`
	Permissions []Permission `gorm:"many2many:role_permission;"` // 角色和权限的多对多关系
}

type Permission struct {
	Base
	Name   string `gorm:"unique;not null"`            // 权限名称，例如 "create_user", "delete_user"
	Path   string `gorm:"not null"`                   // 资源路径，支持通配符如 /product/*
	Method string `gorm:"not null"`                   // 请求方法，支持正则如 GET|POST
	Roles  []Role `gorm:"many2many:role_permission;"` // 权限和角色的多对多关系
}

type UserRole struct {
	UserID int
	RoleID int
}

type RolePermission struct {
	RoleID       int
	PermissionID int
}

type Product struct {
	Base
	Name        string
	Price       float64
	Stock       int
	Image       string
	Description string
}

type CartItem struct {
	Base
	UserID    int
	ProductID int
	Quantity  int     `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID"`
}

type Order struct {
	Base
	UserID     int
	Number     string  `gorm:"unique;not null"`
	TotalPrice float64 `gorm:"not null"`
	// OrderStatus
	Status  int `gorm:"type:varchar(20);not null"`
	Items   []OrderItem
	PayTime *time.Time
}

type OrderItem struct {
	Base
	OrderID      int     // 订单 ID
	ProductID    int     // 商品 ID
	Quantity     int     `gorm:"not null"`             // 商品数量
	Price        float64 `gorm:"not null"`             // 商品单价
	Product      Product `gorm:"foreignKey:ProductID"` // 关联商品
	Order        Order   `gorm:"foreignKey:OrderID"`   // 关联订单
	ProductName  string
	ProductImage string
}

// 黑名单条目模型
type Blacklist struct {
	Base
	Identifier string    `gorm:"type:varchar(255);uniqueIndex;not null"` // 封禁标识（用户ID）
	Reason     string    `gorm:"type:text;not null"`                     // 封禁原因
	ExpiresAt  time.Time `gorm:"index"`                                  // 过期时间（为空表示永久封禁）
}

// 操作类型（用于日志）
type BlacklistOperation string

const (
	AddToBlacklist      BlacklistOperation = "ADD"
	RemoveFromBlacklist BlacklistOperation = "REMOVE"
)

func (User) TableName() string {
	return "user"
}

func (Role) TableName() string {
	return "role"
}

func (Permission) TableName() string {
	return "permission"
}

func (UserRole) TableName() string {
	return "user_role"
}

func (RolePermission) TableName() string {
	return "role_permission"
}

func (CartItem) TableName() string {
	return "cart_item"
}

func (Order) TableName() string {
	return "order"
}

func (OrderItem) TableName() string {
	return "order_item"
}

func (Product) TableName() string {
	return "product"
}

func (Blacklist) TableName() string {
	return "blacklist"
}

// 调用这个来自动调整表结构
func AutoMigrate(db *gorm.DB) {
	if os.Getenv("GO_ENV") != "production" {
		log.Println("进行数据表的Migrate")
		//nolint:errcheck
		db.AutoMigrate(
			&User{},
			&Role{},
			&Permission{},
			&UserRole{},
			&RolePermission{},
			&CartItem{},
			&Order{},
			&OrderItem{},
			&Product{},
			&Blacklist{},
		)
	}
}

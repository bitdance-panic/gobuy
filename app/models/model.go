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

// 定义订单状态常量
type OrderStatus int

const (
	OrderStatusPending   OrderStatus = 1 // 待支付
	OrderStatusPaid      OrderStatus = 2 // 已支付
	OrderStatusCancelled OrderStatus = 3 // 已取消
)

// 贫血模型
type User struct {
	Base
	Email          string `gorm:"unique"`
	PasswordHashed string
	Username       string
	RefreshToken   string

	Roles  []Role `gorm:"many2many:user_role;"` // 用户和角色的多对多关系
	Cart   Cart
	Orders []Order
}

type Role struct {
	Base
	Name        string       `gorm:"unique;not null"`
	Permissions []Permission `gorm:"many2many:role_permission;"` // 角色和权限的多对多关系
}

type Permission struct {
	Base
	Name  string `gorm:"unique;not null"`            // 权限名称，例如 "create_user", "delete_user"
	Roles []Role `gorm:"many2many:role_permission;"` // 权限和角色的多对多关系
}

type UserRole struct {
	UserID uint
	RoleID uint
}

type RolePermission struct {
	RoleID       uint
	PermissionID uint
}

type Product struct {
	Base
	Name        string
	Price       float64
	Stock       int
	Image       string
	Description string
	// 规格
}

type Cart struct {
	Base
	UserID uint       // 用户 ID
	Items  []CartItem // 购物车项
}

type CartItem struct {
	Base
	CartID    uint    // 购物车 ID
	ProductID uint    // 商品 ID
	Quantity  int     `gorm:"not null"`             // 商品数量
	Product   Product `gorm:"foreignKey:ProductID"` // 关联商品
}

type Order struct {
	Base
	UserID      uint    // 用户 ID
	OrderNumber string  `gorm:"unique;not null"` // 订单号
	TotalAmount float64 `gorm:"not null"`        // 订单总金额
	// OrderStatus
	Status int         `gorm:"type:varchar(20);not null"` // 订单状态
	Items  []OrderItem // 订单项
	
}

type OrderItem struct {
	Base
	OrderID     uint    // 订单 ID
	ProductID   uint    // 商品 ID
	Quantity    int     `gorm:"not null"`             // 商品数量
	Price       float64 `gorm:"not null"`             // 商品单价
	Product     Product `gorm:"foreignKey:ProductID"` // 关联商品
	ProductName string
}

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

func (Cart) TableName() string {
	return "cart"
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
			&Cart{},
			&CartItem{},
			&Order{},
			&OrderItem{},
			&Product{},
		)
	}
}

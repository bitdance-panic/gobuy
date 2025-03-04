package dao

import (
	"fmt"
	"github.com/bitdance-panic/gobuy/app/consts"
	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Order = models.Order
type UserAddress = models.UserAddress

func CreateOrder(db *gorm.DB, order *Order) error {
	if err := db.Create(order).Error; err != nil {
		return err
	}
	return nil
}

func SaveOrder(db *gorm.DB, order *Order) error {
	if err := db.Save(order).Error; err != nil {
		return err
	}
	return nil
}

func GetOrderByID(db *gorm.DB, orderID int) (*Order, error) {
	var order Order
	result := db.Preload("Items").First(&order, "id = ?", orderID)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &order, nil
}

func UpdateOrderStatus(db *gorm.DB, order *Order, newStatus consts.OrderStatus) error {
	if order == nil {
		return errors.New("order is nil")
	}
	order.Status = int(newStatus)
	return db.Save(order).Error
}

func UpdateOrderAddress(db *gorm.DB, order *Order, newAddress string) error {
	if order == nil {
		return errors.New("order is nil")
	}
	order.OrderAddress = newAddress
	return db.Save(order).Error
}

func ListUserOrder(db *gorm.DB, userID int, pageNum int, pageSize int) (*[]Order, error) {
	var orders []Order
	err := db.Preload("Items").Limit(pageSize).Offset((pageNum-1)*pageSize).Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return &orders, nil
}

func AdminListOrder(db *gorm.DB, pageNum int, pageSize int) (*[]Order, int64, error) {
	var orders []Order
	err := db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	var count int64
	db.Model(&Order{}).Count(&count)
	return &orders, count, nil
}

func ListPendingOrder(db *gorm.DB) (*[]Order, error) {
	var orders []Order
	err := db.Where("status = ?", consts.OrderStatusPending).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return &orders, nil
}

// 创建订单地址
func CreateUserAddress(db *gorm.DB, userAddress *UserAddress) error {
	if err := db.Create(userAddress).Error; err != nil {
		return err
	}
	return nil
}

func DeleteUserAddress(db *gorm.DB, userID int32) error {
	result := db.Where("user_id = ?", userID).Delete(&UserAddress{})
	if result.RowsAffected == 0 {
		return fmt.Errorf("no records found with user_id %d", userID)
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 更新订单地址
func UpdateUserAddress(db *gorm.DB, userID int32, userAddress string) error {
	var address *UserAddress
	if err := db.Model(&address).Where("user_id = ?", userID).Update("user_address", userAddress).Error; err != nil {
		return err
	}
	return nil
}

// 获取订单地址
func GetUserAddress(db *gorm.DB, userID int32) (*UserAddress, error) {
	var userAddress *UserAddress
	if err := db.Where("user_id = ?", userID).First(&userAddress).Error; err != nil {
		return nil, err
	}
	return userAddress, nil
}

package dao

import (
	"github.com/bitdance-panic/gobuy/app/consts"
	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Order = models.Order
type OrderAddress = models.OrderAddress

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
func CreateOrderAddress(db *gorm.DB, orderAddress *OrderAddress) error {
	if err := db.Create(orderAddress).Error; err != nil {
		return err
	}
	return nil
}

func DeleteOrderAddress(db *gorm.DB, orderID int32) error {
	var orderAddress *OrderAddress
	if err := db.Where("order_id = ?", orderID).Delete(&orderAddress).Error; err != nil {
		return err
	}
	return nil
}

// 更新订单地址
func UpdateOrderAddress(db *gorm.DB, orderID int32, orderAddress string) error {
	var address *OrderAddress
	if err := db.Model(&address).Where("order_id = ?", orderID).Update("order_address", orderAddress).Error; err != nil {
		return err
	}
	return nil
}

// 获取订单地址
func GetOrderAddress(db *gorm.DB, orderID int32) (*OrderAddress, error) {
	var orderAddress *OrderAddress
	if err := db.Where("order_id = ?", orderID).First(&orderAddress).Error; err != nil {
		return nil, err
	}
	return orderAddress, nil
}

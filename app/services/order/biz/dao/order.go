package dao

import (
	"consts"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Order = models.Order

func CreateOrder(db *gorm.DB, order *Order) (*Order, error) {
	if err := db.Create(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func GetOrderByID(db *gorm.DB, orderID int) (*Order, error) {
	var order Order
	result := db.First(&order, "id = ?", orderID)
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

func ListUserOrder(db *gorm.DB, userID int, pageNum int, pageSize int) (*[]models.Order, error) {
	var orders []models.Order
	err := db.Limit(pageSize).Offset((pageNum-1)*pageSize).Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return &orders, nil
}

func AdminListOrder(db *gorm.DB, userID int, pageNum int, pageSize int) (*[]models.Order, error) {
	var orders []models.Order
	err := db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return &orders, nil
}

func ListPendingOrder(db *gorm.DB) (*[]models.Order, error) {
	var orders []models.Order
	err := db.Where("status = ?", consts.OrderStatusPending).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return &orders, nil
}

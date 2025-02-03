package dal

import (
	"github.com/bitdance-panic/gobuy/app/models"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	//tidb.Init()
	//SaveOrder()
	//GetOrderById()
	//UpdateOrderStatus()
}
func SaveOrder(order *models.Order) error {
	result := db.Create(&order)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func GetOrderByUserId(userId int32) ([]models.Order, error) {
	var orders []models.Order
	result := db.Where("user_id = ?", userId).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(orders) == 0 {
		return nil, nil
	}
	return orders, nil
}
func GetOrderById(orderId int32) (*models.Order, error) {
	var order models.Order
	result := db.First(&order, "id = ?", orderId)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &order, nil
}
func UpdateOrderStatus(orderId int32, newStatus int) (*models.Order, error) {
	var order *models.Order
	result := db.Model(&order).Where("id = ?", orderId).Update("status", newStatus)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return order, nil
}

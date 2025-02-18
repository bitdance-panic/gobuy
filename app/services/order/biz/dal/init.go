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

	result := db.Where("user_id = ? AND is_deleted = false", userId).Find(&orders)

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

	result := db.Where("id = ? AND is_deleted = false", orderId).First(&order)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &order, nil
}
func UpdateOrderStatus(orderId int32, newStatus int) (*models.Order, error) {

	var order models.Order

	result := db.Model(&order).Where("id = ?", orderId).Update("status", newStatus)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &order, nil
}

/*
*
安装订单状态获取订单
*/
func GetOrdersByStatus(status models.OrderStatus) ([]models.Order, error) {
	var orders []models.Order
	result := db.Where("status = ?", status).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}
func SoftDeleteOrder(orderId int32) error {
	//更新订单的IsDeleted字段为true
	result := db.Where("order_id = ?", orderId).Update("is_deleted", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil

}

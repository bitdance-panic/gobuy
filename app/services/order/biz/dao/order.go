package dao

import (
	"context"
	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Order = models.Order

func CreateOrder(db *gorm.DB, ctx context.Context, order *models.Order) error {
	return db.WithContext(ctx).Create(order).Error
}

func DeleteOrder(db *gorm.DB, ctx context.Context, order *models.Order) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}
	//删除订单
	result := db.Delete(order)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}
func UpdateOrder(db *gorm.DB, ctx context.Context, order *models.Order) error {
	if order == nil {
		return errors.New("order can't be nil")
	}
	//更新订单状态
	result := db.Model(&order).Updates(map[string]interface{}{
		"status":      order.Status,
		"totalamount": order.TotalAmount,
		"items":       order.Items,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no row updated")
	}
	return nil
}
func GetOrder(db *gorm.DB, ctx context.Context, order *models.Order) error {
	if order == nil {
		return errors.New("order can't be nil")
	}
	//查询订单
	result := db.First(order, "id = ?", order.ID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}
func GetOrderList(db *gorm.DB, ctx context.Context, order *models.Order) error {
	if order == nil {
		return errors.New("order can't be nil")
	}
	//查询订单列表
	var orders []models.Order
	result := db.Where("user_id = ?", order.UserID).Find(&orders)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("order not found")
	}
	*order = orders[0] //假设只需要第一个订单
	return nil
}

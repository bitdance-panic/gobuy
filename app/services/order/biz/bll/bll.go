package bll

import (
	"errors"
	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/order/biz/dal"
)

type OrderBLL struct{}

func NewOrderBLL() *OrderBLL {
	return &OrderBLL{}
}

// 创建一个新订单
func (bll *OrderBLL) CreateOrder(userID int32, items []models.OrderItem) (models.Order, error) {
	order := &models.Order{
		UserID:      uint(userID),
		OrderNumber: "1",
		TotalAmount: 123,
		Status:      1,
		Items:       items,
	}
	err := dal.SaveOrder(order)
	if err != nil {

	}
	return *order, err
}
func (bll *OrderBLL) GetOrder(orderID int32) (*models.Order, error) {
	order, err := dal.GetOrderById(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}
func (bll *OrderBLL) UpdateOrder(orderId int32, newStatus int32) error {
	order, err := dal.UpdateOrderStatus(orderId, int(newStatus))
	if err != nil {
		return err
	}
	return order, nil
}
func (bll *OrderBLL) GetOrdersByUserID(userID int32) ([]models.Order, error) {
	//调用数据访问层来获取订单
	orders, err := dal.GetOrderByUserId(userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

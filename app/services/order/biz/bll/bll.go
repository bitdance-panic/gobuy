package bll

import (
	//"context"
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
	if order == nil || order.IsDeleted {
		return nil, errors.New("order not found")
	}
	return order, nil
}
func (bll *OrderBLL) UpdateOrder(orderId int32, newStatus int32) (*models.Order, error) {
	order, err := dal.UpdateOrderStatus(orderId, int(newStatus))
	if err != nil {
		return nil, err
	}
	return order, nil
}
func (bll *OrderBLL) GetOrdersByUserID(userID int32) ([]models.Order, error) {
	//调用数据访问层来获取订单
	orders, err := dal.GetOrderByUserId(userID)
	if err != nil {
		return nil, err
	}
	//过滤掉已经软删除的订单
	activeOrders := make([]models.Order, 0)
	for _, order := range orders {
		if !order.IsDeleted {
			activeOrders = append(activeOrders, order)
		}
	}
	return activeOrders, nil
}

// 软删除订单
func (bll *OrderBLL) SoftDeleteOrder(orderID int32) error {
	err := dal.SoftDeleteOrder(orderID)
	if err != nil {
		return err
	}
	return nil
}

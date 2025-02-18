package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bitdance-panic/gobuy/app/models"

	rpc_order "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order"
	"github.com/bitdance-panic/gobuy/app/services/order/biz/bll"
)

var orderBll *bll.OrderBLL

func init() {
	orderBll = bll.NewOrderBLL()
}

type OrderServiceImpl struct{}

func (s *OrderServiceImpl) DeleteOrder(ctx context.Context, req *rpc_order.DeleteOrderRequest) (r *rpc_order.DeleteOrderResponse, err error) {
	//TODO implement me
	orderBll.SoftDeleteOrder(req.GetOrderId())
	return &rpc_order.DeleteOrderResponse{
		Success: true,
		Message: "success",
	}, nil
}

func (s *OrderServiceImpl) CreateOrder(ctx context.Context, req *rpc_order.CreateOrderRequest) (*rpc_order.CreateOrderResponse, error) {
	//调用业务逻辑层创建订单
	orderInfo, err := orderBll.CreateOrder(req.UserId, ConvertOrderItems(req.Items))
	if err != nil {
		return nil, err
	}
	// todo 把items弄进来
	rpcOrder := &rpc_order.Order{
		Id:          int32(orderInfo.ID),
		UserId:      int32(orderInfo.UserID),
		OrderNumber: orderInfo.OrderNumber,
		TotalAmount: orderInfo.TotalAmount,
		Status:      rpc_order.OrderStatus(orderInfo.Status),
		Items:       []*rpc_order.OrderItem{},
		IsDeleted:   orderInfo.IsDeleted,
	}
	//返回创建的订单信息
	return &rpc_order.CreateOrderResponse{
		Order: rpcOrder,
	}, nil
}
func (s *OrderServiceImpl) GetUserOrders(ctx context.Context, req *rpc_order.GetUserOrdersRequest) (*rpc_order.GetUserOrdersResponse, error) {
	//调用业务逻辑层获取订单
	orders, err := orderBll.GetOrdersByUserID(req.UserId)
	if err != nil {
		return nil, err
	}
	//输出订单信息
	for _, order := range orders {
		fmt.Printf("Order ID: %d, User ID: %s, Items: %v, Amount: %.2f, Status: %s\n",
			order.ID, order.UserID, order.Items, order.TotalAmount, order.Status)
	}
	return &rpc_order.GetUserOrdersResponse{
		Orders: []*rpc_order.Order{},
	}, nil

	//orderInfo, err := orderBll.GetOrder(r
	//if err != nil {
	//	return nil, err
	//}
	//if orderInfo == nil {
	//	return &rpc_order.CreateOrderResponse{}, nil
	//}
	////返回订单信息
	//return &rpc_order.CreateOrderResponse{
	//	OrderId: orderInfo.ID,
	//	UserId:  orderInfo.UserID,
	//	Items:   orderInfo.Items,
	//	Amount:  orderInfo.Amount,
	//	Status:  orderInfo.Status,
	//}, nil

}
func (s *OrderServiceImpl) UpdateOrder(ctx context.Context, req *rpc_order.UpdateOrderRequest) (*rpc_order.UpdateOrderResponse, error) {
	//调用业务逻辑层更新订单
	order, err := orderBll.UpdateOrder(123, int32(req.Status))
	if err != nil {
		return nil, err
	}
	// todo 把items弄进来
	rpcOrder := &rpc_order.Order{
		Id:          int32(order.ID),
		UserId:      int32(order.UserID),
		OrderNumber: order.OrderNumber,
		TotalAmount: order.TotalAmount,
		Status:      rpc_order.OrderStatus(order.Status),
		Items:       []*rpc_order.OrderItem{},
		IsDeleted:   order.IsDeleted,
	}
	//返回更新结果
	return &rpc_order.UpdateOrderResponse{
		Order: rpcOrder,
	}, nil
}
func (s *OrderServiceImpl) SoftDeleteOrder(ctx context.Context) int {
	orderIdStr := ctx.Value("orderId").(string)
	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		return 0
	}
	err = orderBll.SoftDeleteOrder(int32(orderId))
	if err != nil {
		return 0
	}
	return http.StatusOK
}

func ConvertOrderItems(items []*rpc_order.OrderItem) []models.OrderItem {
	var orderItems []models.OrderItem
	for _, item := range items {
		orderItem := models.OrderItem{
			OrderID:   uint(item.OrderId),
			ProductID: uint(item.ProductId),
			Quantity:  int(item.Quantity),
			Price:     item.Price,
			//Product:     item.Product,
			ProductName: item.ProductName,
		}
		orderItems = append(orderItems, orderItem)
	}
	return orderItems
}

package bll

import (
	"context"
	"errors"

	"github.com/bitdance-panic/gobuy/app/consts"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/cart/biz/clients"
	"github.com/bitdance-panic/gobuy/app/services/cart/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/order/biz/dao"
	"github.com/bitdance-panic/gobuy/app/utils"

	rpc_order "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order"
	rpc_product "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product"
)

type OrderBLL struct{}

func NewOrderBLL() *OrderBLL {
	return &OrderBLL{}
}

// TODO 自动机
// 创建一个新订单
func (bll *OrderBLL) CreateOrder(ctx context.Context, req *rpc_order.CreateOrderReq) (*rpc_order.CreateOrderResp, error) {
	orderNumber, err := utils.GenerateID()
	if err != nil {
		return nil, err
	}
	orderItems := make([]models.OrderItem, len(req.Items))
	totalPrice := 0.0
	// 获取订单项完整信息
	for _, productItem := range req.Items {
		resp, err := clients.ProductClient.GetProductByID(
			ctx,
			&rpc_product.GetProductByIDReq{
				Id: productItem.ProductId,
			},
		)
		if err != nil {
			return nil, err
		}
		p := resp.Product
		if p == nil {
			return nil, errors.New("product not found")
		}
		// 创建订单项
		orderItem := models.OrderItem{
			ProductID:   int(productItem.ProductId),
			ProductName: p.Name,
			Price:       p.Price,
			Quantity:    int(productItem.Quantity),
		}
		// 保存订单项到数据库
		itemWithID, err := dao.CreateOrderItem(tidb.DB, &orderItem)
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, *itemWithID)
		totalPrice += p.Price * float64(productItem.Quantity)
	}

	// 再保存订单
	order := models.Order{
		UserID:     int(req.UserId),
		Number:     orderNumber,
		TotalPrice: totalPrice,
		Status:     int(consts.OrderStatusPending),
		Items:      orderItems,
	}
	orderWithID, err := dao.CreateOrder(tidb.DB, &order)
	if err != nil {
		return nil, err
	}
	protoOrder := convertOrderToProto(orderWithID)
	return &rpc_order.CreateOrderResp{
		Order: protoOrder,
	}, err
}

func (bll *OrderBLL) GetOrder(ctx context.Context, req *rpc_order.GetOrderReq) (*rpc_order.GetOrderResp, error) {
	order, err := dao.GetOrderByID(tidb.DB, int(req.OrderId))
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	protoOrder := convertOrderToProto(order)
	return &rpc_order.GetOrderResp{
		Order: protoOrder,
	}, nil
}

func (bll *OrderBLL) UpdateOrderStatus(ctx context.Context, req *rpc_order.UpdateOrderStatusReq) (*rpc_order.UpdateOrderStatusResp, error) {
	order, err := dao.GetOrderByID(tidb.DB, int(req.OrderId))
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	err = dao.UpdateOrderStatus(tidb.DB, order, consts.OrderStatus(req.Status))
	if err != nil {
		return nil, err
	}
	return &rpc_order.UpdateOrderStatusResp{
		Success: true,
	}, nil
}
func (bll *OrderBLL) ListUserOrder(ctx context.Context, req *rpc_order.ListOrderReq) (*rpc_order.ListOrderResp, error) {
	orders, err := dao.ListUserOrder(tidb.DB, int(req.UserId), int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	protoOrders := make([]*rpc_order.Order, 0, len(*orders))
	for _, order := range *orders {
		protoOrder := convertOrderToProto(&order)
		protoOrders = append(protoOrders, protoOrder)
	}
	return &rpc_order.ListOrderResp{
		Orders: protoOrders,
	}, nil
}

func (bll *OrderBLL) AdminListOrder(ctx context.Context, req *rpc_order.ListOrderReq) (*rpc_order.ListOrderResp, error) {
	orders, err := dao.AdminListOrder(tidb.DB, int(req.UserId), int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	protoOrders := make([]*rpc_order.Order, 0, len(*orders))
	for _, order := range *orders {
		protoOrder := convertOrderToProto(&order)
		protoOrders = append(protoOrders, protoOrder)
	}
	return &rpc_order.ListOrderResp{
		Orders: protoOrders,
	}, nil
}

func convertItemFromProto(items rpc_order.OrderItem) *models.OrderItem {
	return &models.OrderItem{
		OrderID:     int(items.OrderId),
		ProductID:   int(items.ProductId),
		Quantity:    int(items.Quantity),
		Price:       items.Price,
		ProductName: items.ProductName,
	}
}

func convertItemToProto(items models.OrderItem) *rpc_order.OrderItem {
	return &rpc_order.OrderItem{
		OrderId:     int32(items.OrderID),
		ProductId:   int32(items.ProductID),
		Quantity:    int32(items.Quantity),
		Price:       items.Price,
		ProductName: items.ProductName,
	}
}

func convertOrderToProto(order *models.Order) *rpc_order.Order {
	protoItems := make([]*rpc_order.OrderItem, len(order.Items))
	for _, item := range order.Items {
		orderItem := convertItemToProto(item)
		protoItems = append(protoItems, orderItem)
	}
	return &rpc_order.Order{
		Id:         int32(order.ID),
		UserId:     int32(order.UserID),
		Number:     order.Number,
		TotalPrice: order.TotalPrice,
		Status:     int32(order.Status),
		Items:      protoItems,
	}
}

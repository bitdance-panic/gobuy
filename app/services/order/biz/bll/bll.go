package bll

import (
	"context"
	"errors"

	"github.com/bitdance-panic/gobuy/app/consts"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/order/biz/clients"
	"github.com/bitdance-panic/gobuy/app/services/order/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/order/biz/dao"
	"github.com/bitdance-panic/gobuy/app/utils"

	rpc_cart "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart"
	rpc_order "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order"

	"github.com/shopspring/decimal"
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
	orderItems := make([]models.OrderItem, len(req.CartItemIDs))
	totalPrice := 0.0
	for i, cartItemID := range req.CartItemIDs {
		resp, err := clients.CartClient.GetItem(
			ctx,
			&rpc_cart.GetItemReq{ItemId: cartItemID},
		)
		if err != nil {
			return nil, err
		}
		cartItem := resp.Item
		orderItem := models.OrderItem{
			ProductID:    int(cartItem.ProductId),
			ProductName:  cartItem.Name,
			ProductImage: cartItem.Image,
			Price:        cartItem.Price,
			Quantity:     int(cartItem.Quantity),
		}
		// 避免外键报错
		orderItem.Product.ID = int(orderItem.ProductID)
		// itemWithID, err := dao.CreateOrderItem(tidb.DB, &orderItem)
		// if err != nil {
		// 	return nil, err
		// }
		orderItems[i] = orderItem
		p, _ := decimal.NewFromFloat(orderItem.Price).Mul(decimal.NewFromInt(int64(orderItem.Quantity))).Float64()
		// if !exact {
		// 	return nil, errors.New("handle decimal is not exact")
		// }
		totalPrice += p
	}
	// 先保存下订单
	order := models.Order{
		UserID:     int(req.UserId),
		Number:     orderNumber,
		TotalPrice: totalPrice,
		Status:     int(consts.OrderStatusPending),
		PayTime:    nil,
	}

	err = dao.CreateOrder(tidb.DB, &order)
	if err != nil {
		return nil, err
	}
	// 再存订单项
	for i := range len(orderItems) {
		orderItems[i].OrderID = order.ID
	}
	order.Items = orderItems
	// TODO 事务
	err = dao.SaveOrder(tidb.DB, &order)
	if err != nil {
		return nil, err
	}
	protoOrder := convertOrderToProto(&order)
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
	protoOrders := make([]*rpc_order.Order, len(*orders))
	for i, order := range *orders {
		protoOrder := convertOrderToProto(&order)
		protoOrders[i] = protoOrder
	}
	return &rpc_order.ListOrderResp{
		Orders: protoOrders,
	}, nil
}

func (bll *OrderBLL) AdminListOrder(ctx context.Context, req *rpc_order.ListOrderReq) (*rpc_order.ListOrderResp, error) {
	orders, total, err := dao.AdminListOrder(tidb.DB, int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	protoOrders := make([]*rpc_order.Order, 0, len(*orders))
	for _, order := range *orders {
		protoOrder := convertOrderToProto(&order)
		protoOrders = append(protoOrders, protoOrder)
	}
	return &rpc_order.ListOrderResp{
		Orders:     protoOrders,
		TotalCount: total,
	}, nil
}

func convertItemFromProto(items rpc_order.OrderItem) *models.OrderItem {
	return &models.OrderItem{
		OrderID:      int(items.OrderId),
		ProductID:    int(items.ProductId),
		Quantity:     int(items.Quantity),
		Price:        items.Price,
		ProductName:  items.ProductName,
		ProductImage: items.ProductImage,
	}
}

func convertItemToProto(items models.OrderItem) *rpc_order.OrderItem {
	return &rpc_order.OrderItem{
		OrderId:      int32(items.OrderID),
		ProductId:    int32(items.ProductID),
		Quantity:     int32(items.Quantity),
		Price:        items.Price,
		ProductName:  items.ProductName,
		ProductImage: items.ProductImage,
	}
}

func convertOrderToProto(order *models.Order) *rpc_order.Order {
	protoItems := make([]*rpc_order.OrderItem, len(order.Items))
	for i, item := range order.Items {
		orderItem := convertItemToProto(item)
		protoItems[i] = orderItem
	}
	return &rpc_order.Order{
		Id:         int32(order.ID),
		UserId:     int32(order.UserID),
		Number:     order.Number,
		TotalPrice: order.TotalPrice,
		Status:     int32(order.Status),
		Items:      protoItems,
		CreatedAt:  utils.FormatTime(order.CreatedAt),
	}
}

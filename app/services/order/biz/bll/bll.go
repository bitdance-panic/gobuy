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
		UserID:       int(req.UserId),
		Number:       orderNumber,
		TotalPrice:   totalPrice,
		Status:       int(consts.OrderStatusPending),
		PayTime:      nil,
		Phone:        req.Phone,
		OrderAddress: req.OrderAddress,
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
	//TODO 事务
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

func (bll *OrderBLL) UpdateOrderAddress(ctx context.Context, req *rpc_order.UpdateOrderAddressReq) (*rpc_order.UpdateOrderAddressResp, error) {
	order, err := dao.GetOrderByID(tidb.DB, int(req.OrderId))
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	err = dao.UpdateOrderAddress(tidb.DB, order, req.OrderAddress)
	if err != nil {
		return nil, err
	}
	return &rpc_order.UpdateOrderAddressResp{
		Success: true,
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
		Id:           int32(order.ID),
		UserId:       int32(order.UserID),
		Number:       order.Number,
		TotalPrice:   order.TotalPrice,
		Status:       int32(order.Status),
		Items:        protoItems,
		CreatedAt:    utils.FormatTime(order.CreatedAt),
		Phone:        order.Phone,
		OrderAddress: order.OrderAddress,
	}
}

// 创建订单地址
func (bll *OrderBLL) CreateUserAddress(ctx context.Context, req *rpc_order.CreateUserAddressReq) (*rpc_order.CreateUserAddressResp, error) {
	userAddress := models.UserAddress{
		UserID:      int(req.UserId),
		Phone:       req.Phone,
		UserAddress: req.UserAddress,
	}

	// 调用 DAO 层创建订单地址
	err := dao.CreateUserAddress(tidb.DB, &userAddress)
	if err != nil {
		return nil, err
	}

	return &rpc_order.CreateUserAddressResp{
		UserId:  int32(userAddress.UserID),
		Success: true,
	}, nil
}

// 删除订单地址
func (bll *OrderBLL) DeleteOrderAddress(ctx context.Context, req *rpc_order.DeleteUserAddressReq) (*rpc_order.DeleteUserAddressResp, error) {
	// 调用 DAO 层删除订单地址
	err := dao.DeleteUserAddress(tidb.DB, req.UserId)
	if err != nil {
		return nil, err
	}

	return &rpc_order.DeleteUserAddressResp{
		UserId:  req.UserId,
		Success: true,
	}, nil
}

// 更新订单地址
func (bll *OrderBLL) UpdateUserAddress(ctx context.Context, req *rpc_order.UpdateUserAddressReq) (*rpc_order.UpdateUserAddressResp, error) {
	// 调用 DAO 层更新订单地址
	err := dao.UpdateUserAddress(tidb.DB, req.UserId, req.UserAddress)
	if err != nil {
		return nil, err
	}

	return &rpc_order.UpdateUserAddressResp{
		UserAddress: req.UserAddress,
		Success:     true,
	}, nil
}

// 获取订单地址
func (bll *OrderBLL) GetUserAddress(ctx context.Context, req *rpc_order.GetUserAddressReq) (*rpc_order.GetUserAddressResp, error) {
	// 调用 DAO 层获取订单地址
	userAddress, err := dao.GetUserAddress(tidb.DB, req.UserId)
	if err != nil {
		return nil, err
	}

	return &rpc_order.GetUserAddressResp{
		UserAddress: &rpc_order.UserAddress{
			UserId:      int32(userAddress.UserID),
			Phone:       userAddress.Phone,
			UserAddress: userAddress.UserAddress,
		},
	}, nil
}

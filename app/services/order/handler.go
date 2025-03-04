package main

import (
	"context"

	rpc_order "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order"
	"github.com/bitdance-panic/gobuy/app/services/order/biz/bll"
)

var orderBll *bll.OrderBLL

func init() {
	orderBll = bll.NewOrderBLL()
}

type OrderServiceImpl struct{}

func (*OrderServiceImpl) CreateOrder(ctx context.Context, req *rpc_order.CreateOrderReq) (*rpc_order.CreateOrderResp, error) {
	return orderBll.CreateOrder(ctx, req)
}
func (*OrderServiceImpl) ListUserOrder(ctx context.Context, req *rpc_order.ListOrderReq) (*rpc_order.ListOrderResp, error) {
	return orderBll.ListUserOrder(ctx, req)
}

func (*OrderServiceImpl) GetOrder(ctx context.Context, req *rpc_order.GetOrderReq) (*rpc_order.GetOrderResp, error) {
	return orderBll.GetOrder(ctx, req)
}

func (*OrderServiceImpl) UpdateOrderStatus(ctx context.Context, req *rpc_order.UpdateOrderStatusReq) (*rpc_order.UpdateOrderStatusResp, error) {
	return orderBll.UpdateOrderStatus(ctx, req)
}

func (*OrderServiceImpl) UpdateOrderAddress(ctx context.Context, req *rpc_order.UpdateOrderAddressReq) (*rpc_order.UpdateOrderAddressResp, error) {
	return orderBll.UpdateOrderAddress(ctx, req)
}

func (*OrderServiceImpl) AdminListOrder(ctx context.Context, req *rpc_order.ListOrderReq) (*rpc_order.ListOrderResp, error) {
	return orderBll.AdminListOrder(ctx, req)
}

func (*OrderServiceImpl) CreateUserAddress(ctx context.Context, req *rpc_order.CreateUserAddressReq) (*rpc_order.CreateUserAddressResp, error) {
	return orderBll.CreateUserAddress(ctx, req)
}

func (*OrderServiceImpl) DeleteUserAddress(ctx context.Context, req *rpc_order.DeleteUserAddressReq) (*rpc_order.DeleteUserAddressResp, error) {
	return orderBll.DeleteUserAddress(ctx, req)
}
func (*OrderServiceImpl) UpdateUserAddress(ctx context.Context, req *rpc_order.UpdateUserAddressReq) (*rpc_order.UpdateUserAddressResp, error) {
	return orderBll.UpdateUserAddress(ctx, req)
}
func (*OrderServiceImpl) GetUserAddress(ctx context.Context, req *rpc_order.GetUserAddressReq) (*rpc_order.GetUserAddressResp, error) {
	return orderBll.GetUserAddress(ctx, req)
}

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

func (*OrderServiceImpl) AdminListOrder(ctx context.Context, req *rpc_order.ListOrderReq) (*rpc_order.ListOrderResp, error) {
	return orderBll.AdminListOrder(ctx, req)
}

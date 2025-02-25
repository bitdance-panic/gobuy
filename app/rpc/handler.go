package main

import (
	"context"
	order "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order"
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct{}

// CreateOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) CreateOrder(ctx context.Context, req *order.CreateOrderReq) (resp *order.CreateOrderResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateOrderStatus implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) UpdateOrderStatus(ctx context.Context, req *order.UpdateOrderStatusReq) (resp *order.UpdateOrderStatusResp, err error) {
	// TODO: Your code here...
	return
}

// GetOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) GetOrder(ctx context.Context, req *order.GetOrderReq) (resp *order.GetOrderResp, err error) {
	// TODO: Your code here...
	return
}

// ListUserOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) ListUserOrder(ctx context.Context, req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	// TODO: Your code here...
	return
}

// AdminListOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) AdminListOrder(ctx context.Context, req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	// TODO: Your code here...
	return
}

// CreateOrderAddressResp implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) CreateOrderAddressResp(ctx context.Context, req *order.CreateOrderAddressReq) (resp *order.CreateOrderAddressResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteOrderAddressResp implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) DeleteOrderAddressResp(ctx context.Context, req *order.DeleteOrderAddressReq) (resp *order.DeleteOrderAddressResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateOrderAddressResp implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) UpdateOrderAddressResp(ctx context.Context, req *order.UpdateOrderAddressReq) (resp *order.UpdateOrderAddressResp, err error) {
	// TODO: Your code here...
	return
}

// GetOrderAddressResp implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) GetOrderAddressResp(ctx context.Context, req *order.GetOrderAddressReq) (resp *order.GetOrderAddressResp, err error) {
	// TODO: Your code here...
	return
}

// CreateOrderAddress implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) CreateOrderAddress(ctx context.Context, req *order.CreateOrderAddressReq) (resp *order.CreateOrderAddressResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteOrderAddress implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) DeleteOrderAddress(ctx context.Context, req *order.DeleteOrderAddressReq) (resp *order.DeleteOrderAddressResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateOrderAddress implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) UpdateOrderAddress(ctx context.Context, req *order.UpdateOrderAddressReq) (resp *order.UpdateOrderAddressResp, err error) {
	// TODO: Your code here...
	return
}

// GetOrderAddress implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) GetOrderAddress(ctx context.Context, req *order.GetOrderAddressReq) (resp *order.GetOrderAddressResp, err error) {
	// TODO: Your code here...
	return
}

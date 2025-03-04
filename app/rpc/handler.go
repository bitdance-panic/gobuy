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

// UpdateOrderAddressResp implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) UpdateOrderAddressResp(ctx context.Context, req *order.UpdateOrderAddressReq) (resp *order.UpdateOrderAddressResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateOrderAddress implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) UpdateOrderAddress(ctx context.Context, req *order.UpdateOrderAddressReq) (resp *order.UpdateOrderAddressResp, err error) {
	// TODO: Your code here...
	return
}

// CreateUserAddress implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) CreateUserAddress(ctx context.Context, req *order.CreateUserAddressReq) (resp *order.CreateUserAddressResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteUserAddress implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) DeleteUserAddress(ctx context.Context, req *order.DeleteUserAddressReq) (resp *order.DeleteUserAddressResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateUserAddress implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) UpdateUserAddress(ctx context.Context, req *order.UpdateUserAddressReq) (resp *order.UpdateUserAddressResp, err error) {
	// TODO: Your code here...
	return
}

// GetUserAddress implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) GetUserAddress(ctx context.Context, req *order.GetUserAddressReq) (resp *order.GetUserAddressResp, err error) {
	// TODO: Your code here...
	return
}

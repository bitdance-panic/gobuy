// Code generated by Kitex v0.12.1. DO NOT EDIT.

package orderservice

import (
	"context"
	order "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	CreateOrder(ctx context.Context, req *order.CreateOrderReq, callOptions ...callopt.Option) (r *order.CreateOrderResp, err error)
	UpdateOrderStatus(ctx context.Context, req *order.UpdateOrderStatusReq, callOptions ...callopt.Option) (r *order.UpdateOrderStatusResp, err error)
	GetOrder(ctx context.Context, req *order.GetOrderReq, callOptions ...callopt.Option) (r *order.GetOrderResp, err error)
	ListUserOrder(ctx context.Context, req *order.ListOrderReq, callOptions ...callopt.Option) (r *order.ListOrderResp, err error)
	AdminListOrder(ctx context.Context, req *order.ListOrderReq, callOptions ...callopt.Option) (r *order.ListOrderResp, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfoForClient(), options...)
	if err != nil {
		return nil, err
	}
	return &kOrderServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kOrderServiceClient struct {
	*kClient
}

func (p *kOrderServiceClient) CreateOrder(ctx context.Context, req *order.CreateOrderReq, callOptions ...callopt.Option) (r *order.CreateOrderResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.CreateOrder(ctx, req)
}

func (p *kOrderServiceClient) UpdateOrderStatus(ctx context.Context, req *order.UpdateOrderStatusReq, callOptions ...callopt.Option) (r *order.UpdateOrderStatusResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.UpdateOrderStatus(ctx, req)
}

func (p *kOrderServiceClient) GetOrder(ctx context.Context, req *order.GetOrderReq, callOptions ...callopt.Option) (r *order.GetOrderResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetOrder(ctx, req)
}

func (p *kOrderServiceClient) ListUserOrder(ctx context.Context, req *order.ListOrderReq, callOptions ...callopt.Option) (r *order.ListOrderResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ListUserOrder(ctx, req)
}

func (p *kOrderServiceClient) AdminListOrder(ctx context.Context, req *order.ListOrderReq, callOptions ...callopt.Option) (r *order.ListOrderResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.AdminListOrder(ctx, req)
}

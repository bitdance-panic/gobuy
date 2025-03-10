// Code generated by Kitex v0.12.2. DO NOT EDIT.

package productservice

import (
	"context"
	product "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	CreateProduct(ctx context.Context, req *product.CreateProductReq, callOptions ...callopt.Option) (r *product.CreateProductResp, err error)
	UpdateProduct(ctx context.Context, req *product.UpdateProductReq, callOptions ...callopt.Option) (r *product.UpdateProductResp, err error)
	RemoveProduct(ctx context.Context, req *product.RemoveProductReq, callOptions ...callopt.Option) (r *product.RemoveProductResp, err error)
	GetProductByID(ctx context.Context, req *product.GetProductByIDReq, callOptions ...callopt.Option) (r *product.GetProductByIDResp, err error)
	ListProduct(ctx context.Context, req *product.ListProductReq, callOptions ...callopt.Option) (r *product.ListProductResp, err error)
	AdminListProduct(ctx context.Context, req *product.ListProductReq, callOptions ...callopt.Option) (r *product.ListProductResp, err error)
	SearchProducts(ctx context.Context, req *product.SearchProductsReq, callOptions ...callopt.Option) (r *product.SearchProductsResp, err error)
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
	return &kProductServiceClient{
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

type kProductServiceClient struct {
	*kClient
}

func (p *kProductServiceClient) CreateProduct(ctx context.Context, req *product.CreateProductReq, callOptions ...callopt.Option) (r *product.CreateProductResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.CreateProduct(ctx, req)
}

func (p *kProductServiceClient) UpdateProduct(ctx context.Context, req *product.UpdateProductReq, callOptions ...callopt.Option) (r *product.UpdateProductResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.UpdateProduct(ctx, req)
}

func (p *kProductServiceClient) RemoveProduct(ctx context.Context, req *product.RemoveProductReq, callOptions ...callopt.Option) (r *product.RemoveProductResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.RemoveProduct(ctx, req)
}

func (p *kProductServiceClient) GetProductByID(ctx context.Context, req *product.GetProductByIDReq, callOptions ...callopt.Option) (r *product.GetProductByIDResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetProductByID(ctx, req)
}

func (p *kProductServiceClient) ListProduct(ctx context.Context, req *product.ListProductReq, callOptions ...callopt.Option) (r *product.ListProductResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ListProduct(ctx, req)
}

func (p *kProductServiceClient) AdminListProduct(ctx context.Context, req *product.ListProductReq, callOptions ...callopt.Option) (r *product.ListProductResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.AdminListProduct(ctx, req)
}

func (p *kProductServiceClient) SearchProducts(ctx context.Context, req *product.SearchProductsReq, callOptions ...callopt.Option) (r *product.SearchProductsResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.SearchProducts(ctx, req)
}

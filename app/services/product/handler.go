package main

import (
	"context"

	rpc_product "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product"
	"github.com/bitdance-panic/gobuy/app/services/product/biz/bll"
)

type ProductServiceImpl struct{}

func (*ProductServiceImpl) ListProduct(ctx context.Context, req *rpc_product.ListProductReq) (*rpc_product.ListProductResp, error) {
	return bll.ListProduct(ctx, req)
}

func (*ProductServiceImpl) GetProductByID(ctx context.Context, req *rpc_product.GetProductByIDReq) (*rpc_product.GetProductByIDResp, error) {
	return bll.GetProductByID(ctx, req)
}

func (*ProductServiceImpl) CreateProduct(ctx context.Context, req *rpc_product.CreateProductReq) (*rpc_product.CreateProductResp, error) {
	return bll.CreateProduct(ctx, req)
}

func (*ProductServiceImpl) UpdateProduct(ctx context.Context, req *rpc_product.UpdateProductReq) (*rpc_product.UpdateProductResp, error) {
	return bll.UpdateProduct(ctx, req)
}

func (*ProductServiceImpl) RemoveProduct(ctx context.Context, req *rpc_product.RemoveProductReq) (*rpc_product.RemoveProductResp, error) {
	return bll.RemoveProduct(ctx, req)
}

func (*ProductServiceImpl) AdminListProduct(ctx context.Context, req *rpc_product.ListProductReq) (resp *rpc_product.ListProductResp, err error) {
	return bll.AdminListProduct(ctx, req)
}

func (*ProductServiceImpl) SearchProducts(ctx context.Context, req *rpc_product.SearchProductsReq) (*rpc_product.SearchProductsResp, error) {
	return bll.SearchProducts(ctx, req)
}

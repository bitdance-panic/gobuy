package main

import (
	"context"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product"
	"github.com/bitdance-panic/gobuy/app/services/product/biz/bll"
)

var ub *bll.ProductBLL

func init() {
	ub = bll.NewProductBLL()
}

type ProductServiceImpl struct {
}

func (s *ProductServiceImpl) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	return ub.CreateProduct(ctx, req)
}

func (s *ProductServiceImpl) UpdateProduct(ctx context.Context, req *product.UpdateProductRequest) (*product.UpdateProductResponse, error) {
	return ub.UpdateProduct(ctx, req)
}

func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	return ub.DeleteProduct(ctx, req)
}

func (s *ProductServiceImpl) GetProductByID(ctx context.Context, req *product.GetProductByIDRequest) (*product.GetProductByIDResponse, error) {
	return ub.GetProductByID(ctx, req)
}

func (s *ProductServiceImpl) SearchProducts(ctx context.Context, req *product.SearchProductsRequest) (*product.SearchProductsResponse, error) {
	return ub.SearchProducts(ctx, req)
}

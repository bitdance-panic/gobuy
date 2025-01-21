package bll

import (
	"context"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product"
	"github.com/bitdance-panic/gobuy/app/services/product/biz/dao"
)

type ProductBLL struct{}

func NewProductBLL() *ProductBLL {
	return &ProductBLL{}
}

func (s *ProductBLL) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	p := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Stock),
		Image:       req.Img,
	}
	if err := dao.Create(&p); err != nil {
		return nil, err
	}
	protoProduct := convertToProtoProduct(&p)
	return &product.CreateProductResponse{
		Product: protoProduct,
	}, nil
}

func (s *ProductBLL) UpdateProduct(ctx context.Context, req *product.UpdateProductRequest) (*product.UpdateProductResponse, error) {
	p, err := dao.GetByID(uint(req.Id))
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.Description != "" {
		p.Description = req.Description
	}
	if req.Price != 0 {
		p.Price = req.Price
	}
	if req.Stock != 0 {
		p.Stock = int(req.Stock)
	}
	if req.Img != "" {
		p.Image = req.Img
	}
	if err := dao.Update(p); err != nil {
		return nil, err
	}
	protoProduct := convertToProtoProduct(p)
	return &product.UpdateProductResponse{
		Product: protoProduct,
	}, nil
}

func (s *ProductBLL) DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	err := dao.Delete(uint(req.Id))
	if err != nil {
		return &product.DeleteProductResponse{Success: false}, err
	}
	return &product.DeleteProductResponse{Success: true}, nil
}

func (s *ProductBLL) GetProductByID(ctx context.Context, req *product.GetProductByIDRequest) (*product.GetProductByIDResponse, error) {
	p, err := dao.GetByID(uint(req.Id))
	if err != nil {
		return nil, err
	}
	protoProduct := convertToProtoProduct(p)
	return &product.GetProductByIDResponse{
		Product: protoProduct,
	}, nil
}

func (s *ProductBLL) SearchProducts(ctx context.Context, req *product.SearchProductsRequest) (*product.SearchProductsResponse, error) {
	products, err := dao.Search(req.Query)
	if err != nil {
		return nil, err
	}
	protoProducts := make([]*product.Product, 0, len(products))
	for _, p := range products {
		protoProducts = append(protoProducts, convertToProtoProduct(&p))
	}
	return &product.SearchProductsResponse{
		Products: protoProducts,
	}, nil
}

func convertToProtoProduct(p *models.Product) *product.Product {
	return &product.Product{
		Id:          int32(p.ID),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       int32(p.Stock),
		Img:         p.Image,
		CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

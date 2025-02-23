package bll

import (
	"context"

	"github.com/bitdance-panic/gobuy/app/models"
	rpc_product "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product"
	"github.com/bitdance-panic/gobuy/app/services/product/biz/dao"
	"github.com/bitdance-panic/gobuy/app/services/user/biz/dal/tidb"
)

type ProductBLL struct{}

func NewProductBLL() *ProductBLL {
	return &ProductBLL{}
}

func (*ProductBLL) ListProduct(ctx context.Context, req *rpc_product.ListProductReq) (*rpc_product.ListProductResp, error) {
	p, err := dao.List(tidb.DB, int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	protoProducts := make([]*rpc_product.Product, len(*p))
	for i, v := range *p {
		protoProducts[i] = convertToProtoProduct(&v)
	}
	return &rpc_product.ListProductResp{
		Products: protoProducts,
	}, nil
}

func (*ProductBLL) AdminListProduct(ctx context.Context, req *rpc_product.ListProductReq) (*rpc_product.ListProductResp, error) {
	p, err := dao.AdminList(tidb.DB, int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	protoProducts := make([]*rpc_product.Product, len(*p))
	for i, v := range *p {
		protoProducts[i] = convertToProtoProduct(&v)
	}
	return &rpc_product.ListProductResp{
		Products: protoProducts,
	}, nil
}

func (*ProductBLL) GetProductByID(ctx context.Context, req *rpc_product.GetProductByIDReq) (*rpc_product.GetProductByIDResp, error) {
	p, err := dao.GetByID(tidb.DB, int(req.Id))
	if err != nil {
		return nil, err
	}
	if p == nil {
		return &rpc_product.GetProductByIDResp{
			Product: nil,
		}, nil
	}
	protoProduct := convertToProtoProduct(p)
	return &rpc_product.GetProductByIDResp{
		Product: protoProduct,
	}, nil
}

func (*ProductBLL) CreateProduct(ctx context.Context, req *rpc_product.CreateProductReq) (*rpc_product.CreateProductResp, error) {
	p := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Stock),
		Image:       req.Img,
	}
	if err := dao.Create(tidb.DB, &p); err != nil {
		return nil, err
	}
	protoProduct := convertToProtoProduct(&p)
	return &rpc_product.CreateProductResp{
		Product: protoProduct,
	}, nil
}

func (*ProductBLL) UpdateProduct(ctx context.Context, req *rpc_product.UpdateProductReq) (*rpc_product.UpdateProductResp, error) {
	p, err := dao.GetByID(tidb.DB, int(req.Id))
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
	if err := dao.Update(tidb.DB, p); err != nil {
		return nil, err
	}
	protoProduct := convertToProtoProduct(p)
	return &rpc_product.UpdateProductResp{
		Product: protoProduct,
	}, nil
}

func (s *ProductBLL) RemoveProduct(ctx context.Context, req *rpc_product.RemoveProductReq) (*rpc_product.RemoveProductResp, error) {
	err := dao.Remove(tidb.DB, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &rpc_product.RemoveProductResp{Success: true}, nil
}

func (s *ProductBLL) SearchProducts(ctx context.Context, req *rpc_product.SearchProductsReq) (*rpc_product.SearchProductsResp, error) {
	products, err := dao.Search(tidb.DB, req.Query)
	if err != nil {
		return nil, err
	}
	protoProducts := make([]*rpc_product.Product, 0, len(products))
	for _, p := range products {
		protoProducts = append(protoProducts, convertToProtoProduct(&p))
	}
	return &rpc_product.SearchProductsResp{
		Products: protoProducts,
	}, nil
}

func convertToProtoProduct(p *models.Product) *rpc_product.Product {
	return &rpc_product.Product{
		Id:          int32(p.ID),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       int32(p.Stock),
		Img:         p.Image,
		CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

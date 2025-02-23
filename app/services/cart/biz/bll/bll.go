package bll

import (
	"context"
	"errors"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart"
	rpc_cart "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product"
	rpc_product "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product"
	"github.com/bitdance-panic/gobuy/app/services/cart/biz/clients"
	"github.com/bitdance-panic/gobuy/app/services/cart/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/cart/biz/dao"
)

type CartBLL struct{}

func NewCartBLL() *CartBLL {
	return &CartBLL{}
}

func (*CartBLL) CreateItem(ctx context.Context, req *rpc_cart.CreateItemReq) (*rpc_cart.CreateItemResp, error) {
	// 先检查商品是否存在
	getProductResp, err := clients.ProductClient.GetProductByID(ctx, &rpc_product.GetProductByIDReq{Id: int32(req.ProductId)})
	if err != nil {
		return nil, err
	}
	if getProductResp.Product == nil {
		return nil, errors.New("product not found")
	}
	// 再添加
	ok, err := dao.Create(tidb.DB, int(req.UserId), int(req.ProductId))
	if err != nil {
		return nil, err
	}
	return &rpc_cart.CreateItemResp{Success: ok}, nil
}

func (*CartBLL) DeleteItem(ctx context.Context, req *rpc_cart.DeleteItemReq) (*rpc_cart.DeleteItemResp, error) {
	err := dao.Delete(tidb.DB, int(req.ItemId))
	if err != nil {
		return &cart.DeleteItemResp{Success: false}, err
	}
	return &cart.DeleteItemResp{Success: true}, nil
}

func (*CartBLL) UpdateQuantity(ctx context.Context, req *rpc_cart.UpdateQuantityReq) (*rpc_cart.UpdateQuantityResp, error) {
	item, err := dao.GetItemByID(tidb.DB, int(req.ItemId))
	if err != nil {
		return nil, err
	}
	getProductResp, err := clients.ProductClient.GetProductByID(ctx, &product.GetProductByIDReq{Id: int32(item.ProductID)})
	if err != nil {
		return nil, err
	}
	// 1. 商品已被删除
	if getProductResp.Product == nil {
		return nil, errors.New("product not found")
	}
	// 否则直接改，在订单结算时再考虑
	dao.UpdateQuantity(tidb.DB, item, int(req.NewQuantity_))
	return &rpc_cart.UpdateQuantityResp{
		Success: true,
	}, nil
}

func (*CartBLL) ListItem(ctx context.Context, req *rpc_cart.ListItemReq) (*rpc_cart.ListItemResp, error) {
	items, err := dao.ListItemsByUserID(tidb.DB, int(req.UserId), int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	protoItems := make([]*rpc_cart.CartItem, 0, len(*items))
	for _, item := range *items {
		getProductResp, err := clients.ProductClient.GetProductByID(ctx, &rpc_product.GetProductByIDReq{Id: int32(item.ProductID)})
		if err != nil {
			return nil, err
		}
		var protoItem *rpc_cart.CartItem
		if getProductResp.Product == nil {
			protoItem = &rpc_cart.CartItem{
				Valid: false,
			}
		} else {
			protoItem = convertToProtoCartItem(&item, getProductResp.Product)
		}
		protoItems = append(protoItems, protoItem)
	}
	return &rpc_cart.ListItemResp{
		Items: protoItems,
	}, nil
}

func convertToProtoCartItem(item *models.CartItem, p *rpc_product.Product) *rpc_cart.CartItem {
	return &rpc_cart.CartItem{
		Id:       int32(item.ID),
		Name:     p.Name,
		Price:    p.Price,
		Quantity: int32(item.Quantity),
		Img:      p.Img,
		Valid:    true,
	}
}

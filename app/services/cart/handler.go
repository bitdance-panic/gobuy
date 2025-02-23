package main

import (
	"context"

	"github.com/bitdance-panic/gobuy/app/services/cart/biz/bll"

	rpc_cart "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart"
)

var b *bll.CartBLL

func init() {
	b = bll.NewCartBLL()
}

type CartServiceImpl struct{}

func (*CartServiceImpl) CreateItem(ctx context.Context, req *rpc_cart.CreateItemReq) (resp *rpc_cart.CreateItemResp, err error) {
	return b.CreateItem(ctx, req)
}

func (*CartServiceImpl) DeleteItem(ctx context.Context, req *rpc_cart.DeleteItemReq) (resp *rpc_cart.DeleteItemResp, err error) {
	return b.DeleteItem(ctx, req)
}

func (*CartServiceImpl) UpdateQuantity(ctx context.Context, req *rpc_cart.UpdateQuantityReq) (resp *rpc_cart.UpdateQuantityResp, err error) {
	return b.UpdateQuantity(ctx, req)
}

func (*CartServiceImpl) ListItem(ctx context.Context, req *rpc_cart.ListItemReq) (resp *rpc_cart.ListItemResp, err error) {
	return b.ListItem(ctx, req)
}

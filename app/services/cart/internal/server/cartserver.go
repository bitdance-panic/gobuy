package server

import (
	"context"
	"github.com/bitdance-panic/gobuy/app/services/cart/internal/logic"
	"github.com/bitdance-panic/gobuy/app/services/cart/internal/svc"
	"github.com/bitdance-panic/gobuy/app/services/cart/proto/cart"
)

type CartServiceServer struct {
	svcCtx *svc.ServiceContext
	cart.UnimplementedCartServiceServer
}

func NewCartServiceServer(svcCtx *svc.ServiceContext) *CartServiceServer {
	return &CartServiceServer{
		svcCtx: svcCtx,
	}
}

// 封装逻辑处理对象的创建和调用
func (s *CartServiceServer) handleRequest(ctx context.Context, req interface{}, logicFunc func(context.Context, *svc.ServiceContext, interface{}) (interface{}, error)) (interface{}, error) {

	// 获取逻辑处理对象并调用其方法
	resp, err := logicFunc(ctx, s.svcCtx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *CartServiceServer) AddItem(ctx context.Context, in *cart.AddItemReq) (*cart.AddItemResp, error) {
	return s.handleRequest(ctx, in, func(ctx context.Context, svcCtx *svc.ServiceContext, req interface{}) (interface{}, error) {
		l := logic.NewAddItemLogic(ctx, svcCtx)
		return l.AddItem(req.(*cart.AddItemReq))
	}).(*cart.AddItemResp), nil
}

func (s *CartServiceServer) GetCart(ctx context.Context, in *cart.GetCartReq) (*cart.GetCartResp, error) {
	return s.handleRequest(ctx, in, func(ctx context.Context, svcCtx *svc.ServiceContext, req interface{}) (interface{}, error) {
		l := logic.NewGetCartLogic(ctx, svcCtx)
		return l.GetCart(req.(*cart.GetCartReq))
	}).(*cart.GetCartResp), nil
}

func (s *CartServiceServer) EmptyCart(ctx context.Context, in *cart.EmptyCartReq) (*cart.EmptyCartResp, error) {
	return s.handleRequest(ctx, in, func(ctx context.Context, svcCtx *svc.ServiceContext, req interface{}) (interface{}, error) {
		l := logic.NewEmptyCartLogic(ctx, svcCtx)
		return l.EmptyCart(req.(*cart.EmptyCartReq))
	}).(*cart.EmptyCartResp), nil
}

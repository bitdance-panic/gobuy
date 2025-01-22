package server

import (
	"context"
	"mall/service/cart/internal/logic"
	"mall/service/cart/internal/svc"
	"mall/service/cart/proto/cart"
	"github.com/zeromicro/go-zero/core/logx"
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

//封装逻辑处理对象的创建和调用
func (s *CartServiceServer) handleRequest[T any, R any](ctx context.Context, req T, logicFunc func(context.Context, *svc.ServiceContext) interface{}) (*R, error) {
	log := logx.WithContext(ctx)
	log.Infow("Handling request", "request", req)

	// 获取逻辑处理对象并调用其方法
	resp, err := logicFunc(ctx, s.svcCtx).(func(T) (*R, error))(req)
	if err != nil {
		log.Errorw("错误处理请求", "error", err, "request", req)
		return nil, err
	}

	log.Infow("成功处理请求", "response", resp)
	return resp, nil
}

func (s *CartServiceServer) AddItem(ctx context.Context, in *cart.AddItemReq) (*cart.AddItemResp, error) {
	return s.handleRequest[cart.AddItemReq, cart.AddItemResp](ctx, in, func(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
		return func(req *cart.AddItemReq) (*cart.AddItemResp, error) {
			l := logic.NewAddItemLogic(ctx, svcCtx)
			return l.AddItem(req)
		}
	})
}

func (s *CartServiceServer) GetCart(ctx context.Context, in *cart.GetCartReq) (*cart.GetCartResp, error) {
	return s.handleRequest[cart.GetCartReq, cart.GetCartResp](ctx, in, func(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
		return func(req *cart.GetCartReq) (*cart.GetCartResp, error) {
			l := logic.NewGetCartLogic(ctx, svcCtx)
			return l.GetCart(req)
		}
	})
}

func (s *CartServiceServer) EmptyCart(ctx context.Context, in *cart.EmptyCartReq) (*cart.EmptyCartResp, error) {
	return s.handleRequest[cart.EmptyCartReq, cart.EmptyCartResp](ctx, in, func(ctx context.Context, svcCtx *svc.ServiceContext) interface{} {
		return func(req *cart.EmptyCartReq) (*cart.EmptyCartResp, error) {
			l := logic.NewEmptyCartLogic(ctx, svcCtx)
			return l.EmptyCart(req)
		}
	})
}

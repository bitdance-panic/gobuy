package cartservice

import (
	"context"
	"github.com/bitdance-panic/gobuy/app/models" // 添加导入 models 包
	"github.com/bitdance-panic/gobuy/app/services/cart/proto/cart"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AddItemReq    = cart.ProtoAddItemReq
	AddItemResp   = cart.ProtoAddItemResp
	CartItem      = cart.ProtoCartItem
	EmptyCartReq  = models.EmptyCartReq
	EmptyCartResp = models.EmptyCartResp
	GetCartReq    = models.GetCartReq
	GetCartResp   = models.GetCartResp
	CartService   interface {
		AddItem(ctx context.Context, in *AddItemReq, opts ...grpc.CallOption) (*AddItemResp, error)
		GetCart(ctx context.Context, in *GetCartReq, opts ...grpc.CallOption) (*GetCartResp, error)
		EmptyCart(ctx context.Context, in *EmptyCartReq, opts ...grpc.CallOption) (*EmptyCartResp, error)
	}
)

func NewCartService(cli zrpc.Client) CartService {
	return &defaultCartService{
		cli: cli,
	}
}

func (m *defaultCartService) AddItem(ctx context.Context, in *AddItemReq, opts ...grpc.CallOption) (*AddItemResp, error) {
	client := cart.NewCartServiceClient(m.cli.Conn())
	return client.AddItem(ctx, in, opts...)
}

func (m *defaultCartService) GetCart(ctx context.Context, in *GetCartReq, opts ...grpc.CallOption) (*GetCartResp, error) {
	client := cart.NewCartServiceClient(m.cli.Conn())
	return client.GetCart(ctx, in, opts...)
}

func (m *defaultCartService) EmptyCart(ctx context.Context, in *EmptyCartReq, opts ...grpc.CallOption) (*EmptyCartResp, error) {
	client := cart.NewCartServiceClient(m.cli.Conn())
	return client.EmptyCart(ctx, in, opts...)
}

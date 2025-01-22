package cartservice

import (
    "context"
    "github.com/cloudwego/kitex/client"
    "github.com/cloudwego/kitex/pkg/klog"
    "github.com/cloudwego/kitex/pkg/rpcinfo"
    "github.com/kitex-contrib/obs-opentelemetry/tracing"
	"github.com/kitex-contrib/registry-nacos/nacos/resolver"
	"mall/service/cart/proto/cart"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)


type (
	AddItemReq    = cart.AddItemReq
	AddItemResp   = cart.AddItemResp
	CartItem      = cart.CartItem
	EmptyCartReq  = cart.EmptyCartReq
	EmptyCartResp = cart.EmptyCartResp
	GetCartReq    = cart.GetCartReq
	GetCartResp   = cart.GetCartResp

	CartService interface {
		AddItem(ctx context.Context, in *AddItemReq, opts ...grpc.CallOption) (*AddItemResp, error)
		GetCart(ctx context.Context, in *GetCartReq, opts ...grpc.CallOption) (*GetCartResp, error)
		EmptyCart(ctx context.Context, in *EmptyCartReq, opts ...grpc.CallOption) (*EmptyCartResp, error)
	}

	defaultCartService struct {
		cli zrpc.Client
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

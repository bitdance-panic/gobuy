package cart

import (
	"context"
	"io"

	"google.golang.org/grpc"
)

// CartServiceClient is the client API for CartService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CartServiceClient interface {
	GetCart(ctx context.Context, in *GetCartRequest, opts ...grpc.CallOption) (*GetCartResponse, error)
	AddCartItem(ctx context.Context, in *AddCartItemRequest, opts ...grpc.CallOption) (*AddCartItemResponse, error)
}

type cartServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCartServiceClient(cc grpc.ClientConnInterface) CartServiceClient {
	return &cartServiceClient{cc}
}

func (c *cartServiceClient) GetCart(ctx context.Context, in *GetCartRequest, opts ...grpc.CallOption) (*GetCartResponse, error) {
	out := new(GetCartResponse)
	err := c.cc.Invoke(ctx, "/cart.CartService/GetCart", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cartServiceClient) AddCartItem(ctx context.Context, in *AddCartItemRequest, opts ...grpc.CallOption) (*AddCartItemResponse, error) {
	out := new(AddCartItemResponse)
	err := c.cc.Invoke(ctx, "/cart.CartService/AddCartItem", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CartServiceServer is the server API for CartService service.
// All implementations must embed UnimplementedCartServiceServer
// for forward compatibility
type CartServiceServer interface {
	GetCart(context.Context, *GetCartRequest) (*GetCartResponse, error)
	AddCartItem(context.Context, *AddCartItemRequest) (*AddCartItemResponse, error)
	mustEmbedUnimplementedCartServiceServer()
}

// UnimplementedCartServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCartServiceServer struct {
}

func (UnimplementedCartServiceServer) GetCart(context.Context, *GetCartRequest) (*GetCartResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCart not implemented")
}

func (UnimplementedCartServiceServer) AddCartItem(context.Context, *AddCartItemRequest) (*AddCartItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCartItem not implemented")
}

func (UnimplementedCartServiceServer) mustEmbedUnimplementedCartServiceServer() {}

// UnsafeCartServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CartServiceServer will
// result in compilation errors.
type UnsafeCartServiceServer interface {
	mustEmbedUnimplementedCartServiceServer()
}

func RegisterCartServiceServer(s grpc.ServiceRegistrar, srv CartServiceServer) {
	s.RegisterService(&_CartService_serviceDesc, srv)
}

func _CartService_GetCart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CartServiceServer).GetCart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cart.CartService/GetCart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CartServiceServer).GetCart(ctx, req.(*GetCartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CartService_AddCartItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddCartItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CartServiceServer).AddCartItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cart.CartService/AddCartItem",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CartServiceServer).AddCartItem(ctx, req.(*AddCartItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CartService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cart.CartService",
	HandlerType: (*CartServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCart",
			Handler:    _CartService_GetCart_Handler,
		},
		{
			MethodName: "AddCartItem",
			Handler:    _CartService_AddCartItem_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cart.proto",
}

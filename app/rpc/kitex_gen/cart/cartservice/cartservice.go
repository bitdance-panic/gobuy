// Code generated by Kitex v0.12.2. DO NOT EDIT.

package cartservice

import (
	"context"
	"errors"
	cart "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"createItem": kitex.NewMethodInfo(
		createItemHandler,
		newCartServiceCreateItemArgs,
		newCartServiceCreateItemResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"updateQuantity": kitex.NewMethodInfo(
		updateQuantityHandler,
		newCartServiceUpdateQuantityArgs,
		newCartServiceUpdateQuantityResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"deleteItem": kitex.NewMethodInfo(
		deleteItemHandler,
		newCartServiceDeleteItemArgs,
		newCartServiceDeleteItemResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"listItem": kitex.NewMethodInfo(
		listItemHandler,
		newCartServiceListItemArgs,
		newCartServiceListItemResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
}

var (
	cartServiceServiceInfo                = NewServiceInfo()
	cartServiceServiceInfoForClient       = NewServiceInfoForClient()
	cartServiceServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return cartServiceServiceInfo
}

// for stream client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return cartServiceServiceInfoForStreamClient
}

// for client
func serviceInfoForClient() *kitex.ServiceInfo {
	return cartServiceServiceInfoForClient
}

// NewServiceInfo creates a new ServiceInfo containing all methods
func NewServiceInfo() *kitex.ServiceInfo {
	return newServiceInfo(false, true, true)
}

// NewServiceInfo creates a new ServiceInfo containing non-streaming methods
func NewServiceInfoForClient() *kitex.ServiceInfo {
	return newServiceInfo(false, false, true)
}
func NewServiceInfoForStreamClient() *kitex.ServiceInfo {
	return newServiceInfo(true, true, false)
}

func newServiceInfo(hasStreaming bool, keepStreamingMethods bool, keepNonStreamingMethods bool) *kitex.ServiceInfo {
	serviceName := "CartService"
	handlerType := (*cart.CartService)(nil)
	methods := map[string]kitex.MethodInfo{}
	for name, m := range serviceMethods {
		if m.IsStreaming() && !keepStreamingMethods {
			continue
		}
		if !m.IsStreaming() && !keepNonStreamingMethods {
			continue
		}
		methods[name] = m
	}
	extra := map[string]interface{}{
		"PackageName": "cart",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.12.2",
		Extra:           extra,
	}
	return svcInfo
}

func createItemHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*cart.CartServiceCreateItemArgs)
	realResult := result.(*cart.CartServiceCreateItemResult)
	success, err := handler.(cart.CartService).CreateItem(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newCartServiceCreateItemArgs() interface{} {
	return cart.NewCartServiceCreateItemArgs()
}

func newCartServiceCreateItemResult() interface{} {
	return cart.NewCartServiceCreateItemResult()
}

func updateQuantityHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*cart.CartServiceUpdateQuantityArgs)
	realResult := result.(*cart.CartServiceUpdateQuantityResult)
	success, err := handler.(cart.CartService).UpdateQuantity(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newCartServiceUpdateQuantityArgs() interface{} {
	return cart.NewCartServiceUpdateQuantityArgs()
}

func newCartServiceUpdateQuantityResult() interface{} {
	return cart.NewCartServiceUpdateQuantityResult()
}

func deleteItemHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*cart.CartServiceDeleteItemArgs)
	realResult := result.(*cart.CartServiceDeleteItemResult)
	success, err := handler.(cart.CartService).DeleteItem(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newCartServiceDeleteItemArgs() interface{} {
	return cart.NewCartServiceDeleteItemArgs()
}

func newCartServiceDeleteItemResult() interface{} {
	return cart.NewCartServiceDeleteItemResult()
}

func listItemHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*cart.CartServiceListItemArgs)
	realResult := result.(*cart.CartServiceListItemResult)
	success, err := handler.(cart.CartService).ListItem(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newCartServiceListItemArgs() interface{} {
	return cart.NewCartServiceListItemArgs()
}

func newCartServiceListItemResult() interface{} {
	return cart.NewCartServiceListItemResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) CreateItem(ctx context.Context, req *cart.CreateItemReq) (r *cart.CreateItemResp, err error) {
	var _args cart.CartServiceCreateItemArgs
	_args.Req = req
	var _result cart.CartServiceCreateItemResult
	if err = p.c.Call(ctx, "createItem", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) UpdateQuantity(ctx context.Context, req *cart.UpdateQuantityReq) (r *cart.UpdateQuantityResp, err error) {
	var _args cart.CartServiceUpdateQuantityArgs
	_args.Req = req
	var _result cart.CartServiceUpdateQuantityResult
	if err = p.c.Call(ctx, "updateQuantity", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) DeleteItem(ctx context.Context, req *cart.DeleteItemReq) (r *cart.DeleteItemResp, err error) {
	var _args cart.CartServiceDeleteItemArgs
	_args.Req = req
	var _result cart.CartServiceDeleteItemResult
	if err = p.c.Call(ctx, "deleteItem", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) ListItem(ctx context.Context, req *cart.ListItemReq) (r *cart.ListItemResp, err error) {
	var _args cart.CartServiceListItemArgs
	_args.Req = req
	var _result cart.CartServiceListItemResult
	if err = p.c.Call(ctx, "listItem", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

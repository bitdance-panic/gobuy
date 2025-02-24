// Code generated by Kitex v0.12.1. DO NOT EDIT.

package productservice

import (
	"context"
	"errors"
	product "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"createProduct": kitex.NewMethodInfo(
		createProductHandler,
		newProductServiceCreateProductArgs,
		newProductServiceCreateProductResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"updateProduct": kitex.NewMethodInfo(
		updateProductHandler,
		newProductServiceUpdateProductArgs,
		newProductServiceUpdateProductResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"removeProduct": kitex.NewMethodInfo(
		removeProductHandler,
		newProductServiceRemoveProductArgs,
		newProductServiceRemoveProductResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"getProductByID": kitex.NewMethodInfo(
		getProductByIDHandler,
		newProductServiceGetProductByIDArgs,
		newProductServiceGetProductByIDResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"listProduct": kitex.NewMethodInfo(
		listProductHandler,
		newProductServiceListProductArgs,
		newProductServiceListProductResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"adminListProduct": kitex.NewMethodInfo(
		adminListProductHandler,
		newProductServiceAdminListProductArgs,
		newProductServiceAdminListProductResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"searchProducts": kitex.NewMethodInfo(
		searchProductsHandler,
		newProductServiceSearchProductsArgs,
		newProductServiceSearchProductsResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
}

var (
	productServiceServiceInfo                = NewServiceInfo()
	productServiceServiceInfoForClient       = NewServiceInfoForClient()
	productServiceServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return productServiceServiceInfo
}

// for stream client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return productServiceServiceInfoForStreamClient
}

// for client
func serviceInfoForClient() *kitex.ServiceInfo {
	return productServiceServiceInfoForClient
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
	serviceName := "ProductService"
	handlerType := (*product.ProductService)(nil)
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
		"PackageName": "product",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.12.1",
		Extra:           extra,
	}
	return svcInfo
}

func createProductHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*product.ProductServiceCreateProductArgs)
	realResult := result.(*product.ProductServiceCreateProductResult)
	success, err := handler.(product.ProductService).CreateProduct(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newProductServiceCreateProductArgs() interface{} {
	return product.NewProductServiceCreateProductArgs()
}

func newProductServiceCreateProductResult() interface{} {
	return product.NewProductServiceCreateProductResult()
}

func updateProductHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*product.ProductServiceUpdateProductArgs)
	realResult := result.(*product.ProductServiceUpdateProductResult)
	success, err := handler.(product.ProductService).UpdateProduct(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newProductServiceUpdateProductArgs() interface{} {
	return product.NewProductServiceUpdateProductArgs()
}

func newProductServiceUpdateProductResult() interface{} {
	return product.NewProductServiceUpdateProductResult()
}

func removeProductHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*product.ProductServiceRemoveProductArgs)
	realResult := result.(*product.ProductServiceRemoveProductResult)
	success, err := handler.(product.ProductService).RemoveProduct(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newProductServiceRemoveProductArgs() interface{} {
	return product.NewProductServiceRemoveProductArgs()
}

func newProductServiceRemoveProductResult() interface{} {
	return product.NewProductServiceRemoveProductResult()
}

func getProductByIDHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*product.ProductServiceGetProductByIDArgs)
	realResult := result.(*product.ProductServiceGetProductByIDResult)
	success, err := handler.(product.ProductService).GetProductByID(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newProductServiceGetProductByIDArgs() interface{} {
	return product.NewProductServiceGetProductByIDArgs()
}

func newProductServiceGetProductByIDResult() interface{} {
	return product.NewProductServiceGetProductByIDResult()
}

func listProductHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*product.ProductServiceListProductArgs)
	realResult := result.(*product.ProductServiceListProductResult)
	success, err := handler.(product.ProductService).ListProduct(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newProductServiceListProductArgs() interface{} {
	return product.NewProductServiceListProductArgs()
}

func newProductServiceListProductResult() interface{} {
	return product.NewProductServiceListProductResult()
}

func adminListProductHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*product.ProductServiceAdminListProductArgs)
	realResult := result.(*product.ProductServiceAdminListProductResult)
	success, err := handler.(product.ProductService).AdminListProduct(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newProductServiceAdminListProductArgs() interface{} {
	return product.NewProductServiceAdminListProductArgs()
}

func newProductServiceAdminListProductResult() interface{} {
	return product.NewProductServiceAdminListProductResult()
}

func searchProductsHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*product.ProductServiceSearchProductsArgs)
	realResult := result.(*product.ProductServiceSearchProductsResult)
	success, err := handler.(product.ProductService).SearchProducts(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newProductServiceSearchProductsArgs() interface{} {
	return product.NewProductServiceSearchProductsArgs()
}

func newProductServiceSearchProductsResult() interface{} {
	return product.NewProductServiceSearchProductsResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) CreateProduct(ctx context.Context, req *product.CreateProductReq) (r *product.CreateProductResp, err error) {
	var _args product.ProductServiceCreateProductArgs
	_args.Req = req
	var _result product.ProductServiceCreateProductResult
	if err = p.c.Call(ctx, "createProduct", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) UpdateProduct(ctx context.Context, req *product.UpdateProductReq) (r *product.UpdateProductResp, err error) {
	var _args product.ProductServiceUpdateProductArgs
	_args.Req = req
	var _result product.ProductServiceUpdateProductResult
	if err = p.c.Call(ctx, "updateProduct", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) RemoveProduct(ctx context.Context, req *product.RemoveProductReq) (r *product.RemoveProductResp, err error) {
	var _args product.ProductServiceRemoveProductArgs
	_args.Req = req
	var _result product.ProductServiceRemoveProductResult
	if err = p.c.Call(ctx, "removeProduct", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetProductByID(ctx context.Context, req *product.GetProductByIDReq) (r *product.GetProductByIDResp, err error) {
	var _args product.ProductServiceGetProductByIDArgs
	_args.Req = req
	var _result product.ProductServiceGetProductByIDResult
	if err = p.c.Call(ctx, "getProductByID", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) ListProduct(ctx context.Context, req *product.ListProductReq) (r *product.ListProductResp, err error) {
	var _args product.ProductServiceListProductArgs
	_args.Req = req
	var _result product.ProductServiceListProductResult
	if err = p.c.Call(ctx, "listProduct", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) AdminListProduct(ctx context.Context, req *product.ListProductReq) (r *product.ListProductResp, err error) {
	var _args product.ProductServiceAdminListProductArgs
	_args.Req = req
	var _result product.ProductServiceAdminListProductResult
	if err = p.c.Call(ctx, "adminListProduct", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) SearchProducts(ctx context.Context, req *product.SearchProductsReq) (r *product.SearchProductsResp, err error) {
	var _args product.ProductServiceSearchProductsArgs
	_args.Req = req
	var _result product.ProductServiceSearchProductsResult
	if err = p.c.Call(ctx, "searchProducts", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

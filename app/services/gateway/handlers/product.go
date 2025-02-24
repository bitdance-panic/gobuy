package handlers

import (
	"context"
	"strconv"
	"time"

	rpc_product "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product"
	clients "github.com/bitdance-panic/gobuy/app/services/gateway/biz/clients"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client/callopt"
)

// handleProductPut 这是更新商品
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /product [put]
func HandleUpdateProduct(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	var req rpc_product.UpdateProductReq
	if err := c.BindAndValidate(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	req.Id = int32(id)
	resp, err := clients.ProductClient.UpdateProduct(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"newProduct": resp.Product})
}

// handleProductPost 这是创建商品
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /product [post]
func HandleCreateProduct(ctx context.Context, c *app.RequestContext) {
	var req rpc_product.CreateProductReq
	if err := c.Bind(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	resp, err := clients.ProductClient.CreateProduct(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"product": resp.Product})
}

// handleProductDELETE 这是删除商品
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /product [delete]
func HandleRemoveProduct(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	req := rpc_product.RemoveProductReq{
		Id: int32(id),
	}
	resp, err := clients.ProductClient.RemoveProduct(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	if resp.Success {
		utils.Success(c, nil)
	} else {
		utils.Fail(c, "删除失败")
	}
}

// handleProductGet 这是获取一个商品
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /product [get]
func HandleGetProduct(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	req := rpc_product.GetProductByIDReq{
		Id: int32(id),
	}
	resp, err := clients.ProductClient.GetProductByID(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"product": resp.Product})
}

// handleProductSearch 这是模糊查询商品
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /product/search [get]
func HandleSearchProducts(ctx context.Context, c *app.RequestContext) {
	pageNumStr := c.Query("page")
	pageSizeStr := c.Query("size")
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	// var body struct {
	// 	Query string `json:"query"`
	// }
	// err = c.Bind(&body)
	// if err != nil {
	// 	utils.Fail(c, err.Error())
	// 	return
	// }
	query := c.Query("query")
	req := rpc_product.SearchProductsReq{
		Query:    query,
		PageNum:  int32(pageNum),
		PageSize: int32(pageSize),
	}
	resp, err := clients.ProductClient.SearchProducts(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"products": resp.Products, "total_count": resp.TotalCount})
}

// handleProductSearch 这是模糊查询商品
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /product/search [get]
func HandleAdminListProduct(ctx context.Context, c *app.RequestContext) {
	pageNum, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	pageSize, err := strconv.Atoi(c.Query("size"))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	req := rpc_product.ListProductReq{
		PageNum:  int32(pageNum),
		PageSize: int32(pageSize),
	}
	resp, err := clients.ProductClient.AdminListProduct(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"products": resp.Products, "total_count": resp.TotalCount})
}

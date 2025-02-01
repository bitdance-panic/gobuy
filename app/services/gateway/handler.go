package main

import (
	"context"
	"strconv"
	"time"

	rpc_product "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product"
	rpc_user "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client/callopt"
)

func handlePong(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, utils.H{"message": "pong"})
}

// handleLogin 这是一个handler
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /login [get]
func handleLogin(ctx context.Context, c *app.RequestContext) {
	// 通过 /login?email=1234&pass=1234 测试
	email := c.Query("email")
	password := c.Query("pass")
	req := rpc_user.LoginReq{
		Email:    email,
		Password: password,
	}
	resp, err := userservice.Login(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{"message": err.Error()})
		return
	}
	if resp.Success {
		c.JSON(consts.StatusOK, utils.H{"userID": resp.UserId})
	} else {
		c.JSON(consts.StatusOK, utils.H{"message": "登录失败"})
	}
}

// handleProductPut 这是更新商品
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /product [put]
func handleProductPut(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusOK, utils.H{"message": "参数错误"})
		return
	}
	var req rpc_product.UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(consts.StatusOK, utils.H{"message": "参数错误"})
		return
	}
	req.Id = int32(id)
	// todo 打印日志
	resp, err := productservice.UpdateProduct(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{"message": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, resp.Product)
}

// handleProductPost 这是创建商品
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /product [post]
func handleProductPost(ctx context.Context, c *app.RequestContext) {
	var req rpc_product.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(consts.StatusOK, utils.H{"message": "参数错误"})
		return
	}
	// todo 打印日志
	resp, err := productservice.CreateProduct(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{"message": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, resp.Product)
}

// handleProductDELETE 这是删除商品
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /product [delete]
func handleProductDELETE(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusOK, utils.H{"message": "参数错误"})
		return
	}
	req := rpc_product.DeleteProductRequest{
		Id: int32(id),
	}
	resp, err := productservice.DeleteProduct(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{"message": err.Error()})
		return
	}
	if resp.Success {
		c.JSON(consts.StatusOK, utils.H{"message": "删除成功"})
	} else {
		c.JSON(consts.StatusOK, utils.H{"message": "删除失败"})
	}
}

// handleProductGet 这是获取一个商品
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /product [get]
func handleProductGet(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusOK, utils.H{"message": "参数错误"})
		return
	}
	req := rpc_product.GetProductByIDRequest{
		Id: int32(id),
	}
	resp, err := productservice.GetProductByID(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{"message": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, resp.Product)
}

// handleProductSearch 这是模糊查询商品
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /product/search [get]
func handleProductSearch(ctx context.Context, c *app.RequestContext) {
	query := c.Query("query")
	req := rpc_product.SearchProductsRequest{
		Query: query,
	}
	resp, err := productservice.SearchProducts(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{"message": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, resp.Products)
}

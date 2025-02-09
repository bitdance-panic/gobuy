package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	rpc_product "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product"
	rpc_user "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user"

	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dao"
	"github.com/bitdance-panic/gobuy/app/services/gateway/middleware"
	gutils "github.com/bitdance-panic/gobuy/app/services/gateway/utils"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client/callopt"
)

func handlePong(ctx context.Context, c *app.RequestContext) {
	utils.Success(c, utils.H{"message": "pong"})
}

// handleLogin 这是一个handler
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /login [get]

// RefreshTokenHandler 处理 Token 续期
func RefreshTokenHandler(ctx context.Context, c *app.RequestContext) {
	// 从请求头获取 Access Token
	accessToken := string(c.GetHeader("Authorization"))
	if accessToken == "" {
		utils.FailFull(c, consts.StatusUnauthorized, "Missing Access Token.", nil)
		return
	}

	// 验证 access_token 是否有效
	claims, err := middleware.AuthMiddleware.GetClaimsFromJWT(ctx, c)
	if err == nil {
		// access_token 仍然有效，不需要刷新
		utils.FailFull(c, consts.StatusOK, "Access Token is still valid", nil)
		return
	}

	// access_token 过期，提取 uid
	str := fmt.Sprintf("%v", claims["uid"])
	userID, err := strconv.Atoi(str)

	// 查询数据库中的 refresh_token
	storedRefreshToken, err := dao.GetRefreshTokenByUserID(tidb.DB, userID)
	if err != nil {
		hlog.Warnf("Refresh Token not found for user %d", userID)
		utils.FailFull(c, consts.StatusUnauthorized, "Refresh Token Not Found", nil)
		return
	}

	// 验证 Refresh Token 是否过期
	if gutils.IsRefreshTokenExpired(storedRefreshToken) {
		utils.FailFull(c, consts.StatusUnauthorized, "Refresh Token Expired, Please Login Again", nil)
		return
	}

	// 返回新的 access_token
	tokenString, _, err := middleware.AuthMiddleware.TokenGenerator(userID)
	if err != nil {
		utils.FailFull(c, consts.StatusInternalServerError, "Token Generation Failed", nil)
		return
	}
	utils.Success(c, map[string]interface{}{
		"access_token": tokenString,
	})
}

func handleLogin(ctx context.Context, c *app.RequestContext) {
	req := rpc_user.LoginReq{}

	// 从请求体中绑定参数并验证
	if err := c.BindAndValidate(&req); err != nil {
		hlog.Warnf("Login failed for email: %s, validation error: %v", req.Email, err)
		utils.Fail(c, err.Error())
		return
	}

	resp, err := userservice.Login(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	if resp.Success {
		utils.Success(c, utils.H{"userID": resp.UserId})
	} else {
		utils.FailFull(c, consts.StatusUnauthorized, "Login failed. Invalid email or password.", nil)
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
		utils.Fail(c, "参数错误")
		return
	}
	var req rpc_product.UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	req.Id = int32(id)
	// todo 打印日志
	resp, err := productservice.UpdateProduct(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"product": resp.Product})
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
		utils.Fail(c, "参数错误")
		return
	}
	// todo 打印日志
	resp, err := productservice.CreateProduct(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
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
func handleProductDelete(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	req := rpc_product.DeleteProductRequest{
		Id: int32(id),
	}
	resp, err := productservice.DeleteProduct(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
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
func handleProductGet(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	req := rpc_product.GetProductByIDRequest{
		Id: int32(id),
	}
	resp, err := productservice.GetProductByID(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
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
func handleProductSearch(ctx context.Context, c *app.RequestContext) {
	query := c.Query("query")
	req := rpc_product.SearchProductsRequest{
		Query: query,
	}
	resp, err := productservice.SearchProducts(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"products": resp.Products})
}

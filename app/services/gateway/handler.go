package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	rpc_payment "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/payment"
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
	"github.com/hertz-contrib/jwt"
)

func handlePong(ctx context.Context, c *app.RequestContext) {
	utils.Success(c, utils.H{"message": "pong"})
}

// 用户注册
func RegisterHandler(ctx context.Context, c *app.RequestContext) {
	req := rpc_user.RegisterReq{}

	// 从请求体中绑定参数并验证
	if err := c.BindAndValidate(&req); err != nil {
		hlog.Warnf("User register failed, validation error: %v", err)
		utils.Fail(c, err.Error())
		return
	}

	resp, err := userservice.Register(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	} else {
		utils.Success(c, utils.H{"userID": resp.UserId})
	}
}

// 封禁用户
func DeleteUserHandler(ctx context.Context, c *app.RequestContext) {
	req := rpc_user.DeleteUserReq{}

	// 从请求体中绑定参数并验证
	if err := c.BindAndValidate(&req); err != nil {
		hlog.Warnf("User deletion failed for user id: %s, validation error: %v", req.UserId, err)
		utils.Fail(c, err.Error())
		return
	}

	resp, err := userservice.DeleteUser(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	if resp.Success {
		utils.Success(c, utils.H{"userID": req.UserId})
	} else {
		utils.FailFull(c, consts.StatusInternalServerError, "User deletion failed", nil)
	}
}

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
	userID := int(claims["uid"].(float64))

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

// 获取用户信息
func GetUserHandler(ctx context.Context, c *app.RequestContext) {
	claims := jwt.ExtractClaims(ctx, c)
	// userID := int(claims["uid"].(float64))
	userID := fmt.Sprintf("%v", claims["uid"].(float64))
	req := rpc_user.GetUserReq{UserId: userID}

	resp, err := userservice.GetUser(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	if resp.Success {
		utils.Success(c, utils.H{
			"userID":   resp.UserId,
			"email":    resp.Email,
			"username": resp.Username,
		})
	} else {
		utils.FailFull(c, consts.StatusInternalServerError, "Get user failed", nil)
	}
}

// 更新用户信息
func UpdateUserHandler(ctx context.Context, c *app.RequestContext) {
	claims := jwt.ExtractClaims(ctx, c)
	// userID := int(claims["uid"].(float64))
	userID := fmt.Sprintf("%v", claims["uid"].(float64))

	req := rpc_user.UpdateUserReq{UserId: userID}
	if err := c.BindAndValidate(&req); err != nil {
		hlog.Warnf("User update failed for user id: %s, validation error: %v", req.UserId, err)
		utils.Fail(c, err.Error())
		return
	}

	resp, err := userservice.UpdateUser(context.Background(), &req)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	if resp.Success {
		utils.Success(c, utils.H{"userID": req.UserId})
	} else {
		utils.FailFull(c, consts.StatusInternalServerError, "User update failed", nil)
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

// handleCreatePayment 创建支付记录
// @Summary 创建支付记录
// @Description 创建一个新的支付记录并返回相关信息
// @Accept application/json
// @Produce application/json
// @Router /payment [post]
func handleCreatePayment(ctx context.Context, c *app.RequestContext) {
	req := rpc_payment.CreatePaymentRequest{}

	// 从请求体中绑定参数并验证
	if err := c.BindAndValidate(&req); err != nil {
		hlog.Warnf("CreatePayment failed , validation error: %v", err)
		utils.Fail(c, err.Error())
		return
	}

	resp, err := paymentservice.CreatePayment(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.FailFull(c, consts.StatusInternalServerError, fmt.Sprintf("Create payment failed: %v", err.Error()), nil)
		return
	} else {
		utils.Success(c, utils.H{"payment": resp.Payment})
	}
}

// handleGetPayment 获取支付记录
// @Summary 获取支付记录
// @Description 根据支付记录ID获取支付记录详情
// @Accept application/json
// @Produce application/json
// @Router /payment/{id} [get]
func handleGetPayment(ctx context.Context, c *app.RequestContext) {
	req := rpc_payment.GetPaymentRequest{}

	// 从请求体中绑定参数并验证
	if err := c.BindAndValidate(&req); err != nil {
		hlog.Warnf("GetPayment failed , validation error: %v", err)
		utils.Fail(c, err.Error())
		return
	}

	resp, err := paymentservice.GetPayment(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.FailFull(c, consts.StatusNotFound, fmt.Sprintf("Payment not found: %v", err.Error()), nil)
		return
	} else {
		utils.Success(c, utils.H{"payment": resp.Payment})
	}
}

// handleUpdatePayment 更新支付记录
// @Summary 更新支付记录
// @Description 更新指定支付记录的状态或其他信息
// @Accept application/json
// @Produce application/json
// @Router /payment/{id} [put]
func handleUpdatePayment(ctx context.Context, c *app.RequestContext) {
	req := rpc_payment.UpdatePaymentRequest{}

	// 从请求体中绑定参数并验证
	if err := c.BindAndValidate(&req); err != nil {
		hlog.Warnf("UpdatePayment failed , validation error: %v", err)
		utils.Fail(c, err.Error())
		return
	}

	resp, err := paymentservice.UpdatePayment(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.FailFull(c, consts.StatusInternalServerError, fmt.Sprintf("UpdatePayment failed: %v", err.Error()), nil)
		return
	} else {
		utils.Success(c, utils.H{"payment": resp.Payment})
	}
}

// handleDeletePayment 删除支付记录
// @Summary 删除支付记录
// @Description 删除指定ID的支付记录
// @Accept application/json
// @Produce application/json
// @Router /payment/{id} [delete]
func handleDeletePayment(ctx context.Context, c *app.RequestContext) {
	req := rpc_payment.DeletePaymentRequest{}

	// 从请求体中绑定参数并验证
	if err := c.BindAndValidate(&req); err != nil {
		hlog.Warnf("DeletePayment failed , validation error: %v", err)
		utils.Fail(c, err.Error())
		return
	}

	_, err := paymentservice.DeletePayment(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.FailFull(c, consts.StatusInternalServerError, fmt.Sprintf("DeletePayment failed: %v", err.Error()), nil)
		return
	} else {
		utils.Success(c, utils.H{"payment_id": req.PaymentId})
	}
}

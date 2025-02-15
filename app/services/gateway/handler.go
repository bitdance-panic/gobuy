package main

import (
	"context"
	"net/http"
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

// 获取用户信息
func GetUserHandler(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("userid") // 获取 URL 参数中的 userid
	if userID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  "User ID is required",
		})
		return
	}

	// 调用 user 服务的 GetUser RPC
	req := rpc_user.GetUserReq{UserId: userID}
	resp, err := userservice.GetUser(context.Background(), &req)

	if err != nil || !resp.Success {
		c.JSON(http.StatusNotFound, map[string]interface{}{
			"code": http.StatusNotFound,
			"msg":  "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":     http.StatusOK,
		"msg":      "User found",
		"user_id":  resp.UserId,
		"email":    resp.Email,
		"username": resp.Username,
	})
}

// 用户信息更新
func UpdateUserHandler(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("userid") // 获取 URL 参数中的 userid
	if userID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  "User ID is required",
		})
		return
	}

	var req rpc_user.UpdateUserReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  "Invalid input",
		})
		return
	}

	// 解析 Query 参数
	email := c.Query("email")
	username := c.Query("username")

	// 赋值到 `req`
	req.UserId = userID
	req.Email = &email
	req.Username = &username

	// 赋值 user_id
	req.UserId = userID

	// 调用 user 服务的 UpdateUser RPC
	resp, err := userservice.UpdateUser(context.Background(), &req)

	if err != nil || !resp.Success {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code": http.StatusInternalServerError,
			"msg":  "Update failed",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg":  "User updated successfully",
	})
}

// 封禁用户
func DeleteUserHandler(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("userid") // 获取 URL 参数中的 userid
	if userID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  "User ID is required",
		})
		return
	}

	// 构造请求
	req := rpc_user.DeleteUserReq{UserId: userID}

	// 调用 user 服务的 DeleteUser RPC
	resp, err := userservice.DeleteUser(context.Background(), &req)

	if err != nil || !resp.Success {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code": http.StatusInternalServerError,
			"msg":  "User deletion failed",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg":  "User deleted successfully",
	})
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
	// 解析 Query 参数
	userId := c.Query("user_id")
	orderId := c.Query("order_id")
	amount := c.Query("amount")

	// 转换 userId 和 orderId 为 int32
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  "Invalid user_id",
		})
		return
	}

	orderIdInt, err := strconv.Atoi(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  "Invalid order_id",
		})
		return
	}

	// 转换 amount 为 float64
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  "Invalid amount",
		})
		return
	}

	// 创建 req 并赋值
	req := rpc_payment.CreatePaymentRequest{
		UserId:  int32(userIdInt),  // 转换为 int32
		OrderId: int32(orderIdInt), // 转换为 int32
		Amount:  amountFloat,       // 转换为 float64
	}

	// 调用支付服务创建支付记录
	resp, err := paymentservice.CreatePayment(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}

	// 返回创建的支付记录信息
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"msg":     "成功创建支付订单",
		"payment": resp.Payment,
	})
}

// handleGetPayment 获取支付记录
// @Summary 获取支付记录
// @Description 根据支付记录ID获取支付记录详情
// @Accept application/json
// @Produce application/json
// @Router /payment/{id} [get]
func handleGetPayment(ctx context.Context, c *app.RequestContext) {
	payment_id := c.Param("paymentId")
	if payment_id == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  "Payment ID is required",
		})
		return
	}

	// 转换 payment_id 为 int32
	paymentIdInt, err := strconv.Atoi(payment_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  "Invalid payment_id",
		})
		return
	}

	// 调用支付服务获取支付记录
	req := rpc_payment.GetPaymentRequest{PaymentId: int32(paymentIdInt)}
	resp, err := paymentservice.GetPayment(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]interface{}{
			"code": http.StatusNotFound,
			"msg":  "Payment not found",
		})
		return
	}

	// 返回支付记录详情
	utils.Success(c, utils.H{"payment": resp.Payment})
}

// handleUpdatePayment 更新支付记录
// @Summary 更新支付记录
// @Description 更新指定支付记录的状态或其他信息
// @Accept application/json
// @Produce application/json
// @Router /payment/{id} [put]
func handleUpdatePayment(ctx context.Context, c *app.RequestContext) {
	payment_id := c.Param("paymentId")
	status := c.Query("status")

	if payment_id == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  "Payment ID is required",
		})
		return
	}

	paymentIdInt, err := strconv.Atoi(payment_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  "Invalid payment_id",
		})
		return
	}

	statusInt64, err := strconv.ParseInt(status, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  "Invalid status",
		})
		return
	}

	// 调用支付服务获取支付记录
	req := rpc_payment.UpdatePaymentRequest{
		PaymentId: int32(paymentIdInt),
		Status:    rpc_payment.PaymentStatus(statusInt64), // 进行类型转换
	}

	// 调用支付服务更新支付记录
	resp, err := paymentservice.UpdatePayment(context.Background(), &req)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	// 返回更新后的支付记录信息
	utils.Success(c, utils.H{"payment": resp.Payment})
}

// handleDeletePayment 删除支付记录
// @Summary 删除支付记录
// @Description 删除指定ID的支付记录
// @Accept application/json
// @Produce application/json
// @Router /payment/{id} [delete]
func handleDeletePayment(ctx context.Context, c *app.RequestContext) {
	payment_id := c.Param("paymentId")
	if payment_id == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  "Payment ID is required",
		})
		return
	}

	// 转换 payment_id 为 int32
	paymentIdInt, err := strconv.Atoi(payment_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  "Invalid payment_id",
		})
		return
	}

	// 调用支付服务获取支付记录
	req := rpc_payment.DeletePaymentRequest{PaymentId: int32(paymentIdInt)}
	resp, err := paymentservice.DeletePayment(context.Background(), &req)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	if resp.Success {
		utils.Success(c, nil)
	} else {
		utils.Fail(c, "删除支付记录失败")
	}
}

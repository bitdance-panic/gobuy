package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/bitdance-panic/gobuy/app/consts"
	rpc_order "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order"
	clients "github.com/bitdance-panic/gobuy/app/services/gateway/biz/clients"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client/callopt"
)

func HandleCreateOrder(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt(consts.CONTEXT_UID_KEY)
	var body struct {
		ItemIDs      []int  `json:"itemIDs"`
		Phone        string `json:"phone"`
		OrderAddress string `json:"order_address"`
	}
	if err := c.Bind(&body); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	itemIDs := make([]int32, len(body.ItemIDs))
	for i, itemID := range body.ItemIDs {
		itemIDs[i] = int32(itemID)
	}
	req := rpc_order.CreateOrderReq{
		UserId:       int32(userID),
		CartItemIDs:  itemIDs,
		Phone:        body.Phone,
		OrderAddress: body.OrderAddress,
	}
	resp, err := clients.OrderClient.CreateOrder(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"order": resp.Order})
}

func HandleUpdateOrderAddress(ctx context.Context, c *app.RequestContext) {
	// 1. 绑定请求参数
	var body struct {
		OrderID      int    `json:"order_id"`
		OrderAddress string `json:"order_address" binding:"required"`
	}
	if err := c.Bind(&body); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	// 2. 参数验证
	if body.OrderID <= 0 {
		utils.Fail(c, "无效的订单ID")
		return
	}
	if body.OrderAddress == "" {
		utils.Fail(c, "地址不能为空")
		return
	}

	// 3. 构建RPC请求
	req := rpc_order.UpdateOrderAddressReq{
		OrderId:      int32(body.OrderID),
		OrderAddress: body.OrderAddress,
	}

	// 4. 调用服务
	resp, err := clients.OrderClient.UpdateOrderAddress(
		context.Background(),
		&req,
		callopt.WithRPCTimeout(3*time.Second),
	)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	// 5. 返回响应
	utils.Success(c, utils.H{
		"success": resp.Success,
		"message": "地址更新成功",
	})
}

func HandleGetOrder(ctx context.Context, c *app.RequestContext) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	req := rpc_order.GetOrderReq{
		OrderId: int32(orderID),
	}
	resp, err := clients.OrderClient.GetOrder(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"order": resp.Order})
}

func HandleListUserOrder(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt(consts.CONTEXT_UID_KEY)
	req := rpc_order.ListOrderReq{
		UserId:   int32(userID),
		PageNum:  1,
		PageSize: 1,
	}
	resp, err := clients.OrderClient.ListUserOrder(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"orders": resp.Orders})
}

func HandleAdminListOrder(ctx context.Context, c *app.RequestContext) {
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
	req := rpc_order.ListOrderReq{
		PageNum:  int32(pageNum),
		PageSize: int32(pageSize),
	}
	resp, err := clients.OrderClient.AdminListOrder(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"orders": resp.Orders, "total_count": resp.TotalCount})
}

func HandleCreateUserAddress(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt(consts.CONTEXT_UID_KEY)
	var body struct {
		UserId      int32  `json:"userID"`
		Phone       string `json:"phone"`
		UserAddress string `json:"userAddress"`
	}
	if err := c.Bind(&body); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	req := rpc_order.CreateUserAddressReq{
		UserId:      int32(userID),
		Phone:       body.Phone,
		UserAddress: body.UserAddress,
	}

	// 调用 OrderClient 的 CreateOrderAddress 方法
	resp, err := clients.OrderClient.CreateUserAddress(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	utils.Success(c, utils.H{"userAddress": resp.UserId, "success": resp.Success})
}

func HandleDeleteUserAddress(ctx context.Context, c *app.RequestContext) {
	var body struct {
		UserID int32 `json:"userID"`
	}
	if err := c.Bind(&body); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	req := rpc_order.DeleteUserAddressReq{
		UserId: body.UserID,
	}

	// 调用 OrderClient 的 DeleteOrderAddress 方法
	resp, err := clients.OrderClient.DeleteUserAddress(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	utils.Success(c, utils.H{"userID": resp.UserId, "success": resp.Success})
}

func HandleUpdateUserAddress(ctx context.Context, c *app.RequestContext) {
	var body struct {
		UserID      int32  `json:"userID"`
		UserAddress string `json:"userAddress"`
	}
	if err := c.Bind(&body); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	req := rpc_order.UpdateUserAddressReq{
		UserId:      body.UserID,
		UserAddress: body.UserAddress,
	}

	resp, err := clients.OrderClient.UpdateUserAddress(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	utils.Success(c, utils.H{"userAddress": resp.UserAddress, "success": resp.Success})
}

func HandleGetUserAddress(ctx context.Context, c *app.RequestContext) {
	var body struct {
		UserID int32 `json:"userID"`
	}
	if err := c.Bind(&body); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	req := rpc_order.GetUserAddressReq{
		UserId: body.UserID,
	}

	// 调用 OrderClient 的 GetOrderAddress 方法
	resp, err := clients.OrderClient.GetUserAddress(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	utils.Success(c, utils.H{"userAddress": resp.UserAddress})
}

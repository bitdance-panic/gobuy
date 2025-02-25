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
		ItemIDs []int `json:"itemIDs"`
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
		UserId:      int32(userID),
		CartItemIDs: itemIDs,
	}
	resp, err := clients.OrderClient.CreateOrder(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"order": resp.Order})
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

func HandleCreateOrderAddress(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt(consts.CONTEXT_UID_KEY)
	var body struct {
		Phone        string `json:"phone"`
		OrderID      int32  `json:"orderID"`
		OrderAddress string `json:"orderAddress"`
	}
	if err := c.Bind(&body); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	req := rpc_order.CreateOrderAddressReq{
		UserId:       int32(userID),
		Phone:        body.Phone,
		OrderId:      body.OrderID,
		OrderAddress: body.OrderAddress,
	}

	// 调用 OrderClient 的 CreateOrderAddress 方法
	resp, err := clients.OrderClient.CreateOrderAddress(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	utils.Success(c, utils.H{"orderAddress": resp.OrderId, "success": resp.Success})
}

func HandleDeleteOrderAddress(ctx context.Context, c *app.RequestContext) {
	var body struct {
		OrderID int32 `json:"orderID"`
	}
	if err := c.Bind(&body); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	req := rpc_order.DeleteOrderAddressReq{
		OrderId: body.OrderID,
	}

	// 调用 OrderClient 的 DeleteOrderAddress 方法
	resp, err := clients.OrderClient.DeleteOrderAddress(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	utils.Success(c, utils.H{"orderID": resp.OrderId, "success": resp.Success})
}

func HandleUpdateOrderAddress(ctx context.Context, c *app.RequestContext) {
	var body struct {
		OrderID      int32  `json:"orderID"`
		OrderAddress string `json:"orderAddress"`
	}
	if err := c.Bind(&body); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	req := rpc_order.UpdateOrderAddressReq{
		OrderId:      body.OrderID,
		OrderAddress: body.OrderAddress,
	}

	// 调用 OrderClient 的 UpdateOrderAddress 方法
	resp, err := clients.OrderClient.UpdateOrderAddress(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	utils.Success(c, utils.H{"orderAddress": resp.OrderAddress, "success": resp.Success})
}

func HandleGetOrderAddress(ctx context.Context, c *app.RequestContext) {
	var body struct {
		OrderID int32 `json:"orderID"`
	}
	if err := c.Bind(&body); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	req := rpc_order.GetOrderAddressReq{
		OrderId: body.OrderID,
	}

	// 调用 OrderClient 的 GetOrderAddress 方法
	resp, err := clients.OrderClient.GetOrderAddress(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	utils.Success(c, utils.H{"orderAddress": resp.OrderAddress})
}

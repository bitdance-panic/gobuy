package handlers

import (
	"context"
	"strconv"
	"time"

	rpc_order "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order"
	clients "github.com/bitdance-panic/gobuy/app/services/gateway/biz/clients"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client/callopt"
)

func HandleCreateOrder(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	req := rpc_order.CreateOrderReq{
		UserId: 1,
		Items: []*rpc_order.OrderProductItem{
			{
				ProductId: int32(id),
				Quantity:  1,
			},
		},
	}
	resp, err := clients.OrderClient.CreateOrder(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"order": resp.Order})
}

func HandleGetOrder(ctx context.Context, c *app.RequestContext) {
	req := rpc_order.GetOrderReq{
		OrderId: 1,
	}
	resp, err := clients.OrderClient.GetOrder(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"order": resp.Order})
}

func HandleListUserOrder(ctx context.Context, c *app.RequestContext) {
	req := rpc_order.ListOrderReq{
		UserId:   1,
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
	utils.Success(c, utils.H{"orders": resp.Orders})
}

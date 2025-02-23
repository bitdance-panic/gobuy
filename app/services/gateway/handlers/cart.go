package handlers

import (
	"context"
	"strconv"
	"time"

	rpc_cart "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart"
	clients "github.com/bitdance-panic/gobuy/app/services/gateway/biz/clients"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client/callopt"
)

// HandleAddToCart 这是更新商品
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /product [put]
func HandleCreateCartItem(ctx context.Context, c *app.RequestContext) {
	// productID, err := strconv.Atoi(c.Param("id"))
	var productID int32
	// if err != nil {
	// 	utils.Fail(c, "参数错误")
	// 	return
	// }
	req := rpc_cart.CreateItemReq{
		UserId:    1,
		ProductId: productID,
	}
	resp, err := clients.CartClient.CreateItem(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil || !resp.Success {
		utils.Fail(c, err.Error())
		return
	}
	// 返回是否成功即可
	utils.Success(c, utils.H{})
}

// handleProductPost 这是创建商品
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /product [post]
func HandleListCartItem(ctx context.Context, c *app.RequestContext) {
	req := rpc_cart.ListItemReq{
		UserId: 1,
	}
	resp, err := clients.CartClient.ListItem(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"cartItems": resp.Items})
}

// handleProductDELETE 这是删除商品
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /product [delete]
func HandleDeleteCartItem(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	req := rpc_cart.DeleteItemReq{
		ItemId: int32(id),
	}
	resp, err := clients.CartClient.DeleteItem(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
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
func HandleUpdateCartItemQuantity(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	req := rpc_cart.UpdateQuantityReq{
		ItemId:       int32(id),
		NewQuantity_: 1,
	}
	resp, err := clients.CartClient.UpdateQuantity(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"success": resp.Success})
}

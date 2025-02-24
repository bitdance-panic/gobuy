package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	rpc_product "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product"
	rpc_user "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/clients"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client/callopt"
)

// 用户注册
func HandleRegister(ctx context.Context, c *app.RequestContext) {
	req := rpc_user.RegisterReq{}

	// 从请求体中绑定参数并验证
	if err := c.BindAndValidate(&req); err != nil {
		hlog.Warnf("User register failed, validation error: %v", err)
		utils.Fail(c, err.Error())
		return
	}

	resp, err := clients.UserClient.Register(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil || !resp.Success {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"userID": resp.UserId})
}

// 搜索首页商品
func HandleListIndexProduct(ctx context.Context, c *app.RequestContext) {
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
	if err := c.BindAndValidate(&req); err != nil {
		hlog.Warnf("ListIndexProduct failed, validation error: %v", err)
		utils.Fail(c, err.Error())
		return
	}
	resp, err := clients.ProductClient.ListProduct(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	fmt.Println(resp, err)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"products": resp.Products})
}

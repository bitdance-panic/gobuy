package handlers

import (
	"context"
	"time"

	rpc_user "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user"
	"github.com/bitdance-panic/gobuy/app/services/user/biz/clients"
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

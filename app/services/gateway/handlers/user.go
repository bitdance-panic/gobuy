package handlers

import (
	"context"
	"time"

	rpc_user "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/clients"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/hertz-contrib/jwt"
)

// 获取用户信息
func HandleGetUser(ctx context.Context, c *app.RequestContext) {
	claims := jwt.ExtractClaims(ctx, c)
	userID := claims["uid"].(int)
	req := rpc_user.GetUserReq{UserId: int32(userID)}
	resp, err := clients.UserClient.GetUser(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	if !resp.Success {
		utils.FailFull(c, consts.StatusInternalServerError, "Get user failed", nil)
	}
	utils.Success(c, utils.H{
		"userID":   resp.UserId,
		"email":    resp.Email,
		"username": resp.Username,
	})
}

// 更新用户信息
func HandleUpdateUser(ctx context.Context, c *app.RequestContext) {
	claims := jwt.ExtractClaims(ctx, c)
	userID := claims["uid"].(int)

	req := rpc_user.UpdateUserReq{UserId: int32(userID)}
	if err := c.BindAndValidate(&req); err != nil {
		hlog.Warnf("User update failed for user id: %s, validation error: %v", req.UserId, err)
		utils.Fail(c, err.Error())
		return
	}
	resp, err := clients.UserClient.UpdateUser(context.Background(), &req)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	if !resp.Success {
		utils.FailFull(c, consts.StatusInternalServerError, "User update failed", nil)
		return
	}
	utils.Success(c, utils.H{"userID": req.UserId})
}

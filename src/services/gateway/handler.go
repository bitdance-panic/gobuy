package main

import (
	"context"
	rpc_user "rpc/kitex_gen/user"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client/callopt"
)

func handlePong(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, utils.H{"message": "pong"})
}

// handleLogin 测试 handler
// @Summary 测试 Summary
// @Description 测试 Description
// @Accept application/json
// @Produce application/json
// @Router /login [get]
func handleLogin(ctx context.Context, c *app.RequestContext) {
	// 通过 /login?email=1234&pass=1234 测试
	email := c.Query("email")
	password := c.Query("pass")
	req := rpc_user.LoginReq{
		Email:    email,
		Password: password,
	}
	resp, err := cli.Login(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{"message": err.Error()})
		return
	}
	if resp.Success {
		c.JSON(consts.StatusOK, utils.H{"userID": resp.UserId})
	} else {
		c.JSON(consts.StatusOK, utils.H{"message": "登录失败"})
	}
}

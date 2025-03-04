package handlers

import (
	"context"
	"time"

	"github.com/bitdance-panic/gobuy/app/consts"
	rpc_agent "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/agent"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/clients"
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
func HandleAskAgent(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt(consts.CONTEXT_UID_KEY)
	var req rpc_agent.AskReq
	if err := c.BindAndValidate(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	req.UserId = int32(userID)
	resp, err := clients.AgentClient.Ask(context.Background(), &req, callopt.WithRPCTimeout(30*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"content": resp.Content})
}

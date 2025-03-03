package main

import (
	"context"

	"github.com/bitdance-panic/gobuy/app/services/agent/biz/bll"

	rpc_agent "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/agent"
)

type AgentServiceImpl struct{}

func (*AgentServiceImpl) Ask(ctx context.Context, req *rpc_agent.AskReq) (resp *rpc_agent.AskResp, err error) {
	return bll.Ask(ctx, req)
}

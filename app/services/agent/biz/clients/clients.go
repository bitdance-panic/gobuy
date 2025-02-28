package clients

import (
	"log"

	"github.com/bitdance-panic/gobuy/app/common/clientsuite"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/agent/agentservice"
	"github.com/bitdance-panic/gobuy/app/services/agent/conf"
	"github.com/cloudwego/kitex/client"
)

func Init() (AgentClient agentservice.Client) {
	var err error

	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: "agent",
			RegistryAddr:       conf.GetConf().Registry.RegistryAddress[0],
		}),
	}
	AgentClient, err = agentservice.NewClient(
		"agent",
		opts...,
	)
	if err != nil {
		log.Fatalf("faild to new agent client: %v", err)
	}

	return AgentClient
}

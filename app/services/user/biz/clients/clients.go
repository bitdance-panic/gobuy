package clients

import (
	"log"

	"github.com/bitdance-panic/gobuy/app/common/clientsuite"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user/userservice"
	"github.com/bitdance-panic/gobuy/app/services/user/conf"
	"github.com/cloudwego/kitex/client"
)

var UserClient userservice.Client

func Init() {
	var err error

	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: conf.GetConf().Kitex.Service,
			RegistryAddr:       conf.GetConf().Registry.RegistryAddress[0],
		}),
	}
	UserClient, err = userservice.NewClient(
		"user",
		opts...,
	)
	if err != nil {
		log.Fatalf("faild to new user client: %v", err)
	}
}

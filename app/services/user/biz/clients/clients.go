package clients

import (
	"log"

	"github.com/bitdance-panic/gobuy/app/common/clientsuite"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product/productservice"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user/userservice"
	pclient "github.com/bitdance-panic/gobuy/app/services/product/biz/clients"
	"github.com/bitdance-panic/gobuy/app/services/user/conf"
	"github.com/cloudwego/kitex/client"
)

var ProductClient productservice.Client

func NewClients() {
	ProductClient = pclient.Init()
}

func Init() (UserClient userservice.Client) {
	var err error

	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: "user",
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

	return UserClient
}

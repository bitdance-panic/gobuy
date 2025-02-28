package clients

import (
	"log"

	"github.com/bitdance-panic/gobuy/app/common/clientsuite"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart/cartservice"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product/productservice"
	"github.com/bitdance-panic/gobuy/app/services/cart/conf"
	pclient "github.com/bitdance-panic/gobuy/app/services/product/biz/clients"
	"github.com/cloudwego/kitex/client"
)

var ProductClient productservice.Client

func NewClients() {
	ProductClient = pclient.Init()
}

func Init() (CartClient cartservice.Client) {
	var err error

	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: "cart",
			RegistryAddr:       conf.GetConf().Registry.RegistryAddress[0],
		}),
	}
	CartClient, err = cartservice.NewClient(
		"cart",
		opts..., // 注入监控配置
	)
	if err != nil {
		log.Fatalf("faild to new cart client: %v", err)
	}

	return CartClient
}

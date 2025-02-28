package clients

import (
	"log"

	"github.com/bitdance-panic/gobuy/app/common/clientsuite"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart/cartservice"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order/orderservice"
	cclient "github.com/bitdance-panic/gobuy/app/services/cart/biz/clients"
	"github.com/bitdance-panic/gobuy/app/services/order/conf"
	"github.com/cloudwego/kitex/client"
)

var CartClient cartservice.Client

func NewClients() {
	CartClient = cclient.Init()
}

func Init() (OrderClient orderservice.Client) {
	var err error

	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: "order",
			RegistryAddr:       conf.GetConf().Registry.RegistryAddress[0],
		}),
	}
	OrderClient, err = orderservice.NewClient(
		"order",
		opts..., // 注入监控配置
	)
	if err != nil {
		log.Fatalf("faild to new order client: %v", err)
	}

	return OrderClient
}

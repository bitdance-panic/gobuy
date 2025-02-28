package clients

import (
	"log"

	"github.com/bitdance-panic/gobuy/app/common/clientsuite"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product/productservice"
	"github.com/bitdance-panic/gobuy/app/services/product/conf"
	"github.com/cloudwego/kitex/client"
)

func Init() (ProductClient productservice.Client) {
	var err error

	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: "product",
			RegistryAddr:       conf.GetConf().Registry.RegistryAddress[0],
		}),
	}
	ProductClient, err = productservice.NewClient(
		"product",
		opts..., // 注入监控配置
	)
	if err != nil {
		log.Fatalf("faild to new product client: %v", err)
	}

	return ProductClient
}

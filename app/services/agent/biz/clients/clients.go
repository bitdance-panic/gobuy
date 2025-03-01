package clients

import (
	"github.com/bitdance-panic/gobuy/app/clients"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order/orderservice"
	"github.com/bitdance-panic/gobuy/app/services/agent/conf"
)

var OrderClient orderservice.Client

func init() {
	registryConf := conf.GetConf().Registry
	OrderClient = clients.NewOrderClient(registryConf.RegistryAddress[0])
}

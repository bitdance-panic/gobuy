package main

import (
	// "common/utils"
	"log"
	"net"
	"strings"

	"github.com/bitdance-panic/gobuy/app/common/mtl"
	"github.com/bitdance-panic/gobuy/app/common/serversuite"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart/cartservice"
	"github.com/bitdance-panic/gobuy/app/services/cart/biz/clients"
	"github.com/bitdance-panic/gobuy/app/services/cart/biz/dal"
	"github.com/bitdance-panic/gobuy/app/services/cart/conf"
	"github.com/bitdance-panic/gobuy/app/utils"

	"github.com/cloudwego/kitex/server"
)

var (
	ServiceName  = "cart" // conf.GetConf().Kitex.Service
	RegistryAddr = conf.GetConf().Registry.RegistryAddress[0]
)

func kitexInit() (opts []server.Option) {
	// address
	address := conf.GetConf().Kitex.Address
	if strings.HasPrefix(address, ":") {
		localIp := utils.MustGetLocalIPv4()
		// localIp := "0.0.0.0"
		address = localIp + address
	}
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr), server.WithSuite(serversuite.CommonServerSuite{
		CurrentServiceName: ServiceName,
		RegistryAddr:       RegistryAddr,
	}))
	return
}

func main() {
	// 初始化指标监控
	mtl.InitMetric(ServiceName, conf.GetConf().Kitex.MetricsPort, RegistryAddr)
	dal.Init()

	clients.NewClients()

	opts := kitexInit()
	svr := cartservice.NewServer(new(CartServiceImpl), opts...)
	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}

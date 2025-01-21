package main

import (
	"log"
	"net"
	"strings"

	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product/productservice"

	// "common/utils"

	"github.com/bitdance-panic/gobuy/app/services/product/biz/dal"
	"github.com/bitdance-panic/gobuy/app/services/product/conf"
	"github.com/cloudwego/kitex/server"
)

func kitexInit() (opts []server.Option) {
	// address
	address := conf.GetConf().Kitex.Address
	if strings.HasPrefix(address, ":") {
		// localIp := utils.MustGetLocalIPv4()
		localIp := "0.0.0.0"
		address = localIp + address
	}
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr))
	// opts = append(opts, server.WithServiceAddr(addr), server.WithSuite(serversuite.CommonServerSuite{CurrentServiceName: serviceName, RegistryAddr: conf.GetConf().Registry.RegistryAddress[0]}))
	return
}

func main() {
	dal.Init()
	opts := kitexInit()
	svr := productservice.NewServer(new(ProductServiceImpl), opts...)
	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}

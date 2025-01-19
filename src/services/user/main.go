package main

import (
	// "common/utils"
	"log"
	"net"
	"rpc/kitex_gen/user/userservice"
	"strings"
	"user/biz/dal"
	"user/conf"

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
	svr := userservice.NewServer(new(UserServiceImpl), opts...)
	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}

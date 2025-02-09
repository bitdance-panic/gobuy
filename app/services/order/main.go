package main

import (
	"log"
	"net"
	"strings"

	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order/orderservice"
	"github.com/bitdance-panic/gobuy/app/services/order/biz/dal"
	"github.com/bitdance-panic/gobuy/app/services/order/conf"

	"github.com/cloudwego/kitex/server"
)

func kitexInit() (opts []server.Option) {
	address := conf.GetConf().Kitex.Address
	if strings.HasPrefix(address, ":") {
		localIp := "0.0.0.0"
		address = localIp + address
	}
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr))
	return
}
func main() {
	dal.Init()
	opts := kitexInit()
	svr := orderservice.NewServer(new(OrderServiceImpl), opts...)
	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
	//创建一个示例订单

	//保存订单

	//获取订单

	//获取订单列表

	//更新订单状态

}

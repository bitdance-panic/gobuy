package main

import (
	"log"
	"net"
	"strings"

	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order/orderservice"
	"github.com/bitdance-panic/gobuy/app/services/order/biz/crontask"
	"github.com/bitdance-panic/gobuy/app/services/order/biz/dal"
	"github.com/bitdance-panic/gobuy/app/services/order/conf"
	"github.com/robfig/cron/v3"

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
	startExpireCron()
}

func startExpireCron() {
	// 创建一个新的cron实例
	c := cron.New()
	// 添加定时任务，每分钟执行一次
	_, err := c.AddFunc("@every 1m", crontask.CheckAndUpdateExpiredOrders)
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}
	// 启动cron
	c.Start()
	// 保持主程序运行
	select {}
}

package main

import (
	"log"
	"net"
	"strings"

	"github.com/cloudwego/kitex/server"

	// 假设此处是你生成的 paymentservice 包路径
	payment "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/payment/paymentservice"

	// 假设你在 payment/biz/dal 包中定义了数据库初始化方法
	"github.com/bitdance-panic/gobuy/app/services/payment/biz/dal"

	// 假设你在 payment/conf 包中定义了 GetConf() 获取配置
	"github.com/bitdance-panic/gobuy/app/services/payment/conf"
)

// kitexInit 负责初始化 Kitex 的 server 选项，例如服务监听地址等
func kitexInit() (opts []server.Option) {
	address := conf.GetConf().Kitex.Address // 从配置中读取监听地址，比如 ":8888"

	// 如果只写了端口号，自动补上 0.0.0.0
	if strings.HasPrefix(address, ":") {
		localIP := "0.0.0.0"
		address = localIP + address // 0.0.0.0:8888
	}

	// 将字符串地址解析为 TCPAddr
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		panic(err)
	}

	// 设置服务监听地址
	opts = append(opts, server.WithServiceAddr(addr))
	return
}

func main() {
	// 1. 初始化数据库等资源
	dal.Init()

	// 2. 初始化 Kitex 服务的监听地址、其他选项
	opts := kitexInit()

	// 3. 创建 PaymentServiceImpl 实例
	svr := payment.NewServer(new(PaymentServiceImpl), opts...)

	// 4. 启动服务
	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}

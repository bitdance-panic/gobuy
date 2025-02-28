package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/bitdance-panic/gobuy/app/services/payment/client"
	"github.com/bitdance-panic/gobuy/app/services/payment/conf"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/smartwalle/alipay/v3"
	"github.com/smartwalle/xid"
)

func configServer() *server.Hertz {
	address := conf.GetConf().Hertz.Address
	s := fmt.Sprintf("localhost%s", address)
	return server.New(server.WithHostPorts(s))
}

var (
	notifyURL   = conf.GetConf().Alipay.ServerDomain + "/alipay/notify"
	callbackURL = conf.GetConf().Alipay.ServerDomain + "/alipay/callback"
)

func handleGetUrl(ctx context.Context, c *app.RequestContext) {
	// 从请求中获取系统订单ID
	orderID := c.Query("order_id")
	// 生成支付宝交易号
	tradeNo := fmt.Sprintf("alipay_%d", xid.Next())
	// TODO 根据ID查询订单详情
	p := alipay.TradePagePay{}
	p.NotifyURL = notifyURL
	p.ReturnURL = callbackURL
	// TODO 和订单一致
	p.Subject = "支付订单:" + "qwerqwer"
	p.OutTradeNo = tradeNo
	// TODO 和订单一致
	p.TotalAmount = "20.00"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	// 使用 PassbackParams 传递系统订单ID
	// 需要进行URL编码，因为这个参数会通过URL传递
	p.PassbackParams = url.QueryEscape(orderID)
	url, err := client.Client.TradePagePay(p)
	if err != nil {
		utils.Fail(c, "创建支付链接失败")
		return
	}
	payURL := url.String()
	log.Printf("创建支付链接成功, 系统订单ID: %s, 支付宝交易号: %s", orderID, tradeNo)
	fmt.Println(payURL)
	// 重定向到支付页面
	utils.Success(c, utils.H{"payURL": payURL})
}

func handleCallback(ctx context.Context, c *app.RequestContext) {
	type Req struct {
		OutTradeNo     string `form:"out_trade_no"`
		PassbackParams string `form:"passback_params"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		log.Println("参数错误", err.Error())
		return
	}
	// 获取系统订单ID (PassbackParams)
	systemOrderID, err := url.QueryUnescape(req.PassbackParams)
	if err != nil {
		log.Println("解析系统订单ID发生错误", err.Error())
		systemOrderID = ""
		return
	}
	// 验证签名
	// if err := client.Client.VerifySign(c.Request.Form); err != nil {
	// 	log.Println("回调验证签名发生错误", err)
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "回调验证签名发生错误"})
	// 	return
	// }

	log.Println("回调验证签名通过")

	// 查询订单状态
	p := alipay.TradeQuery{}
	p.OutTradeNo = req.OutTradeNo

	rsp, err := client.Client.TradeQuery(ctx, p)
	if err != nil {
		log.Printf("验证订单 %s 信息发生错误: %s", req.OutTradeNo, err.Error())
		return
	}
	if rsp.IsFailure() {
		log.Printf("验证订单 %s 信息发生错误: %s-%s", req.OutTradeNo, rsp.Msg, rsp.SubMsg)
		return
	}
	// 在这里处理系统订单状态更新
	if systemOrderID != "" {
		log.Printf("系统订单 %s 支付成功，支付宝交易号: %s", systemOrderID, req.OutTradeNo)
		// TODO: 在这里更新你系统中的订单状态
		// updateOrderStatus(systemOrderID, "PAID")
	}
}

func main() {
	client.NewClients()

	h := configServer()
	{
		// http://127.0.0.1:8280/alipay/pay?order_id=12
		h.GET("/alipay/pay", handleGetUrl)
		h.GET("/alipay/callback", handleCallback)
		// h.POST("/alipay/notify", handleGetUrl)
	}

	h.Spin()
}

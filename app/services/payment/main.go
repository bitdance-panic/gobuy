package main

import (
	"github.com/bitdance-panic/gobuy/app/services/payment/conf"
	"github.com/gin-gonic/gin"
)

var (
	notifyURL   = conf.GetConf().Alipay.ServerDomain + "/alipay/notify"
	callbackURL = conf.GetConf().Alipay.ServerDomain + "/alipay/callback"
)

func main() {
	//路由函数
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	// http://127.0.0.1:8280/alipay/pay?order_id=150003
	r.GET("/alipay/pay", handlePayURL)
	r.GET("/alipay/callback", handleCallback)
	r.POST("/alipay/notify", handleNotify)
	hertzConf := conf.GetConf().Hertz
	r.Run(hertzConf.Address)
}

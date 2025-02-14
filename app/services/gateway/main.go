package main

import (
	"fmt"
	"log"

	// 引入 product 和 user 服务的客户端
	productservice_ "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product/productservice"
	userservice_ "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user/userservice"

	// 引入 payment 服务的客户端
	paymentservice_ "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/payment/paymentservice"

	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal"
	"github.com/bitdance-panic/gobuy/app/services/gateway/conf"
	_ "github.com/bitdance-panic/gobuy/app/services/gateway/docs"
	"github.com/bitdance-panic/gobuy/app/services/gateway/middleware"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
)

var (
	userservice    userservice_.Client
	productservice productservice_.Client
	// 定义 paymentservice 客户端
	paymentservice paymentservice_.Client
)

// @title userservice
// @version 1.0
// @description API Doc for user service.

// @contact.name hertz-contrib
// @contact.url https://github.com/hertz-contrib

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8888
// @BasePath /
// @schemes http
func main() {
	// 初始化数据库等
	dal.Init()
	address := conf.GetConf().Hertz.Address
	s := fmt.Sprintf("localhost%s", address)
	h := server.New(server.WithHostPorts(s))
	h.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有来源
		AllowMethods:     []string{"*"}, // 允许所有方法
		AllowHeaders:     []string{"*"}, // 允许所有头信息
		ExposeHeaders:    []string{"*"}, // 暴露所有头信息
		AllowCredentials: true,          // 允许携带凭证（如 cookies）
	}))

	// 初始化 userservice 客户端
	c, err := userservice_.NewClient("user", client.WithHostPorts("0.0.0.0:8881"))
	if err != nil {
		log.Fatal(err)
	}
	userservice = c
	middleware.UserClient = userservice

	// 初始化 productservice 客户端
	cp, errp := productservice_.NewClient("product", client.WithHostPorts("0.0.0.0:8882"))
	if errp != nil {
		log.Fatal(err)
	}
	productservice = cp

	// 初始化 paymentservice 客户端
	cpmt, errpmt := paymentservice_.NewClient("payment", client.WithHostPorts("0.0.0.0:8883"))
	if errpmt != nil {
		log.Fatal(errpmt)
	}
	paymentservice = cpmt

	// 初始化中间件
	middleware.InitAuth()
	registerRoutes(h)

	// 启动 Swagger 文档服务
	url := swagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", s))
	h.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, url))
	h.Spin()
}

func registerRoutes(h *server.Hertz) {
	// 用户相关路由
	h.POST("/login", middleware.AuthMiddleware.LoginHandler) // 用户登录

	adminGroup := h.Group("/auth")
	{
		adminGroup.POST("/refresh", RefreshTokenHandler) // 令牌刷新
	}

	product := h.Group("/product")
	{
		product.GET("/search", handleProductSearch)
		product.GET("/:id", handleProductGet)
		// product.DELETE("/:id", handleProductDELETE)
		product.PUT("/:id", handleProductPut)
		product.POST("", handleProductPost)
	}

	user := h.Group("/user")
	{
		user.PUT("/:userid", UpdateUserHandler)
		user.GET("/:userid", GetUserHandler)
		user.POST("/:userid", DeleteUserHandler)
	}

	// 处理与支付相关的路由
	payment := h.Group("/payment")
	{
		payment.POST("/create", handleCreatePayment)
		payment.GET("/:paymentId", handleGetPayment)
		payment.PUT("/:paymentId", handleUpdatePayment)
		payment.DELETE("/:paymentId", handleDeletePayment)
	}
}

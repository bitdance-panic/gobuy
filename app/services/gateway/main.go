package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal/redis"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/gateway/casbin"
	"github.com/bitdance-panic/gobuy/app/services/gateway/conf"
	_ "github.com/bitdance-panic/gobuy/app/services/gateway/docs"
	"github.com/bitdance-panic/gobuy/app/services/gateway/handlers"
	"github.com/bitdance-panic/gobuy/app/services/gateway/middleware"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/hertz-contrib/registry/consul"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"

	consulapi "github.com/hashicorp/consul/api"
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

var (
	address string
	h       *server.Hertz
)

func registerToConsul() {
	// build a consul client
	config := consulapi.DefaultConfig()
	config.Address = conf.GetConf().Registry.RegistryAddress[0] // "localhost:8500"
	consulclient, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalf("failed to build a consul client: %v", err)
		return
	}
	// build a consul register with the consul client
	r := consul.NewConsulRegister(consulclient)

	// 解析服务的 IP 和端口
	address = conf.GetConf().Hertz.Address
	if strings.HasPrefix(address, ":") {
		localIp := utils.MustGetLocalIPv4()
		address = localIp + address
	}
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		panic(err)
	}

	// // 生成唯一服务 ID（例如使用 IP:Port）
	// parts := strings.Split(address, ":")
	// if len(parts) != 2 {
	// 	log.Fatalf("地址格式错误: %s", address)
	// }
	// ip := parts[0]
	// port, _ := strconv.Atoi(parts[1])
	// serviceID := fmt.Sprintf("gateway-%s-%d", ip, port)

	// run Hertz with the consul register
	h = server.Default(
		server.WithHostPorts(address),
		server.WithRegistry(r, &registry.Info{
			ServiceName: "gateway",
			Addr:        addr,
			Weight:      10,
			Tags:        nil,
		}),
	)
}

func main() {
	// 初始化数据库
	dal.Init()

	// 不要每次都初始化Casbin
	if err := casbin.InitCasbin(tidb.DB); err != nil {
		hlog.Fatalf("Casbin初始化失败: %v", err)
	}
	// dao.AddUserRole(tidb.DB, 540001, 1)

	// 同步黑名单到Redis
	redis.SyncBlacklistToRedis()

	// 启动自动清理任务
	middleware.StartRedisCleanupTask()

	// 创建Hertz实例
	registerToConsul()

	// 中间件链
	h.Use(
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"}, // 允许所有来源
			AllowMethods:     []string{"*"}, // 允许所有方法
			AllowHeaders:     []string{"*"}, // 允许所有头信息
			ExposeHeaders:    []string{"*"}, // 暴露所有头信息
			AllowCredentials: true,          // 允许携带凭证（如 cookies）
		}),
		// 白名单放行接口
		middleware.WhiteListMiddleware(),
		middleware.ConditionalAuthMiddleware(),
		middleware.AddUidMiddleware(),
		// 黑名单检查
		middleware.BlacklistMiddleware(),
		// 用户权限检查
		middleware.CasbinMiddleware(),
	)
	// 注册路由
	registerRoutes(h)
	// 注册Swagger
	registerSwagger(h, address)
	h.Spin()
}

// 作为占位
func TODOHandler(ctx context.Context, c *app.RequestContext) {}

func registerRoutes(h *server.Hertz) {
	noAuthGroup := h.Group("")
	{
		// TODO 登陆
		noAuthGroup.POST("/login", middleware.AuthMiddleware.LoginHandler)
		// 注册
		noAuthGroup.POST("/register", handlers.HandleRegister)
		// 未登录时获取首页商品
		noAuthGroup.GET("/product/search", handlers.HandleSearchProducts)
		// 让前端移除token就行，这里废弃
		// noAuthGroup.POST("/logout", middleware.AuthMiddleware.LogoutHandler)
	}
	authGroup := h.Group("/auth")
	{
		// TODO 刷新token
		authGroup.POST("/refresh", middleware.AuthMiddleware.RefreshHandler)
	}
	userGroup := h.Group("/user")
	{
		// 自己获取自己信息
		userGroup.GET("", handlers.HandleGetUser)
		// 更新个人信息
		userGroup.PUT("", handlers.HandleUpdateUser)
	}
	productGroup := h.Group("/product")
	{
		// 获取单个商品详情
		productGroup.GET("/:id", handlers.HandleGetProduct)
	}
	cartGroup := h.Group("/cart")
	{
		// 获取用户购物车
		cartGroup.GET("", handlers.HandleListCartItem)
		// 将商品放入购物车
		cartGroup.POST("/:productID", handlers.HandleCreateCartItem)
		// 从购物车移除单个商品
		cartGroup.DELETE("/:itemID", handlers.HandleDeleteCartItem)
		cartGroup.PUT("/:itemID", handlers.HandleUpdateCartItemQuantity)
	}
	orderGroup := h.Group("/order")
	{
		orderGroup.POST("", handlers.HandleCreateOrder)
		orderGroup.GET("/:id", handlers.HandleGetOrder)
		orderGroup.GET("/user", handlers.HandleListUserOrder)
	}
	paymentGroup := h.Group("/payment")
	{
		// TODO 只需要处理支付操作，应该是个回调的接口
		paymentGroup.POST("/:orderID", TODOHandler)
	}
	agentGroup := h.Group("/agent")
	{
		agentGroup.POST("/ask", handlers.HandleAskAgent)
	}
	adminGroup := h.Group("/admin")
	{
		adminProductGroup := adminGroup.Group("/product")
		{
			// 创建商品
			adminProductGroup.POST("", handlers.HandleCreateProduct)
			// 更新商品
			adminProductGroup.PUT("/:id", handlers.HandleUpdateProduct)
			// 移除商品
			adminProductGroup.DELETE("/:id", handlers.HandleRemoveProduct)
			// 获取所有商品
			adminProductGroup.GET("/list", handlers.HandleAdminListProduct)
		}
		adminUserGroup := adminGroup.Group("/user")
		{
			// 获取所有的用户信息
			adminUserGroup.GET("/list", handlers.HandleAdminListUser)
			// TODO 封禁用户
			adminUserGroup.POST("/block", handlers.HandleBlockUser)
			// TODO 解封
			adminUserGroup.DELETE("/unblock/:identifier", handlers.HandleUnblockUser)
			// 移除用户
			adminUserGroup.DELETE("/:userID", handlers.HandleRemoveUser)
		}
		adminOrderGroup := adminGroup.Group("/order")
		{
			// TODO 获取所有的订单(分页)（订单包括支付信息）
			adminOrderGroup.GET("/list", handlers.HandleAdminListOrder)
		}
	}
}

func registerSwagger(h *server.Hertz, addr string) {
	url := swagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", addr))
	h.GET("/swagger/*any",
		swagger.WrapHandler(swaggerFiles.Handler,
			swagger.DefaultModelsExpandDepth(-1), // 隐藏模型定义
			url,
		),
	)
}

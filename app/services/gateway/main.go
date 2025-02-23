package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	userservice_ "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user/userservice"

	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal/redis"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/gateway/casbin"
	"github.com/bitdance-panic/gobuy/app/services/gateway/conf"
	_ "github.com/bitdance-panic/gobuy/app/services/gateway/docs"
	"github.com/bitdance-panic/gobuy/app/services/gateway/handlers"
	"github.com/bitdance-panic/gobuy/app/services/gateway/middleware"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
)

var (
	userservice userservice_.Client
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
	// 初始化数据库
	dal.Init()

	// 初始化Casbin
	if err := casbin.InitCasbin(tidb.DB); err != nil {
		hlog.Fatalf("Casbin初始化失败: %v", err)
	}
	// dao.AddUserRole(tidb.DB, 540001, 1)

	// 同步黑名单到Redis
	redis.SyncBlacklistToRedis()

	// 启动自动清理任务
	middleware.StartRedisCleanupTask()

	// 初始化黑名单
	// middleware.InitBlacklistCache()
	// middleware.StartBlacklistCleanupTask()

	// 创建Hertz实例
	address := conf.GetConf().Hertz.Address
	s := fmt.Sprintf("localhost%s", address)
	h := server.New(server.WithHostPorts(s))

	// 中间件链
	h.Use(
		// 黑名单检查
		middleware.BlacklistMiddleware(),
		// 白名单放行接口
		middleware.WhiteListMiddleware(),
		conditionalAuthMiddleware(),
		// 用户权限检查
		middleware.CasbinMiddleware(),
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"}, // 允许所有来源
			AllowMethods:     []string{"*"}, // 允许所有方法
			AllowHeaders:     []string{"*"}, // 允许所有头信息
			ExposeHeaders:    []string{"*"}, // 暴露所有头信息
			AllowCredentials: true,          // 允许携带凭证（如 cookies）
		}),
	)
	// 注册路由
	registerRoutes(h)
	// 注册Swagger
	registerSwagger(h, s)
	h.Spin()
}

// 作为占位
func TODOHandler(ctx context.Context, c *app.RequestContext) {}

func registerRoutes(h *server.Hertz) {
	noAuthGroup := h.Group("")
	{
		// 登陆
		noAuthGroup.POST("/login", middleware.AuthMiddleware.LoginHandler)
		// 注册
		noAuthGroup.POST("/register", handlers.HandleRegister)
		// 获取首页商品
		noAuthGroup.GET("/index/products", handlers.HandleListIndexProduct)
		// 让前端移除token就行，这里废弃
		// noAuthGroup.POST("/logout", middleware.AuthMiddleware.LogoutHandler)
	}
	authGroup := h.Group("/auth")
	{
		//TODO
		authGroup.POST("/refresh", middleware.AuthMiddleware.RefreshHandler)
	}
	userGroup := h.Group("/user")
	{
		userGroup.GET("", handlers.HandleGetUser)
		userGroup.PUT("", handlers.HandleUpdateUser)
	}
	productGroup := h.Group("/product")
	{
		productGroup.GET("/search", handlers.HandleSearchProducts)
		//获取单个商品详情
		productGroup.GET("/:id", handlers.HandleGetProduct)
	}
	cartGroup := h.Group("/cart")
	{
		cartGroup.GET("", handlers.HandleListCartItem)
		cartGroup.POST("/:productID", handlers.HandleCreateCartItem)
		cartGroup.DELETE("/:itemID", handlers.HandleDeleteCartItem)
		cartGroup.PUT("/:itemID", handlers.HandleUpdateCartItemQuantity)
	}
	orderGroup := h.Group("/order")
	{
		// 创建订单
		orderGroup.POST("", handlers.HandleCreateOrder)
		// 获取单个订单详情
		orderGroup.GET("/:id", handlers.HandleGetOrder)
		// 获取用户的所有订单
		orderGroup.GET("/user", handlers.HandleListUserOrder)
	}
	paymentGroup := h.Group("/payment")
	{
		//TODO 只需要处理支付操作，应该是个回调的接口
		paymentGroup.POST("/:orderID", TODOHandler)
	}
	agentGroup := h.Group("/agent")
	{
		//TODO 根据用户输入获取对应商品
		agentGroup.POST("/product/search", TODOHandler)
		//TODO 根据用户输入获取对应订单
		agentGroup.POST("/order/search", TODOHandler)
	}
	adminGroup := h.Group("/admin")
	{
		adminProductGroup := adminGroup.Group("/product")
		{
			adminProductGroup.POST("", handlers.HandleCreateProduct)
			adminProductGroup.PUT("/:id", handlers.HandleUpdateProduct)
			adminProductGroup.DELETE("/:id", handlers.HandleRemoveProduct)
			adminProductGroup.GET("", handlers.HandleAdminListProduct)
		}
		adminUserGroup := adminGroup.Group("/users")
		{
			// 获取所有的用户信息
			adminUserGroup.GET("", handlers.HandleAdminListUser)
			adminUserGroup.GET("/block/:userID", handlers.HandleBlockUser)
			adminUserGroup.GET("/unblock/:userID", handlers.HandleUnblockUser)
			adminUserGroup.DELETE("/:userID", handlers.HandleRemoveUser)
		}
		adminOrderGroup := adminGroup.Group("/orders")
		{
			// 获取所有的订单(分页)（订单包括支付信息）
			adminOrderGroup.GET("", handlers.HandleAdminListOrder)
		}
	}
}

func conditionalAuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if skip, exists := c.Get("skip_auth"); exists && skip.(bool) {
			c.Next(ctx) // 跳过认证
			return
		}
		middleware.AuthMiddleware.MiddlewareFunc()(ctx, c) // 执行认证
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

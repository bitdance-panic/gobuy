package main

import (
	"context"
	"fmt"

	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal/redis"
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

func addUidMiddleware() app.HandlerFunc {
	// claims := jwt.ExtractClaims(ctx, c)
	// userID := claims["uid"].(int)
	return func(ctx context.Context, c *app.RequestContext) {
		fmt.Println("设置UID")
		c.Set("uid", 450002)
		c.Next(ctx)
	}
}

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

	// 不要每次都初始化Casbin
	// if err := casbin.InitCasbin(tidb.DB); err != nil {
	// 	hlog.Fatalf("Casbin初始化失败: %v", err)
	// }
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
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"}, // 允许所有来源
			AllowMethods:     []string{"*"}, // 允许所有方法
			AllowHeaders:     []string{"*"}, // 允许所有头信息
			ExposeHeaders:    []string{"*"}, // 暴露所有头信息
			AllowCredentials: true,          // 允许携带凭证（如 cookies）
		}),
		// 黑名单检查
		// middleware.BlacklistMiddleware(),
		// 白名单放行接口
		middleware.WhiteListMiddleware(),
		// conditionalAuthMiddleware(),
		addUidMiddleware(),
		// 用户权限检查
		// middleware.CasbinMiddleware(),
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
		// TODO 根据用户输入获取对应商品
		agentGroup.POST("/product/search", TODOHandler)
		// TODO 根据用户输入获取对应订单
		agentGroup.POST("/order/search", TODOHandler)
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
			adminUserGroup.GET("/block/:userID", handlers.HandleBlockUser)
			// TODO 解封
			adminUserGroup.GET("/unblock/:userID", handlers.HandleUnblockUser)
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

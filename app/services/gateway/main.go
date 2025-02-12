package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	productservice_ "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product/productservice"
	userservice_ "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user/userservice"

	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal"
	"github.com/bitdance-panic/gobuy/app/services/gateway/conf"
	_ "github.com/bitdance-panic/gobuy/app/services/gateway/docs"
	"github.com/bitdance-panic/gobuy/app/services/gateway/middleware"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
)

var (
	userservice    userservice_.Client
	productservice productservice_.Client
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
	dal.Init()

	// 初始化Casbin
	// if err := casbin.InitCasbin(tidb.DB); err != nil {
	// 	hlog.Fatalf("Casbin初始化失败: %v", err)
	// }

	// 创建Hertz实例
	address := conf.GetConf().Hertz.Address
	s := fmt.Sprintf("localhost%s", address)
	h := server.New(server.WithHostPorts(s))

	c, err := userservice_.NewClient("user", client.WithHostPorts("0.0.0.0:8881"))
	if err != nil {
		hlog.Fatal(err)
	}
	userservice = c
	middleware.UserClient = userservice
	cp, errp := productservice_.NewClient("product", client.WithHostPorts("0.0.0.0:8882"))
	if errp != nil {
		hlog.Fatal(err)
	}
	productservice = cp

	// 初始化中间件
	middleware.InitAuth()

	// 中间件链
	h.Use(
		// middleware.CasbinMiddleware(), // 权限中间件
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

func registerRoutes(h *server.Hertz) {
	// 用户路由
	user := h.Group("/")
	user.Use(
		middleware.WhiteListMiddleware(),
		conditionalAuthMiddleware())
	{
		user.POST("/login", middleware.AuthMiddleware.LoginHandler)
	}

	// 需要认证的路由
	adminGroup := h.Group("/auth")
	adminGroup.Use(middleware.RBACMiddleware("admin"))
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

	// // 受保护的业务 API
	// authGroup := h.Group("/api")
	// authGroup.Use(middleware.AuthMiddleware.MiddlewareFunc()) // JWT 认证
	// {
	// 	authGroup.GET("/profile", handler.ProfileHandler) // 获取用户信息
	// 	authGroup.POST("/update", handler.UpdateProfile)  // 更新用户信息
	// }
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

func conditionalAuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if skip, exists := c.Get("skip_auth"); exists && skip.(bool) {
			c.Next(ctx) // 跳过认证
			return
		}
		middleware.AuthMiddleware.MiddlewareFunc()(ctx, c) // 执行认证
	}
}

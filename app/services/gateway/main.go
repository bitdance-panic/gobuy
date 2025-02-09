package main

import (
	"fmt"
	"log"

	productservice_ "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product/productservice"
	userservice_ "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user/userservice"

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
	// middleware.InitCasbin()
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

	c, err := userservice_.NewClient("user", client.WithHostPorts("0.0.0.0:8881"))
	if err != nil {
		log.Fatal(err)
	}
	userservice = c
	middleware.UserClient = userservice
	cp, errp := productservice_.NewClient("product", client.WithHostPorts("0.0.0.0:8882"))
	if errp != nil {
		log.Fatal(err)
	}
	productservice = cp

	// 初始化中间件
	middleware.InitAuth()
	registerRoutes(h)
	url := swagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", s))
	h.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, url))
	h.Spin()
}

func registerRoutes(h *server.Hertz) {
	// h.GET("/ping", handlePong)
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
	// // 受保护的业务 API
	// authGroup := h.Group("/api")
	// authGroup.Use(middleware.AuthMiddleware.MiddlewareFunc()) // JWT 认证
	// {
	// 	authGroup.GET("/profile", handler.ProfileHandler) // 获取用户信息
	// 	authGroup.POST("/update", handler.UpdateProfile)  // 更新用户信息
	// }

	// // 角色权限管理（需要 RBAC 控制）
	// adminGroup := h.Group("/admin")
	// adminGroup.Use(middleware.AuthMiddleware.MiddlewareFunc(), middleware.CasbinMiddleware())
	// {
	//  adminGroup.POST("/refresh", RefreshTokenHandler) // 令牌刷新
	// 	adminGroup.POST("/users", handler.CreateUser)       // 创建用户
	// 	adminGroup.PUT("/users/:id", handler.UpdateUser)    // 更新用户信息
	// 	adminGroup.DELETE("/users/:id", handler.DeleteUser) // 删除用户
	// }

	// 健康检查
	// h.GET("/health", handler.HealthCheck)
}

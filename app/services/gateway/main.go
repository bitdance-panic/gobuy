package main

import (
	"fmt"
	"log"

	productservice_ "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product/productservice"
	userservice_ "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user/userservice"

	"github.com/bitdance-panic/gobuy/app/services/gateway/conf"
	_ "github.com/bitdance-panic/gobuy/app/services/gateway/docs"

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
	// middleware.RegsterAuth(h)
	c, err := userservice_.NewClient("user", client.WithHostPorts("0.0.0.0:8881"))
	if err != nil {
		log.Fatal(err)
	}
	userservice = c
	cp, errp := productservice_.NewClient("product", client.WithHostPorts("0.0.0.0:8882"))
	if errp != nil {
		log.Fatal(err)
	}
	productservice = cp

	h.GET("/ping", handlePong)
	h.GET("/login", handleLogin)
	product := h.Group("/product")
	{
		product.GET("/search", handleProductSearch)
		product.GET("/:id", handleProductGet)
		product.DELETE("/:id", handleProductDELETE)
		product.PUT("/:id", handleProductPut)
		product.POST("", handleProductPost)
	}
	// The url pointing to swagger API definition
	url := swagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", s))
	h.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, url))
	h.Spin()
}

package main

import (
	"fmt"
	"log"

	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user/userservice"

	_ "github.com/bitdance-panic/gobuy/app/services/gateway/docs"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
)

var (
	cli userservice.Client
)

// @title HertzTest
// @version 1.0
// @description This is a demo using Hertz.

// @contact.name hertz-contrib
// @contact.url https://github.com/hertz-contrib

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8888
// @BasePath /
// @schemes http
func main() {
	s := "0.0.0.0:8888"
	h := server.New(server.WithHostPorts(s))
	h.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有来源
		AllowMethods:     []string{"*"}, // 允许所有方法
		AllowHeaders:     []string{"*"}, // 允许所有头信息
		ExposeHeaders:    []string{"*"}, // 暴露所有头信息
		AllowCredentials: true,          // 允许携带凭证（如 cookies）
	}))
	// middleware.RegsterAuth(h)
	c, err := userservice.NewClient("user", client.WithHostPorts("0.0.0.0:9999"))
	if err != nil {
		log.Fatal(err)
	}
	cli = c
	h.GET("/ping", handlePong)
	h.GET("/login", handleLogin)
	// The url pointing to swagger API definition
	url := swagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", s))
	h.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, url))
	h.Spin()
}

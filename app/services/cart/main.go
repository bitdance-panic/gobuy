package main

import (
	"app/services/cart/biz/bll"
	"app/services/cart/biz/dal"
	"app/services/user/biz/dal/postgres"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	dal.Init()
	cartDAL := dal.NewCartDAL(postgres.DB)
	cartBLL = bll.NewCartBLL(cartDAL)

	s := "0.0.0.0:8889"
	h := server.New(server.WithHostPorts(s))

	h.GET("/cart", GetCartHandler)
	h.POST("/cart/add", AddItemToCartHandler)

	h.Spin()
}

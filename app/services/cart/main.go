package main

import (
	"github.com/bitdance-panic/gobuy/app/services/cart/biz/bll"
	"github.com/bitdance-panic/gobuy/app/services/cart/biz/dal"
	//"github.com/bitdance-panic/gobuy/hertz/pkg/app/server"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/cloudwego/kitex/server"
)

// 初始化 TiDB 连接
func initTiDB() (*gorm.DB, error) {
	//TiDB 的连接信息

dsn:
	"%s:%s@tcp(%s:%s)/gobuy?tls=tidb&charset=utf8mb4"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	// 初始化 TiDB 连接
	tidbDB, err := initTiDB()
	if err != nil {
		panic("Failed to connect to TiDB: " + err.Error())
	}

	cartDAL := dal.NewCartDAL(tidbDB)
	cartBLL := bll.NewCartBLL(cartDAL)

	s := "0.0.0.0:8889"
	h := server.New(server.WithHostPorts(s))

	h.GET("/cart", GetCartHandler)
	h.POST("/cart/add", AddItemToCartHandler)

	h.Spin()
}

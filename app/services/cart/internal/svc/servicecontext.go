package svc

import (
	"github.com/bitdance-panic/gobuy/app/services/cart/proto/cart"
	"github.com/cloudwego/thriftgo/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type ServiceContext struct {
	Config   config.Config
	DB       *gorm.DB
	CartRepo cart.CartRepo
}

func NewServiceContext(c config.Config) *ServiceContext {
	ctx := &ServiceContext{Config: c}
	dsn := "DJnAaVe8x5ioexY.root:3NeJC53vJdFYSqz8@tcp(gateway01.ap-southeast-1.prod.aws.tidbcloud.com:4000)/@gateway01.ap-southeast-1.prod.aws.tidbcloud.com?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return ctx
	}
	sqlDB, err := db.DB()
	if err != nil {
	} else {
		sqlDB.SetMaxOpenConns(16)
		sqlDB.SetMaxIdleConns(8)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	ctx.DB = db
	ctx.CartRepo = cart.NewCartRepo(db)

	return ctx
}

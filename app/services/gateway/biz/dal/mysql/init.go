package mysql

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/bitdance-panic/gobuy/app/services/gateway/conf"
	TLS "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

// Init 初始化TiDB连接
func Init() {
	cfg := conf.GetConf().Tidb
	// 配置TLS
	TLS.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: cfg.Host,
	})

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=tidb&charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		"gobuy", // 数据库名
	)

	var initErr error
	dbOnce.Do(func() {
		instance, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			initErr = fmt.Errorf("failed to connect to TiDB: %v", err)
			return
		}

		sqlDB, _ := instance.DB()
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			initErr = fmt.Errorf("TiDB heartbeat check failed: %v", err)
			return
		}

		db = instance
	})

	if initErr != nil {
		panic(initErr)
	}
}

// DB 获取数据库实例
func DB() *gorm.DB {
	return db
}

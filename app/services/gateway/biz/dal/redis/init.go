package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bitdance-panic/gobuy/app/services/gateway/conf"
	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once   sync.Once
)

// Init 初始化Redis连接
func Init() {
	cfg := conf.GetConf().Redis

	var initErr error
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     cfg.Address,
			Password: cfg.Password,
			DB:       cfg.DB,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := client.Ping(ctx).Result(); err != nil {
			initErr = fmt.Errorf("failed to connect to Redis: %v", err)
		}
	})

	if initErr != nil {
		panic(initErr)
	}
}

// Client 获取Redis客户端
func Client() *redis.Client {
	return client
}

package redis

import (
	"context"
	"fmt"
	"time"
	"web/setting"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

func Init(cfg *setting.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	_, err = rdb.Ping().Result()
	return
}

func Close() {
	_ = rdb.Close()
}

// 设置key-value
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rdb.Set(key, value, expiration).Err()
}

// 获取key-value
func Get(ctx context.Context, key string) (string, error) {
	return rdb.Get(key).Result()
}

// 删除key-value
func Del(ctx context.Context, key string) error {
	return rdb.Del(key).Err()
}

// 判断key是否存在
func Exist(ctx context.Context, key string) bool {
	n, _ := rdb.Exists(key).Result()
	return n > 0
}

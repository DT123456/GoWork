package core

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis() {
	addr := viper.GetString("redis.addr")
	password := viper.GetString("redis.password")

	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	// 测试连接
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		ZapRedis.Fatal("Redis连接失败",
			zap.String("addr", addr),
			zap.Error(err),
		)
	}

	ZapRedis.Info("Redis连接成功",
		zap.String("addr", addr),
	)
}

// Close 关闭 Redis 连接
func CloseRedis() error {
	if RDB != nil {
		return RDB.Close()
	}
	return nil
}

// RedisLog 记录 Redis 操作日志（可选，用于调试）
func RedisLog(operation, key string, err error) {
	if err != nil {
		ZapRedis.Warn(fmt.Sprintf("Redis %s 失败", operation),
			zap.String("key", key),
			zap.Error(err),
		)
	}
}

package configs

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// NewRedisClient 创建一个新的Redis客户端
var RedisClient *redis.Client

func NewRedisClient(config *RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// Test the connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return client, nil
}

// 初始化Redis配置
func ConnectRedis() error {
	redisConfig := &RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "", // 如果有密码，请设置
		DB:       0,
	}

	client, err := NewRedisClient(redisConfig)
	if err != nil {
		return fmt.Errorf("redis连接失败: %v", err)
	}

	RedisClient = client
	fmt.Println("Redis连接成功")
	return nil
}

func CheckRedisConnection() error {
	if RedisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis connection error: %v", err)
	}
	return nil
}

// 设置缓存
func SetCache(key string, value interface{}, expiration time.Duration) error {
	// 检查连接
	if err := CheckRedisConnection(); err != nil {
		return fmt.Errorf("redis connection check failed: %v", err)
	}

	ctx := context.Background()
	err := RedisClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %v", err)
	}

	// 验证缓存是否设置成功
	_, err = RedisClient.Get(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("cache verification failed: %v", err)
	}

	return nil
}

// 获取缓存
func GetCache(key string) (string, error) {
	ctx := context.Background()
	return RedisClient.Get(ctx, key).Result()
}

// 删除缓存
func DeleteCache(key string) error {
	ctx := context.Background()
	return RedisClient.Del(ctx, key).Err()
}

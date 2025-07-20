package utils

import (
	"context"
	"fmt"
	"go-project/configs"
	"time"
)

// CacheKeys 用于管理资源的缓存键
type CacheKeys struct {
	resource   string
	ListPrefix string
	TotalKey   string
}

const (
	DefaultCacheTime = 360 * time.Hour // 默认缓存时间为360小时
)

// NewCacheKeys 创建新的缓存键管理器
func NewCacheKeys(resource string) CacheKeys {
	return CacheKeys{
		resource:   resource,
		ListPrefix: resource + ":list",
		TotalKey:   resource + ":total",
	}
}

// 生成列表缓存key
func GenListCacheKey(prefix string, limit, offset int) string {
	return fmt.Sprintf("%s:limit_%d:offset_%d", prefix, limit, offset)
}

// 生成总数缓存key
func GenTotalCacheKey(prefix string) string {
	return fmt.Sprintf("%s:total", prefix)
}

// 生成分页缓存key
func GenPageCacheKey(business, resource string, limit, offset int) string {
	return fmt.Sprintf("%s:%s:page:limit_%d:offset_%d", business, resource, limit, offset)
}

// 获取详情缓存键
func (ck CacheKeys) GetDetailKey(id interface{}) string {
	return fmt.Sprintf("%s:detail:%v", ck.resource, id)
}

// ClearListCache 通用的列表缓存清理函数
func ClearListCache(Cache CacheKeys) {
	ctx := context.Background()

	// 清除总数缓存
	configs.RedisClient.Del(ctx, Cache.TotalKey)

	// 清除列表缓存（使用通配符）
	pattern := fmt.Sprintf("%s*", Cache.ListPrefix)
	iter := configs.RedisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		configs.RedisClient.Del(ctx, iter.Val())
	}
}

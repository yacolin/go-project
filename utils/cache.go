package utils

import (
	"context"
	"fmt"
	"time"

	"go-project/configs"
)

const (
	AlbumListCachePrefix = "album:list"
	AlbumTotalCacheKey   = "album:total"
	TeamListCachePrefix  = "team:list"
	TeamTotalCacheKey    = "team:total"
	DefaultCacheTime     = 5 * time.Minute
)

// 生成列表缓存key
func GenerateListCacheKey(prefix string, limit, offset int) string {
	return fmt.Sprintf("%s:limit_%d:offset_%d", prefix, limit, offset)
}

// 生成总数缓存key
func GenerateTotalCacheKey(prefix string) string {
	return fmt.Sprintf("%s:total", prefix)
}

// 生成分页缓存key
func GeneratePageCacheKey(business, resource string, limit, offset int) string {
	return fmt.Sprintf("%s:%s:page:limit_%d:offset_%d", business, resource, limit, offset)
}

// 清理专辑相关的缓存
func ClearAlbumCache() {
	ctx := context.Background()

	// 清除总数缓存
	configs.RedisClient.Del(ctx, AlbumTotalCacheKey)

	// 清除列表缓存（使用通配符）
	pattern := fmt.Sprintf("%s*", AlbumListCachePrefix)
	iter := configs.RedisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		configs.RedisClient.Del(ctx, iter.Val())
	}
}

// 清理团队相关的缓存
func ClearTeamCache() {
	ctx := context.Background()

	// 清除总数缓存
	configs.RedisClient.Del(ctx, TeamTotalCacheKey)

	// 清除列表缓存（使用通配符）
	pattern := fmt.Sprintf("%s*", TeamListCachePrefix)
	iter := configs.RedisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		configs.RedisClient.Del(ctx, iter.Val())
	}
}

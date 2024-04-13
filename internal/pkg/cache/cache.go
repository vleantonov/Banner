package cache

import (
	"banner/internal/repository/ttl"
	"github.com/ReneKroon/ttlcache"
	"go.uber.org/zap"
	"time"
)

func SetupCache(t time.Duration, l *zap.Logger) *ttl.BannerContent {
	cache := ttlcache.NewCache()
	cache.SetTTL(t)
	cache.SetNewItemCallback(func(key string, value interface{}) {
		l.Info("new cache element", zap.String("key", key))
	})
	cache.SetExpirationCallback(func(key string, value interface{}) {
		l.Info("cache element has been expired", zap.String("key", key))
	})
	return ttl.New(
		cache,
	)
}

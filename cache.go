package cache

import (
	"context"
	"github.com/jiebutech/config"
	"sync"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) string
	GetAndUnmarshal(ctx context.Context, key string, container interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error
	HSet(ctx context.Context, key, field string, value interface{}) error
	HDelete(ctx context.Context, key string, field ...string) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context)
}

type CacheGenerator func(conf config.CacheConfig) Cache

// RegisterDriver 驱动注册
func RegisterDriver(driverName string, generator CacheGenerator) {
	cacheGeneratorStash[driverName] = generator
}

var cacheGeneratorStash = map[string]CacheGenerator{}

var cacheInstance = sync.Map{}

// New 实例化缓存
func New(cfg config.CacheConfig) Cache {
	name := cfg.Type
	ins, ok := cacheInstance.Load(name)
	if !ok {
		generator, ok := cacheGeneratorStash[name]
		if !ok {
			panic(name + " cache driver not register")
		}
		ins = generator(cfg)
	}
	return ins.(Cache)
}

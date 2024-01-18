package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jiebutech/config"
	"github.com/jiebutech/log"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

const DriverName = "redis"

func init() {
	cache.RegisterDriver(DriverName, func(conf config.CacheConfig) cache.Cache {
		return &redisCache{
			cli: redis.NewClient(&redis.Options{
				Addr:     conf.Host + ":" + conf.Password,
				Password: conf.Password,
				DB:       conf.DB,
			}),
			logErr: conf.Debug,
		}
	})

}

type redisCache struct {
	cli    *redis.Client
	logErr bool
}

func (r *redisCache) Get(ctx context.Context, key string) string {
	result, err := r.cli.Get(ctx, key).Result()
	if err != nil && r.logErr {
		log.Error(ctx, "get value from redis cache error", zap.Error(err), zap.String("key", key))
	}
	return result
}

func (r *redisCache) GetAndUnmarshal(ctx context.Context, key string, container interface{}) error {
	s := r.Get(ctx, key)
	if s == "" {
		return errors.New("cache value is empty")
	}
	err := json.Unmarshal([]byte(s), container)
	return err
}

func (r *redisCache) Set(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error {
	expire := time.Duration(redis.KeepTTL)
	if len(ttl) > 0 {
		expire = ttl[0]
	}
	str, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = r.cli.Set(ctx, key, string(str), expire).Result()
	if err != nil && r.logErr {
		log.Error(ctx, "set value to redis cache error", zap.Error(err), zap.String("key", key), zap.Any("value", value))
	}
	return err
}

func (r *redisCache) HSet(ctx context.Context, key, field string, value interface{}) error {
	_, err := r.cli.HSet(ctx, key, field, value).Result()
	if err != nil && r.logErr {
		log.Error(ctx, "set hash value to redis cache error", zap.Error(err), zap.String("key", key), zap.String("field", field), zap.Any("value", value))
	}
	return err
}

func (r *redisCache) HDelete(ctx context.Context, key string, field ...string) error {
	_, err := r.cli.HDel(ctx, key, field...).Result()
	if err != nil && r.logErr {
		log.Error(ctx, "set hash value to redis cache error", zap.Error(err), zap.String("key", key), zap.Any("field", field))
	}
	return err
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	_, err := r.cli.Del(ctx, key).Result()
	if err != nil && r.logErr {
		log.Error(ctx, "delete cache from redis error", zap.Error(err), zap.String("key", key))
	}
	return err
}

func (r *redisCache) Clear(ctx context.Context) {
	panic("The redis cache not support clear all cache")
}

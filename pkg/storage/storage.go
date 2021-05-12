package storage

import (
	"github.com/Walker-PI/edgex-gateway/conf"
	"github.com/go-redis/redis/v8"
)

var (
	// RedisClient ...
	RedisClient *redis.Client
)

func InitStorage() {
	initRedisClient()
}

func initRedisClient() {
	redisOpt := &redis.Options{
		Addr:     conf.RedisConf.Address,
		Password: conf.RedisConf.Password,
		DB:       conf.RedisConf.DB,
	}
	RedisClient = redis.NewClient(redisOpt)
}

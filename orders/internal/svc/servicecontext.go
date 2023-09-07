package svc

import (
	"lab-zero/orders/internal/config"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/syncx"
)

type ServiceContext struct {
	Config       config.Config
	RedisCli     *redis.Redis
	SingleFlight syncx.SingleFlight
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		RedisCli:     redis.MustNewRedis(c.RedisConf),
		SingleFlight: syncx.NewSingleFlight(),
	}
}

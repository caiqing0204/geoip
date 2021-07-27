package cache

import (
	"geoip/pkg/logger"
	"github.com/go-redis/redis/v7"
)

type redisCache struct {
	options Options
	Client  *redis.Client
}

func newRedisCache(opts ...Option) Cache {
	r := new(redisCache)
	_ = r.Init(opts...)
	return r
}

func (r *redisCache) Close() error {
	return r.Client.Close()
}

func (r *redisCache) Init(opts ...Option) error {
	for _, o := range opts {
		o(&r.options)
	}
	nodes := r.options.Nodes

	if len(nodes) == 0 {
		nodes = []string{"localhost:6379"}
	}

	redisOptions := &redis.Options{
		Addr:     nodes[0],
		Password: r.options.Password,
		DB:       r.options.DataBase,
	}
	r.Client = redis.NewClient(redisOptions)

	pong, err := r.Client.Ping().Result()
	if err != nil {
		logger.Fatalf("init redis error:%s", err)
		return err
	}
	logger.Infof("init redis success cmd:PING result:%s", pong)
	return nil
}

func (r *redisCache) Set(name string, value interface{}) error {
	return r.Client.Set(name, value, 0).Err()
}

func (r *redisCache) Get(name string) ([]byte, error) {
	val, err := r.Client.Get(name).Bytes()
	return val, err
}

func (r *redisCache) Del(name string) error {
	return r.Client.Del(name).Err()
}

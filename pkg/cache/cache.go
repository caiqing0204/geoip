package cache

import "context"

type Options struct {
	Nodes    []string
	DataBase int
	Password string
	Context  context.Context
}

type Option func(o *Options)

func WithNode(nodes ...string) Option {
	return func(o *Options) {
		o.Nodes = nodes
	}
}

func WithDB(db int) Option {
	return func(o *Options) {
		o.DataBase = db
	}
}

func WithPassword(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

type Cache interface {
	Init(...Option) error
	Set(name string, value interface{}) error
	Get(name string) ([]byte, error)
	Del(name string) error
}

func NewCache(opts ...Option) Cache {
	return newRedisCache(opts...)
}

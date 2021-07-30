package geoip

import (
	"context"
	"fmt"
	"geoip/pkg/cache"
	"geoip/pkg/codec"
	"geoip/pkg/logger"
)

type geoipCache struct {
	codec   codec.Marshaler
	client  cache.Cache
	context context.Context
	options Options
}

func WithCache(cache cache.Cache) Option {
	return func(o *Options) {
		o.cache = cache
	}
}

func NewGeoIpCache(opts ...Option) Interface {
	options := Options{
		codec:    codec.DefaultMarshaler,
		cache:    nil,
		regionId: DefaultRegionId,
		Context:  context.TODO(),
	}
	for _, o := range opts {
		o(&options)
	}
	if options.cache == nil {
		options.cache = cache.NewCache()
	}

	g := &geoipCache{
		options: options,
		context: context.TODO(),
		client:  options.cache,
		codec:   options.codec,
	}
	return g
}

func (g geoipCache) Init(option ...Option) error {
	return nil
}

func (g geoipCache) format(ipv4 IPV4) string {
	return fmt.Sprintf("geoip/%s", ipv4)
}

func (g *geoipCache) Select(ipv4 IPV4) *IPLocation {
	name := g.format(ipv4)
	b, err := g.client.Get(name)
	if err != nil {
		return nil
	}
	ip := &IPLocation{}
	err = g.codec.Unmarshal(b, ip)
	if err != nil {
		logger.Errorf("geoip cache unmarshal ip:%s err:%s", ipv4.String(), err)
		return nil
	}
	return ip
}

func (g *geoipCache) Update(ipv4 IPV4, location *IPLocation) bool {
	name := g.format(ipv4)
	value, err := g.codec.Marshal(location)
	if err != nil {
		return false
	}
	err = g.client.Set(name, value)
	if err != nil {
		logger.Errorf("geoip cache set ip:%s err:%s", ipv4.String(), err)
		return false
	}
	return true
}

func (g *geoipCache) Delete(ipv4 IPV4) bool {
	name := g.format(ipv4)
	err := g.client.Del(name)
	if err != nil {
		logger.Errorf("geoip cache del ip:%s err:%s", ipv4.String(), err)
		return false
	}
	return true
}

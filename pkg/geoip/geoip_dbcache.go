package geoip

import (
	"geoip/pkg/logger"
)

type geoipDBCache struct {
	cache Interface
	db    Interface
}

// 当不知道起什么名字的时候 就handler

func NewGeoIpHandler(db, cache Interface) Interface {

	g := geoipDBCache{
		db:    db,
		cache: cache,
	}

	return &g
}

func (g geoipDBCache) Init(option ...Option) error {
	return nil
}

func (g geoipDBCache) Select(ipv4 IPV4) *IPLocation {
	ip := g.cache.Select(ipv4)
	if ip == nil {
		ip = g.db.Select(ipv4)
		if ip != nil {
			g.Update(ip.Ip, ip)
		}
	}
	if ip != nil {
		logger.Infof("HIT ip:%s", ip.Ip.String())
	}
	return ip
}

func (g geoipDBCache) Update(ipv4 IPV4, location *IPLocation) bool {
	return g.cache.Update(ipv4, location)
}

func (g geoipDBCache) Delete(ipv4 IPV4) bool {
	return g.cache.Delete(ipv4)
}

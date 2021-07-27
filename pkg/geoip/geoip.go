package geoip

import (
	"context"
	"geoip/pkg/cache"
	"geoip/pkg/codec"
)

var (
	DefaultRegionId = "cn-hangzhou"
)

type IPV4 string

func (ip IPV4) String() string {
	return string(ip)
}

type IPLocation struct {
	Ip       IPV4   `json:"ip"`        // IP 地址
	Country  string `json:"country"`   // 国家
	Province string `json:"province" ` // 省份
	City     string `json:"city" `     // 城市
	County   string `json:"county" `   // 县
	Isp      string `json:"isp"`       // 运营商
}

type Options struct {
	regionId        string
	accessKeyId     string
	accessKeySecret string
	cache           cache.Cache
	codec           codec.Marshaler
	Context         context.Context
}

type Option func(*Options)

func WithAccessKeySecret(ak, sk string) Option {
	return func(o *Options) {
		o.accessKeyId = ak
		o.accessKeySecret = sk
	}
}

type Interface interface {
	Init(...Option) error
	Select(ipv4 IPV4) *IPLocation
	Update(ipv4 IPV4, location *IPLocation) bool
	Delete(ipv4 IPV4) bool
}

func NewGeoIP(opt ...Option) Interface {
	return NewGeoIpDB(opt...)
}

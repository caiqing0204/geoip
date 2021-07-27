package geoip

import (
	"context"
	"geoip/pkg/codec"
	"geoip/pkg/logger"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/geoip"
)

type geoIPDBAliyun struct {
	client  *geoip.Client
	options Options
}

func NewGeoIpDB(opts ...Option) Interface {
	options := Options{
		regionId: DefaultRegionId,
		codec:    codec.DefaultMarshaler,
		cache:    nil,
		Context:  context.TODO(),
	}
	for _, o := range opts {
		o(&options)
	}

	c, err := geoip.NewClientWithAccessKey(options.regionId, options.accessKeyId, options.accessKeySecret)
	if err != nil {
		logger.Fatalf("init geoip aliyun error:%s", err)
		return nil
	}

	g := &geoIPDBAliyun{
		options: options,
		client:  c,
	}
	logger.Info("init geoip aliyun success")
	return g
}

func (g *geoIPDBAliyun) Init(opts ...Option) error {
	return nil
}

func (g *geoIPDBAliyun) Select(ipv4 IPV4) *IPLocation {
	if g.client == nil {
		return nil
	}
	request := geoip.CreateDescribeIpv4LocationRequest()
	request.Ip = string(ipv4)
	request.Lang = "en"
	request.AcceptFormat = "json"
	response, err := g.client.DescribeIpv4Location(request)
	if err != nil {
		logger.Errorf("select ip from geoip aliyun error:%s", err)
		return nil
	}
	location := &IPLocation{}
	location.Ip = IPV4(response.Ip)
	location.City = response.City
	location.Country = response.Country
	location.Isp = response.Isp
	location.Province = response.Province
	return location
}

func (g *geoIPDBAliyun) Update(ipv4 IPV4, location *IPLocation) bool {
	return true
}

func (g *geoIPDBAliyun) Delete(ipv4 IPV4) bool {
	return true
}

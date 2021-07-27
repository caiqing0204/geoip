package geoip

import (
	"testing"
)

func TestGeoIPCache_Update(t *testing.T) {
	ip := &IPLocation{Ip: "192.168.11.1", City: "beiJing", Country: "china", Province: "beijing", Isp: "liantong"}
	g := NewGeoIpCache()
	ok := g.Update(ip.Ip, ip)
	t.Log(ok)
}
func TestGeoipCache_Select(t *testing.T) {
	ip := "192.168.11.1"
	g := NewGeoIpCache()
	result := g.Select(IPV4(ip))
	t.Log(result)
}

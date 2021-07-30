package config

type RedisConfig struct {
	Host     string `json:"host"`
	DB       int    `json:"db"`
	Password string `json:"password"`
}

type GeoIpConfig struct {
	Ak string `json:"ak"`
	Sk string `json:"sk"`
}

type GlobalConfig struct {
	Port  int         `json:"port"`
	Redis RedisConfig `json:"redis"`
	GeoIp GeoIpConfig `json:"geoip"`
}

package main

import (
	"fmt"
	"geoip/pkg/cache"
	"geoip/pkg/config"
	"geoip/pkg/geoip"
	"geoip/pkg/logger"
	"geoip/pkg/logger/zerolog"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"regexp"
	"time"
)

func IpMatch(ip string) bool {
	reg := regexp.MustCompile("^[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}$")
	return reg.Match([]byte(ip))
}

type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type IPLocationResponse struct {
	BaseResponse
	Data geoip.IPLocation `json:"data"`
}

func HttpLogger() gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		logger.Infof("code:%d method:%s latency:%d x-forwared-for:%s url:%s", statusCode, method, latency.String(), clientIP, c.Request.RequestURI)
	}
}

func daemon(conf *config.GlobalConfig) {
	// 日志初始化
	logger.DefaultLogger = zerolog.NewLogger(
		logger.WithLevel(logger.InfoLevel),
		logger.WithOutputer(logger.NewOutputer("geoip", "")),
		logger.WithCallerSkipCount(4),
		zerolog.WithTimeFormat(time.RFC3339))

	logger.Info("logger init success")

	geoIpdb := geoip.NewGeoIpDB(
		geoip.WithAccessKeySecret(conf.GeoIp.Ak, conf.GeoIp.Sk))

	geoIpCache := geoip.NewGeoIpCache(
		geoip.WithCache(
			cache.NewCache(
				cache.WithNode(conf.Redis.Host),
				cache.WithDB(conf.Redis.DB),
				cache.WithPassword(conf.Redis.Password),
			),
		))

	handler := geoip.NewGeoIpHandler(geoIpdb, geoIpCache)

	gin.SetMode(gin.ReleaseMode)
	service := gin.New()
	service.Use(HttpLogger())
	service.GET("/geoip/api/v1/ip", func(c *gin.Context) {
		ipQuery := c.Query("ip")
		if ok := IpMatch(ipQuery); ok {
			result := handler.Select(geoip.IPV4(ipQuery))
			c.JSON(http.StatusOK, &IPLocationResponse{BaseResponse{Code: 0, Msg: "success"}, *result})
			return
		} else {
			c.JSON(http.StatusOK, &BaseResponse{Code: 1, Msg: "IP不合法或IP参数未传"})
			return
		}
	})
	logger.Infof("http listen :%d", conf.Port)
	err := service.Run(fmt.Sprintf(":%d", conf.Port))
	logger.Fatalf("http stop message:%s", err)
}

func main() {
	app := &cli.App{
		Name: "geoip service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "the configuration file",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			cfgFileName := c.String("config")
			var conf config.GlobalConfig
			err := configor.New(&configor.Config{Debug: false}).Load(&conf, cfgFileName)
			if err != nil {
				logger.Fatalf("parse config error:%s", err)
				return err
			}
			daemon(&conf)
			return nil
		},
	}
	app.Version = "1.0.0"
	err := app.Run(os.Args)
	if err != nil {
		logger.Fatalf("start service err: %s", err)
	}
}

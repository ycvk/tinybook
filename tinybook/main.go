package main

import (
	"context"
	"geek_homework/tinybook/ioc"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"net/http"
	"time"
)

func init() {
	initViper()         // 初始化配置
	ioc.InitSnowflake() // 初始化雪花算法
	initPrometheus()    // 初始化prometheus
}

func main() {
	otel := ioc.InitOTEL() // 初始化otel
	defer func() {
		timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
		defer cancelFunc()
		otel(timeout) // 服务器关闭时, 超时控制去关闭otel
	}()
	app := InitWebServer()         // 初始化web服务
	for i := range app.consumers { // 启动kafka消费者
		app.consumers[i].Start()
	}
	err := app.server.Run(":8081") // 启动web服务
	if err != nil {
		panic(err)
	}
}

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8082", nil)
	}()
}

func initViper() {
	//s := pflag.String("config", "config/dev.yaml", "配置文件路径")
	//pflag.Parse()
	//viper.SetConfigFile(*s)
	//err := viper.ReadInConfig()
	//if err != nil {
	//	panic(err)
	//}
	//err := viper.AddRemoteProvider("etcd3", "tinybook-etcd:2381", "tinybook")
	err := viper.AddRemoteProvider("etcd3", "127.0.0.1:32379", "tinybook-local")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

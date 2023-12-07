package main

import (
	"geek_homework/tinybook/ioc"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {
	initViper()
	ioc.InitSnowflake()
	app := InitWebServer()
	for i := range app.consumers {
		app.consumers[i].Start()
	}
	err := app.server.Run(":8081")
	if err != nil {
		panic(err)
	}
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

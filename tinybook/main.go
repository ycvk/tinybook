package main

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {
	initViper()
	engine := InitWebServer()
	err := engine.Run(":8081")
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
	err := viper.AddRemoteProvider("etcd3", "http://127.0.0.1:12379", "tinybook")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

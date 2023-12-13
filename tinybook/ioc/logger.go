package ioc

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	err := viper.UnmarshalKey("log", &config)
	if err != nil {
		panic(err)
	}
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	return logger
}

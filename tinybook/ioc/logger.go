package ioc

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	err := viper.UnmarshalKey("log", &config) // 从配置文件中读取配置
	if err != nil {
		panic(err)
	}
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 配置日志级别颜色
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	return logger
}

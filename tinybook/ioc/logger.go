package ioc

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func InitLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	//err := viper.UnmarshalKey("log", &config) // 从配置文件中读取配置
	//if err != nil {
	//	panic(err)
	//}
	//newConfig := prom.NewConfig()
	//core, err := newConfig.Build()
	//if err != nil {
	//	panic(err)
	//}
	//l := zap.New(core)

	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 配置日志级别颜色
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // 配置时间格式
	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(config.EncoderConfig),
		//zapcore.Lock(os.Stdout),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(zap.InfoLevel),
	))
	//logger, err := config.Build()
	//if err != nil {
	//	panic(err)
	//}
	zap.ReplaceGlobals(logger)
	return logger
}

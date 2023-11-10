package ioc

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

var (
	gormDB *gorm.DB
	once   sync.Once
)

func InitDB(zipLog *zap.Logger) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}

	once.Do(func() {
		gormDB, err = gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
			Logger: logger.New(gormLoggerFunc(func(msg string, data ...interface{}) {
				zipLog.Info(msg, zap.Any("data", data))
			}), logger.Config{
				SlowThreshold: 0,           // 慢查询阈值 0 表示打印所有sql
				LogLevel:      logger.Info, // 日志级别
			}),
		})
	})
	if err != nil {
		panic(err)
	}
	// TODO 为了方便测试，每次启动都会重新创建表 仅供测试使用
	CreateTable(gormDB)
	return gormDB
}

// gormLoggerFunc gorm日志
type gormLoggerFunc func(msg string, data ...interface{})

// Printf 实现gorm日志接口
func (f gormLoggerFunc) Printf(format string, args ...interface{}) {
	f(format, args...)
}

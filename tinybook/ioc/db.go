package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	gormDB *gorm.DB
	once   sync.Once
)

func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}

	once.Do(func() {
		gormDB, err = gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
		gormDB = gormDB.Debug() // 开启debug模式 会打印sql语句 便于调试
		if err != nil {
			panic(err)
		}
	})
	// TODO 为了方便测试，每次启动都会重新创建表 仅供测试使用
	CreateTable(gormDB)
	return gormDB
}

package ioc

import (
	"geek_homework/tinybook/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	gormDB *gorm.DB
	once   sync.Once
)

func InitDB() *gorm.DB {
	once.Do(func() {
		dsn := config.Config.DB.Host
		var err error
		gormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		gormDB = gormDB.Debug()
		if err != nil {
			panic(err)
		}
	})
	return gormDB
}

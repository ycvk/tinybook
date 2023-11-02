package ioc

import (
	"geek_homework/tinybook/internal/repository/dao"
	"gorm.io/gorm"
)

func CreateTable(db *gorm.DB) {
	err := db.AutoMigrate(&dao.User{}, &dao.Code{})
	if err != nil {
		panic(err)
	}
}

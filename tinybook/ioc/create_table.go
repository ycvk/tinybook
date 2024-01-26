package ioc

import (
	"gorm.io/gorm"
	dao2 "tinybook/tinybook/article/repository/dao"
	"tinybook/tinybook/internal/repository/dao"
)

func CreateTable(db *gorm.DB) {
	err := db.AutoMigrate(
		&dao.User{},
		&dao.SMS{},
		&dao2.Article{},
		&dao2.PublishedArticle{},
		&dao.Job{},
	)
	if err != nil {
		panic(err)
	}
}

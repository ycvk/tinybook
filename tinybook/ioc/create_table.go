package ioc

import (
	"gorm.io/gorm"
	"tinybook/tinybook/internal/repository/dao"
)

func CreateTable(db *gorm.DB) {
	err := db.AutoMigrate(
		&dao.User{},
		&dao.SMS{},
		&dao.Article{},
		&dao.PublishedArticle{},
		&dao.Job{},
	)
	if err != nil {
		panic(err)
	}
}

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
		&dao.Interactive{},
		&dao.LikeRecord{},
		&dao.CollectRecord{},
		&dao.Job{},
	)
	if err != nil {
		panic(err)
	}
}

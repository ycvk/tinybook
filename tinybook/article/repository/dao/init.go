package dao

import "gorm.io/gorm"

func CreateTableForArticle(db *gorm.DB) {
	err := db.AutoMigrate(
		&Article{},
		&PublishedArticle{},
	)
	if err != nil {
		panic(err)
	}
}

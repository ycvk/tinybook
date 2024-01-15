package dao

import "gorm.io/gorm"

func CreateTableForInteractive(db *gorm.DB) {
	err := db.AutoMigrate(
		&Interactive{},
		&LikeRecord{},
		&CollectRecord{},
	)
	if err != nil {
		panic(err)
	}
}

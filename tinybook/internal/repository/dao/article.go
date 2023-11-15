package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, article Article) (int64, error)
}

type Article struct {
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement;not null" json:"id"`
	Title    string `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Content  string `gorm:"column:content;type:BLOB;not null" json:"content"`
	AuthorId int64  `gorm:"index;column:author_id;not null" json:"author_id"`
	Ctime    int64  `gorm:"column:ctime" json:"ctime"`
	Utime    int64  `gorm:"column:utime" json:"utime"`
}

type GormArticleDAO struct {
	db *gorm.DB
}

func NewGormArticleDAO(db *gorm.DB) ArticleDAO {
	return &GormArticleDAO{db: db}
}

func (g *GormArticleDAO) Insert(ctx context.Context, article Article) (int64, error) {
	now := time.Now().Unix()
	article.Ctime, article.Utime = now, now
	err := g.db.WithContext(ctx).Create(&article).Error
	return article.ID, err
}

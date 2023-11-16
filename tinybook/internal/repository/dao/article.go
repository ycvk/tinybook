package dao

import (
	"context"
	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, article Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
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

func (g *GormArticleDAO) UpdateById(ctx context.Context, article Article) error {
	updates := g.db.WithContext(ctx).Model(&article).
		Where("id = ? AND author_id = ?", article.ID, article.AuthorId).
		Updates(map[string]any{
			"title":   article.Title,
			"content": article.Content,
			"utime":   time.Now().Unix(),
		})
	if updates.Error != nil {
		return updates.Error
	}
	if updates.RowsAffected == 0 {
		return errors.New("作者ID与文章ID不匹配")
	}
	return nil
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

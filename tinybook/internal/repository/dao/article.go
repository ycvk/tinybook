package dao

import (
	"context"
	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, article Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
	Sync(ctx context.Context, dao Article) (int64, error)
}

type Article struct {
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement;not null" json:"id"`
	Title    string `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Content  string `gorm:"column:content;type:BLOB;not null" json:"content"`
	AuthorId int64  `gorm:"index;column:author_id;not null" json:"author_id"`
	Ctime    int64  `gorm:"column:ctime" json:"ctime"`
	Utime    int64  `gorm:"column:utime" json:"utime"`
}

type PublishedArticle Article

type GormArticleDAO struct {
	db *gorm.DB
}

func (g *GormArticleDAO) Sync(ctx context.Context, article Article) (int64, error) {
	txErr := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		articleDAO := NewGormArticleDAO(tx) // 事务中的DAO
		if article.ID > 0 {                 // 更新
			if err := articleDAO.UpdateById(ctx, article); err != nil {
				return err
			}
		} else { // 新增
			id, err := articleDAO.Insert(ctx, article)
			if err != nil {
				return err
			}
			article.ID = id
		}
		now := time.Now().Unix()
		publishedArticle := PublishedArticle(article)
		publishedArticle.Ctime, publishedArticle.Utime = now, now // 更新发布时间 & 更新时间

		err := tx.Clauses(clause.OnConflict{
			// id更新冲突时，只更新title、content、utime字段
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":   publishedArticle.Title,
				"content": publishedArticle.Content,
				"utime":   now,
			}),
		}).Create(&publishedArticle).Error
		return err
	})
	if txErr != nil {
		return 0, txErr
	}
	return article.ID, nil
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

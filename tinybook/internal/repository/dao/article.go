package dao

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/godruoyi/go-snowflake"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, article Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
	Sync(ctx context.Context, dao Article) (int64, error)
	SyncStatus(ctx context.Context, dao Article, u uint8) error
	GetArticlesByAuthor(ctx context.Context, uid int64, limit int, offset int) ([]Article, error)
}

type Article struct {
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement;not null" json:"id" bson:"id,omitempty"`
	Title    string `gorm:"column:title;type:varchar(255);not null" json:"title" bson:"title,omitempty"`
	Content  string `gorm:"column:content;type:BLOB;not null" json:"content" bson:"content,omitempty"`
	AuthorId int64  `gorm:"index;column:author_id;not null" json:"author_id" bson:"author_id,omitempty"`
	Status   uint8  `gorm:"column:status;type:tinyint(1);not null" json:"status" bson:"status,omitempty"`
	Ctime    int64  `gorm:"column:ctime" json:"ctime" bson:"ctime,omitempty"`
	Utime    int64  `gorm:"column:utime" json:"utime" bson:"utime,omitempty"`
}

type PublishedArticle Article

type GormArticleDAO struct {
	db *gorm.DB
}

func (g *GormArticleDAO) GetArticlesByAuthor(ctx context.Context, uid int64, limit int, offset int) ([]Article, error) {
	var articles []Article
	err := g.db.WithContext(ctx).
		Where("author_id = ?", uid).
		Order("utime desc").
		Limit(limit).
		Offset(offset).
		Find(&articles).
		Error
	return articles, err
}

func NewGormArticleDAO(db *gorm.DB) ArticleDAO {
	return &GormArticleDAO{db: db}
}

func (g *GormArticleDAO) SyncStatus(ctx context.Context, dao Article, u uint8) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().Unix()
		updates := tx.Model(&Article{}).
			Where("id = ? and author_id = ?", dao.ID, dao.AuthorId). // 判断作者ID与文章ID是否匹配
			Updates(map[string]any{
				"status": u,
				"utime":  now,
			})
		if updates.Error != nil {
			return updates.Error
		}
		if updates.RowsAffected == 0 {
			return errors.New("作者ID与文章ID不匹配")
		}
		return tx.Model(&PublishedArticle{}).
			Where("id = ?", dao.ID). // 前面已经判断了author_id，这里不需要再判断
			Updates(map[string]any{
				"status": u,
				"utime":  now,
			}).Error
	})
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
			DoUpdates: clause.Assignments(map[string]any{
				"title":   publishedArticle.Title,
				"content": publishedArticle.Content,
				"status":  publishedArticle.Status,
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
			"status":  article.Status,
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

func (g *GormArticleDAO) Insert(ctx context.Context, article Article) (int64, error) {
	now := time.Now().Unix()
	article.Ctime, article.Utime = now, now
	err := g.db.WithContext(ctx).Create(&article).Error
	return article.ID, err
}

type MongoDBArticleDAO struct {
	db            *qmgo.Database
	coll          *qmgo.Collection
	publishedColl *qmgo.Collection
}

func NewMongoDBArticleDAO(db *qmgo.Database) ArticleDAO {
	return &MongoDBArticleDAO{
		db:            db,
		coll:          db.Collection("articles"),
		publishedColl: db.Collection("published_articles"),
	}
}

func (m *MongoDBArticleDAO) GetArticlesByAuthor(ctx context.Context, uid int64, limit int, offset int) ([]Article, error) {
	var articles []Article
	err := m.coll.Find(ctx, bson.M{"author_id": uid}).Skip(int64(offset)).Limit(int64(limit)).All(&articles)
	return articles, err
}

func (m *MongoDBArticleDAO) Insert(ctx context.Context, article Article) (int64, error) {
	article.ID = int64(snowflake.ID())
	now := time.Now().Unix()
	article.Ctime, article.Utime = now, now
	_, err := m.coll.InsertOne(ctx, &article)
	return article.ID, err
}

func (m *MongoDBArticleDAO) UpdateById(ctx context.Context, article Article) error {
	now := time.Now().Unix()
	err := m.coll.UpdateOne(ctx,
		bson.M{"id": article.ID, "author_id": article.AuthorId}, // 判断作者ID与文章ID是否匹配
		bson.M{
			"$set": bson.M{
				"title":   article.Title,
				"content": article.Content,
				"status":  article.Status,
				"utime":   now,
			},
		})
	return err
}

func (m *MongoDBArticleDAO) Sync(ctx context.Context, dao Article) (int64, error) {
	now := time.Now().Unix()
	var err error
	article := PublishedArticle(dao)
	if dao.ID > 0 { // 更新
		err = m.UpdateById(ctx, dao)
	} else {
		article.Ctime, article.Utime = now, now // 更新发布时间 & 更新时间
		dao.ID, err = m.Insert(ctx, Article(article))
	}
	if err != nil { // 更新或发布文章失败
		return dao.ID, err
	}
	// 更新发布文章
	_, err = m.publishedColl.Upsert(ctx, bson.M{"id": dao.ID, "author_id": dao.AuthorId},
		bson.M{
			"id":        dao.ID, // 保证id不变
			"title":     article.Title,
			"content":   article.Content,
			"status":    article.Status,
			"author_id": article.AuthorId,
			"utime":     now,
		})
	return dao.ID, err
}

func (m *MongoDBArticleDAO) SyncStatus(ctx context.Context, dao Article, u uint8) error {
	now := time.Now().Unix()
	err := m.UpdateById(ctx, dao)
	if err != nil {
		return err
	}
	return m.publishedColl.UpdateOne(ctx, bson.M{"id": dao.ID, "author_id": dao.AuthorId},
		bson.M{
			"$set": bson.M{
				"status": u,
				"utime":  now,
			}},
	)
}

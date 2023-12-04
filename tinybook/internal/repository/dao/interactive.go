package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Interactive struct {
	Id    int64  `gorm:"column:id;primaryKey;autoIncrement;not null"`
	BizId int64  `gorm:"column:biz_id;not null;uniqueIndex:idx_biz_type_id"`
	Biz   string `gorm:"column:biz;not null;uniqueIndex:idx_biz_type;type:varchar(32)"`

	ReadCount    int64 `gorm:"column:read_count"`
	LikeCount    int64 `gorm:"column:like_count"`
	CollectCount int64 `gorm:"column:collect_count"`
	Utime        int64 `gorm:"column:utime;not null"`
	Ctime        int64 `gorm:"column:ctime;not null"`
}

type InteractiveDAO interface {
	IncreaseReadCount(ctx context.Context, biz string, bizId int64) error
}

type GormInteractiveDAO struct {
	db *gorm.DB
}

func NewGormInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &GormInteractiveDAO{db: db}
}

func (g *GormInteractiveDAO) IncreaseReadCount(ctx context.Context, biz string, bizId int64) error {
	now := time.Now().Unix()
	return g.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"read_count": gorm.Expr("read_count + ?", 1),
			"utime":      now,
		}),
	}).Create(&Interactive{
		BizId:     bizId,
		Biz:       biz,
		ReadCount: 1,
		Ctime:     now,
		Utime:     now,
	}).Error
}

package dao

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var DuplicatePhoneError = errors.New("手机号码已经在等待发送队列中")

type CodeDAO interface {
	Insert(ctx context.Context, code ...*Code) error
	Delete(ctx context.Context, code Code) error
}

type Code struct {
	Id     int64  `gorm:"column:id;primaryKey;autoIncrement;not null"`
	Number string `gorm:"unique;column:number;not null"`
	Ctime  int64  `gorm:"column:ctime;not null"`
	Utime  int64  `gorm:"column:utime;not null"`
}

type GormCodeDAO struct {
	db *gorm.DB
}

func NewGormCodeDAO(db *gorm.DB) CodeDAO {
	return &GormCodeDAO{
		db: db,
	}
}

func (g *GormCodeDAO) Insert(ctx context.Context, code ...*Code) error {
	now := time.Now().UnixMilli()
	for _, c := range code {
		c.Ctime, c.Utime = now, now
	}
	err := g.db.WithContext(ctx).Create(code).Error
	var my *mysql.MySQLError
	if errors.As(err, &my) {
		// 如果是重复的手机号码，说明此号码已经在等待发送队列中了，返回错误
		if my.Number == 1062 {
			return DuplicatePhoneError
		}
	}
	return err
}

func (g *GormCodeDAO) Delete(ctx context.Context, code Code) error {
	return g.db.WithContext(ctx).Delete(&code).Error
}

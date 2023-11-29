package dao

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var DuplicatePhoneError = errors.New("手机号码已经在等待发送队列中")

type SMSDAO interface {
	Insert(ctx context.Context, code ...*SMS) error
	Delete(ctx context.Context, code ...*SMS) error
}

type SMS struct {
	Id    int64  `gorm:"column:id;primaryKey;autoIncrement;not null"`
	Phone string `gorm:"unique;column:phone;not null"`
	Code  string `gorm:"column:code;not null"`
	TplId string `gorm:"column:tpl_id"`
	Ctime int64  `gorm:"column:ctime;not null"`
	Utime int64  `gorm:"column:utime;not null"`
}

type GormSMSDAO struct {
	db *gorm.DB
}

func NewGormSMSDAO(db *gorm.DB) SMSDAO {
	return &GormSMSDAO{db: db}
}

func (g *GormSMSDAO) Insert(ctx context.Context, sms ...*SMS) error {
	now := time.Now().UnixMilli()
	for _, s := range sms {
		s.Ctime, s.Utime = now, now
	}
	err := g.db.WithContext(ctx).Create(sms).Error
	var my *mysql.MySQLError
	if errors.As(err, &my) {
		// 如果是重复的手机号码，说明此号码已经在等待发送队列中了，返回错误
		if my.Number == 1062 {
			return DuplicatePhoneError
		}
	}
	return err
}

func (g *GormSMSDAO) Delete(ctx context.Context, sms ...*SMS) error {
	strings := make([]string, 0, len(sms))
	for _, sm := range sms {
		strings = append(strings, sm.Phone)
	}
	return g.db.WithContext(ctx).Delete(&sms, "phone in (?)", strings).Error
}

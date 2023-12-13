package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	duplicateEmailError = errors.New("邮箱已存在")
	ErrUserNotFound     = gorm.ErrRecordNotFound.Error()
	ErrUserDuplicate    = gorm.ErrDuplicatedKey.Error()
)

type UserDAO interface {
	Insert(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindByWechatOpenId(ctx context.Context, openId string) (User, error)
	Update(ctx context.Context, user User) error
}

type GormUserDAO struct {
	db *gorm.DB
}

// User 用户表
type User struct {
	Id            int64          `gorm:"column:id;primaryKey;autoIncrement;not null"`
	Email         sql.NullString `gorm:"unique;column:email"`
	Phone         sql.NullString `gorm:"unique;column:phone"`
	Password      string         `gorm:"column:password"`
	Ctime         int64          `gorm:"column:ctime"`
	Utime         int64          `gorm:"column:utime"`
	Nickname      string         `gorm:"column:nickname"`
	Birthday      string         `gorm:"column:birthday"`
	AboutMe       string         `gorm:"column:about_me"`
	WechatOpenId  sql.NullString `gorm:"unique;column:wechat_open_id"`
	WechatUnionId sql.NullString `gorm:"column:wechat_union_id"`
}

func NewGormUserDAO(db *gorm.DB) UserDAO {
	return &GormUserDAO{
		db: db,
	}
}

// Insert 插入用户
func (dao *GormUserDAO) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.Ctime = now
	user.Utime = now
	err := dao.db.WithContext(ctx).Create(&user).Error
	var my *mysql.MySQLError
	if errors.As(err, &my) {
		if my.Number == 1062 {
			return duplicateEmailError
		}
	}
	return err
}

// FindByEmail 根据邮箱查找用户
func (dao *GormUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

// FindById 根据id查找用户
func (dao *GormUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return user, err
}

// Update 调用gorm的更新方法，更新用户信息
func (dao *GormUserDAO) Update(ctx context.Context, user User) error {
	err := dao.db.WithContext(ctx).Updates(&user).Error
	return err
}

// FindByPhone 根据手机号查找用户
func (dao *GormUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	return user, err
}

// FindByWechatOpenId 根据微信openId查找用户
func (dao *GormUserDAO) FindByWechatOpenId(ctx context.Context, openId string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("wechat_open_id = ?", openId).First(&user).Error
	return user, err
}

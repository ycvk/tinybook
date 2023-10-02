package dao

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var duplicateEmailError = errors.New("邮箱已存在")
var ErrUserNotFound = gorm.ErrRecordNotFound.Error()

type UserDAO struct {
	db *gorm.DB
}

type User struct {
	Id       int64  `gorm:"column:id;primaryKey;autoIncrement;not null"`
	Email    string `gorm:"unique;column:email"`
	Password string `gorm:"column:password"`
	Ctime    int64  `gorm:"column:ctime"`
	Utime    int64  `gorm:"column:utime"`
	Nickname string `gorm:"column:nickname"`
	Birthday string `gorm:"column:birthday"`
	AboutMe  string `gorm:"column:about_me"`
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, user User) error {
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

func (dao *UserDAO) FindByEmail(ctx *gin.Context, email string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (dao *UserDAO) FindById(ctx *gin.Context, id int64) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return user, err
}

func (dao *UserDAO) Update(ctx *gin.Context, user User) error {
	err := dao.db.WithContext(ctx).Updates(&user).Error
	return err
}

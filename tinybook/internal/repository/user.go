package repository

import (
	"context"
	"geek_homework/tinybook/internal/domain"
	"geek_homework/tinybook/internal/repository/dao"
	"github.com/gin-gonic/gin"
	"time"
)

var ErrUserNotFound = dao.ErrUserNotFound

type UserRepository struct {
	userDao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		userDao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, user domain.User) error {
	return repo.userDao.Insert(ctx, dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}

func (repo *UserRepository) FindByEmail(ctx *gin.Context, email string) (domain.User, error) {
	user, err := repo.userDao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (repo *UserRepository) UpdateById(ctx *gin.Context, id int64, birthday string, nickname string, me string) error {
	user, err := repo.userDao.FindById(ctx, id)
	if err != nil {
		return err
	}
	user.Birthday = birthday
	user.Nickname = nickname
	user.AboutMe = me
	user.Utime = time.Now().UnixMilli()

	err = repo.userDao.Update(ctx, user)
	return err
}

func (repo *UserRepository) FindById(ctx *gin.Context, id int64) (domain.User, error) {
	byId, err := repo.userDao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Email:    byId.Email,
		Nickname: byId.Nickname,
		Birthday: byId.Birthday,
		AboutMe:  byId.AboutMe,
	}, nil
}

package repository

import (
	"context"
	"database/sql"
	"geek_homework/tinybook/internal/domain"
	"geek_homework/tinybook/internal/repository/cache"
	"geek_homework/tinybook/internal/repository/dao"
	"log/slog"
	"time"
)

var (
	ErrUserNotFound = dao.ErrUserNotFound
	ErrorUserExist  = dao.ErrUserDuplicate
)

type UserRepository struct {
	userDao   *dao.UserDAO
	userCache cache.UserCache
}

// NewUserRepository 构建UserRepository
func NewUserRepository(dao *dao.UserDAO, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		userDao:   dao,
		userCache: *cache,
	}
}

// Create 创建用户
func (repo *UserRepository) Create(ctx context.Context, user domain.User) error {
	return repo.userDao.Insert(ctx, dao.User{
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Password: user.Password,
	})
}

// FindByEmail 根据邮箱查找用户
func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := repo.userDao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       user.Id,
		Email:    user.Email.String,
		Password: user.Password,
	}, nil
}

// UpdateById 根据id更新用户信息
func (repo *UserRepository) UpdateById(ctx context.Context, id int64, birthday string, nickname string, me string) error {
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

// FindById 根据id查找用户
func (repo *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// 先从缓存中查找
	cacheById, err := repo.userCache.GetById(ctx, id)
	// 封装一个从数据库中查找的方法, 定义后不会立即执行, 只有在调用时才会执行
	databaseById := func() (domain.User, error) {
		byId, err := repo.userDao.FindById(ctx, id)
		if err != nil {
			return domain.User{}, err
		}
		// 封装成domain.User
		user := domain.User{
			Id:       byId.Id,
			Email:    byId.Email.String,
			Nickname: byId.Nickname,
			Birthday: byId.Birthday,
			AboutMe:  byId.AboutMe,
		}
		return user, nil
	}
	// 根据err判断缓存中是否有数据
	switch err {
	case nil:
		// 缓存中有则直接返回
		slog.Info("从缓存中获取用户信息", "id", id)
		return cacheById, nil
	case cache.ErrUserNotFound:
		// 缓存中没有则从数据库中查找
		slog.Info("从数据库中获取用户信息", "id", id)
		byId, err := databaseById()
		if err != nil {
			return domain.User{}, err
		}
		go func() {
			// 将查找到的用户信息存入缓存
			err = repo.userCache.SetById(ctx, byId)
			if err != nil {
				slog.Error("缓存用户信息失败", "err", err)
			}
		}()
		return byId, nil
	default:
		// redis有问题，降级处理
		//return domain.User{}, err

		//redis有问题，直接从数据库中查找
		return databaseById()
	}
}

func (repo *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	byPhone, err := repo.userDao.FindByPhone(ctx, phone)
	if err != nil {
		slog.Error("根据手机号查找用户失败", "phone", phone)
		return domain.User{}, err
	}
	return domain.User{
		Id:    byPhone.Id,
		Email: byPhone.Email.String,
		Phone: byPhone.Phone.String,
	}, nil
}

package service

import (
	"errors"
	"geek_homework/tinybook/internal/domain"
	"geek_homework/tinybook/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

var (
	ErrUserNotFound = repository.ErrUserNotFound
	ErrorUserExist  = repository.ErrorUserExist
)

type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService 构建UserService
func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepository,
	}
}

// Signup 注册
func (userService *UserService) Signup(ctx *gin.Context, user domain.User) error {
	password := user.ValidatePassword()
	email := user.ValidateEmail()
	if !email {
		slog.Error("邮箱格式不正确", "email", user.Email)
		return errors.New("邮箱格式不正确")
	}
	if !password {
		slog.Error("密码格式或长度不正确", "password", user.Password)
		return errors.New("密码格式或长度不正确, 长度6-16位, 且只能包含数字和字母和特殊字符")
	}
	generateFromPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("密码加密失败", "password", user.Password)
		return errors.New("密码加密失败")
	}
	user.Password = string(generateFromPassword)
	return userService.userRepo.Create(ctx, user)
}

// Login 登录
func (userService *UserService) Login(ctx *gin.Context, email string, password string) (domain.User, error) {
	byEmail, err := userService.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(byEmail.Password), []byte(password))
	if err != nil {
		slog.Error("密码不正确", "password", password)
		return domain.User{}, errors.New("密码不正确")
	}
	return byEmail, nil
}

// Edit 编辑
func (userService *UserService) Edit(ctx *gin.Context, user domain.User) error {
	birthday := user.ValidateBirthday()
	nickname := user.ValidateNickname()
	aboutMe := user.ValidateAboutMe()
	if !birthday {
		slog.Error("生日格式不正确", "birthday", user.Birthday)
		return errors.New("生日格式不正确 (格式为yyyy-MM-dd) 或输入日期超出当前日期")
	}
	if !nickname {
		slog.Error("昵称格式不正确", "nickname", user.Nickname)
		return errors.New("昵称长度不正确, 长度1-10位")
	}
	if !aboutMe {
		slog.Error("个人简介格式不正确", "aboutMe", user.AboutMe)
		return errors.New("个人简介长度不正确, 长度不能超过200位")
	}
	err := userService.userRepo.UpdateById(ctx, user.Id, user.Birthday, user.Nickname, user.AboutMe)
	return err
}

// Profile 个人信息
func (userService *UserService) Profile(ctx *gin.Context, userId int64) (domain.User, error) {
	return userService.userRepo.FindById(ctx, userId)
}

// LoginOrSignup 登录或注册
func (userService *UserService) LoginOrSignup(ctx *gin.Context, phone string) (domain.User, error) {
	byPhone, err := userService.userRepo.FindByPhone(ctx, phone)
	if err != nil {
		if err.Error() == ErrUserNotFound {
			// 用户不存在, 注册
			user := domain.User{
				Phone: phone,
			}
			err := userService.userRepo.Create(ctx, user)
			if err != nil {
				// 注册失败, 可能是手机号已存在, 也可能是其他原因
				if err.Error() == ErrorUserExist {
					return domain.User{}, errors.New("手机号已存在")
				}
				return domain.User{}, err
			}
			// 注册成功后再次查询
			byPhone, err = userService.userRepo.FindByPhone(ctx, phone)
			if err != nil {
				return domain.User{}, err
			}
			return byPhone, nil
		}
		return domain.User{}, err
	}
	return byPhone, nil
}

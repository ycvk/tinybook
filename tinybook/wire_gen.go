// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"geek_homework/tinybook/internal/repository"
	"geek_homework/tinybook/internal/repository/cache"
	"geek_homework/tinybook/internal/repository/dao"
	"geek_homework/tinybook/internal/service"
	"geek_homework/tinybook/internal/web"
	"geek_homework/tinybook/internal/web/jwt"
	"geek_homework/tinybook/ioc"
	"github.com/gin-gonic/gin"
)

import (
	_ "github.com/spf13/viper/remote"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := ioc.InitRedis()
	handler := jwt.NewRedisJWTHandler(cmdable)
	logger := ioc.InitLogger()
	v := ioc.InitHandlerFunc(cmdable, handler, logger)
	db := ioc.InitDB()
	userDAO := dao.NewGormUserDAO(db)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository, logger)
	theineCache := ioc.InitLocalCache()
	codeCache := cache.NewLocalCodeCache(theineCache)
	codeRepository := repository.NewCachedCodeRepository(codeCache)
	smsdao := dao.NewGormSMSDAO(db)
	smsRepository := repository.NewGormSMSRepository(smsdao)
	smsService := ioc.InitSMSService(cmdable, smsRepository)
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService, handler)
	wechatService := ioc.InitWechatService()
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, userService, handler)
	engine := ioc.InitWebServer(v, userHandler, oAuth2WechatHandler)
	return engine
}

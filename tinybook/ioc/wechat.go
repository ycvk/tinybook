package ioc

import (
	"geek_homework/tinybook/internal/service/oauth2/wechat"
	"os"
)

func InitWechatService() wechat.Service {
	env, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("env WECHAT_APP_ID not found")
	}
	return wechat.NewService(env)
}

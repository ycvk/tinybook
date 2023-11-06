package ioc

import (
	"geek_homework/tinybook/internal/service/oauth2/wechat"
	"os"
)

func InitWechatService() wechat.Service {
	env, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		//panic("env WECHAT_APP_ID not found")
		env = "wx4f4bc4dec97d474b"
	}
	lookupEnv, b := os.LookupEnv("WECHAT_APP_SECRET")
	if !b {
		//panic("env WECHAT_APP_SECRET not found")
		lookupEnv = "wx4f4bc4dec97d474b"
	}
	return wechat.NewService(env, lookupEnv)
}

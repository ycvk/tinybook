package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"tinybook/tinybook/internal/service/sms"
)

type SMSService struct {
	svc sms.Service
	key []byte
}

type SMSClaims struct {
	jwt.RegisteredClaims
	tpl string //模板ID
}

func NewSMSService(svc sms.Service) *SMSService {
	return &SMSService{
		svc: svc,
	}
}

func (s SMSService) Send(ctx context.Context, tplToken string, args []string, numbers ...string) error {
	var claims SMSClaims
	// 解析token
	_, err := jwt.ParseWithClaims(tplToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	// 如果解析失败，那么就返回错误
	if err != nil {
		return err
	}
	// 如果解析成功，那么就发送短信
	return s.svc.Send(ctx, claims.tpl, args, numbers...)
}

package localsms

import (
	"context"
	"geek_homework/tinybook/internal/service/sms"
	"log"
)

type Service struct {
}

func NewService() sms.Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	log.Println("验证码是", args)
	return nil
}

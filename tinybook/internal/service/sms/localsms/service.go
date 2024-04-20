package localsms

import (
	"context"
	"log"
	"tinybook/tinybook/internal/service/sms"
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

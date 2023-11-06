package sms

import (
	"context"
)

// Service 短信服务
type Service interface {
	Send(ctx context.Context, tplId string, args []string, numbers ...string) error
}

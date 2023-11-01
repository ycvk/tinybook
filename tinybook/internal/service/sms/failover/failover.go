package failover

import (
	"context"
	"geek_homework/tinybook/internal/service/sms"
	"github.com/cockroachdb/errors"
	"log/slog"
)

var ErrAllServicesFailed = errors.New("all send sms services failed")

type FailoverSMSService struct {
	services []sms.Service
}

func NewFailoverSMSService(services ...sms.Service) *FailoverSMSService {
	return &FailoverSMSService{
		services: services,
	}
}

func (f FailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	for _, service := range f.services {
		err := service.Send(ctx, tplId, args, numbers...)
		if err == nil { // 如果没有错误，说明发送成功，直接返回
			return nil
		}
		// 如果是短信服务的错误，说明短信服务出错了，继续尝试下一个短信服务
		err = errors.AssertionFailedf("Send sms failed\n%+v", err)
		slog.Error(err.Error())
	}
	return ErrAllServicesFailed
}

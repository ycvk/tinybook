package tencent

import (
	"context"
	"github.com/samber/lo"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
	"log/slog"
)

type Service struct {
	client   *sms.Client
	appId    *string
	signName *string
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	request := sms.NewSendSmsRequest()
	request.SetContext(ctx)
	request.SmsSdkAppId = s.appId
	request.SignName = s.signName
	request.TemplateId = &tplId
	request.TemplateParamSet = s.toPtrSlice(args)
	request.PhoneNumberSet = s.toPtrSlice(numbers)
	response, err := s.client.SendSms(request)
	if err != nil {
		slog.Error("发送短信失败", "error", err)
	}
	statuses := response.Response.SendStatusSet
	for _, status := range statuses {
		if *status.Code != "Ok" {
			slog.Error("发送短信失败", "code", status.Code, "error", status.Message)
			return err
		}
	}
	return nil
}

// toPtrSlice 将字符串切片转换为指针切片
func (s *Service) toPtrSlice(data []string) []*string {
	return lo.Map(data, func(src string, idx int) *string {
		return &src
	})
}

func NewService(client *sms.Client, appId string, signName string) *Service {
	return &Service{
		client:   client,
		appId:    &appId,
		signName: &signName,
	}
}

package wechat

import (
	"context"
	"fmt"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/url"
)

var redirectURI = url.PathEscape("https://tinybook.ycvk.app/oauth2/wechat/callback")

type Service interface {
	AuthURL(ctx context.Context) (string, error)
}

type service struct {
	appId string
}

func NewService(appId string) Service {
	return &service{
		appId: appId,
	}
}

func (s *service) AuthURL(ctx context.Context) (string, error) {
	u := uuid.New()
	const authURL = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	return fmt.Sprintf(authURL, s.appId, redirectURI, u), nil
}

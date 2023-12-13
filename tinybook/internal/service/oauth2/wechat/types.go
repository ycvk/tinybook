package wechat

import (
	"context"
	"fmt"
	"geek_homework/tinybook/internal/domain"
	"github.com/bytedance/sonic"
	"net/http"
	"net/url"
)

var redirectURI = url.PathEscape("https://tinybook.ycvk.app/oauth2/wechat/callback")

type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	Verify(ctx context.Context, code string) (domain.WechatInfo, error)
}

type service struct {
	appId     string
	appSecret string
	client    *http.Client
}

type wechatAccessToken struct {
	AccessToken  string `json:"access_token"`  // 接口调用凭证
	ExpiresIn    int64  `json:"expires_in"`    // access_token接口调用凭证超时时间，单位（秒）
	RefreshToken string `json:"refresh_token"` // 用户刷新access_token
	OpenId       string `json:"openid"`        // 授权用户唯一标识
	Scope        string `json:"scope"`         // 用户授权的作用域，使用逗号（,）分隔
	Unionid      string `json:"unionid"`       // 仅在该应用已获得该用户的userinfo授权时，才会出现该字段。

	// 错误码
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func NewService(id string, secret string) Service {
	return &service{
		appId:     id,
		appSecret: secret,
		client:    http.DefaultClient,
	}
}

func (s *service) Verify(ctx context.Context, code string) (domain.WechatInfo, error) {
	accessTokenUrl := fmt.Sprintf(`https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code`, s.appId, s.appSecret, code)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, accessTokenUrl, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	response, err := s.client.Do(request)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	defer response.Body.Close()
	var accessToken wechatAccessToken
	// 解析json
	err = sonic.ConfigDefault.NewDecoder(response.Body).Decode(&accessToken)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	if accessToken.ErrCode != 0 {
		return domain.WechatInfo{},
			fmt.Errorf("获取access_token失败，错误码：%d，错误信息：%s", accessToken.ErrCode, accessToken.ErrMsg)
	}
	return domain.WechatInfo{
		OpenId:  accessToken.OpenId,
		UnionId: accessToken.Unionid,
	}, nil
}

func (s *service) AuthURL(ctx context.Context, state string) (string, error) {
	const authURL = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	return fmt.Sprintf(authURL, s.appId, redirectURI, state), nil
}

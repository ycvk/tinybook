package web

import (
	"geek_homework/tinybook/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OAuth2WechatHandler struct {
	service wechat.Service
}

func NewOAuth2WechatHandler(service wechat.Service) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		service: service,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(engine *gin.Engine) {
	engine.Group("/oauth2/wechat")
	engine.GET("/authurl", h.Auth2URL)
	engine.Any("/callback", h.Callback)

}

func (h *OAuth2WechatHandler) Auth2URL(context *gin.Context) {
	url, err := h.service.AuthURL(context)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "构造跳转链接失败",
		})
		return
	}
	context.JSON(200, Result{
		Code: 200,
		Data: url,
	})
}

func (h *OAuth2WechatHandler) Callback(context *gin.Context) {

}

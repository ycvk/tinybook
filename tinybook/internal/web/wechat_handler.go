package web

import (
	"geek_homework/tinybook/internal/service"
	"geek_homework/tinybook/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OAuth2WechatHandler struct {
	JWTHandler
	wechatService wechat.Service
	userService   service.UserService
}

func NewOAuth2WechatHandler(service wechat.Service, userService service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		wechatService: service,
		userService:   userService,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(engine *gin.Engine) {
	group := engine.Group("/oauth2/wechat")
	group.GET("/authurl", h.Auth2URL)
	group.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) Auth2URL(context *gin.Context) {
	url, err := h.wechatService.AuthURL(context)
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
	code := context.Query("code")
	wechatInfo, err := h.wechatService.Verify(context, code)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "授权码验证失败",
		})
		return
	}
	byWechat, err := h.userService.LoginOrSignupByWechat(context, wechatInfo)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "登录失败",
		})
		return
	}
	jwtToken, err := h.GetJWTToken(context, byWechat)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "登录失败",
		})
		return
	}
	context.Header("X-Jwt-Token", jwtToken)
	context.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "登录成功",
	})
}

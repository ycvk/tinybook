package web

import (
	"geek_homework/tinybook/internal/service"
	"geek_homework/tinybook/internal/service/oauth2/wechat"
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
)

var ErrInvalidState = errors.New("state不匹配")

type OAuth2WechatHandler struct {
	jwtHandler      *JWTHandler
	wechatService   wechat.Service
	userService     service.UserService
	stateCookieName string
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string `json:"state"`
}

func NewOAuth2WechatHandler(service wechat.Service, userService service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		wechatService:   service,
		userService:     userService,
		jwtHandler:      NewJWTHandler(),
		stateCookieName: "jwt-state",
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(engine *gin.Engine) {
	group := engine.Group("/oauth2/wechat")
	group.GET("/authurl", h.Auth2URL)
	group.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) Auth2URL(context *gin.Context) {
	state := uuid.New()
	url, err := h.wechatService.AuthURL(context, state)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "构造跳转链接失败",
		})
		return
	}
	stateErr := h.SetStateCookie(context, state)
	if stateErr != nil {
		context.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "系统错误",
		})
		return
	}
	context.JSON(200, Result{
		Code: 200,
		Data: url,
	})
}

func (h *OAuth2WechatHandler) Callback(context *gin.Context) {
	verifyErr := h.VerifyStateCookie(context)
	if verifyErr != nil {
		context.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "非法请求",
		})
		return
	}
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
	verifyErr = h.jwtHandler.SetJWTToken(context, byWechat)
	if verifyErr != nil {
		context.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "登录失败",
		})
		return
	}
	context.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "登录成功",
	})
}

func (h *OAuth2WechatHandler) SetStateCookie(context *gin.Context, state string) error {
	claims := StateClaims{
		State: state,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(JWTKey))
	if err != nil {
		return err
	}
	context.SetCookie(h.stateCookieName, signedString, 600,
		"/oauth2/wechat/callback", // 限制只能在/oauth2/wechat/callback下访问
		"", false, true)
	return nil
}

func (h *OAuth2WechatHandler) VerifyStateCookie(ctx *gin.Context) error {
	state := ctx.Query("state")
	cookie, err := ctx.Cookie(h.stateCookieName)
	if err != nil {
		return err
	}
	claims := StateClaims{}
	_, err = jwt.ParseWithClaims(cookie, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTKey), nil
	})
	if err != nil {
		return err
	}
	if claims.State != state {
		return ErrInvalidState
	}
	return nil
}

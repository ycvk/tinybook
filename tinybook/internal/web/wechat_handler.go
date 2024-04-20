package web

import (
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"tinybook/tinybook/internal/service"
	"tinybook/tinybook/internal/service/oauth2/wechat"
	jwt2 "tinybook/tinybook/internal/web/jwt"
)

const JWTKey = "MK7z43qKmUkY5sy9w3rQ8CygFpOSN90W"

var ErrInvalidState = errors.New("state不匹配")

type OAuth2WechatHandler struct {
	jwtHandler      jwt2.Handler
	wechatService   wechat.Service
	userService     service.UserService
	stateCookieName string
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string `json:"state"`
}

func NewOAuth2WechatHandler(service wechat.Service, userService service.UserService, handler jwt2.Handler) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		wechatService:   service,
		userService:     userService,
		jwtHandler:      handler,
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
	// 验证state 防止csrf攻击
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
	// 调用service层的LoginOrSignup方法
	byWechat, err := h.userService.LoginOrSignupByWechat(context, wechatInfo)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "登录失败",
		})
		return
	}
	// 生成jwt 包含refresh token和jwt token
	verifyErr = h.jwtHandler.SetLoginToken(context, byWechat)
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

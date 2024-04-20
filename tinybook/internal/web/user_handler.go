package web

import (
	"geek_homework/tinybook/internal/domain"
	"geek_homework/tinybook/internal/service"
	jwt2 "geek_homework/tinybook/internal/web/jwt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

var ErrUserNotFound = service.ErrUserNotFound

const (
	bizLogin = "login"
)

type UserHandler struct {
	jwtHandler  jwt2.Handler
	userService service.UserService
	codeService service.CodeService
}

func NewUserHandler(userService service.UserService, codeService service.CodeService, handler jwt2.Handler) *UserHandler {
	return &UserHandler{
		userService: userService,
		codeService: codeService,
		jwtHandler:  handler,
	}
}

// SignUp 注册
func (userHandler *UserHandler) SignUp(ctx *gin.Context) {
	type Sign struct {
		Password        string `json:"password"`
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var sign Sign
	if err := ctx.Bind(&sign); err != nil {
		ctx.JSON(http.StatusOK, "格式不正确")
		ctx.Error(err)
		return
	}
	if strings.Compare(sign.Password, sign.ConfirmPassword) != 0 {
		ctx.JSON(http.StatusOK, "两次密码不一致")
		return
	}
	// 调用service层的Signup方法
	err := userHandler.userService.Signup(ctx.Request.Context(), domain.User{
		Email:    sign.Email,
		Password: sign.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, err.Error())
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, "注册成功")
}

// LoginJWT 登录
func (userHandler *UserHandler) LoginJWT(ctx *gin.Context) {
	type Login struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	var login Login
	if err := ctx.Bind(&login); err != nil {
		ctx.JSON(http.StatusOK, "格式不正确")
		ctx.Error(err)
		return
	}
	// 调用service层的Login方法
	user, err := userHandler.userService.Login(ctx, login.Email, login.Password)
	if err != nil {
		if err.Error() == ErrUserNotFound {
			ctx.JSON(http.StatusOK, "用户不存在")
			return
		}
		ctx.JSON(http.StatusOK, "密码不正确")
		return
	}
	err = userHandler.jwtHandler.SetLoginToken(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusOK, "登录失败")
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, "登录成功")
}

// Login 登录
func (userHandler *UserHandler) Login(ctx *gin.Context) {
	type Login struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	var login Login
	if err := ctx.Bind(&login); err != nil {
		ctx.JSON(http.StatusOK, "格式不正确")
		return
	}
	// 调用service层的Login方法
	user, err := userHandler.userService.Login(ctx, login.Email, login.Password)
	if err != nil {
		if err.Error() == ErrUserNotFound {
			ctx.JSON(http.StatusOK, "用户不存在")
			return
		}
		ctx.JSON(http.StatusOK, "密码不正确")
		ctx.Error(err)
		return
	}
	// 设置session 保存用户id 有效时间1小时
	session := sessions.Default(ctx)
	session.Set("userId", user.Id)
	session.Options(sessions.Options{
		MaxAge: 60 * 60 * 1,
	})
	err = session.Save()
	if err != nil {
		ctx.JSON(http.StatusOK, "登录失败")
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, "登录成功")
}

// Edit 编辑
func (userHandler *UserHandler) Edit(ctx *gin.Context) {
	type Edit struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var edit Edit
	if err := ctx.Bind(&edit); err != nil {
		ctx.JSON(http.StatusOK, "格式不正确")
		ctx.Error(err)
		return
	}
	// 获取session中的userId
	//session := sessions.Default(ctx)
	//userId := session.Get("userId")
	//if userId == nil {
	//	ctx.JSON(http.StatusOK, gin.H{"msg": "用户未登录，请先登录"})
	//	return
	//}
	// 使用jwt token获取userId
	claims, ok := (ctx.MustGet("userClaims")).(jwt2.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": "用户未登录, 请先登录"})
		return
	}

	// 调用service层的Edit方法
	err := userHandler.userService.Edit(ctx, domain.User{
		Id:       claims.Uid,
		Nickname: edit.Nickname,
		Birthday: edit.Birthday,
		AboutMe:  edit.AboutMe,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": err.Error()})
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "编辑成功"})
}

// Profile 获取个人信息
func (userHandler *UserHandler) Profile(ctx *gin.Context) {
	// 获取session中的userId
	//session := sessions.Default(ctx)
	//userId := session.Get("userId")
	// 使用jwt token获取userId
	claims, ok := (ctx.MustGet("userClaims")).(jwt2.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": "用户未登录, 请先登录"})
		return
	}
	// 调用service层的Profile方法
	user, err := userHandler.userService.Profile(ctx, claims.Uid)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": err.Error()})
		ctx.Error(err)
		return
	}

	// 返回个人信息 我发现前端页面设计时没有设计手机号码字段，所以这里写死
	ctx.JSON(200, gin.H{
		"Email":    user.Email,
		"Phone":    "18011111111",
		"Nickname": user.Nickname,
		"Birthday": user.Birthday,
		"AboutMe":  user.AboutMe,
	})
}

// SendSMSLoginCode 发送登录验证码
func (userHandler *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	type Send struct {
		Phone string `json:"phone"`
	}
	var send Send
	if err := ctx.Bind(&send); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  "格式不正确",
		})
		return
	}
	// 调用service层的Send方法
	err := userHandler.codeService.Send(ctx, bizLogin, send.Phone, "600")
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  err.Error(),
		})
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "发送成功",
	})
}

// LoginSMS 短信登录
func (userHandler *UserHandler) LoginSMS(ctx *gin.Context) {
	type Login struct {
		Code  string `json:"code"`
		Phone string `json:"phone"`
	}
	var login Login
	if err := ctx.Bind(&login); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  "格式不正确",
		})
		return
	}
	// 调用service层的Verify方法
	verify, err := userHandler.codeService.Verify(ctx, bizLogin, login.Phone, login.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  err.Error(),
		})
		ctx.Error(err)
		return
	}
	if !verify {
		ctx.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  "验证码错误, 请重试",
		})
		return
	}
	// 调用service层的LoginOrSignup方法
	user, err := userHandler.userService.LoginOrSignupByPhone(ctx, login.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  err.Error(),
		})
		ctx.Error(err)
		return
	}
	// 生成与设置jwt token
	err = userHandler.jwtHandler.SetLoginToken(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  "登录失败",
		})
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "登录成功",
	})
}

// RefreshToken 刷新token 通过refresh token 刷新jwt token
func (userHandler *UserHandler) RefreshToken(ctx *gin.Context) {
	// 从header中获取refresh token
	authorization := userHandler.jwtHandler.ExtractAuthorization(ctx)
	var refreshClaims jwt2.RefreshClaims
	token, err := jwt.ParseWithClaims(authorization, &refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTKey), nil
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		ctx.Error(err)
		return
	}
	// 判断refresh token是否存在于redis
	err = userHandler.jwtHandler.CheckToken(ctx, refreshClaims.Ssid)
	if err != nil {
		// refresh token存在于redis 或者redis崩溃了
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if token == nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 生成新的jwt token 和 refresh token
	err = userHandler.jwtHandler.SetLoginToken(ctx, domain.User{
		Id: refreshClaims.Uid,
	})
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "ok",
	})
}

func (userHandler *UserHandler) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Options(sessions.Options{
		MaxAge: -1,
	})
	err := session.Save()
	if err != nil {
		ctx.Error(err)
		return
	}
}

func (userHandler *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := userHandler.jwtHandler.DeregisterToken(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "ok",
	})
}

// RegisterRoutes 注册路由
func (userHandler *UserHandler) RegisterRoutes(engine *gin.Engine) {
	group := engine.Group("/users")

	group.POST("/signup", userHandler.SignUp)
	//group.POST("/login", userHandler.Login)
	group.POST("/login", userHandler.LoginJWT)
	group.POST("/edit", userHandler.Edit)
	group.GET("/profile", userHandler.Profile)
	group.GET("/refresh_token", userHandler.RefreshToken)
	group.GET("/logout", userHandler.Logout)
	group.POST("/login_sms/code/send", userHandler.SendSMSLoginCode)
	group.POST("/login_sms", userHandler.LoginSMS)
}

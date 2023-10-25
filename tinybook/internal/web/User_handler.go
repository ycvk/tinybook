package web

import (
	"geek_homework/tinybook/internal/domain"
	"geek_homework/tinybook/internal/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

var ErrUserNotFound = service.ErrUserNotFound

const (
	JWTKey   = "MK7z43qKmUkY5sy9w3rQ8CygFpOSN90W"
	bizLogin = "login"
)

type UserHandler struct {
	userService service.UserService
	codeService service.CodeService
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64  `json:"uid"`
	UserAgent string `json:"userAgent"`
}

func NewUserHandler(userService service.UserService, codeService service.CodeService) *UserHandler {
	return &UserHandler{
		userService: userService,
		codeService: codeService,
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
		return
	}
	if strings.Compare(sign.Password, sign.ConfirmPassword) != 0 {
		ctx.JSON(http.StatusOK, "两次密码不一致")
		return
	}
	// 调用service层的Signup方法
	err := userHandler.userService.Signup(ctx, domain.User{
		Email:    sign.Email,
		Password: sign.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, err.Error())
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

	tokenStr, jwtErr := userHandler.GetJWTToken(ctx, user)
	if jwtErr != nil {
		ctx.JSON(http.StatusOK, "系统错误")
		return
	}
	ctx.Header("X-Jwt-Token", tokenStr)
	ctx.JSON(http.StatusOK, "登录成功")
}

// GetJWTToken 获取jwt token
func (userHandler *UserHandler) GetJWTToken(ctx *gin.Context, user domain.User) (string, error) {
	userClaims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
		},
		Uid:       user.Id,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaims) //生成token
	tokenStr, err := token.SignedString([]byte(JWTKey))
	return tokenStr, err
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
	claims, ok := (ctx.MustGet("userClaims")).(UserClaims)
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
	claims, ok := (ctx.MustGet("userClaims")).(UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": "用户未登录, 请先登录"})
		return
	}
	// 调用service层的Profile方法
	user, err := userHandler.userService.Profile(ctx, claims.Uid)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": err.Error()})
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
	user, err := userHandler.userService.LoginOrSignup(ctx, login.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}
	// 生成jwt token
	tokenStr, jwtErr := userHandler.GetJWTToken(ctx, user)
	if jwtErr != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  "系统错误",
		})
		return
	}
	ctx.Header("X-Jwt-Token", tokenStr)
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "登录成功",
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
	group.POST("/login_sms/code/send", userHandler.SendSMSLoginCode)
	group.POST("/login_sms", userHandler.LoginSMS)
}

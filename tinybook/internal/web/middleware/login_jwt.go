package middleware

import (
	"geek_homework/tinybook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
	"net/http"
)

type LoginJWTMiddlewareBuilder struct {
	jwtHandler web.JWTHandler
}

// Build 构建登录中间件
func (builder *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	pathList := make([]string, 10)
	pathList = append(pathList,
		"/users/login", "/users/signup",
		"/users/login_sms/code/send", "/users/login_sms",
		"/oauth2/wechat/authurl", "/oauth2/wechat/callback")
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		// 登录和注册不需要经过登录中间件, 可以在未经认证的情况下访问
		if lo.Contains(pathList, path) {
			return
		}
		// 从header中提取jwt token
		jwtToken := builder.jwtHandler.ExtractAuthorization(ctx)
		var claims web.UserClaims
		withClaims, err := jwt.ParseWithClaims(jwtToken, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(web.JWTKey), nil
		})
		if err != nil {
			// 解析失败 token过期或者token不合法
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !withClaims.Valid {
			// token不合法
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			// ua不一致 正常用户不会进入此分支 说明token被盗用
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//expiresAt := claims.ExpiresAt              // 过期时间
		//if expiresAt.Sub(time.Now()) < time.Hour { // 过期时间小于1小时 刷新token
		//	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour)) // 设置过期时间为1小时后
		//	tokenStr, err := withClaims.SignedString([]byte(web.JWTKey))     // 重新签名
		//	if err != nil {
		//		// 刷新失败
		//		slog.Info("刷新token失败", "err", err)
		//	}
		//	ctx.Header("X-Jwt-Token", tokenStr)
		//}
		ctx.Set("userClaims", claims) // 设置userId到上下文
	}
}

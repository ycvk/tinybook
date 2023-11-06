package web

import (
	"geek_homework/tinybook/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const JWTKey = "MK7z43qKmUkY5sy9w3rQ8CygFpOSN90W"

type JWTHandler struct {
}

// GetJWTToken 获取jwt token
func (jwtHandler *JWTHandler) GetJWTToken(ctx *gin.Context, user domain.User) (string, error) {
	userClaims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)), //过期时间
		},
		Uid:       user.Id,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaims) //生成token
	tokenStr, err := token.SignedString([]byte(JWTKey))            //签名
	return tokenStr, err
}

// SetJWTToken 设置jwt token
func (jwtHandler *JWTHandler) SetJWTToken(ctx *gin.Context, token string) {
	ctx.Header("X-Jwt-Token", token)
}

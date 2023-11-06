package web

import (
	"geek_homework/tinybook/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"strings"
	"time"
)

const JWTKey = "MK7z43qKmUkY5sy9w3rQ8CygFpOSN90W"

type JWTHandler struct {
	signMethod jwt.SigningMethod
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid int64 `json:"uid"`
}

func NewJWTHandler() *JWTHandler {
	return &JWTHandler{
		signMethod: jwt.SigningMethodHS512,
	}
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
	token := jwt.NewWithClaims(jwtHandler.signMethod, userClaims) //生成token
	tokenStr, err := token.SignedString([]byte(JWTKey))           //签名
	return tokenStr, err
}

// SetJWTToken 设置jwt token
func (jwtHandler *JWTHandler) SetJWTToken(ctx *gin.Context, user domain.User) error {
	// 设置refresh token
	err := jwtHandler.SetRefreshToken(ctx, user.Id)
	if err != nil {
		slog.Error("设置refresh token失败", "err", err)
		return err
	}
	tokenStr, jwtErr := jwtHandler.GetJWTToken(ctx, user)
	if jwtErr != nil {
		return jwtErr
	}
	ctx.Header("X-Jwt-Token", tokenStr)
	return nil
}

func (jwtHandler *JWTHandler) SetRefreshToken(ctx *gin.Context, uid int64) error {
	refreshClaims := RefreshClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), //过期时间
		},
	}
	refreshToken := jwt.NewWithClaims(jwtHandler.signMethod, refreshClaims) //生成token
	signedString, err := refreshToken.SignedString([]byte(JWTKey))          //签名
	if err != nil {
		return err
	}
	ctx.Header("X-Refresh-Token", signedString)
	return nil
}

// ExtractAuthorization 从header中提取jwt token
func (jwtHandler *JWTHandler) ExtractAuthorization(ctx *gin.Context) string {
	// 从header中获取jwt token
	authCode := ctx.GetHeader("Authorization")
	if jwtToken := strings.Split(authCode, " "); len(jwtToken) != 2 {
		return ""
	} else {
		return jwtToken[1]
	}
}

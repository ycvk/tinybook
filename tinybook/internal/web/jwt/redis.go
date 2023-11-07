package jwt

import (
	"fmt"
	"geek_homework/tinybook/internal/domain"
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"strings"
	"time"
)

const JWTKey = "MK7z43qKmUkY5sy9w3rQ8CygFpOSN90W"

type RedisJWTHandler struct {
	signMethod jwt.SigningMethod
	cmd        redis.Cmdable
	expire     time.Duration
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid  int64  `json:"uid"`
	Ssid string `json:"ssid"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64  `json:"uid"`
	Ssid      string `json:"ssid"`
	UserAgent string `json:"userAgent"`
}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		signMethod: jwt.SigningMethodHS512,
		expire:     time.Hour * 24 * 7,
		cmd:        cmd,
	}
}

func (jwtHandler *RedisJWTHandler) SetLoginToken(ctx *gin.Context, user domain.User) error {
	// 生成ssid
	ssid := uuid.New().String()
	// 设置refresh token
	err := jwtHandler.SetRefreshToken(ctx, user.Id, ssid)
	if err != nil {
		slog.Error("设置refresh token失败", "err", err)
		return err
	}
	// 设置jwt token
	err = jwtHandler.SetJWTToken(ctx, ssid, user)
	if err != nil {
		slog.Error("设置jwt token失败", "err", err)
		return err
	}
	return nil
}

// GetJWTToken 获取jwt token
func (jwtHandler *RedisJWTHandler) getJWTToken(ctx *gin.Context, ssid string, user domain.User) (string, error) {
	userClaims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)), //过期时间
		},
		Uid:       user.Id,
		Ssid:      ssid,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwtHandler.signMethod, userClaims) //生成token
	tokenStr, err := token.SignedString([]byte(JWTKey))           //签名
	return tokenStr, err
}

// SetJWTToken 设置jwt token
func (jwtHandler *RedisJWTHandler) SetJWTToken(ctx *gin.Context, ssid string, user domain.User) error {
	// 设置refresh token
	err := jwtHandler.SetRefreshToken(ctx, user.Id, ssid)
	if err != nil {
		slog.Error("设置refresh token失败", "err", err)
		return err
	}
	tokenStr, jwtErr := jwtHandler.getJWTToken(ctx, ssid, user)
	if jwtErr != nil {
		return jwtErr
	}
	ctx.Header("X-Jwt-Token", tokenStr)
	return nil
}

func (jwtHandler *RedisJWTHandler) SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	refreshClaims := RefreshClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtHandler.expire)), //过期时间
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
func (jwtHandler *RedisJWTHandler) ExtractAuthorization(ctx *gin.Context) string {
	// 从header中获取jwt token
	authCode := ctx.GetHeader("Authorization")
	if jwtToken := strings.Split(authCode, " "); len(jwtToken) != 2 {
		return ""
	} else {
		return jwtToken[1]
	}
}

func (jwtHandler *RedisJWTHandler) DeregisterToken(ctx *gin.Context) error {
	ctx.Header("X-Jwt-Token", "")
	ctx.Header("X-Refresh-Token", "")
	claims := ctx.MustGet("userClaims").(UserClaims)
	//把目前的refresh token添加到redis，这样就说明此token已经失效了
	return jwtHandler.cmd.Set(ctx, fmt.Sprintf("refresh_token:ssid:%s", claims.Ssid), "", jwtHandler.expire).Err()
}

func (jwtHandler *RedisJWTHandler) CheckToken(ctx *gin.Context, ssid string) error {
	result, err := jwtHandler.cmd.Exists(ctx, fmt.Sprintf("refresh_token:ssid:%s", ssid)).Result()
	if err != nil || result > 0 {
		// refresh token存在于redis 或者redis崩溃了
		//ctx.AbortWithStatus(http.StatusUnauthorized)
		return errors.New("不存在的token")
	}
	return nil
}

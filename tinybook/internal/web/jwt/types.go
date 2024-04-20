package jwt

import (
	"github.com/gin-gonic/gin"
	"tinybook/tinybook/internal/domain"
)

type Handler interface {
	SetLoginToken(ctx *gin.Context, user domain.User) error
	SetJWTToken(ctx *gin.Context, ssid string, user domain.User) error
	SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error
	ExtractAuthorization(ctx *gin.Context) string
	DeregisterToken(ctx *gin.Context) error
	CheckToken(ctx *gin.Context, ssid string) error
}

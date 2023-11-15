package web

import (
	"geek_homework/tinybook/internal/domain"
	"geek_homework/tinybook/internal/service"
	"geek_homework/tinybook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type ArticleHandler struct {
	service service.ArticleService
	l       *zap.Logger
}

func NewArticleHandler(service service.ArticleService, l *zap.Logger) *ArticleHandler {
	return &ArticleHandler{service: service, l: l}
}

func (h *ArticleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  "参数错误",
		})
		return
	}
	claims := (ctx.MustGet("userClaims")).(jwt.UserClaims)
	id, err := h.service.Save(ctx, domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Author:  domain.Author{ID: claims.Uid},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "服务器错误",
		})
		h.l.Error("保存文章失败, 作者ID: "+strconv.FormatInt(claims.Uid, 10), zap.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "保存成功",
		Data: id,
	})
}

func (h *ArticleHandler) RegisterRoutes(engine *gin.Engine) {
	group := engine.Group("/articles")
	group.POST("/edit", h.Edit)
}

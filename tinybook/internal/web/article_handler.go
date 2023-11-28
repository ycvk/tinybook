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
		Id      int64  `json:"id"`
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
		ID:      req.Id,
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

func (h *ArticleHandler) Publish(ctx *gin.Context) {
	type Req struct {
		Id      int64  `json:"id"`
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
	id, err := h.service.Publish(ctx, domain.Article{
		ID:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author:  domain.Author{ID: claims.Uid},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "服务器错误",
		})
		h.l.Error("发表文章失败, 作者ID: "+strconv.FormatInt(claims.Uid, 10), zap.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "发表成功",
		Data: id,
	})
}

func (h *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64 `json:"id"`
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
	err := h.service.Withdraw(ctx, domain.Article{
		ID: req.Id,
		Author: domain.Author{
			ID: claims.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "服务器错误",
		})
		h.l.Error("撤回文章失败, 作者ID: "+
			strconv.FormatInt(claims.Uid, 10)+
			" 文章ID: "+strconv.FormatInt(req.Id, 10),
			zap.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "撤回成功",
		Data: req.Id,
	})
}

func (h *ArticleHandler) Detail(context *gin.Context) {
	param := context.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  "参数错误",
		})
		return
	}
	claims := (context.MustGet("userClaims")).(jwt.UserClaims)
	article, err := h.service.GetArticleById(context, id)
	if strconv.FormatInt(claims.Uid, 10) != article.Author {
		context.JSON(http.StatusOK, Result{
			Code: 401,
			Msg:  "无权限",
		})
	}
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "服务器错误",
		})
		h.l.Error("获取文章详情失败, 文章ID: "+strconv.FormatInt(id, 10), zap.Error(err))
		return
	}
	context.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "获取成功",
		Data: article,
	})
}

func (h *ArticleHandler) List(context *gin.Context) {
	var page Page
	if err := context.Bind(&page); err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  "参数错误",
		})
		return
	}
	claims := (context.MustGet("userClaims")).(jwt.UserClaims)
	articles, err := h.service.GetArticlesByAuthor(context, claims.Uid, page.Limit, page.Offset)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "服务器错误",
		})
		h.l.Error("获取文章列表失败, 作者ID: "+strconv.FormatInt(claims.Uid, 10), zap.Error(err))
		return
	}
	context.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "获取成功",
		Data: articles,
	})
}

func (h *ArticleHandler) PubDetail(context *gin.Context) {
	param := context.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  "参数错误",
		})
		return
	}
	article, err := h.service.GetPubArticleById(context, id)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "服务器错误",
		})
		h.l.Error("读者获取文章详情失败, 文章ID: "+strconv.FormatInt(id, 10), zap.Error(err))
		return
	}
	context.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "获取成功",
		Data: article,
	})
}

func (h *ArticleHandler) RegisterRoutes(engine *gin.Engine) {
	group := engine.Group("/articles")
	group.POST("/edit", h.Edit)               // 编辑文章
	group.POST("/publish", h.Publish)         // 发表文章
	group.POST("/withdraw", h.Withdraw)       // 撤回文章
	group.POST("/list", h.List)               // 文章列表
	group.GET("/detail/:id", h.Detail)        // 文章详情
	group.GET("/pub/detail/:id", h.PubDetail) // 读者查看文章详情
}

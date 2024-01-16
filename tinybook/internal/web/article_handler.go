package web

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	intrv1 "tinybook/tinybook/api/proto/gen/intr/v1"
	"tinybook/tinybook/internal/domain"
	"tinybook/tinybook/internal/service"
	"tinybook/tinybook/internal/web/jwt"
)

type ArticleHandler struct {
	articleService     service.ArticleService
	interactiveService intrv1.InteractiveServiceClient
	l                  *zap.Logger
	biz                string
}

func NewArticleHandler(artService service.ArticleService, interService intrv1.InteractiveServiceClient, l *zap.Logger) *ArticleHandler {
	return &ArticleHandler{
		articleService:     artService,
		interactiveService: interService,
		l:                  l,
		biz:                "article",
	}
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
	id, err := h.articleService.Save(ctx, domain.Article{
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
	id, err := h.articleService.Publish(ctx, domain.Article{
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
	err := h.articleService.Withdraw(ctx, domain.Article{
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
	article, err := h.articleService.GetArticleById(context, id)
	if strconv.FormatInt(claims.Uid, 10) != article.Author {
		context.JSON(http.StatusOK, Result{
			Code: 401,
			Msg:  "无权限",
		})
		if err != nil {
			h.l.Error("获取文章详情失败, 文章ID: "+strconv.FormatInt(id, 10), zap.Error(err))
		}
		return
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
	articles, err := h.articleService.GetArticlesByAuthor(context, claims.Uid, page.Limit, page.Offset)
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

	var (
		eg          errgroup.Group
		article     domain.ArticleVo
		interactive *intrv1.GetInteractiveResponse
	)
	claims := (context.MustGet("userClaims")).(jwt.UserClaims)

	eg.Go(func() error {
		var articleErr error
		article, articleErr = h.articleService.GetPubArticleById(context, id, claims.Uid)
		return articleErr
	})

	eg.Go(func() error {
		var interactiveErr error
		claims := (context.MustGet("userClaims")).(jwt.UserClaims)
		interactive, interactiveErr = h.interactiveService.GetInteractive(context, &intrv1.GetInteractiveRequest{
			Biz:   h.biz,
			BizId: id,
			Uid:   claims.Uid,
		})
		return interactiveErr
	})

	err = eg.Wait()
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "服务器错误",
		})
		h.l.Error("读者获取文章详情失败, 文章ID: "+strconv.FormatInt(id, 10), zap.Error(err))
		return
	}
	// 增加阅读数 不在这里处理了，放到kafka处理了
	//go func() {
	//	err = h.interactiveService.IncreaseReadCount(context, h.biz, id)
	//	if err != nil {
	//		h.l.Warn("文章阅读数增加失败, 文章ID: "+strconv.FormatInt(id, 10), zap.Error(err))
	//	}
	//}()

	article.BizId = interactive.Interactive.BizId
	article.Biz = interactive.Interactive.Biz
	article.ReadCount = interactive.Interactive.ReadCount
	article.LikeCount = interactive.Interactive.LikeCount
	article.CollectCount = interactive.Interactive.CollectCount
	article.Liked = interactive.Interactive.Liked
	article.Collected = interactive.Interactive.Collected

	context.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "获取成功",
		Data: article,
	})
}

func (h *ArticleHandler) Like(context *gin.Context) {
	type Req struct {
		Id   int64 `json:"id"`
		Like bool  `json:"like"`
	}
	var req Req
	if err := context.Bind(&req); err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  "参数错误",
		})
		return
	}
	claims := (context.MustGet("userClaims")).(jwt.UserClaims)
	var err error
	if req.Like {
		_, err = h.interactiveService.Like(context, &intrv1.LikeRequest{
			Biz:   h.biz,
			BizId: req.Id,
			Uid:   claims.Uid,
		})
	} else {
		_, err = h.interactiveService.Unlike(context, &intrv1.UnlikeRequest{
			Biz:   h.biz,
			BizId: req.Id,
			Uid:   claims.Uid,
		})
	}
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "服务器错误",
		})
		h.l.Error("点赞失败, 文章ID: "+strconv.FormatInt(req.Id, 10)+" 用户ID: "+strconv.FormatInt(claims.Uid, 10), zap.Error(err))
		return
	}
	var msg string
	if req.Like {
		msg = "点赞成功"
	} else {
		msg = "取消点赞成功"
	}
	context.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  msg,
	})
}

func (h *ArticleHandler) Collect(ctx *gin.Context) {
	type Req struct {
		Id  int64 `json:"id"`
		Cid int64 `json:"cid"`
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
	_, err := h.interactiveService.Collect(ctx, &intrv1.CollectRequest{
		Biz:   h.biz,
		BizId: req.Id,
		Cid:   req.Cid,
		Uid:   claims.Uid,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "服务器错误",
		})
		h.l.Error("收藏失败, 文章ID: "+strconv.FormatInt(req.Id, 10)+" 用户ID: "+strconv.FormatInt(claims.Uid, 10), zap.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "收藏成功",
	})
}

func (h *ArticleHandler) Rank(context *gin.Context) {
	param := context.Param("id")
	num, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  "参数错误",
		})
		return
	}
	if num <= 0 || num > 100 {
		context.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  "非法参数",
		})
		return
	}
	ranks, err := h.interactiveService.GetLikeRanks(context, &intrv1.GetLikeRanksRequest{
		Biz: h.biz,
		Num: num,
	})
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "服务器错误",
		})
		h.l.Error("获取点赞排行榜失败, 数量: "+strconv.FormatInt(num, 10), zap.Error(err))
		return
	}
	context.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "获取成功",
		Data: ranks.GetArticles(),
	})
}

func (h *ArticleHandler) RegisterRoutes(engine *gin.Engine) {
	group := engine.Group("/articles")
	group.POST("/edit", h.Edit)         // 编辑文章
	group.POST("/publish", h.Publish)   // 发表文章
	group.POST("/withdraw", h.Withdraw) // 撤回文章
	group.POST("/list", h.List)         // 文章列表
	group.GET("/detail/:id", h.Detail)  // 文章详情
	group.GET("/pub/:id", h.PubDetail)  // 读者查看文章详情
	group.POST("/like", h.Like)         // 点赞
	group.POST("/collect", h.Collect)   // 收藏
	group.GET("/rank/:id", h.Rank)      // 点赞排行榜
}

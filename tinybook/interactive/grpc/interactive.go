package grpc

import (
	"context"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"tinybook/tinybook/api/proto/gen/intr/v1"
	"tinybook/tinybook/interactive/domain"
	"tinybook/tinybook/interactive/service"
)

type InteractiveServiceServer struct {
	intrv1.UnimplementedInteractiveServiceServer
	interactiveSvc service.InteractiveService
}

func NewInteractiveServiceServer(interactiveSvc service.InteractiveService) *InteractiveServiceServer {
	return &InteractiveServiceServer{interactiveSvc: interactiveSvc}
}

// Register 注册服务
func (i *InteractiveServiceServer) Register(server *grpc.Server) {
	intrv1.RegisterInteractiveServiceServer(server, i)
}

func (i *InteractiveServiceServer) IncreaseReadCount(ctx context.Context, request *intrv1.IncreaseReadCountRequest) (*intrv1.IncreaseReadCountResponse, error) {
	err := i.interactiveSvc.IncreaseReadCount(ctx, request.GetBiz(), request.GetBizId())
	return &intrv1.IncreaseReadCountResponse{}, err
}

func (i *InteractiveServiceServer) Like(ctx context.Context, request *intrv1.LikeRequest) (*intrv1.LikeResponse, error) {
	err := i.interactiveSvc.Like(ctx, request.GetBiz(), request.GetBizId(), request.GetUid())
	return &intrv1.LikeResponse{}, err
}

func (i *InteractiveServiceServer) Unlike(ctx context.Context, request *intrv1.UnlikeRequest) (*intrv1.UnlikeResponse, error) {
	err := i.interactiveSvc.Unlike(ctx, request.GetBiz(), request.GetBizId(), request.GetUid())
	return &intrv1.UnlikeResponse{}, err
}

func (i *InteractiveServiceServer) Collect(ctx context.Context, request *intrv1.CollectRequest) (*intrv1.CollectResponse, error) {
	err := i.interactiveSvc.Collect(ctx, request.GetBiz(), request.GetBizId(), request.GetCid(), request.GetUid())
	return &intrv1.CollectResponse{}, err
}

func (i *InteractiveServiceServer) GetInteractive(ctx context.Context, request *intrv1.GetInteractiveRequest) (*intrv1.GetInteractiveResponse, error) {
	interactive, err := i.interactiveSvc.GetInteractive(ctx, request.GetBiz(), request.GetBizId(), request.GetUid())
	if err != nil {
		return nil, err
	}
	return &intrv1.GetInteractiveResponse{
		Interactive: i.interactiveToDTO(interactive),
	}, nil
}

func (i *InteractiveServiceServer) GetLikeRanks(ctx context.Context, request *intrv1.GetLikeRanksRequest) (*intrv1.GetLikeRanksResponse, error) {
	ranks, err := i.interactiveSvc.GetLikeRanks(ctx, request.GetBiz(), request.GetNum())
	if err != nil {
		return nil, err
	}
	vos := lo.Map[domain.ArticleVo, *intrv1.ArticleVo](ranks, func(item domain.ArticleVo, index int) *intrv1.ArticleVo {
		return i.articleVoToDTO(item)
	})
	return &intrv1.GetLikeRanksResponse{
		Articles: vos,
	}, nil
}

func (i *InteractiveServiceServer) GetByIds(ctx context.Context, request *intrv1.GetByIdsRequest) (*intrv1.GetByIdsResponse, error) {
	ids, err := i.interactiveSvc.GetByIds(ctx, request.GetBiz(), request.GetIds())
	if err != nil {
		return nil, err
	}
	m := make(map[int64]*intrv1.Interactive, len(ids))
	for id := range ids {
		m[id] = i.interactiveToDTO(ids[id])
	}
	return &intrv1.GetByIdsResponse{
		Interactives: m,
	}, nil
}

func (i *InteractiveServiceServer) interactiveToDTO(interactive domain.Interactive) *intrv1.Interactive {
	return &intrv1.Interactive{
		Biz:          interactive.Biz,
		BizId:        interactive.BizId,
		ReadCount:    interactive.ReadCount,
		LikeCount:    interactive.LikeCount,
		CollectCount: interactive.CollectCount,
		Liked:        interactive.Liked,
		Collected:    interactive.Collected,
	}
}

func (i *InteractiveServiceServer) articleVoToDTO(articleVo domain.ArticleVo) *intrv1.ArticleVo {
	return &intrv1.ArticleVo{
		Id:         articleVo.ID,
		Title:      articleVo.Title,
		Content:    articleVo.Content,
		Abstract:   articleVo.Abstract,
		Author:     articleVo.Author,
		AuthorName: articleVo.AuthorName,
		Status:     articleVo.Status,
		Ctime:      articleVo.Ctime,
		Utime:      articleVo.Utime,
		// 以下是Interactive字段，因为ArticleVo是Article和Interactive的组合，但proto不能嵌套组合，所以这里需要手动赋值
		Biz:          articleVo.Biz,
		BizId:        articleVo.BizId,
		ReadCount:    articleVo.ReadCount,
		LikeCount:    articleVo.LikeCount,
		CollectCount: articleVo.CollectCount,
		Liked:        articleVo.Liked,
		Collected:    articleVo.Collected,
	}
}

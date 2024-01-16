package client

import (
	"context"
	"google.golang.org/grpc"
	intrv1 "tinybook/tinybook/api/proto/gen/intr/v1"
	"tinybook/tinybook/interactive/domain"
	"tinybook/tinybook/interactive/service"
)

type LocalInteractiveServiceAdapter struct {
	svc service.InteractiveService
}

func NewLocalInteractiveServiceAdapter(svc service.InteractiveService) *LocalInteractiveServiceAdapter {
	return &LocalInteractiveServiceAdapter{svc: svc}
}

func (l *LocalInteractiveServiceAdapter) IncreaseReadCount(ctx context.Context, in *intrv1.IncreaseReadCountRequest, opts ...grpc.CallOption) (*intrv1.IncreaseReadCountResponse, error) {
	err := l.svc.IncreaseReadCount(ctx, in.GetBiz(), in.GetBizId())
	return &intrv1.IncreaseReadCountResponse{}, err
}

func (l *LocalInteractiveServiceAdapter) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	err := l.svc.Like(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	return &intrv1.LikeResponse{}, err
}

func (l *LocalInteractiveServiceAdapter) Unlike(ctx context.Context, in *intrv1.UnlikeRequest, opts ...grpc.CallOption) (*intrv1.UnlikeResponse, error) {
	err := l.svc.Unlike(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	return &intrv1.UnlikeResponse{}, err
}

func (l *LocalInteractiveServiceAdapter) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	err := l.svc.Collect(ctx, in.GetBiz(), in.GetBizId(), in.GetCid(), in.GetUid())
	return &intrv1.CollectResponse{}, err
}

func (l *LocalInteractiveServiceAdapter) GetInteractive(ctx context.Context, in *intrv1.GetInteractiveRequest, opts ...grpc.CallOption) (*intrv1.GetInteractiveResponse, error) {
	interactive, err := l.svc.GetInteractive(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	if err != nil {
		return nil, err
	}
	return &intrv1.GetInteractiveResponse{
		Interactive: l.interactiveToDTO(interactive),
	}, nil
}

func (l *LocalInteractiveServiceAdapter) GetLikeRanks(ctx context.Context, in *intrv1.GetLikeRanksRequest, opts ...grpc.CallOption) (*intrv1.GetLikeRanksResponse, error) {
	ranks, err := l.svc.GetLikeRanks(ctx, in.GetBiz(), in.GetNum())
	if err != nil {
		return nil, err
	}
	var dtos []*intrv1.ArticleVo
	for _, rank := range ranks {
		dtos = append(dtos, l.articleVoToDTO(rank))
	}
	return &intrv1.GetLikeRanksResponse{
		Articles: dtos,
	}, nil
}

func (l *LocalInteractiveServiceAdapter) GetByIds(ctx context.Context, in *intrv1.GetByIdsRequest, opts ...grpc.CallOption) (*intrv1.GetByIdsResponse, error) {
	articles, err := l.svc.GetByIds(ctx, in.GetBiz(), in.GetIds())
	if err != nil {
		return nil, err
	}
	m := make(map[int64]*intrv1.Interactive)
	for i, interactive := range articles {
		m[i] = l.interactiveToDTO(interactive)
	}
	return &intrv1.GetByIdsResponse{
		Interactives: m,
	}, nil
}

func (l *LocalInteractiveServiceAdapter) interactiveToDTO(interactive domain.Interactive) *intrv1.Interactive {
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

func (l *LocalInteractiveServiceAdapter) articleVoToDTO(articleVo domain.ArticleVo) *intrv1.ArticleVo {
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

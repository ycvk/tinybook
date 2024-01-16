package client

import (
	"context"
	"google.golang.org/grpc"
	"log/slog"
	"math/rand"
	"strconv"
	"sync/atomic"
	intrv1 "tinybook/tinybook/api/proto/gen/intr/v1"
)

type InteractiveClient struct {
	remote intrv1.InteractiveServiceClient
	local  intrv1.InteractiveServiceClient

	threshold int32
}

func NewInteractiveClient(remote intrv1.InteractiveServiceClient, local intrv1.InteractiveServiceClient, threshold int32) *InteractiveClient {
	return &InteractiveClient{
		remote:    remote,
		local:     local,
		threshold: atomic.LoadInt32(&threshold),
	}
}

func (i *InteractiveClient) SetThreshold(th int32) {
	atomic.StoreInt32(&i.threshold, th)
}

func (i *InteractiveClient) selectClient() intrv1.InteractiveServiceClient {
	n := rand.Int31n(100)
	if n < i.threshold {
		slog.Info("n 为 " + strconv.FormatInt(int64(n), 10) + "，使用远程服务")
		return i.remote
	}
	slog.Info("n 为 " + strconv.FormatInt(int64(n), 10) + "，使用本地服务")
	return i.local
}
func (i *InteractiveClient) IncreaseReadCount(ctx context.Context, in *intrv1.IncreaseReadCountRequest, opts ...grpc.CallOption) (*intrv1.IncreaseReadCountResponse, error) {
	return i.selectClient().IncreaseReadCount(ctx, in, opts...)
}

func (i *InteractiveClient) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	return i.selectClient().Like(ctx, in, opts...)
}

func (i *InteractiveClient) Unlike(ctx context.Context, in *intrv1.UnlikeRequest, opts ...grpc.CallOption) (*intrv1.UnlikeResponse, error) {
	return i.selectClient().Unlike(ctx, in, opts...)
}

func (i *InteractiveClient) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	return i.selectClient().Collect(ctx, in, opts...)
}

func (i *InteractiveClient) GetInteractive(ctx context.Context, in *intrv1.GetInteractiveRequest, opts ...grpc.CallOption) (*intrv1.GetInteractiveResponse, error) {
	return i.selectClient().GetInteractive(ctx, in, opts...)
}

func (i *InteractiveClient) GetLikeRanks(ctx context.Context, in *intrv1.GetLikeRanksRequest, opts ...grpc.CallOption) (*intrv1.GetLikeRanksResponse, error) {
	return i.selectClient().GetLikeRanks(ctx, in, opts...)
}

func (i *InteractiveClient) GetByIds(ctx context.Context, in *intrv1.GetByIdsRequest, opts ...grpc.CallOption) (*intrv1.GetByIdsResponse, error) {
	return i.selectClient().GetByIds(ctx, in, opts...)
}

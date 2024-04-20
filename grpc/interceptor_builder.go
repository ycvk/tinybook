package tinybook

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	intrv1 "tinybook/tinybook/api/proto/gen/intr/v1"
	"tinybook/tinybook/pkg/limiter"
)

type InterceptorBuilder struct {
	limiter limiter.Limiter
	key     string
}

// BuildServerUnaryInterceptor 构建一个服务端的拦截器 用于限流 整个服务的限流
func (b *InterceptorBuilder) BuildServerUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		limit, err := b.limiter.Limit(ctx, b.key)
		if err != nil { //如果出现错误，直接返回
			return nil, status.Errorf(codes.ResourceExhausted, "rate limit error: %v", err)
		}
		if limit { //如果限流了，返回错误
			return nil, status.Errorf(codes.ResourceExhausted, "rate limit error: %v", err)
		}
		return handler(ctx, req)
	}
}

// BuildServerUnaryInterceptorBiz 构建一个服务端的拦截器 用于限流 服务的某个方法的限流
func (b *InterceptorBuilder) BuildServerUnaryInterceptorBiz() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// 类型断言 如果是 GetByIdsRequest 类型的请求，就限流
		if request, ok := req.(intrv1.GetByIdsRequest); ok {
			key := fmt.Sprintf("limter:user:get_by_ids:%v", request.Ids)
			limit, err := b.limiter.Limit(ctx, key)
			if err != nil { //如果出现错误，直接返回
				return nil, status.Errorf(codes.ResourceExhausted, "rate limit error: %v", err)
			}
			if limit { //如果限流了，返回错误
				return nil, status.Errorf(codes.ResourceExhausted, "rate limit error: %v", err)
			}
		}
		return handler(ctx, req)
	}
}

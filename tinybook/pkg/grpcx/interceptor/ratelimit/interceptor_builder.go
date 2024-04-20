package ratelimit

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"tinybook/tinybook/pkg/limiter"
)

type InterceptorBuilder struct {
	limiter limiter.Limiter
	key     string
}

func NewInterceptorBuilder(limiter limiter.Limiter, key string) *InterceptorBuilder {
	return &InterceptorBuilder{
		limiter: limiter,
		key:     key,
	}
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

// BuildServerUnaryInterceptorService 构建一个服务端的拦截器 用于限流 服务的某个方法的限流
func (b *InterceptorBuilder) BuildServerUnaryInterceptorService() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if strings.HasPrefix(info.FullMethod, "/UserService") {
			limit, err := b.limiter.Limit(ctx, b.key)
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

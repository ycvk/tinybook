package circuitbreaker

import (
	"context"
	"github.com/go-kratos/aegis/circuitbreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InterceptorBuilder struct {
	breaker circuitbreaker.CircuitBreaker
}

// BuildServerUnaryInterceptor 构建一个服务端的拦截器 用于熔断 整个服务的熔断
func (b *InterceptorBuilder) BuildServerUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		err := b.breaker.Allow()
		if err != nil {
			b.breaker.MarkFailed()
			return nil, status.Errorf(codes.Unavailable, "circuit breaker error: %s", err.Error())
		}
		resp, err := handler(ctx, req)
		if err != nil {
			// 这里是业务逻辑错误，不是熔断错误 所以需要再细粒度的判断是否需要熔断
			b.breaker.MarkFailed()
			return nil, status.Errorf(codes.Unavailable, "circuit breaker error: %s", err.Error())
		} else {
			b.breaker.MarkSuccess()
		}
		return resp, err
	}
}

.PHONY: docker mock
docker:
	# wire 依赖注入
	@ wire
	# 把上次编译的二进制文件删除
	@rm -f tinybook || true
	@go mod tidy
	# 编译二进制文件
	@GOOS=linux GOARCH=arm64 go build -tags=k8s -o tinybook .
	# 删除上次编译的镜像
	@docker rmi -f ycvk/tinybook:v0.0.1 || true
	# 构建镜像
	@docker build -t ycvk/tinybook:v0.0.1 .

mock:
	# 生成user service mock文件
	@mockgen -source=internal/service/user.go -package=mocks -destination=internal/service/mocks/user_mock.go
	# 生成code service mock文件
	@mockgen -source=internal/service/code.go -package=mocks -destination=internal/service/mocks/code_mock.go
	# 生成sms service mock文件
	@mockgen -source=internal/service/sms/types.go -package=smsmocks -destination=internal/service/sms/mocks/sms_mock.go
	# 生成user repository mock文件
	@mockgen -source=internal/repository/user.go -package=repomocks -destination=internal/repository/mocks/user_mock.go
	# 生成code repository mock文件
	@mockgen -source=internal/repository/code.go -package=repomocks -destination=internal/repository/mocks/code_mock.go
	# 生成sms repository mock文件
	@mockgen -source=internal/repository/sms.go -package=repomocks -destination=internal/repository/mocks/sms_mock.go
	# 生成limiter mock文件
	@mockgen -source=pkg/limiter/types.go -package=limitermocks -destination=pkg/limiter/mocks/limiter_mock.go
	# 生成sms async retry mock文件
	@mockgen -source=internal/service/sms/failover/retry/retry_task.go -package=retry -destination=internal/service/sms/failover/retry/async_retry_mock.go
	# 生成sms error monitor mock文件
	@mockgen -source=internal/service/sms/failover/retry/error_rate_monitor.go -package=retry -destination=internal/service/sms/failover/retry/error_rate_monitor_mock.go
	# 解决依赖
	@go mod tidy

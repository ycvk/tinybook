.PHONY: docker
docker:
	# 把上次编译的二进制文件删除
	@rm -f tinybook || true
	@go mod tidy
	# 编译二进制文件
	@GOOS=linux GOARCH=arm64 go build -o tinybook .
	# 删除上次编译的镜像
	@docker rmi -f ycvk/tinybook:v0.0.1 || true
	# 构建镜像
	@docker build -t ycvk/tinybook:v0.0.1 .

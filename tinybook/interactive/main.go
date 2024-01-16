package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"tinybook/tinybook/ioc"
	"tinybook/tinybook/pkg/grpcx"
)

func init() {
	initViper()      // 初始化配置
	initPrometheus() // 初始化prometheus
}

func main() {
	otel := ioc.InitOTEL() // 初始化otel
	defer func() {
		timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
		defer cancelFunc()
		otel(timeout) // 服务器关闭时, 超时控制去关闭otel
	}()

	// 初始化服务
	app := InitInteractiveApp()

	// 启动kafka消费者
	for i := range app.consumers {
		app.consumers[i].Start()
	}

	// 启动grpc服务
	server := app.server
	go func() {
		if err := server.Serve(); err != nil {
			panic(err)
		}
	}()

	// 监听项目退出
	exit(server)
}

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8082", nil)
	}()
}

func initViper() {
	//err := viper.AddRemoteProvider("etcd3", "tinybook-etcd:2381", "tinybook")
	err := viper.AddRemoteProvider("etcd3", "127.0.0.1:32379", "tinybook-local")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

// 监听退出
func exit(engine *grpcx.Server) {
	sigs := make(chan os.Signal, 1)
	quit := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-sigs
		fmt.Println("收到 interactive 退出信号: ", sig)
		// 退出grpc服务
		engine.GracefulStop()
		quit <- true
	}()
	<-quit
	fmt.Println("interactive 服务 PID 为: ", os.Getpid())
	fmt.Println("interactive 服务已退出")
	// 查杀
	exec.Command("killall", "main", strconv.Itoa(os.Getpid())).Run()
	// 自杀
	exec.Command("kill", "-9", strconv.Itoa(os.Getpid())).Run()
}

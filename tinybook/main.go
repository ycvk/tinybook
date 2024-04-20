package main

import (
	"context"
	"errors"
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
)

func init() {
	initViper()         // 初始化配置
	ioc.InitSnowflake() // 初始化雪花算法
	initPrometheus()    // 初始化prometheus
}

func main() {
	otel := ioc.InitOTEL() // 初始化otel
	defer func() {
		timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
		defer cancelFunc()
		otel(timeout) // 服务器关闭时, 超时控制去关闭otel
	}()

	app := InitWebServer() // 初始化web服务

	// 启动kafka消费者
	for i := range app.consumers {
		app.consumers[i].Start()
	}

	// 启动定时任务
	app.cron.Start()
	defer func() {
		// 等待退出
		<-app.cron.Stop().Done()
	}()

	// 启动web服务
	server := &http.Server{
		Addr:    ":8081",
		Handler: app.server,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("web服务启动失败: ", err)
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
	//s := pflag.String("config", "config/dev.yaml", "配置文件路径")
	//pflag.Parse()
	//viper.SetConfigFile(*s)
	//err := viper.ReadInConfig()
	//if err != nil {
	//	panic(err)
	//}
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
	viper.WatchConfig()

}

// 监听退出
func exit(engine *http.Server) {
	sigs := make(chan os.Signal, 1)
	quit := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-sigs
		fmt.Println("收到退出信号: ", sig)
		// 退出web服务
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := engine.Shutdown(ctx); err != nil {
			fmt.Println("web服务退出失败: ", err)
		}
		quit <- true
	}()
	<-quit
	fmt.Println("服务 PID 为: ", os.Getpid())
	fmt.Println("服务已退出")
	// 查杀
	exec.Command("killall", "main", strconv.Itoa(os.Getpid())).Run()
	// 自杀
	exec.Command("kill", "-9", strconv.Itoa(os.Getpid())).Run()
}

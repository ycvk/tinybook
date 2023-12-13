package ioc

import (
	"geek_homework/tinybook/pkg/gormx"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"
	"sync"
)

var (
	gormDB *gorm.DB
	once   sync.Once
)

func InitDB(zipLog *zap.Logger) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}

	once.Do(func() {
		gormDB, err = gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
			Logger: logger.New(gormLoggerFunc(func(msg string, data ...interface{}) {
				zipLog.Info(msg, zap.Any("data", data))
			}), logger.Config{
				SlowThreshold: 0,           // 慢查询阈值 0 表示打印所有sql
				LogLevel:      logger.Info, // 日志级别
			}),
		})
	})
	if err != nil {
		panic(err)
	}
	err = gormDB.Use(prometheus.New(prometheus.Config{
		DBName:          "tinybook",
		RefreshInterval: 15, // 指标刷新频率，单位秒
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"Threads_running"},
			},
		},
	}))
	if err != nil {
		panic(err)
	}
	callbacks := gormx.NewCallbacks(prometheus2.SummaryOpts{
		Namespace: "tinybook",
		Subsystem: "mysql",
		Name:      "gorm_db",
		Help:      "统计gorm的sql执行时间",
		ConstLabels: map[string]string{
			"instance_id": "my_instance",
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})
	// 注册Prometheus插件
	err = gormDB.Use(callbacks)
	if err != nil {
		panic(err)
	}
	// TODO 为了方便测试，每次启动都会重新创建表 仅供测试使用
	CreateTable(gormDB)
	return gormDB
}

// gormLoggerFunc gorm日志
type gormLoggerFunc func(msg string, data ...interface{})

// Printf 实现gorm日志接口
func (f gormLoggerFunc) Printf(format string, args ...interface{}) {
	f(format, args...)
}

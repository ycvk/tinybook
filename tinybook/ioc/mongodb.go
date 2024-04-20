package ioc

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/here-Leslie-Lau/mongo-plus/mongo"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	mgoptions "go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"sync"
)

func InitMongoDB(zipLog *zap.Logger) *qmgo.Database {
	//const (
	//	Ip         = "127.0.0.1"
	//	Port       = "27017"
	//	UserName   = "root"
	//	Password   = "123456"
	//	DBName     = "tinybook"
	//	AuthSource = "admin"
	//)

	type Config struct {
		Ip         string `yaml:"ip"`
		Port       string `yaml:"port"`
		UserName   string `yaml:"username"`
		Password   string `yaml:"password"`
		DBName     string `yaml:"dbname"`
		AuthSource string `yaml:"authSource"`
	}
	var cfg Config
	viperErr := viper.UnmarshalKey("mongodb", &cfg)
	if viperErr != nil {
		panic(viperErr)
	}

	var (
		ConnectTimeoutMS = int64(1000)    // 连接超时时间
		MaxPoolSize      = uint64(100)    // 最大连接数
		MinPoolSize      = uint64(0)      // 最小连接数
		mLog             = zipLog.Sugar() // 日志
	)

	ctx := context.Background()
	// 拼接MongoDB Url
	var mongoUrl string
	if cfg.Password != "" {
		mongoUrl = "mongodb://" + cfg.UserName + ":" + cfg.Password + "@" +
			cfg.Ip + ":" + cfg.Port + "/" + cfg.DBName +
			"?authSource=" + cfg.AuthSource
	} else {
		mongoUrl = "mongodb://" + cfg.Ip + ":" + cfg.Port
	}

	// 创建cmdMonitor，用于打印SQL
	//startedCommands := make(map[int64]bson.Raw)
	startedCommands := sync.Map{}        // map[int64]bson.Raw
	cmdMonitor := &event.CommandMonitor{ // CommandMonitor 用于监控命令
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			startedCommands.Store(evt.RequestID, evt.Command)
			//startedCommands[evt.RequestID] = evt.Command
		},
		Succeeded: func(_ context.Context, evt *event.CommandSucceededEvent) {
			//log.Printf("Command: %v Reply: %v\n",
			//	startedCommands[evt.RequestID],
			//	evt.Reply,
			//)
			var commands bson.Raw
			v, ok := startedCommands.Load(evt.RequestID)
			if ok {
				commands = v.(bson.Raw)
			}
			defer mLog.Sync()
			mLog.Infof("\n[MongoDB] [%.3fms] [%v] %v \n", float64(evt.DurationNanos)/1e6, commands, evt.Reply)
		},
		Failed: func(_ context.Context, evt *event.CommandFailedEvent) {
			//log.Printf("Command: %v Failure: %v\n",
			//	startedCommands[evt.RequestID],
			//	evt.Failure,
			//)
			var commands bson.Raw
			v, ok := startedCommands.Load(evt.RequestID)
			if ok {
				commands = v.(bson.Raw)
			}
			defer mLog.Sync()
			mLog.Fatalf("\n[MongoDB] [%.3fms] [%v] \n %v \n", float64(evt.DurationNanos)/1e6, commands, evt.Failure)
		},
	}
	// 创建options
	ops := options.ClientOptions{ClientOptions: &mgoptions.ClientOptions{}}
	ops.SetMonitor(cmdMonitor)

	// 创建数据库链接
	client, err := qmgo.NewClient(ctx, &qmgo.Config{
		Uri:              mongoUrl,
		ConnectTimeoutMS: &ConnectTimeoutMS,
		MaxPoolSize:      &MaxPoolSize,
		MinPoolSize:      &MinPoolSize,
	}, ops)

	if err != nil {
		err = errors.New("MongoDB连接异常：" + err.Error())
		panic(err)
	}
	// 选择数据库
	db := client.Database(cfg.DBName)
	// 在初始化成功后，测试使用完毕请defer来关闭连接
	//defer func() {
	//	if err = client.Close(ctx); err != nil {
	//		panic(err)
	//	}
	//}()
	// 创建索引
	err = createIndex(db)
	return db
}

func createIndex(db *qmgo.Database) error {
	// 创建索引
	articleColl := db.Collection("articles")
	err := articleColl.CreateIndexes(context.Background(), []options.IndexModel{
		{
			Key:          []string{"id"},
			IndexOptions: mgoptions.Index().SetUnique(true),
		},
		{
			Key: []string{"author_id"},
		},
	})
	if err != nil {
		return err
	}
	publishedArticleColl := db.Collection("published_articles")
	err = publishedArticleColl.CreateIndexes(context.Background(), []options.IndexModel{
		{
			Key:          []string{"id"},
			IndexOptions: mgoptions.Index().SetUnique(true),
		},
		{
			Key: []string{"author_id"},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func InitMongoDBV2() *mongo.Conn {
	type Config struct {
		Ip         string `yaml:"ip"`
		Port       string `yaml:"port"`
		UserName   string `yaml:"username"`
		Password   string `yaml:"password"`
		DBName     string `yaml:"dbname"`
		AuthSource string `yaml:"authSource"`
	}
	var cfg Config
	viperErr := viper.UnmarshalKey("mongodb", &cfg)
	if viperErr != nil {
		panic(viperErr)
	}

	opts := []mongo.Option{
		// Database to Connect
		mongo.WithDatabase(cfg.DBName),
		// Maximum Connection Pool Size
		mongo.WithMaxPoolSize(10),
		// Username
		mongo.WithUsername(cfg.UserName),
		// Password
		mongo.WithPassword(cfg.Password),
		// Connection URL
		mongo.WithAddr(cfg.Ip + ":" + cfg.Port),
	}
	conn, _ := mongo.NewConn(opts...)
	//defer f()
	return conn
}

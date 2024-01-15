package job

import (
	"context"
	"github.com/bsm/redislock"
	"github.com/lithammer/shortuuid/v4"
	"go.uber.org/zap"
	"math/rand"
	"time"
	"tinybook/tinybook/internal/service"
	"tinybook/tinybook/pkg/hashring"
)

type RankingJob struct {
	log        *zap.Logger
	rankingSvc service.RankingService
	time       time.Duration
	RedisLock  *redislock.Client
	key        string

	// TODO:这是用来模拟负载的，实际使用时需要删除
	load      int32                    // 节点负载
	threshold int32                    // 节点负载阈值
	Id        string                   // 节点ID
	hashRing  *hashring.ConsistentHash // 一致性哈希环
}

func NewRankingJob(rankingSvc service.RankingService, t time.Duration, lock *redislock.Client, l *zap.Logger) *RankingJob {
	r := &RankingJob{
		rankingSvc: rankingSvc,
		time:       t,
		RedisLock:  lock,
		log:        l,
		threshold:  85, // 节点负载阈值
	}

	// TODO:这是用来模拟负载的，实际使用时需要删除
	uuid := shortuuid.New()
	r.Id = uuid
	ticker := time.NewTicker(time.Second) // 每秒钟更改一个负载 0-100 模拟负载
	go func() {
		for {
			select {
			case <-ticker.C:
				r.setLoad(rand.Int31n(100))
			}
		}
	}()
	// 初始化一致性哈希环
	ch := &hashring.ConsistentHash{
		Ring:  make(hashring.HashRing, 0),
		Nodes: make(hashring.NodeMap),
	}
	// 添加节点
	ch.AddNode(hashring.Node{ID: r.Id, Load: r.load})
	// 定时(10s)更新哈希环中的节点负载
	ch.AutoUpdateLoadByFunc(r.Id, time.Second*10, r.GetLoad)
	r.hashRing = ch
	return r
}

// setLoad TODO:这是用来设置负载的，实际使用时需要删除
func (r *RankingJob) setLoad(load int32) {
	r.load = load
}

// GetLoad TODO:这是用来获取负载的，实际使用时需要删除
func (r *RankingJob) GetLoad() (int32, error) {
	return r.load, nil
}

func (r *RankingJob) Name() string {
	return "ranking"
}

func (r *RankingJob) Run() error {
	node := r.hashRing.GetNode(string(rand.Int31())) // 获取负载极可能最小的节点
	if node.Load > r.threshold {                     // 如果获取负载的节点也都已经负载超过阈值，那就暂不执行，打印日志
		r.log.Warn("ranking job load is too high", zap.String("node", node.ID), zap.Int32("load", node.Load))
		return nil
	}
	if node.ID != r.Id { // 如果不是自己，就不执行
		return nil
	}

	timeout, c := context.WithTimeout(context.Background(), time.Second*3)
	defer c()
	lock, err := r.RedisLock.Obtain(timeout, r.key, r.time,
		&redislock.Options{
			RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(time.Millisecond*100), 3), // 重试3次，每次间隔100ms
		})
	if err != nil {
		r.log.Error("ranking job lock failed", zap.Error(err))
		return err
	}
	defer func() {
		withTimeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
		defer cancelFunc()
		err2 := lock.Release(withTimeout)
		if err2 != nil {
			r.log.Error("ranking job unlock failed", zap.Error(err2))
		}
	}()
	// 执行前再次检查节点负载 如果目前节点突然已经超过阈值，就不执行，释放此分布式锁，打印日志
	if r.load > r.threshold {
		r.log.Warn("ranking job load is too high", zap.String("node", node.ID), zap.Int32("load", node.Load))
		return nil
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), r.time)
	defer cancelFunc()
	return r.rankingSvc.TopN(ctx)
}

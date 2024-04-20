package wrr

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"sync"
)

const Name = "custom_weighted_round_robin"

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &PickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

// PickerBuilder 是一个负载均衡器 Picker 的构造器
type PickerBuilder struct {
}

func (p *PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	var conns = make([]*weightConn, 0, len(info.ReadySCs))
	for sc, scInfo := range info.ReadySCs {
		// 从 SubConnInfo 中获取权重信息
		weight := scInfo.Address.Metadata.(map[string]any)["weight"].(float64)
		conns = append(conns, &weightConn{
			SubConn:       sc,
			weight:        int(weight),
			currentWeight: 0,
		})
	}
	return &Picker{conns: conns}
}

// Picker 是一个负载均衡器的实现
type Picker struct {
	conns []*weightConn
	lock  sync.Mutex
}

// weightConn 是 SubConn 的包装，增加了权重信息
type weightConn struct {
	balancer.SubConn
	weight        int
	currentWeight int
	// 用于指数退避算法
	failureCount int
	successCount int
}

func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if len(p.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	// 总权重
	total := 0
	// 选择最大权重的连接
	maxWeightConn := p.conns[0]
	for _, conn := range p.conns {
		// 计算总权重
		total += conn.weight
		// 每次加上自己的权重
		conn.currentWeight += conn.weight
		// 选择最大权重的连接
		if conn.currentWeight > maxWeightConn.currentWeight {
			maxWeightConn = conn
		}
	}

	// 选择后减去总权重
	maxWeightConn.currentWeight -= total

	pickedConn := maxWeightConn // 已选择的连接
	return balancer.PickResult{
		SubConn: maxWeightConn.SubConn,
		Done: func(doneInfo balancer.DoneInfo) {
			// 可以在这里对连接的权重进行调整, 例如根据连接的成功率进行调整
			p.adjustWeight(pickedConn, doneInfo.Err)
		},
	}, nil
}

// 调整权重的方法 指数退避算法
func (p *Picker) adjustWeight(conn *weightConn, err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	const MaxIncrease = 10 // 最大增加幅度
	const MaxDecrease = 10 // 最大减少幅度

	if err != nil {
		// 如果有错误，增加失败次数，根据失败次数来进行指数退避减少权重
		//但不能超过最大减少幅度
		if conn.failureCount > 3 {
			return
		}
		decreaseFactor := min(1<<conn.failureCount, MaxDecrease)
		conn.weight -= decreaseFactor
		conn.failureCount++
		conn.successCount = 0
	} else {
		// 如果成功，增加权重
		if conn.failureCount > 0 {
			conn.failureCount--
		}
		increaseFactor := min(1<<conn.successCount, MaxIncrease)
		conn.weight += increaseFactor
		conn.successCount++
		//conn.failureCount = 0
	}
}

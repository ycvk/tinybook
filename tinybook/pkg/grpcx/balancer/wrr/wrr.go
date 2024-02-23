package wrr

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"
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
			SubConn:         sc,
			weight:          int(weight),
			effectiveWeight: int(weight),
			isAvailable:     true,
			currentWeight:   0,
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
	weight          int  // 权重
	effectiveWeight int  // 有效权重
	currentWeight   int  // 当前权重
	isAvailable     bool // 是否可用
}

// OnInvokeSuccess 提权操作
func (conn *weightConn) OnInvokeSuccess() {
	if conn.effectiveWeight < conn.weight {
		conn.effectiveWeight++
	}
}

// OnInvokeFault 降权操作
func (conn *weightConn) OnInvokeFault() {
	conn.effectiveWeight = max(1, conn.effectiveWeight-1)
}

// OnInvokeUnavailable 处理熔断逻辑
func (conn *weightConn) OnInvokeUnavailable() {
	conn.isAvailable = false
	// 可以在这里启动一个 goroutine 进行健康检查
	go conn.healthCheck()
}

// healthCheck 定期进行健康检查
func (conn *weightConn) healthCheck() {
	// 定义健康检查逻辑
	for {
		time.Sleep(3 * time.Second) //健康检查间隔
		// 发送健康检查请求...
		// 如果节点恢复，更新状态并跳出循环
		//if 节点已恢复 {
		//	conn.isAvailable = true
		//	break
		//}
	}
}

func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if len(p.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	// 总有效权重
	totalEffectiveWeight := 0
	for _, conn := range p.conns {
		totalEffectiveWeight += conn.effectiveWeight
	}

	var maxWeightConn *weightConn
	for _, conn := range p.conns {
		// 增加当前权重
		conn.currentWeight += conn.effectiveWeight

		// 选择当前权重最大的可用连接
		if maxWeightConn == nil || (conn.currentWeight > maxWeightConn.currentWeight && conn.isAvailable) {
			maxWeightConn = conn
		}
	}

	// 选中节点后，减去总有效权重
	if maxWeightConn != nil {
		maxWeightConn.currentWeight -= totalEffectiveWeight
		// 检查节点是否被熔断
		if !maxWeightConn.isAvailable {
			return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
		}
	}

	return balancer.PickResult{
		SubConn: maxWeightConn.SubConn,
		Done: func(doneInfo balancer.DoneInfo) {
			// 可以在这里对连接的权重进行调整, 例如根据连接的成功率进行调整
			p.lock.Lock()
			defer p.lock.Unlock()

			if err := doneInfo.Err; err != nil {
				if errStatus, ok := status.FromError(err); ok {
					switch errStatus.Code() {
					case codes.Unavailable: // 不可用 熔断
						maxWeightConn.OnInvokeUnavailable()
					case codes.ResourceExhausted: // 资源耗尽 降级
						maxWeightConn.OnInvokeFault()
					default:
						return
					}
				}
			} else {
				maxWeightConn.OnInvokeSuccess()
			}
		},
	}, nil
}

package hashring

import (
	uuid "github.com/lithammer/shortuuid/v4"
	"math/rand"
	"testing"
)

func BenchmarkConsistentHash_AddNode(b *testing.B) {
	// 准备随机数据
	var nodes []Node
	for i := 0; i < 1000; i++ {
		node := Node{
			ID:   uuid.New(),
			Load: rand.Int31n(MaxLoad),
		}
		nodes = append(nodes, node)
	}

	b.ResetTimer() // 重置计时器，开始测试

	for i := 0; i < b.N; i++ {
		ch := &ConsistentHash{
			Ring:  make(HashRing, 0),
			Nodes: make(NodeMap),
		}

		for _, node := range nodes {
			ch.AddNode(node) // 使用预先准备的数据进行测试
		}
	}
}

func BenchmarkConsistentHash_RemoveNode(b *testing.B) {
	// 预先填充 ConsistentHash 实例
	originalCH := &ConsistentHash{
		Ring:  make(HashRing, 0),
		Nodes: make(NodeMap),
	}
	nodeIDs := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		nodeID := uuid.New()
		nodeIDs[i] = nodeID
		originalCH.AddNode(Node{ID: nodeID, Load: rand.Int31n(MaxLoad)})
	}

	b.ResetTimer() // 重置计时器，开始测试
	for i := 0; i < b.N; i++ {
		// 深度复制 originalCH 以保证每次迭代的环境一致
		ch := originalCH.DeepCopy()

		// 在每次迭代中随机选择一个节点进行移除
		nodeID := nodeIDs[rand.Intn(len(nodeIDs))]
		ch.RemoveNode(nodeID)
	}
}

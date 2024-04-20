package hashring

import (
	uuid "github.com/lithammer/shortuuid/v4"
	"math/rand"
	"strconv"
	"testing"
	"time"
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
		ch := NewHashRing()
		for _, node := range nodes {
			ch.AddNode(node) // 使用预先准备的数据进行测试
		}
	}
}

func BenchmarkConsistentHash_RemoveNode(b *testing.B) {
	// 预先填充 ConsistentHash 实例
	originalCH := NewHashRing()
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

func TestAddNode(t *testing.T) {
	ring := NewHashRing()
	ring.AddNode(Node{ID: "node1", Load: 100})
	ring.AddNode(Node{ID: "node2", Load: 70})
	ring.AddNode(Node{ID: "node3", Load: 50})
	ring.AddNode(Node{ID: "node4", Load: 30})
	ring.AddNode(Node{ID: "node5", Load: 10})
	m := make(map[string]int)
	for i := 0; i < 1000; i++ {
		m[ring.GetNode(strconv.Itoa(i)).ID]++
	}
	for k, v := range m {
		t.Logf("node: %s, count: %d", k, v)
	}
}

func TestRemoveNode(t *testing.T) {
	ring := NewHashRing()
	ring.AddNode(Node{ID: "node1", Load: 50})
	ring.RemoveNode("node1")
	if ring.GetNode("key1").ID != "" {
		t.Errorf("Expected empty, got %s", ring.GetNode("key1").ID)
	}
}

func TestUpdateLoad(t *testing.T) {
	ring := NewHashRing()
	ring.AddNode(Node{ID: "node1", Load: 50})
	ring.UpdateLoad("node1", 10)
	if ring.GetNode("key1").Load != 10 {
		t.Errorf("Expected 10, got %d", ring.GetNode("key1").Load)
	}
}

func TestAutoUpdateLoadByFunc(t *testing.T) {
	ring := NewHashRing()
	ring.AddNode(Node{ID: "node1", Load: 50})
	ring.AutoUpdateLoadByFunc("node1", time.Second, func() (int32, error) {
		return 20, nil
	})
	time.Sleep(2 * time.Second)
	if ring.GetNode("key1").Load != 20 {
		t.Errorf("Expected 20, got %d", ring.GetNode("key1").Load)
	}
}

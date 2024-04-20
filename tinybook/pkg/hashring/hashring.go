package hashring

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
	"time"
)

const MaxLoad = 100

// Node 代表哈希环中的一个节点
type Node struct {
	ID   string
	Load int32 // 节点负载
}

// HashRing 一致性哈希环
type HashRing []uint32

// NodeMap 节点映射表
type NodeMap map[uint32]Node

// ConsistentHash 一致性哈希结构体
type ConsistentHash struct {
	Ring  HashRing
	Nodes NodeMap
	sync.RWMutex
}

// AddNode 添加节点到哈希环
func (c *ConsistentHash) AddNode(node Node) {
	c.Lock()
	defer c.Unlock()

	if node.Load > MaxLoad {
		node.Load = MaxLoad
	}
	virtualNodes := MaxLoad - node.Load // 虚拟节点数量 = 最大负载 - 节点负载, 负载越高 节点越少
	keys := make([]uint32, virtualNodes)

	for i := 0; i < int(virtualNodes); i++ {
		key := c.hashKey(node.ID + strconv.Itoa(i))
		c.Nodes[key] = node
		keys[i] = key
	}

	c.Ring = append(c.Ring, keys...)
	sort.Slice(c.Ring, func(i, j int) bool {
		return c.Ring[i] < c.Ring[j]
	})
}

// RemoveNode 从哈希环中删除节点
func (c *ConsistentHash) RemoveNode(nodeID string) {
	c.Lock()
	defer c.Unlock()
	var newRing HashRing
	for i := 0; i < MaxLoad; i++ {
		key := c.hashKey(nodeID + strconv.Itoa(i))
		if _, exists := c.Nodes[key]; exists {
			delete(c.Nodes, key)
			continue
		}
		newRing = append(newRing, c.Ring[i])
	}
	c.Ring = newRing
}

// GetNode 返回负载极可能最小的节点
func (c *ConsistentHash) GetNode(key string) Node {
	c.RLock()
	defer c.RUnlock()
	if len(c.Ring) == 0 {
		return Node{}
	}
	hash := c.hashKey(key)
	idx := c.search(hash)
	return c.Nodes[c.Ring[idx]]
}

// DeepCopy 方法用于深度复制 ConsistentHash 实例
func (c *ConsistentHash) DeepCopy() *ConsistentHash {
	newCH := &ConsistentHash{
		Ring:  make(HashRing, len(c.Ring)),
		Nodes: make(NodeMap, len(c.Nodes)),
	}
	copy(newCH.Ring, c.Ring)
	for k, v := range c.Nodes {
		newCH.Nodes[k] = v
	}
	return newCH
}

// hashKey 为字符串键生成哈希值
func (c *ConsistentHash) hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

// search 查找哈希环中某个哈希值的索引
func (c *ConsistentHash) search(hash uint32) int {
	idx := sort.Search(len(c.Ring), func(i int) bool {
		return c.Ring[i] >= hash
	})
	if idx == len(c.Ring) {
		idx = 0
	}
	return idx
}

// UpdateLoad 更新节点负载
func (c *ConsistentHash) UpdateLoad(nodeID string, load int32) {
	c.Lock()
	defer c.Unlock()
	c.RemoveNode(nodeID)
	node := Node{ID: nodeID, Load: load}
	c.AddNode(node)
}

// AutoUpdateLoadByFunc 定时更新节点负载
func (c *ConsistentHash) AutoUpdateLoadByFunc(nodeID string, duration time.Duration, loadFunc func() (int32, error)) {
	ticker := time.NewTicker(duration)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				load, err := loadFunc()
				if err != nil {
					c.UpdateLoad(nodeID, MaxLoad) // 如果获取负载失败，就设置为最大负载,节点将不会被选中
				}
				c.UpdateLoad(nodeID, load)
			}
		}
	}()
}

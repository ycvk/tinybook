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

func (c *ConsistentHash) AddNode(node Node) {
	c.Lock()
	defer c.Unlock()
	if node.Load > MaxLoad {
		node.Load = MaxLoad
	}
	virtualNodes := MaxLoad - node.Load // 虚拟节点数量 = 最大负载 - 节点负载, 负载越高 节点越少
	for i := 0; i < int(virtualNodes); i++ {
		key := c.hashKey(node.ID + strconv.Itoa(i))
		c.Nodes[key] = node
	}
}

func (c *ConsistentHash) RemoveNode(nodeID string) {
	c.Lock()
	defer c.Unlock()
	for i := 0; i < MaxLoad; i++ {
		key := c.hashKey(nodeID + strconv.Itoa(i))
		if _, exists := c.Nodes[key]; exists {
			delete(c.Nodes, key)
		}
	}
}

func (c *ConsistentHash) GetNode(key string) Node {
	c.RLock()
	defer c.RUnlock()
	if len(c.Ring) == 0 {
		return Node{}
	}
	hash := c.hashKey(key)
	idx := sort.Search(len(c.Ring), func(i int) bool { return c.Ring[i] >= hash })
	if idx == len(c.Ring) {
		idx = 0
	}
	return c.Nodes[c.Ring[idx]]
}

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

func (c *ConsistentHash) hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *ConsistentHash) UpdateLoad(nodeID string, load int32) {
	c.Lock()
	defer c.Unlock()
	c.RemoveNode(nodeID)
	node := Node{ID: nodeID, Load: load}
	c.AddNode(node)
}

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

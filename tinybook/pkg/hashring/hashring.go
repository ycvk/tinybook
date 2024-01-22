package hashring

import (
	"github.com/cespare/xxhash/v2"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"
)

const MaxLoad = 100

// ConsistentHashRing 一致性哈希环接口
type ConsistentHashRing interface {
	AddNode(node Node)
	RemoveNode(nodeID string)
	GetNode(key string) Node
	DeepCopy() ConsistentHashRing
	UpdateLoad(nodeID string, load int32)
	AutoUpdateLoadByFunc(nodeID string, duration time.Duration, loadFunc func() (int32, error))
}

// Node 代表哈希环中的一个节点
type Node struct {
	ID   string
	Load int32 // 节点负载
}

// HashRing 哈希环
type HashRing []uint64

// ConsistentHash 一致性哈希结构体
type ConsistentHash struct {
	Ring  HashRing
	Nodes sync.Map
	sync.RWMutex
	hash *xxhash.Digest
}

func NewHashRing() ConsistentHashRing {
	return &ConsistentHash{
		Ring:  make(HashRing, 0),
		Nodes: sync.Map{},
		hash:  xxhash.New(),
	}
}

func (c *ConsistentHash) AddNode(node Node) {
	c.Lock()
	defer c.Unlock()
	c.updateVirtualNodes(node.ID, node.Load)
	c.sortRing() // 添加节点后，可能需要重新排序
}

func (c *ConsistentHash) RemoveNode(nodeID string) {
	c.Lock()
	defer c.Unlock()
	c.removeVirtualNodes(nodeID)
	c.sortRing() // 移除节点后，可能需要重新排序
}

// GetNode 根据 key 随机获取一个节点
func (c *ConsistentHash) GetNode(key string) Node {
	c.RLock()
	defer c.RUnlock()
	if len(c.Ring) == 0 {
		return Node{}
	}
	int31 := rand.Int31()
	hashKey := c.hashKey(key + strconv.Itoa(int(int31))) // 生成哈希值
	idx := sort.Search(len(c.Ring), func(i int) bool { return c.Ring[i] >= hashKey })
	if idx == len(c.Ring) {
		idx = 0
	}
	node, ok := c.Nodes.Load(c.Ring[idx]) // 从节点映射表中获取节点
	if !ok {
		return Node{} // 如果类型断言失败，返回空节点
	}
	return node.(Node)
}

func (c *ConsistentHash) DeepCopy() ConsistentHashRing {
	c.RLock()
	defer c.RUnlock()
	newCH := &ConsistentHash{
		Ring:  make(HashRing, len(c.Ring)),
		Nodes: sync.Map{},
		hash:  xxhash.New(), // 复制hash字段
	}
	copy(newCH.Ring, c.Ring)
	c.Nodes.Range(func(key, value interface{}) bool {
		newCH.Nodes.Store(key, value)
		return true
	})
	return newCH
}

func (c *ConsistentHash) UpdateLoad(nodeID string, load int32) {
	c.Lock()
	defer c.Unlock()
	// 移除所有与该节点ID关联的虚拟节点
	c.removeVirtualNodes(nodeID)
	// 根据新的负载添加虚拟节点
	c.updateVirtualNodes(nodeID, load)
	// 由于添加或移除了节点，需要重新排序哈希环
	c.sortRing()
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
					load = MaxLoad // 如果获取负载失败，设置为最大负载
				}
				c.UpdateLoad(nodeID, load)
			}
		}
	}()
}

func (c *ConsistentHash) updateVirtualNodes(nodeID string, load int32) {
	// 确保负载不超过最大值
	if load > MaxLoad {
		load = MaxLoad
	}
	virtualNodes := MaxLoad - load // 虚拟节点数量 = 最大负载 - 节点负载, 负载越高 节点越少
	for i := 0; i < int(virtualNodes); i++ {
		virtualNodeKey := nodeID + strconv.Itoa(i)
		hashedKey := c.hashKey(virtualNodeKey)
		c.Nodes.Store(hashedKey, Node{ID: nodeID, Load: load})
		c.Ring = append(c.Ring, hashedKey) // 将虚拟节点的哈希值添加到哈希环中
	}
}

func (c *ConsistentHash) removeVirtualNodes(nodeID string) {
	for i := 0; i < MaxLoad; i++ {
		key := c.hashKey(nodeID + strconv.Itoa(i))
		c.Nodes.Delete(key)
	}
	newRing := make(HashRing, 0)
	for _, v := range c.Ring {
		if node, ok := c.Nodes.Load(v); ok && node.(Node).ID != nodeID {
			newRing = append(newRing, v)
		}
	}
	c.Ring = newRing
}

func (c *ConsistentHash) sortRing() {
	sort.Slice(c.Ring, func(i, j int) bool { return c.Ring[i] < c.Ring[j] })
}

func (c *ConsistentHash) hashKey(key string) uint64 {
	c.hash.Reset()
	c.hash.Write([]byte(key))
	return c.hash.Sum64()
}

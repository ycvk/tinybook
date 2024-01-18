package priorityqueue

import (
	"container/heap"
	"sync"
)

// HeapKind 指定堆类型 - 最小或最大
type HeapKind int

const (
	MinHeap HeapKind = iota //小顶堆
	MaxHeap                 //大顶堆
)

// Item 代表优先级队列中的一个元素
type Item[T comparable, V int64 | float64] struct {
	// 元素值
	Value T
	// 优先级
	Priority V
	// 更新时需要索引 - 用于在更新元素优先级时, 快速定位元素在堆中的位置
	index int
}

// PriorityQueue 基于容器/堆的优先级队列实现
type PriorityQueue[T comparable, V int64 | float64] struct {
	lock      sync.RWMutex
	items     []*Item[T, V]
	lookupMap map[T]int // 使用map来查找元素在堆中的位置，减少遍历时间
	kind      HeapKind
}

// New 创建一个新的优先级队列, 包含类型为T,优先级为V的项
func New[T comparable, V int64 | float64](kind HeapKind) *PriorityQueue[T, V] {
	pq := &PriorityQueue[T, V]{
		items:     make([]*Item[T, V], 0),
		lookupMap: make(map[T]int),
		kind:      kind,
	}
	heap.Init(pq) // 初始化堆
	return pq
}

// Len implements sort.Interface
func (pq *PriorityQueue[T, V]) Len() int {
	return len(pq.items)
}

// Less implements sort.Interface
func (pq *PriorityQueue[T, V]) Less(i, j int) bool {
	if pq.kind == MinHeap {
		return pq.items[i].Priority < pq.items[j].Priority
	}
	return pq.items[i].Priority > pq.items[j].Priority
}

// Swap implements sort.Interface
func (pq *PriorityQueue[T, V]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}

// Push implements heap.Interface
func (pq *PriorityQueue[T, V]) Push(x any) {
	n := len(pq.items)
	item := x.(*Item[T, V])
	item.index = n
	pq.items = append(pq.items, item)
}

// Pop implements heap.Interface
func (pq *PriorityQueue[T, V]) Pop() any {
	old := pq.items
	n := len(old)
	item := old[n-1]
	pq.lookupMap[item.Value] = -1 // 标记元素已被删除，避免内存泄漏
	delete(pq.lookupMap, item.Value)
	return item
}

// Put 将指定优先级的值添加到优先级队列中
func (pq *PriorityQueue[T, V]) Put(value T, priority V) {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	if _, ok := pq.lookupMap[value]; !ok { // 如果元素已存在，则不添加
		item := &Item[T, V]{
			Value:    value,
			Priority: priority,
		}
		pq.lookupMap[value] = len(pq.items)
		heap.Push(pq, item)
	}
}

// Get 返回优先级队列中的下一个元素
func (pq *PriorityQueue[T, V]) Get() *Item[T, V] {
	if pq.IsEmpty() {
		return nil
	}
	pq.lock.Lock()
	defer pq.lock.Unlock()
	item := heap.Pop(pq).(*Item[T, V])
	return item
}

// IsEmpty 返回一个布尔值,表示优先级队列是否为空
func (pq *PriorityQueue[T, V]) IsEmpty() bool {
	pq.lock.RLock()
	defer pq.lock.RUnlock()
	return len(pq.items) == 0
}

// Update 更新与给定值关联的优先级
func (pq *PriorityQueue[T, V]) Update(value T, priority V) bool {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	if index, ok := pq.lookupMap[value]; ok { // 如果元素存在，则更新优先级并调整堆
		pq.items[index].Priority = priority
		heap.Fix(pq, index)
	}
	return false
}

// Clear 清空优先级队列
func (pq *PriorityQueue[T, V]) Clear() {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	pq.items = nil
	pq.lookupMap = make(map[T]int)
}

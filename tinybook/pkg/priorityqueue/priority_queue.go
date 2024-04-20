package priorityqueue

import (
	"container/heap"
	"sync"
	"time"
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
	// 时间戳字段 - 用于记录元素的创建时间, 也用于在优先级相同时比较元素, 确保先入队的元素先出队
	time int64
}

// PriorityQueue 基于容器/堆的优先级队列实现
type PriorityQueue[T comparable, V int64 | float64] struct {
	lock      sync.RWMutex
	items     []*Item[T, V]
	lookupMap map[T]*Item[T, V]
	kind      HeapKind
}

// New 创建一个新的优先级队列, 包含类型为T,优先级为V的项
func New[T comparable, V int64 | float64](kind HeapKind) *PriorityQueue[T, V] {
	pq := &PriorityQueue[T, V]{
		items:     make([]*Item[T, V], 0),
		lookupMap: make(map[T]*Item[T, V]),
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
	// 当优先级相同的时候，比较时间戳
	if pq.items[i].Priority == pq.items[j].Priority {
		return pq.items[i].time < pq.items[j].time
	}
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
	old[n-1] = nil  // 最后一个元素置空 - 避免内存泄漏
	item.index = -1 // for safety
	pq.items = old[0 : n-1]
	return item
}

// Put 将指定优先级的值添加到优先级队列中
func (pq *PriorityQueue[T, V]) Put(value T, priority V) {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	item := &Item[T, V]{
		Value:    value,
		Priority: priority,
		time:     time.Now().UnixNano(),
	}
	pq.lookupMap[value] = item

	heap.Push(pq, item)
}

// PutItem 将指定优先级的值添加到优先级队列中
func (pq *PriorityQueue[T, V]) PutItem(item *Item[T, V]) {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	item.time = time.Now().UnixNano()
	pq.lookupMap[item.Value] = item
	heap.Push(pq, item)
}

// Get 返回优先级队列中的下一个元素
func (pq *PriorityQueue[T, V]) Get() *Item[T, V] {
	if pq.IsEmpty() {
		return nil
	}
	pq.lock.Lock()
	defer pq.lock.Unlock()
	item := heap.Pop(pq).(*Item[T, V])
	delete(pq.lookupMap, item.Value)
	return item
}

// IsEmpty 返回一个布尔值,表示优先级队列是否为空
func (pq *PriorityQueue[T, V]) IsEmpty() bool {
	pq.lock.RLock()
	defer pq.lock.RUnlock()
	return pq.Len() == 0
}

// Update 更新与给定值关联的优先级
func (pq *PriorityQueue[T, V]) Update(value T, priority V) bool {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	item, ok := pq.lookupMap[value]
	if ok {
		item.Priority = priority
		heap.Fix(pq, item.index)
	}
	return ok
}

// Clear 清空优先级队列
func (pq *PriorityQueue[T, V]) Clear() {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	pq.items = pq.items[:0]
	pq.lookupMap = make(map[T]*Item[T, V])
}

// Destroy 销毁优先级队列
func (pq *PriorityQueue[T, V]) Destroy() {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	pq.items = nil
	pq.lookupMap = nil
}

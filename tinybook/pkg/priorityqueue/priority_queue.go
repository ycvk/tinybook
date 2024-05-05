package priorityqueue

import (
	"container/heap"
	"golang.org/x/exp/constraints"
	"sync"
)

type PriorityQueue[T any, P constraints.Ordered] interface {
	Len() int
	Pop() any
	Put(value T, priority P)
	Get() *Item[T, P]
	GetAndPop() *Item[T, P]
	IsEmpty() bool
	Update(value T, priority P)
	Clear()
}

// HeapType 指定堆类型 - 最小或最大
type HeapType int

const (
	MinHeap HeapType = iota // 小顶堆
	MaxHeap                 // 大顶堆
)

// Item 代表优先级队列中的一个元素
type Item[T any, P constraints.Ordered] struct {
	Value    T // 元素值
	Priority P // 元素优先级
}

// HeapPriorityQueue 基于容器/堆的优先级队列实现
type HeapPriorityQueue[T comparable, P constraints.Ordered] struct {
	lock     sync.RWMutex
	items    []*Item[T, P]
	lookup   map[T]int // 元素值到索引的映射
	heapType HeapType
}

// New 创建一个新的优先级队列
func New[T comparable, P constraints.Ordered](kind HeapType) PriorityQueue[T, P] {
	pq := &HeapPriorityQueue[T, P]{
		items:    make([]*Item[T, P], 0),
		lookup:   make(map[T]int),
		heapType: kind,
	}
	return pq
}

// Len implements heap.Interface
func (pq *HeapPriorityQueue[T, P]) Len() int {
	return len(pq.items)
}

// Less implements heap.Interface
func (pq *HeapPriorityQueue[T, P]) Less(i, j int) bool {
	switch pq.heapType {
	case MinHeap:
		return pq.items[i].Priority < pq.items[j].Priority
	case MaxHeap:
		return pq.items[i].Priority > pq.items[j].Priority
	}
	return false // Should never reach here
}

// Swap implements heap.Interface
func (pq *HeapPriorityQueue[T, P]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.lookup[pq.items[i].Value] = i
	pq.lookup[pq.items[j].Value] = j
}

// Push implements heap.Interface
func (pq *HeapPriorityQueue[T, P]) Push(x any) {
	n := len(pq.items)
	item := x.(*Item[T, P])
	pq.lookup[item.Value] = n
	pq.items = append(pq.items, item)
}

// Pop implements heap.Interface 弹出最后一个元素
func (pq *HeapPriorityQueue[T, P]) Pop() any {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	n := len(pq.items)
	item := pq.items[n-1]
	pq.items = pq.items[:n-1]     // 直接截断 slice
	delete(pq.lookup, item.Value) // 立即从 lookupMap 中删除元素
	return item
}

// Put 将元素添加到优先级队列中
func (pq *HeapPriorityQueue[T, P]) Put(value T, priority P) {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	item := &Item[T, P]{Value: value, Priority: priority}
	heap.Push(pq, item)
	pq.lookup[value] = len(pq.items) - 1 // 更新 lookupMap
}

// Get 返回优先级队列中的下一个元素而不移除它
func (pq *HeapPriorityQueue[T, P]) Get() *Item[T, P] {
	pq.lock.RLock()
	defer pq.lock.RUnlock()
	return pq.items[0] // 直接返回堆顶元素
}

// GetAndPop 移除并返回优先级队列中的下一个元素
func (pq *HeapPriorityQueue[T, P]) GetAndPop() *Item[T, P] {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	return heap.Pop(pq).(*Item[T, P])
}

// IsEmpty 检查优先级队列是否为空
func (pq *HeapPriorityQueue[T, P]) IsEmpty() bool {
	pq.lock.RLock()
	defer pq.lock.RUnlock()
	return len(pq.items) == 0
}

// Update 更新元素的优先级
func (pq *HeapPriorityQueue[T, P]) Update(value T, priority P) {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	if index, ok := pq.lookup[value]; ok {
		pq.items[index].Priority = priority
		heap.Fix(pq, index) // 调整堆
	}
}

// Clear 清空优先级队列
func (pq *HeapPriorityQueue[T, P]) Clear() {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	pq.items = pq.items[:0] // 清空 slice 但保留其底层数组
	for k := range pq.lookup {
		delete(pq.lookup, k) // 清空 lookupMap
	}
}

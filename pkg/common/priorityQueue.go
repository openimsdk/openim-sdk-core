package common

import (
	"container/heap"
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// Event 定义事件结构
type Event struct {
	Priority int // 优先级，数字越大优先级越高
	Data     any // 事件内容
	Created  int64
	index    int // heap内部用
}

// eventHeap 实现heap.Interface
type eventHeap []*Event

func (h eventHeap) Len() int { return len(h) }
func (h eventHeap) Less(i, j int) bool {
	if h[i].Priority == h[j].Priority {
		return h[i].Created < h[j].Created // 同优先级时，Created 小的优先（先入先出）
	}
	return h[i].Priority > h[j].Priority
}
func (h eventHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i]; h[i].index = i; h[j].index = j }
func (h *eventHeap) Push(x any)   { *h = append(*h, x.(*Event)) }
func (h *eventHeap) Pop() any     { old := *h; n := len(old); x := old[n-1]; *h = old[0 : n-1]; return x }

// PriorityQueue 通用的线程安全优先队列
type PriorityQueue struct {
	mu           sync.Mutex
	cond         *sync.Cond
	events       eventHeap
	closed       bool
	capacity     int
	createdCount int64
}

// NewPriorityQueue 新建一个队列，capacity = 0 表示不限制
func NewPriorityQueue(capacity int) *PriorityQueue {
	q := &PriorityQueue{
		capacity: capacity,
	}
	q.cond = sync.NewCond(&q.mu)
	heap.Init(&q.events)
	return q
}

// Push 插入一个事件
func (q *PriorityQueue) Push(event *Event) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return errors.New("priority queue closed")
	}
	if q.capacity > 0 && len(q.events) >= q.capacity {
		return errors.New("priority queue is full")
	}

	event.Created = atomic.AddInt64(&q.createdCount, 1)

	heap.Push(&q.events, event)
	q.cond.Signal()
	return nil
}

// UpdatePriority 更新某个事件的优先级（并自动重新排序）
func (q *PriorityQueue) UpdatePriority(event *Event, newPriority int) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return errors.New("priority queue closed")
	}

	if event.index < 0 || event.index >= len(q.events) {
		return errors.New("invalid event index")
	}

	// 更新优先级
	event.Priority = newPriority
	// 重新维护堆
	heap.Fix(&q.events, event.index)
	q.cond.Broadcast()
	return nil
}

// PopWithContext 取出最高优先级的事件，支持ctx cancel
func (q *PriorityQueue) PopWithContext(ctx context.Context) (*Event, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.events) == 0 && !q.closed {
		// 使用 cond.Wait() 前，先检查 ctx 是否结束
		done := make(chan struct{})
		go func() {
			select {
			case <-ctx.Done():
				q.mu.Lock()
				q.cond.Signal()
				q.mu.Unlock()
			case <-done:
			}
		}()
		q.cond.Wait()
		close(done)

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	if q.closed && len(q.events) == 0 {
		return nil, errors.New("priority queue closed")
	}
	return heap.Pop(&q.events).(*Event), nil
}

// Close 关闭队列
func (q *PriorityQueue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.closed = true
	q.cond.Broadcast()
}

// Len 当前队列长度
func (q *PriorityQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.events)
}

// IsClosed 是否已关闭
func (q *PriorityQueue) IsClosed() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.closed
}

package common

import (
	"container/heap"
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Event defines the event structure
type Event struct {
	Priority int   // Priority: the higher the number, the higher the priority
	Data     any   // Event data
	Created  int64 // The incremented order in which the event was added to the queue
	index    int   // Internal use by heap to track event position
}

// eventHeap implements heap.Interface
type eventHeap []*Event

func (h eventHeap) Len() int { return len(h) }
func (h eventHeap) Less(i, j int) bool {
	if h[i].Priority == h[j].Priority {
		return h[i].Created < h[j].Created // FIFO for events with the same priority
	}
	return h[i].Priority > h[j].Priority
}
func (h eventHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i]; h[i].index = i; h[j].index = j }
func (h *eventHeap) Push(x any)   { *h = append(*h, x.(*Event)) }
func (h *eventHeap) Pop() any     { old := *h; n := len(old); x := old[n-1]; *h = old[0 : n-1]; return x }

// PriorityQueue represents a thread-safe priority queue
type PriorityQueue struct {
	mu           sync.Mutex
	cond         *sync.Cond
	events       eventHeap
	closed       bool
	capacity     int
	createdCount int64
}

// NewPriorityQueue creates a new priority queue with the specified capacity
// capacity = 0 means no limit on the number of events in the queue
func NewPriorityQueue(capacity int) *PriorityQueue {
	q := &PriorityQueue{
		capacity: capacity,
	}
	q.cond = sync.NewCond(&q.mu)
	heap.Init(&q.events)
	return q
}

// Push adds an event to the queue
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

func (q *PriorityQueue) PushWithContext(ctx context.Context, event *Event) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	for {
		if q.closed {
			return errors.New("priority queue closed")
		}
		if q.capacity == 0 || len(q.events) < q.capacity {
			event.Created = atomic.AddInt64(&q.createdCount, 1)
			heap.Push(&q.events, event)
			q.cond.Signal()
			return nil
		}

		timer := time.AfterFunc(50*time.Millisecond, func() {
			q.cond.Signal() // Wake up periodically
		})
		q.cond.Wait()
		timer.Stop()

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
}

// UpdatePriority updates the priority of an event and reorders the heap
func (q *PriorityQueue) UpdatePriority(event *Event, newPriority int) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return errors.New("priority queue closed")
	}

	if event.index < 0 || event.index >= len(q.events) {
		return errors.New("invalid event index")
	}

	// Update the event's priority
	event.Priority = newPriority
	// Re-heapify the queue to maintain the correct order
	heap.Fix(&q.events, event.index)
	q.cond.Broadcast()
	return nil
}

// PopWithContext removes and returns the highest-priority event, with context cancellation support
func (q *PriorityQueue) PopWithContext(ctx context.Context) (*Event, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.events) == 0 && !q.closed {
		// Wait for an event to be pushed into the queue, but check if context is canceled
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

// Close closes the priority queue, preventing any further operations
func (q *PriorityQueue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.closed = true
	q.cond.Broadcast()
}

// Len returns the current length of the queue
func (q *PriorityQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.events)
}

// IsClosed returns whether the queue has been closed
func (q *PriorityQueue) IsClosed() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.closed
}

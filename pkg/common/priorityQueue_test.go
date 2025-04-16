package common

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPriorityQueue_PushAndBlockingPop(t *testing.T) {
	q := NewPriorityQueue(0)

	_ = q.Push(&Event{Priority: 1, Data: "low1"})
	_ = q.Push(&Event{Priority: 1, Data: "low2"})
	_ = q.Push(&Event{Priority: 10, Data: "high1"})
	_ = q.Push(&Event{Priority: 12, Data: "high2"})
	_ = q.Push(&Event{Priority: 1, Data: "low3"})
	_ = q.Push(&Event{Priority: 9, Data: "mid"})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		event, err := q.PopWithContext(ctx)
		if err != nil {
			t.Logf("Exit pop loop, reason: %v", err)
			break
		}
		t.Logf("Got event: priority=%d, data=%v", event.Priority, event.Data)
	}
}

func TestPriorityQueue_FIFOOrder(t *testing.T) {
	q := NewPriorityQueue(0)

	_ = q.Push(&Event{Priority: 5, Data: "first"})
	_ = q.Push(&Event{Priority: 5, Data: "second"})
	_ = q.Push(&Event{Priority: 5, Data: "third"})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	e1, _ := q.PopWithContext(ctx)
	e2, _ := q.PopWithContext(ctx)
	e3, _ := q.PopWithContext(ctx)

	assert.Equal(t, "first", e1.Data)
	assert.Equal(t, "second", e2.Data)
	assert.Equal(t, "third", e3.Data)
}

func TestPriorityQueue_UpdatePriority(t *testing.T) {
	q := NewPriorityQueue(0)

	e1 := &Event{Priority: 1, Data: "low"}
	e2 := &Event{Priority: 5, Data: "mid"}
	_ = q.Push(e1)
	_ = q.Push(e2)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 提高低优先级事件的优先级
	err := q.UpdatePriority(e1, 10)
	assert.NoError(t, err)

	// 现在 e1 应该优先出来
	event, _ := q.PopWithContext(ctx)
	assert.Equal(t, "low", event.Data)

	event, _ = q.PopWithContext(ctx)
	assert.Equal(t, "mid", event.Data)
}

func TestPriorityQueue_ContextCancel(t *testing.T) {
	q := NewPriorityQueue(0)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// 没有数据，应该 ctx 超时退出
	event, err := q.PopWithContext(ctx)
	assert.Nil(t, event)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestPriorityQueue_Close(t *testing.T) {
	q := NewPriorityQueue(0)

	_ = q.Push(&Event{Priority: 1, Data: "data"})
	q.Close()

	assert.True(t, q.IsClosed())

	// 后续 push 会报错
	err := q.Push(&Event{Priority: 2, Data: "should fail"})
	assert.Error(t, err)

	// Pop 正常拿完后，再 pop 会报错
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	event, err := q.PopWithContext(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "data", event.Data)

	event, err = q.PopWithContext(ctx)
	assert.Nil(t, event)
	assert.Error(t, err)
}

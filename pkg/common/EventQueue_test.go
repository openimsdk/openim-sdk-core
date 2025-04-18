package common

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventQueue_ProduceAndConsume(t *testing.T) {
	var (
		mu           sync.Mutex
		handledOrder []string
	)

	handler := func(ctx context.Context, event *Event) {
		mu.Lock()
		handledOrder = append(handledOrder, event.Data.(string))
		mu.Unlock()
	}

	logCallback := func(msg string, fields ...any) {

	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eq := NewEventQueue(10)

	_, _ = eq.Produce("low1", 1)
	_, _ = eq.Produce("medium", 5)
	_, _ = eq.Produce("low2", 1)
	_, _ = eq.Produce("high", 10)

	go eq.ConsumeLoop(ctx, handler, logCallback)

	time.Sleep(500 * time.Millisecond)
	cancel()

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []string{"high", "medium", "low1", "low2"}, handledOrder)
}

func TestEventQueue_CancelContext(t *testing.T) {
	called := false

	handler := func(ctx context.Context, event *Event) {
		called = true
	}

	logCallback := func(msg string, fields ...any) {}

	ctx, cancel := context.WithCancel(context.Background())
	eq := NewEventQueue(1)

	go eq.ConsumeLoop(ctx, handler, logCallback)

	cancel()
	time.Sleep(100 * time.Millisecond)

	assert.False(t, called)
}

func TestEventQueue_UpdatePriority(t *testing.T) {
	var mu sync.Mutex
	var results []string

	handler := func(ctx context.Context, event *Event) {
		mu.Lock()
		results = append(results, event.Data.(string))
		mu.Unlock()
	}

	logCallback := func(msg string, fields ...any) {}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eq := NewEventQueue(5)
	go eq.ConsumeLoop(ctx, handler, logCallback)

	e1, _ := eq.Produce("a", 1)
	_, _ = eq.Produce("b", 2)

	time.Sleep(50 * time.Millisecond)
	_ = eq.UpdatePriority(e1, 5)

	time.Sleep(500 * time.Millisecond)
	cancel()

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []string{"a", "b"}, results)
}

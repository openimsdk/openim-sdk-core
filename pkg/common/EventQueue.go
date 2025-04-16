package common

import (
	"context"
)

type EventHandler func(ctx context.Context, event *Event)

type LogCallback func(msg string, fields ...any)

// EventQueue 负责并发消费 PriorityQueue 的事件
type EventQueue struct {
	queue *PriorityQueue
}

// NewEventQueue 创建 worker pool 实例
func NewEventQueue(capacity int) *EventQueue {
	return &EventQueue{
		queue: NewPriorityQueue(capacity),
	}
}

func (e *EventQueue) Produce(data any, priority int) (*Event, error) {
	event := &Event{Data: data, Priority: priority}
	err := e.queue.Push(event)
	return event, err
}

func (e *EventQueue) UpdatePriority(event *Event, newPriority int) error {
	return e.queue.UpdatePriority(event, newPriority)
}

func (e *EventQueue) ConsumeLoop(ctx context.Context, handle EventHandler, log LogCallback) {
	for {
		event, err := e.queue.PopWithContext(ctx)
		if err != nil {
			select {
			case <-ctx.Done():
				log("ctx canceled", err)
				return
			default:
				log("pop error", err)
				continue
			}
		}
		handle(ctx, event)
	}
}

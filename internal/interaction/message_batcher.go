package interaction

import (
	"context"
	"sync"
	"time"

	"github.com/openimsdk/protocol/sdkws"
)

const (
	maxBatchMessages     = 400
	minAggregationDelay  = 50 * time.Millisecond
	maxAggregationDelay  = time.Second
	lowLoadWindow        = 10 * time.Second
	lowLoadMessageLimit  = 20
	highLoadMessageLimit = 200
)

type arrivalRecord struct {
	ts    time.Time
	count int
}

type MessageBatcher struct {
	mutex         sync.Mutex
	buffer        *sdkws.PushMessages
	contexts      []context.Context
	handler       func([]context.Context, *sdkws.PushMessages)
	flushTimer    *time.Timer
	nextFlushAt   time.Time
	arrivals      []arrivalRecord
	recentTotal   int
	firstBuffered time.Time
}

func NewMessageBatcher(handler func([]context.Context, *sdkws.PushMessages)) *MessageBatcher {
	return &MessageBatcher{handler: handler}
}

// Close flushes buffered data (if any) and detaches the handler. Call when shutting down.
func (b *MessageBatcher) Close() {
	var (
		pending *sdkws.PushMessages
		ctxs    []context.Context
		handler func([]context.Context, *sdkws.PushMessages)
	)
	b.mutex.Lock()
	pending, ctxs = b.consumeLocked()
	b.cancelTimerLocked()
	handler = b.handler
	b.handler = nil
	b.mutex.Unlock()
	dispatch(handler, ctxs, pending)
}

func (b *MessageBatcher) dispatch(ctxs []context.Context, messages *sdkws.PushMessages) {
	if messages == nil || b.handler == nil {
		return
	}
	if len(messages.Msgs) == 0 && len(messages.NotificationMsgs) == 0 {
		return
	}
	b.handler(ctxs, messages)
}

func dispatch(handler func([]context.Context, *sdkws.PushMessages), ctxs []context.Context, messages *sdkws.PushMessages) {
	if handler == nil || messages == nil {
		return
	}
	if len(messages.Msgs) == 0 && len(messages.NotificationMsgs) == 0 {
		return
	}
	handler(ctxs, messages)
}

func (b *MessageBatcher) mergeLocked(batch *sdkws.PushMessages) int {
	if batch == nil {
		return 0
	}
	if b.buffer == nil {
		b.buffer = &sdkws.PushMessages{}
	}
	total := 0
	total += b.mergePullsLocked(batch.Msgs, true)
	total += b.mergePullsLocked(batch.NotificationMsgs, false)
	return total
}

func (b *MessageBatcher) mergePullsLocked(source map[string]*sdkws.PullMsgs, isMessage bool) int {
	if len(source) == 0 {
		return 0
	}

	var destination map[string]*sdkws.PullMsgs
	if isMessage {
		if b.buffer.Msgs == nil {
			b.buffer.Msgs = make(map[string]*sdkws.PullMsgs, len(source))
		}
		destination = b.buffer.Msgs
	} else {
		if b.buffer.NotificationMsgs == nil {
			b.buffer.NotificationMsgs = make(map[string]*sdkws.PullMsgs, len(source))
		}
		destination = b.buffer.NotificationMsgs
	}

	added := 0
	for conversationID, pulls := range source {
		if existing, ok := destination[conversationID]; ok {
			existing.Msgs = append(existing.Msgs, pulls.Msgs...)
			existing.IsEnd = pulls.IsEnd
			existing.EndSeq = pulls.EndSeq
		} else {
			destination[conversationID] = pulls
		}
		added += len(pulls.Msgs)
	}
	return added
}

func (b *MessageBatcher) consumeLocked() (*sdkws.PushMessages, []context.Context) {
	if b.buffer == nil {
		return nil, nil
	}
	pending := b.buffer
	ctxs := b.contexts
	b.buffer = nil
	b.contexts = nil
	b.firstBuffered = time.Time{}
	return pending, ctxs
}

func (b *MessageBatcher) pendingCountLocked() int {
	if b.buffer == nil {
		return 0
	}
	total := 0
	for _, pulls := range b.buffer.Msgs {
		total += len(pulls.Msgs)
	}
	for _, pulls := range b.buffer.NotificationMsgs {
		total += len(pulls.Msgs)
	}
	return total
}

func (b *MessageBatcher) ensureTimerLocked(target time.Time) {
	delay := time.Until(target)
	if delay <= 0 {
		delay = time.Millisecond
		target = time.Now().Add(delay)
	}
	if b.flushTimer == nil {
		b.nextFlushAt = target
		b.flushTimer = time.AfterFunc(delay, b.onTimer)
		return
	}
	if target.Equal(b.nextFlushAt) {
		return
	}
	b.flushTimer.Stop()
	b.nextFlushAt = target
	b.flushTimer = time.AfterFunc(delay, b.onTimer)
}

func (b *MessageBatcher) cancelTimerLocked() {
	if b.flushTimer == nil {
		return
	}
	b.flushTimer.Stop()
	b.flushTimer = nil
	b.nextFlushAt = time.Time{}
}

func (b *MessageBatcher) onTimer() {
	b.mutex.Lock()
	pending, ctxs := b.consumeLocked()
	b.flushTimer = nil
	b.nextFlushAt = time.Time{}
	b.mutex.Unlock()
	b.dispatch(ctxs, pending)
}

func (b *MessageBatcher) Enqueue(ctx context.Context, batch *sdkws.PushMessages) {
	var (
		toFlush *sdkws.PushMessages
		toCtxs  []context.Context
	)

	b.mutex.Lock()
	now := time.Now()
	addedCount := countMessages(batch)
	recent := b.recordArrivalLocked(now, addedCount)

	if recent < lowLoadMessageLimit {
		toFlush, toCtxs = b.consumeLocked()
		b.cancelTimerLocked()
		b.mutex.Unlock()
		b.dispatch(toCtxs, toFlush)
		b.dispatch([]context.Context{ctx}, batch)
		return
	}

	b.contexts = append(b.contexts, ctx)

	b.mergeLocked(batch)
	if b.buffer != nil && b.firstBuffered.IsZero() {
		b.firstBuffered = now
	}

	pendingCount := b.pendingCountLocked()
	if pendingCount >= maxBatchMessages && recent < highLoadMessageLimit {
		toFlush, toCtxs = b.consumeLocked()
		b.cancelTimerLocked()
	} else {
		if b.firstBuffered.IsZero() {
			b.firstBuffered = now
		}
		elapsed := now.Sub(b.firstBuffered)
		totalDelay := b.computeDelayLocked(recent)
		targetFlush := b.firstBuffered.Add(totalDelay)
		if elapsed >= maxAggregationDelay || elapsed >= totalDelay {
			toFlush, toCtxs = b.consumeLocked()
			b.cancelTimerLocked()
		} else {
			if now.After(targetFlush) {
				toFlush, toCtxs = b.consumeLocked()
				b.cancelTimerLocked()
			} else {
				b.ensureTimerLocked(targetFlush)
			}
		}
	}
	b.mutex.Unlock()

	b.dispatch(toCtxs, toFlush)
}

func (b *MessageBatcher) recordArrivalLocked(now time.Time, count int) int {
	if count <= 0 {
		return b.recentTotal
	}
	cutoff := now.Add(-lowLoadWindow)
	idx := 0
	for idx < len(b.arrivals) && b.arrivals[idx].ts.Before(cutoff) {
		b.recentTotal -= b.arrivals[idx].count
		idx++
	}
	if idx > 0 {
		copy(b.arrivals, b.arrivals[idx:])
		b.arrivals = b.arrivals[:len(b.arrivals)-idx]
	}
	b.arrivals = append(b.arrivals, arrivalRecord{ts: now, count: count})
	b.recentTotal += count
	return b.recentTotal
}

func (b *MessageBatcher) computeDelayLocked(recent int) time.Duration {
	if recent >= highLoadMessageLimit {
		return maxAggregationDelay
	}
	if recent <= lowLoadMessageLimit {
		return minAggregationDelay
	}
	span := highLoadMessageLimit - lowLoadMessageLimit
	scale := float64(recent-lowLoadMessageLimit) / float64(span)
	delay := minAggregationDelay + time.Duration(scale*float64(maxAggregationDelay-minAggregationDelay))
	return clampDuration(delay, minAggregationDelay, maxAggregationDelay)
}

func clampDuration(val, min, max time.Duration) time.Duration {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func countMessages(batch *sdkws.PushMessages) int {
	if batch == nil {
		return 0
	}
	total := 0
	for _, pulls := range batch.Msgs {
		total += len(pulls.Msgs)
	}
	for _, pulls := range batch.NotificationMsgs {
		total += len(pulls.Msgs)
	}
	return total
}

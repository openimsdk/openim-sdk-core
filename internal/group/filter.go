package group

import (
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2/simplelru"
)

// NotificationFilter deduplicates events by UUID, ensuring that
// the same UUID is only processed once within the specified timeout.
type NotificationFilter struct {
	lock    sync.Mutex
	data    *simplelru.LRU[string, time.Time]
	timeout time.Duration
}

// NewNotificationFilter creates a NotificationFilter with the given
// capacity (size) and timeout duration.
func NewNotificationFilter(size int, timeout time.Duration) *NotificationFilter {
	lru, err := simplelru.NewLRU[string, time.Time](size, nil)
	if err != nil {
		panic(err)
	}
	return &NotificationFilter{
		data:    lru,
		timeout: timeout,
	}
}

// ShouldExecute returns true if the UUID has not been processed
// within the timeout period. It also records the current time for this UUID.
// If the UUID was processed less than timeout ago, it returns false.
func (f *NotificationFilter) ShouldExecute(uuid string) bool {
	f.lock.Lock()
	defer f.lock.Unlock()

	now := time.Now()
	if last, exists := f.data.Get(uuid); exists && now.Sub(last) <= f.timeout {
		return false
	}
	f.data.Add(uuid, now)
	return true
}

// ExecuteIfNew calls fn only if the UUID has not been processed
// within the timeout period. Otherwise, it does nothing.
func (f *NotificationFilter) ExecuteIfNew(uuid string, fn func()) {
	if f.ShouldExecute(uuid) {
		fn()
	}
}

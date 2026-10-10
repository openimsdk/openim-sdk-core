package progress

import (
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"time"
)

func NewBar(name string, now, total int, ifRemove bool) *Bar {
	return &Bar{
		name:        name,
		now:         now,
		total:       total,
		ifRemove:    ifRemove,
		delayRemove: config.BarRemoveWaiteSec * time.Second,
	}
}

func NewRemoveBar(name string, now, total int) *Bar {
	return &Bar{
		name:        name,
		now:         now,
		total:       total,
		ifRemove:    true,
		delayRemove: config.BarRemoveWaiteSec * time.Second,
	}
}

type Bar struct {
	name         string
	now          int
	total        int
	completeTime time.Time
	delayRemove  time.Duration
	ifRemove     bool
}

func (b *Bar) shouldRemove() bool {
	if !b.ifRemove || b.now != b.total {
		return false
	}
	if b.completeTime.IsZero() {
		// first complete
		b.completeTime = time.Now()
	}
	if time.Since(b.completeTime) >= b.delayRemove {
		return true
	}
	return false
}

func (b *Bar) isDone() bool {
	return b.now == b.total
}

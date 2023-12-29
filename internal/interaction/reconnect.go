package interaction

import (
	"time"
)

type ReconnectStrategy interface {
	GetSleepInterval() time.Duration
	Reset()
}

type ExponentialRetry struct {
	attempts []int
	index    int
}

func NewExponentialRetry() *ExponentialRetry {
	return &ExponentialRetry{
		attempts: []int{1, 2, 4, 8, 16},
		index:    -1,
	}
}

func (rs *ExponentialRetry) GetSleepInterval() time.Duration {
	rs.index++
	interval := rs.index % len(rs.attempts)
	return time.Second * time.Duration(rs.attempts[interval])
}

func (rs *ExponentialRetry) Reset() {
	rs.index = -1
}

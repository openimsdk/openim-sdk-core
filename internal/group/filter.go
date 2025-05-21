package group

import (
	"crypto/md5"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2/simplelru"
)

type filterKey struct {
	ContentType int32
	Sum         [md5.Size]byte
}

func newNotificationFilter() *notificationFilter {
	data, err := simplelru.NewLRU[filterKey, time.Time](1024, nil)
	if err != nil {
		panic(err)
	}
	return &notificationFilter{
		data:    data,
		timeout: time.Second * 30,
	}
}

type notificationFilter struct {
	lock    sync.Mutex
	data    *simplelru.LRU[filterKey, time.Time]
	timeout time.Duration
}

func (x *notificationFilter) AlreadyExecuted(contentType int32, data []byte) bool {
	key := filterKey{
		ContentType: contentType,
		Sum:         md5.Sum(data),
	}
	x.lock.Lock()
	defer x.lock.Unlock()
	now := time.Now()
	if last, ok := x.data.Get(key); ok {
		duration := now.Sub(last)
		if duration >= 0 && duration <= x.timeout {
			return true
		}
	}
	x.data.Add(key, now)
	return false
}

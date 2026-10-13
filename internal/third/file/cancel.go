package file

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
)

func newCancelMap() *cancelMap {
	return &cancelMap{data: make(map[string]map[uint64]context.CancelCauseFunc)}
}

type cancelMap struct {
	mu   sync.Mutex
	data map[string]map[uint64]context.CancelCauseFunc
	incr atomic.Uint64
}

func (m *cancelMap) WithCancel(ctx context.Context, cancelID string) (context.Context, context.CancelFunc) {
	value := m.incr.Add(1)
	ctx, cancel := context.WithCancelCause(ctx)
	m.mu.Lock()
	if m.data[cancelID] == nil {
		m.data[cancelID] = make(map[uint64]context.CancelCauseFunc)
	}
	m.data[cancelID][value] = cancel
	m.mu.Unlock()
	return ctx, func() {
		cancel(context.Canceled)
		m.mu.Lock()
		defer m.mu.Unlock()
		delete(m.data[cancelID], value)
		if len(m.data[cancelID]) == 0 {
			delete(m.data, cancelID)
		}
	}
}

func (m *cancelMap) CancelID(cancelID string, err error) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	items := m.data[cancelID]
	for _, cancel := range items {
		cancel(err)
	}
	delete(m.data, cancelID)
	return len(items)
}

func (f *File) WithUploadCancel(ctx context.Context, cancelID string) (context.Context, context.CancelFunc) {
	if cancelID == "" {
		return ctx, func() {}
	}
	return f.cancel.WithCancel(ctx, cancelID)
}

func (f *File) CancelUpload(_ context.Context, cancelID string) int {
	if cancelID == "" {
		return 0
	}
	return f.cancel.CancelID(cancelID, sdkerrs.ErrCancelUploadingFile.Wrap())
}

package utils

import "sync"

// resource manager
type LockPool struct {
	maxLocks int
	locks    map[string]struct{}

	globalMutex sync.Mutex
	cond        *sync.Cond
}

// Usage

func NewLockPool(maxLocks int) *LockPool {
	lp := &LockPool{
		maxLocks: maxLocks,
		locks:    make(map[string]struct{}),
	}

	lp.cond = sync.NewCond(&lp.globalMutex)
	return lp
}

func (lp *LockPool) Lock(key string) {
	lp.globalMutex.Lock()
	defer lp.globalMutex.Unlock()

	for {
		// If the key exists in the map, then wait.
		if _, exists := lp.locks[key]; exists {
			lp.cond.Wait()
		} else if len(lp.locks) >= lp.maxLocks {
			// If the number of locks reaches the maximum value, then wait.
			lp.cond.Wait()
		} else {
			lp.locks[key] = struct{}{}
			return
		}
	}
}

func (lp *LockPool) Unlock(key string) {
	lp.globalMutex.Lock()
	defer lp.globalMutex.Unlock()

	delete(lp.locks, key)
	// wakes all goroutines waiting
	lp.cond.Broadcast()
}

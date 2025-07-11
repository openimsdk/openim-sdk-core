package group

import (
	"sync"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	var (
		wg   sync.WaitGroup
		lock sync.Mutex
	)

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !lock.TryLock() {
				t.Log("TryLock=false")
				return
			}
			t.Log("TryLock=true")
			time.Sleep(time.Second)
			lock.Unlock()
		}()
	}
	wg.Wait()
}

// Package reerrgroup is a rewrite of errgroup
package reerrgroup

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

const (
	checkTaskDoneMilSec = 100 // Check task completion interval (unit: Milliseconds)
)

// A Group is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// A zero Group is valid, has no limit on the number of active goroutines,
// and does not cancel on error.
type Group struct {
	cancel func(error)
	wg     sync.WaitGroup // worker wg

	taskChan  chan func() error
	taskCount atomic.Int64

	errOnce sync.Once
	err     error

	beforeTasks []func() error
	afterTasks  []func() error
}

func (g *Group) done() {

	g.wg.Done()
}

// WithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func WithContext(ctx context.Context, workerCount int) (*Group, context.Context) {
	g := &Group{}
	g.initialize(workerCount)
	ctx, cancel := context.WithCancelCause(ctx)
	g.cancel = cancel
	return g, ctx
}

func (g *Group) initialize(workerCount int) {
	g.taskChan = make(chan func() error, workerCount) // Initialize the channel with the provided buffer size.
	// Start multiple goroutines based on the specified workerCount.
	for i := 0; i < workerCount; i++ {
		g.wg.Add(1)
		go func() {
			defer g.wg.Done()
			for task := range g.taskChan {
				doTask := func() error {
					defer g.taskCount.Add(-1)
					tasks := append(append(g.beforeTasks, task), g.afterTasks...)
					for _, t := range tasks { // Execute the function
						if g.err != nil {
							return nil
						}
						if err := t(); err != nil {
							return err
						}
					}
					return nil
				}
				if err := doTask(); err != nil {
					g.errOnce.Do(func() {
						g.err = err
					})
					if g.cancel != nil {
						g.cancel(g.err)
					}
				}

			}
		}()
	}
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	close(g.taskChan)

	g.wg.Wait()

	if g.cancel != nil {
		g.cancel(g.err)
	}
	return g.err
}

// WaitTaskDone only wait all task done without cancel ctx and close taskChan.
func (g *Group) WaitTaskDone() {
	ticker := time.NewTicker(checkTaskDoneMilSec * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if g.taskCount.Load() == 0 {
				return
			}
		}
	}
}

// Go calls the given function in a new goroutine.
// It blocks until the new goroutine can be added without the number of
// active goroutines in the group exceeding the configured limit.
//
// The first call to return a non-nil error cancels the group's context, if the
// group was created by calling WithContext. The error will be returned by Wait.
func (g *Group) Go(f func() error) {
	if g.err != nil {
		return
	}

	select {
	case g.taskChan <- f:
		g.taskCount.Add(1)
		return
	}
}

// TryGo calls the given function in a new goroutine only if the number of
// active goroutines in the group is currently below the configured limit.
//
// The return value reports whether the goroutine was started.
func (g *Group) TryGo(f func() error) bool {
	if g.err != nil {
		return false
	}

	select {
	case g.taskChan <- f:
		g.taskCount.Add(1)
		return true
	default:
		return false
	}
}

func (g *Group) SetBeforeTasks(f ...func() error) {
	g.beforeTasks = append(g.beforeTasks, f...)
}

func (g *Group) SetAfterTasks(f ...func() error) {
	g.afterTasks = append(g.afterTasks, f...)
}

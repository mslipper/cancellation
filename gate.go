package cancellation

import (
	"sync"
	"sync/atomic"
)

// Gate is an object whose Cancel method blocks until
// all actions enqueued with the Do method return.
// For example:
//
//     g.Do(func() {
//	       // do some long running task
//     })
//     g.Cancel() // will block until the g.Do above returns.
//
// This is useful in cases where a producer of values continues
// producing values even after they aren't needed anymore.
type Gate struct {
	mtx  sync.RWMutex
	done uint32
}

// Do calls the function f unless the Gate is cancelled.
// Calling Do after the Gate is cancelled is a no-op.
func (c *Gate) Do(f func()) {
	if atomic.LoadUint32(&c.done) == 1 {
		return
	}

	c.mtx.RLock()
	if c.done == 0 {
		f()
	}
	c.mtx.RUnlock()
}

// Cancel cancels the Gate. Cancel will wait for all
// currently-executing actions enqueued by Do to return
// before marking the Gate as cancelled.
func (c *Gate) Cancel() {
	if atomic.LoadUint32(&c.done) == 1 {
		return
	}

	c.mtx.Lock()
	atomic.StoreUint32(&c.done, 1)
	c.mtx.Unlock()
}

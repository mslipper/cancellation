package cancellation

import "sync"

// Once represents a cancellable action that can be waited on.
type Once struct {
	ch   chan error
	once sync.Once
}

// NewOnce returns a new Once instance.
func NewOnce() *Once {
	return &Once{
		ch: make(chan error, 1),
	}
}

// Cancel cancels the Once. It can be called multiple times,
// however all calls after the first are no-ops.
func (c *Once) Cancel(err error) {
	c.once.Do(func() {
		c.ch <- err
	})
}

// Wait blocks until the Once is cancelled.
func (c *Once) Wait() error {
	return <-c.ch
}

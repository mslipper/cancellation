package cancellation

import (
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

func TestGroup_Cancelled(t *testing.T) {
	var g Group
	var callCount uint32
	g.Do(func(killed <-chan struct{}) {
		atomic.AddUint32(&callCount, 1)
		<-killed
	})
	g.Do(func(killed <-chan struct{}) {
		atomic.AddUint32(&callCount, 1)
	})
	g.Do(func(killed <-chan struct{}) {
		<-killed
		time.Sleep(10 * time.Millisecond)
		atomic.AddUint32(&callCount, 1)
	})

	g.Cancel()
	g.Wait()
	assert.EqualValues(t, 3, atomic.LoadUint32(&callCount))
}

func TestGroup_AllExited(t *testing.T) {
	var g Group
	var callCount uint32
	g.Do(func(killed <-chan struct{}) {
		atomic.AddUint32(&callCount, 1)
	})
	g.Do(func(killed <-chan struct{}) {
		time.Sleep(10 * time.Millisecond)
		atomic.AddUint32(&callCount, 1)
	})

	g.Wait()
	assert.EqualValues(t, 2, atomic.LoadUint32(&callCount))
}

func TestGroup_Panics(t *testing.T) {
	var g Group
	var callCount uint32
	g.Do(func(killed <-chan struct{}) {
		atomic.AddUint32(&callCount, 1)
	})

	g.Cancel()
	assert.Panics(t, func() {
		g.Cancel()
	})
	assert.Panics(t, func() {
		g.Do(func(killed <-chan struct{}) {
			atomic.AddUint32(&callCount, 1)
		})
	})

	assert.EqualValues(t, 1, atomic.LoadUint32(&callCount))
}

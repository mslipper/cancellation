package cancellation

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGate(t *testing.T) {
	var g Gate
	var callCount uint32
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go g.Do(func() {
			time.Sleep(10 * time.Millisecond)
			atomic.AddUint32(&callCount, 1)
			wg.Done()
		})
	}
	wg.Wait()
	wg.Add(10)
	g.Cancel()
	for i := 0; i < 10; i++ {
		go func() {
			g.Do(func() {
				atomic.AddUint32(&callCount, 1)
			})
			wg.Done()
		}()
	}
	assert.EqualValues(t, 10, atomic.LoadUint32(&callCount))
}

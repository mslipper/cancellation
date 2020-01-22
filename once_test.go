package cancellation

import (
	"math/rand"
	"testing"
)

func TestOnce(t *testing.T) {
	once := NewOnce()
	doneCh := make(chan struct{})
	go func() {
		for i := 0; i < rand.Intn(100); i++ {
			once.Cancel()
			doneCh <- struct{}{}
		}
	}()

	once.Wait()
	<-doneCh
}
